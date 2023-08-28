package uf

import (
	"github.com/gocql/gocql"
)

// Routes
var MahaGlogV1 = "/v1/glog.create"
var MahaGlogsGetByUfIdV1 = "/v1/glogs/get.byUfId"

type GlogStruct struct {
	GlogId    gocql.UUID
	ToddId    gocql.UUID
	UfId      gocql.UUID
	ChatId    gocql.UUID
	SurferId  gocql.UUID
	Code      string
	Level     int
	Msg       string
	Service   string
	Time      int64
	Interface interface{}
}

func GlogCreate(
	gibs *Gibs,
	glog GlogStruct,
) UfResponse {
	ures := MakeRequest(
		gibs,
		gibs.Conf.MahaAddress,
		MahaGlogV1,
		glog,
		&EmptyStruct{},
	)

	return ures
}
