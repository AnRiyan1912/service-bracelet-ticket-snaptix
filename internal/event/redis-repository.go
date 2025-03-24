package event

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"

	"bracelet-ticket-system-be/internal/constan"
	"bracelet-ticket-system-be/internal/domain"
	"bracelet-ticket-system-be/pkg/xlogger"
)

type RedisEventRepository struct {
	client redis.UniversalClient
}

func NewRedisEventRepository(client redis.UniversalClient) domain.RedisEventRepository {
	return RedisEventRepository{client: client}
}

// InsertTotalAndTotalCheckInBraceletTicketByEventId implements domain.RedisEventRepository.
func (r RedisEventRepository) InsertTotalAndTotalCheckInBraceletTicketByEventId(eventId string, totalBraceletTicket int, totalCheckInBraceletTicket int) error {
	xlogger := xlogger.Logger
	ctx := context.Background()

	keyTotalBraceletTicket := constan.BRACELET_STOCK_EVENT + eventId
	keyTotalCheckInBraceletTicket := constan.BRACELET_CHECK_IN + eventId

	_, err := r.client.Set(ctx, keyTotalBraceletTicket, totalBraceletTicket, 0).Result()
	if err != nil {
		xlogger.Error().Err(err).Msg("Failed to insert total bracelet ticket")
		return err
	}
	_, err = r.client.Set(ctx, keyTotalCheckInBraceletTicket, totalCheckInBraceletTicket, 0).Result()
	if err != nil {
		xlogger.Error().Err(err).Msg("Failed to insert total bracelet ticket")
		return err
	}
	return nil
}

// FindTotalAndTotalCheckInBraceletTicketByEventId implements domain.RedisEventRepository.
func (r RedisEventRepository) FindTotalAndTotalCheckInBraceletTicketByEventId(eventId string) (*domain.EventTotalBraceletAndCheckInBraceletTicket, error) {
	xlogger := xlogger.Logger

	keyTotalBraceletTicket := constan.BRACELET_STOCK_EVENT + eventId
	keyTotalCheckInBraceletTicket := constan.BRACELET_CHECK_IN + eventId

	totalBraceletTicketStr, err := r.client.Get(context.Background(), keyTotalBraceletTicket).Result()
	if err != nil {
		xlogger.Error().Err(err).Msg("Failed to find total bracelet ticket")
		return nil, err
	}
	totalBraceletTicket, err := strconv.Atoi(totalBraceletTicketStr)
	if err != nil {
		xlogger.Error().Err(err).Msg("Failed to convert total bracelet ticket to int")
		return nil, err
	}

	totalCheckInBraceletTicketStr, err := r.client.Get(context.Background(), keyTotalCheckInBraceletTicket).Result()
	if err != nil {
		xlogger.Error().Err(err).Msg("Failed to find total check in bracelet ticket")
		return nil, err
	}
	totalCheckInBraceletTicket, err := strconv.Atoi(totalCheckInBraceletTicketStr)
	if err != nil {
		xlogger.Error().Err(err).Msg("Failed to convert total check in bracelet ticket to int")
		return nil, err
	}

	return &domain.EventTotalBraceletAndCheckInBraceletTicket{
		TotalBraceletTicket:        totalBraceletTicket,
		TotalCheckInBraceletTicket: totalCheckInBraceletTicket,
	}, nil

}

// UpdateTotalCheckInBraceletTicketByEventId implements domain.RedisEventRepository.
func (r RedisEventRepository) UpdateTotalCheckInBraceletTicketByEventId(eventId string, totalCheckInBraceletTicket int) error {
	xlogger := xlogger.Logger
	ctx := context.Background()

	keyTotalCheckIn := constan.BRACELET_CHECK_IN + eventId
	pubSubChannel := "update_check_in_bracelet_ticket:" + eventId

	err := r.client.Watch(ctx, func(tx *redis.Tx) error {
		currentCheckInStr, err := tx.Get(ctx, keyTotalCheckIn).Result()
		if err != nil && err != redis.Nil {
			xlogger.Error().Err(err).Msg("Failed to get total check-in bracelet ticket")
			return err
		}

		currentCheckIn := 0
		if currentCheckInStr != "" {
			currentCheckIn, _ = strconv.Atoi(currentCheckInStr)
		}

		newTotalCheckIn := currentCheckIn + totalCheckInBraceletTicket

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, keyTotalCheckIn, newTotalCheckIn, 0)

			pipe.Publish(ctx, pubSubChannel, newTotalCheckIn)
			return nil
		})

		return err
	}, keyTotalCheckIn)

	if err != nil {
		xlogger.Error().Err(err).Msg("Failed to update total check-in bracelet ticket")
		return err
	}

	return nil
}
