package uf

import (
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/valyala/fasthttp"
)

type Pile struct {
	Conf      Config
	Ctx       *fasthttp.RequestCtx
	Pool      *pgxpool.Pool
	TaxTable  TaxTable
	LogMailer *Mailer
}

func (pile *Pile) InitFields(ufName string) {
	if ufName == "order" {
		pile.TaxTable = createTaxTable()
	}

	pile.Conf.InterpretDefaults(ufName)

	pile.LogMailer = &logMailer
	pile.LogMailer.smtpHost = pile.Conf.SmtpHost
	pile.LogMailer.smtpPort = pile.Conf.SmtpPort
	pile.LogMailer.smtpUser = pile.Conf.SmtpUser
	pile.LogMailer.smtpPass = pile.Conf.SmtpPass

	var err error = nil
	// want to set local global loggingLevel
	loggingLevel, err = strconv.Atoi(pile.Conf.LoggingLevel)
	if err != nil {
		Panic("failed to convert logging level string")
	}

	if loggingLevel > 50 {
		Panic("LoggingLevel must be less than 50")
	}

	telegramBot.BotId = pile.Conf.TelegramBotId
	telegramBot.ChannelId, err = strconv.ParseInt(
		pile.Conf.TelegramChannelId,
		10,
		64)
	if err != nil {
		Panic("Telegram ChannelId must parse to int64")
	}
}
