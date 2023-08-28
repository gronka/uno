package uf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v73"
	//log "github.com/sirupsen/logrus"
)

type Config struct {
	//-- This config is set at runtime
	ConnectCassandra bool
	Environment      string
	UfName           string
	//-- End of config set at runtime

	//-- This config set programmatically:
	BindTo       string
	JsonLoc      string
	StripePrefix string

	AimAddress    string
	BorderAddress string
	MahaAddress   string
	OrderAddress  string
	PublicAddress string
	UserAddress   string

	AimDomain    string
	BorderDomain string
	MahaDomain   string
	OrderDomain  string
	PublicDomain string
	UserDomain   string

	AimPort    string
	BorderPort string
	MahaPort   string
	OrderPort  string
	PublicPort string
	UserPort   string

	AimPgDsn   string
	MahaPgDsn  string
	OrderPgDsn string
	UserPgDsn  string
	//-- End of programmatic config

	//-- This config is set in Json
	AmazonEmail    string
	AmazonPassword string
	ApiDn          string
	AppDn          string
	FridayyPhone   string
	HushJwt        string

	Mapsac    string
	OpenAiKey string

	SendBlueApiKeyId      string // header to SendBlue
	SendBlueApiSecret     string // header to SendBlue
	SendBlueHookSecret    string // header from SendBlue messages
	LoggingLevel          string
	LoopAuthorization     string
	LoopHookAuthorization string
	LoopSecretKey         string
	SipSecret             string
	SipHookSecret         string
	SipUrl                string

	SmtpHost          string
	SmtpPort          string
	SmtpUser          string
	SmtpPass          string
	StripePublicKey   string
	StripeSecretKey   string
	TelegramBotId     string
	TelegramChannelId string
	UspsUsername      string
	UspsPassword      string
	UfKey             string
	UfScheme          string

	//Loglevel
	//CassandraIps
	//CassandraKeyspace
	ZincCardName            string
	ZincCardNumber          string
	ZincCardSecurityCode    string
	ZincCardExpirationMonth int
	ZincCardExpirationYear  int
	ZincToken               string
	//-- End of config set in Json
}

func (conf *Config) InterpretDefaults(ufName string) {
	conf.UfName = ufName
	conf.Environment = os.Getenv("Environment")
	if conf.Environment != "local" && conf.Environment != "prod" {
		panic(errors.New("invalid conf.Environment: " + conf.Environment))
	}
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	conf.AimPort = "8030"
	conf.BorderPort = "8010"
	conf.MahaPort = "8060"
	conf.OrderPort = "8040"
	conf.PublicPort = "8050"
	conf.UserPort = "8020"

	conf.AimDomain = "uf-aim"
	conf.BorderDomain = "uf-border"
	conf.MahaDomain = "uf-maha"
	conf.OrderDomain = "uf-order"
	conf.PublicDomain = "uf-public"
	conf.UserDomain = "uf-user"

	conf.AimAddress = fmt.Sprintf("%s:%s", conf.AimDomain, conf.AimPort)
	conf.BorderAddress = fmt.Sprintf("%s:%s", conf.BorderDomain, conf.BorderPort)
	conf.MahaAddress = fmt.Sprintf("%s:%s", conf.MahaDomain, conf.MahaPort)
	conf.OrderAddress = fmt.Sprintf("%s:%s", conf.OrderDomain, conf.OrderPort)
	conf.PublicAddress = fmt.Sprintf("%s:%s", conf.PublicDomain, conf.PublicPort)
	conf.UserAddress = fmt.Sprintf("%s:%s", conf.UserDomain, conf.UserPort)

	switch conf.UfName {
	case "aim":
		conf.BindTo = "0.0.0.0:" + conf.AimPort
	case "border":
		conf.BindTo = "0.0.0.0:" + conf.BorderPort
	case "maha":
		conf.BindTo = "0.0.0.0:" + conf.MahaPort
	case "order":
		conf.BindTo = "0.0.0.0:" + conf.OrderPort
	case "public":
		conf.BindTo = "0.0.0.0:" + conf.PublicPort
	case "user":
		conf.BindTo = "0.0.0.0:" + conf.UserPort
	default:
		panic(errors.New("invalid config name"))
	}

	pgUser := "postgres://postgres:postgres@"
	if conf.Environment == "local" {
		conf.readConfFromEnvJson()

		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		conf.AimPgDsn = pgUser + "uf-aim:5430/uf_aim"
		conf.MahaPgDsn = pgUser + "uf-maha:5460/uf_maha"
		conf.OrderPgDsn = pgUser + "uf-order:5440/uf_order"
		conf.UserPgDsn = pgUser + "uf-user:5420/uf_user"
	} else {
		conf.readConfFromEnvironment()

		//TODO change to json when we have a log aggregator
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		conf.AimPgDsn = pgUser + "uf-aim-pg:5432/uf_aim"
		conf.MahaPgDsn = pgUser + "uf-maha-pg:5432/uf_maha"
		conf.OrderPgDsn = pgUser + "uf-order-pg:5432/uf_order"
		conf.UserPgDsn = pgUser + "uf-user-pg:5432/uf_user"
	}

	conf.validateValues()

	Debug(conf.Environment)
	Debug(conf.ApiDn)
	Debug(conf.UfScheme)

	Debug("config loaded for " + conf.UfName)
	Debug("BindTo: " + conf.BindTo)

	stripe.Key = conf.StripeSecretKey
}

