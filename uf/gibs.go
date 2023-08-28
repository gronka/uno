package uf

import (
	"encoding/json"
	"strings"

	"github.com/gocql/gocql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gronka/tg"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

type Gibs struct {
	Ctx  *fasthttp.RequestCtx
	Pile *Pile
	Conf *Config

	// UfAuth is a custom auth header because kubernetes seems to default an
	// empty Authorization header to ":" which is "Og==", so it's probably
	// better for us to create our own authorization space
	UfAuth                string
	IsJwtValid            bool
	IsJwtFromSms          bool
	SendBlueHookSecret    string
	SipHookSecret         string
	LoopAuthorization     string
	LoopSecretKey         string
	LoopHookAuthorization string
	ToddId                gocql.UUID
	ToddType              string
	// UfId is for tracking a request for its lifetime across services
	UfId  gocql.UUID
	UfKey string
}

func (gibs *Gibs) InitGibsForAll(
	ctx *fasthttp.RequestCtx,
	pile *Pile,
	in interface{}) {
	gibs.Ctx = ctx
	gibs.Conf = &pile.Conf
	gibs.Pile = pile

	json.Unmarshal(ctx.PostBody(), &in)

	Trace(ctx.PostBody())

	// If the UfId is set later, it will be overwritten by InitGibs
	gibs.UfId = gocql.TimeUUID()
}

func InitGibs(
	ctx *fasthttp.RequestCtx,
	pile *Pile,
	in interface{},
) (gibs Gibs, rb ResponseBuilder) {
	gibs.InitGibsForAll(ctx, pile, in)

	gibs.UfAuth = string(ctx.Request.Header.Peek("UfAuth"))
	err := checkAndGetJwtClaims(&gibs, &pile.Conf, gibs.UfAuth)
	if err != nil {
		panic(errors.Wrap(err, "failed jwt check"))
	}

	ufId := string(ctx.Request.Header.Peek("UfId"))
	gibs.UfId, err = gocql.ParseUUID(ufId)
	if err != nil {
		// all uf apps except public use UfId to track the request
		if pile.Conf.UfName != "public" {
			Glog(&gibs, GlogStruct{
				Level:     LevelError,
				Code:      "glog.100",
				Interface: strings.Join([]string{"Ill-formed UfId: ", ufId}, ""),
			})
		}
	}

	gibs.UfKey = string(ctx.Request.Header.Peek("UfKey"))
	return
}

func InitGibsForSendBlue(
	ctx *fasthttp.RequestCtx,
	pile *Pile,
	in interface{},
) (gibs Gibs, rb ResponseBuilder) {
	gibs.InitGibsForAll(ctx, pile, in)
	gibs.SendBlueHookSecret = string(ctx.Request.Header.Peek("sb-signing-secret"))
	return
}

func InitGibsForWeb(
	ctx *fasthttp.RequestCtx,
	pile *Pile,
	in interface{},
) (gibs Gibs, rb ResponseBuilder) {
	gibs.InitGibsForAll(ctx, pile, in)
	return
}

func InitGibsForSip(
	ctx *fasthttp.RequestCtx,
	pile *Pile,
	in interface{},
) (gibs Gibs, rb ResponseBuilder) {
	gibs.InitGibsForAll(ctx, pile, in)
	gibs.SipHookSecret = string(ctx.QueryArgs().Peek("secret"))
	return
}

func InitGibsForLoop(
	ctx *fasthttp.RequestCtx,
	pile *Pile,
	in interface{},
) (gibs Gibs, rb ResponseBuilder) {
	gibs.InitGibsForAll(ctx, pile, in)
	gibs.LoopHookAuthorization = string(ctx.Request.Header.Peek("Authorization"))
	return
}

type JwtClaims struct {
	jwt.StandardClaims
	ToddId   string
	ToddType string
}

func CreateJwt(
	conf *Config,
	toddId gocql.UUID,
	issuer string,
) (string, error) {
	if issuer == "email" || issuer == "sip" || issuer == "loop" || issuer == "mock" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"ExpiresAt": InXDaysStamp(30),
			"ToddId":    toddId,
			"ToddType":  "surfer",
			"Issuer":    issuer,
		})

		tokenString, err := token.SignedString([]byte(conf.HushJwt))
		if err != nil {
			return "", errors.Wrap(err, "failed to sign token")
		}
		return tokenString, nil
	} else {
		Fatal("invalid issuer: " + issuer)
	}

	return "", errors.New("invalid jwt issuer: " + issuer)
}

func checkAndGetJwtClaims(gibs *Gibs, conf *Config, authorization string) error {
	gibs.IsJwtValid = false
	if authorization == "" {
		return nil
	}

	tg.Info("checking auth header: " + authorization)

	token, err := jwt.ParseWithClaims(authorization, &JwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(conf.HushJwt), nil
		})
	if err != nil {
		return errors.Wrap(err, "failure parsing jwt")
	}

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		gibs.ToddId, err = gocql.ParseUUID(claims.ToddId)
		if err != nil {
			return errors.Wrap(err, "failure decoding claims.ToddId")
		}
		Info("Token is valid with ToddId " + claims.ToddId)

		if claims.Issuer == "uf_border" {
			gibs.IsJwtFromSms = true
		}

		//gibs.ToddType = claims.ToddType
		//if err != nil {
		//return errors.Wrap(err, "failure decoding claims.ToddType")
		//}
		//Info("Token is valid with ToddType " + claims.ToddType)
	} else {
		return errors.New("jwt token is invalid")
	}

	gibs.IsJwtValid = true
	return nil
}

func (gibs *Gibs) String() string {
	return "UfId: " + gibs.UfId.String() + ", ToddId: " + gibs.ToddId.String()
}
