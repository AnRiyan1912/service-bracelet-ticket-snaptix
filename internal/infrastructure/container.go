package infrastructure

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"

	braceletqrcode "bracelet-ticket-system-be/internal/bracelet-qr-code"
	braceletticket "bracelet-ticket-system-be/internal/bracelet-ticket"
	braceletticketcategory "bracelet-ticket-system-be/internal/bracelet-ticket-category"
	"bracelet-ticket-system-be/internal/config"
	"bracelet-ticket-system-be/internal/domain"
	"bracelet-ticket-system-be/internal/event"
	"bracelet-ticket-system-be/internal/ticket"
	"bracelet-ticket-system-be/pkg/xlogger"
)

var (
	cfg                   config.Config
	braceletTicketService domain.BraceletTicketService
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
	mysqlBraceletTicketRepository := braceletticket.NewMysqlBraceletTicketRepository(db)
	mysqlBraceletCategoryRepository := braceletticketcategory.NewMysqlBraceletTicketCategoryRepository(db)
	mysqlTicketRepository := ticket.NewMysqlTicketRepository(db)
	mysqlEventRepository := event.NewMysqlEventRepository(db)
	braceletQrCodeService := braceletqrcode.BraceletQrCodeService(&cfg)
	braceletTicketService = braceletticket.NewBraceletTicketService(mysqlBraceletTicketRepository, mysqlBraceletCategoryRepository, mysqlTicketRepository, mysqlEventRepository, braceletQrCodeService)
}