// we use env variables via a configMap when deploying to prod
func (conf *Config) readConfFromEnvironment() {
	var err error
	conf.AmazonEmail = os.Getenv("AmazonEmail")
	conf.AmazonPassword = os.Getenv("AmazonPassword")
	conf.ApiDn = os.Getenv("ApiDn")
	conf.AppDn = os.Getenv("AppDn")
	conf.FridayyPhone = os.Getenv("FridayyPhone")
	conf.HushJwt = os.Getenv("HushJwt")
	conf.LoggingLevel = os.Getenv("LoggingLevel")
	conf.LoopAuthorization = os.Getenv("LoopAuthorization")
	conf.LoopHookAuthorization = os.Getenv("LoopHookAuthorization")
	conf.LoopSecretKey = os.Getenv("LoopSecretKey")
	conf.OpenAiKey = os.Getenv("OpenAiKey")

	conf.SendBlueApiSecret = os.Getenv("SendBlueApiSecret")
	conf.SendBlueApiKeyId = os.Getenv("SendBlueApiKeyId")
	conf.SendBlueHookSecret = os.Getenv("SendBlueHookSecret")
	conf.SipSecret = os.Getenv("SipSecret")
	conf.SipHookSecret = os.Getenv("SipHookSecret")
	conf.SipUrl = os.Getenv("SipUrl")

	conf.SmtpHost = os.Getenv("SmtpHost")
	conf.SmtpPort = os.Getenv("SmtpPort")
	conf.SmtpUser = os.Getenv("SmtpUser")
	conf.SmtpPass = os.Getenv("SmtpPass")
	conf.StripePublicKey = os.Getenv("StripePublicKey")
	conf.StripeSecretKey = os.Getenv("StripeSecretKey")
	conf.TelegramBotId = os.Getenv("TelegramBotId")
	conf.TelegramChannelId = os.Getenv("TelegramChannelId")
	conf.UspsUsername = os.Getenv("UspsUsername")
	conf.UspsPassword = os.Getenv("UspsPassword")
	conf.UfKey = os.Getenv("UfKey")
	conf.UfScheme = os.Getenv("UfScheme")

	conf.ZincToken = os.Getenv("ZincToken")
	conf.ZincCardName = os.Getenv("ZincCardName")
	conf.ZincCardNumber = os.Getenv("ZincCardNumber")
	conf.ZincCardExpirationMonth, err = strconv.Atoi(os.Getenv("ZincCardExpirationMonth"))
	if err != nil {
		//msg := "failed to read ZincCardExpirationMonth"
		//Error(msg)
		Error(err)
	}
	conf.ZincCardExpirationYear, err = strconv.Atoi(os.Getenv("ZincCardExpirationYear"))
	if err != nil {
		//msg := "failed to read ZincCardExpirationYear"
		//Error(msg)
		Error(err)
	}
	conf.ZincCardSecurityCode = os.Getenv("ZincCardSecurityCode")
}

