package infrastructure

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/redis/go-redis/v9"

	"bracelet-ticket-system-be/internal/domain"
	"bracelet-ticket-system-be/pkg/xlogger"
)

type Client struct {
	Conn    *websocket.Conn
	EventID string
}

func runWebsocket(fiberServer *fiber.App, mysqlEventRepository domain.MysqlEventRepository, redisEventRepository domain.RedisEventRepository, redisClient redis.UniversalClient) {
	logger := xlogger.Logger
	ctx := context.Background()
	clients := make(map[*websocket.Conn]string)

	fiberServer.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	fiberServer.Get("/ws", websocket.New(func(c *websocket.Conn) {
		logger.Info().Msgf("Client connected: %s", c.RemoteAddr().String())
		defer func() {
			logger.Info().Msgf("Client disconnected: %s", c.RemoteAddr().String())
			delete(clients, c)
			c.Close()
		}()

		queryParams := c.Query("eventId")
		if queryParams == "" {
			logger.Error().Msg("Client must provide eventId")
			c.Close()
			return
		}
		clients[c] = queryParams
		logger.Info().Msgf("Client joined room: %s", queryParams)

		var eventBraceletTicket domain.EventTotalBraceletAndCheckInBraceletTicket

		// Send initial stock
		findBraceletTicketStockInRedis, err := redisEventRepository.FindTotalAndTotalCheckInBraceletTicketByEventId(queryParams)
		if errors.Is(err, redis.Nil) {
			findBraceletTicketStockInMysql, err := mysqlEventRepository.FindTotalAndTotalCheckInBraceletTicketByEventId(queryParams)
			if err != nil {
				logger.Err(err).Msg("Failed to get total bracelet ticket and total checkIn")
			}
			// insert to redis when not yet
			err = redisEventRepository.InsertTotalAndTotalCheckInBraceletTicketByEventId(queryParams, findBraceletTicketStockInMysql.TotalBraceletTicket, findBraceletTicketStockInMysql.TotalCheckInBraceletTicket)
			if err != nil {
				logger.Err(err).Msg("Failed to insert data redis when not yet insert")
			}
			eventBraceletTicket.TotalBraceletTicket = findBraceletTicketStockInMysql.TotalBraceletTicket
			eventBraceletTicket.TotalCheckInBraceletTicket = findBraceletTicketStockInMysql.TotalCheckInBraceletTicket
		} else {
			eventBraceletTicket.TotalBraceletTicket = findBraceletTicketStockInRedis.TotalBraceletTicket
			eventBraceletTicket.TotalCheckInBraceletTicket = findBraceletTicketStockInRedis.TotalCheckInBraceletTicket
		}
		c.WriteJSON(fiber.Map{
			"event":                      "initialData",
			"totalBraceletTicket":        eventBraceletTicket.TotalBraceletTicket,
			"totalCheckInBraceletTicket": eventBraceletTicket.TotalCheckInBraceletTicket,
		})

		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				break
			}
		}
	}))

	// Redis subscriber to broadcast updates
	go func() {
		pubsub := redisClient.Subscribe(ctx, "update_check_in_bracelet_ticket:*")
		defer pubsub.Close()

		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to receive message from Redis Pub/Sub")
				continue
			}

			parts := strings.Split(msg.Channel, ":")
			if len(parts) < 2 {
				continue
			}
			eventId := parts[1]
			stock, _ := strconv.Atoi(msg.Payload)

			logger.Info().Msgf("Broadcasting new stock for event %s: %d", eventId, stock)
			for client, eid := range clients {
				if eid == eventId {
					client.WriteJSON(fiber.Map{"event": "updateTotalCheckIn", "totalCheckIn": stock})
				}
			}
		}
	}()
}
