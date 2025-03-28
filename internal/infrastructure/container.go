package infrastructure

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"

	braceletticket "bracelet-ticket-system-be/internal/bracelet-ticket"
	braceletticketcategory "bracelet-ticket-system-be/internal/bracelet-ticket-category"
	braceletticketexel "bracelet-ticket-system-be/internal/bracelet-ticket-exel"
	"bracelet-ticket-system-be/internal/config"
	"bracelet-ticket-system-be/internal/domain"
	"bracelet-ticket-system-be/internal/event"
	eventbraceletticketcategory "bracelet-ticket-system-be/internal/event-bracelet-ticket-category"
	"bracelet-ticket-system-be/internal/ticket"
	"bracelet-ticket-system-be/pkg/xlogger"
)

var (
	cfg                   config.Config
	braceletTicketService domain.BraceletTicketService
	mysqlEventRepository  domain.MysqlEventRepository
	redisEventRepository  domain.RedisEventRepository
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	xlogger.Setup(cfg)
	db, error := dbSetup()
	if db == nil {
		xlogger.Logger.Error().Err(error).Msg("Failed connect to database!")
	}
	redisSetup()

	mysqlBraceletTicketRepository := braceletticket.NewMysqlBraceletTicketRepository(db)
	mysqlBraceletCategoryRepository := braceletticketcategory.NewMysqlBraceletTicketCategoryRepository(db)
	mysqlTicketRepository := ticket.NewMysqlTicketRepository(db)
	mysqlEventRepository = event.NewMysqlEventRepository(db)
	redisEventRepository = event.NewRedisEventRepository(redisClient)
	mysqlEventBraceletTicketCategoryRepository := eventbraceletticketcategory.NewMysqlEventBraceletTicketCategoryRepository(db)
	mysqlBraceletTicketExelRepository := braceletticketexel.NewMysqlBraceletTicketExelRepository(db)
	braceletTicketExelService := braceletticketexel.NewBraceletTicketExelService(mysqlBraceletTicketExelRepository)
	braceletTicketService = braceletticket.NewBraceletTicketService(mysqlBraceletTicketRepository, mysqlBraceletCategoryRepository, mysqlTicketRepository, mysqlEventRepository, mysqlEventBraceletTicketCategoryRepository, redisEventRepository, braceletTicketExelService, &cfg)
}