// we use env variables via a configMap when deploying to prod
func (conf *Config) readConfFromEnvJson() {
	wd, _ := os.Getwd()
	//parent := filepath.Dir(wd)
	conf.JsonLoc = wd + "/.env"
	Debug("Reading config from " + conf.JsonLoc)

	jsonFile, err := os.Open(conf.JsonLoc)
	if err != nil {
		msg := "failed to open json file"
		Error(msg)
		//panic(errors.Wrap(err, msg))
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		msg := "failed to read json file"
		Error(msg)
		//panic(errors.Wrap(err, msg))
	}

	err = json.Unmarshal(byteValue, &conf)
	if err != nil {
		msg := "failed to unmarshal json file"
		Error(msg)
		//panic(errors.Wrap(err, msg))
	}
}

func (conf *Config) validateValues() {
	if conf.AmazonEmail == "" {
		Error("failed to parse Conf['AmazonEmail']")
	}
	if conf.AmazonPassword == "" {
		Error("failed to parse Conf['AmazonPassword']")
	}

	if conf.ApiDn == "" {
		Error("failed to parse Conf['ApiDn']")
	}

	if conf.AppDn == "" {
		Error("failed to parse Conf['AppDn']")
	}

	if conf.FridayyPhone == "" {
		conf.FridayyPhone = "fake"
		Error("failed to parse Conf['FridayyPhone']")
	}

	if conf.HushJwt == "" {
		conf.HushJwt = "fake"
		Error("failed to parse Conf['HushJwt']")
	}

	if conf.OpenAiKey == "" {
		Error("failed to parse Conf['OpenAiKey']")
	}

	if conf.SendBlueApiSecret == "" {
		Error("failed to parse Conf['SendBlueApiSecret']")
	}

	if conf.SendBlueApiKeyId == "" {
		Error("failed to parse Conf['SendBlueApiKeyId']")
	}

	if conf.SendBlueHookSecret == "" {
		Error("failed to parse Conf['SendBlueHookSecret']")
	}

	if conf.SmtpHost == "" {
		Error("failed to parse Conf['SmtpHost']")
	}
	if conf.SmtpPort == "" {
		Error("failed to parse Conf['SmtpPort']")
	}
	if conf.SmtpUser == "" {
		Error("failed to parse Conf['SmtpUser']")
	}
	if conf.SmtpPass == "" {
		Error("failed to parse Conf['SmtpPass']")
	}

	if conf.StripePublicKey == "" {
		Error("failed to parse Conf['StripePublicKey']")
	}

	if conf.StripeSecretKey == "" {
		Error("failed to parse Conf['StripeSecretKey']")
	}

	if conf.TelegramBotId == "" {
		Error("failed to parse Conf['TelegramBotId']")
	}
	if conf.TelegramChannelId == "" {
		Error("failed to parse Conf['TelegramChannelId']")
	}

	if conf.UspsUsername == "" {
		Error("failed to parse Conf['UspsUsername']")
	}
	if conf.UspsPassword == "" {
		Error("failed to parse Conf['UspsPassword']")
	}

	if conf.UfKey == "" {
		conf.UfKey = "fake"
		Error("failed to parse Conf['UfKey']")
	}

	if conf.ZincToken == "" {
		Error("failed to parse Conf['ZincToken']")
	}
	if conf.ZincCardName == "" {
		Error("failed to parse Conf['ZincCardName']")
	}
	if conf.ZincCardNumber == "" {
		Error("failed to parse Conf['ZincCardNumber']")
	}
	if conf.ZincCardExpirationMonth == 0 {
		Error("failed to parse Conf['ZincCardExpirationMonth']")
	}
	if conf.ZincCardExpirationYear == 0 {
		Error("failed to parse Conf['ZincCardExpirationYear']")
	}
	if conf.ZincCardSecurityCode == "" {
		Error("failed to parse Conf['ZincCardSecurityCode']")
	}
}

func (config *Config) IsProd() bool {
	if config.Environment == "prod" {
		return true
	}
	return false
}
