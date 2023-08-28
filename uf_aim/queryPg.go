package uf_aim

import (
	"sort"
	"strings"

	"github.com/gocql/gocql"
	"gitlab.com/textfridayy/uno/uf"
)

type QueryPg struct {
	QueryId   gocql.UUID
	SurferId  gocql.UUID
	QueryHash gocql.UUID

	// details
	CurrentProductId    gocql.UUID
	SkippedProductIds   []gocql.UUID
	CurrentProductPrice int
	AbsoluteMaxPrice    int
	AbsoluteMinPrice    int
	Count               int

	// taxonomy
	Noun                 string
	Category             string
	AdjectivesSorted     []string
	AntiAdjectivesSorted []string

	TimeCreated int64
	TimeUpdated int64

	// not in db
	MaxPrice int
	MinPrice int
}

func (query *QueryPg) AsQueryString() string {
	return strings.Join(query.AdjectivesSorted[:], ", ") + " " + query.Noun
}

func (aimInfo *AimInfo) LoadQueryFromUno(gibs *uf.Gibs) error {
	query, err := LQueryGetById(gibs, aimInfo.ActiveQueryId)
	aimInfo.Query = &query
	return err
}

func (query *QueryPg) AddAdjectives(adj []string) {
	query.AdjectivesSorted = append(query.AdjectivesSorted, adj...)
	sort.Strings(query.AdjectivesSorted)
}

func (query *QueryPg) ComputeQueryHash() gocql.UUID {
	var sb strings.Builder
	sb.WriteString(query.Noun)
	for _, adj := range query.AdjectivesSorted {
		sb.WriteString(";" + adj)
	}
	return uf.StringToUUID(sb.String())
}

func (aimInfo *AimInfo) CreateNewQuery() {
	now := uf.NowStamp()
	//AdjectivesSorted:     aimInfo.MakoResponse.AdjectivesSorted(),
	uf.Trace("about to create QueryPg")
	uf.Trace(uf.RandomUUID())
	uf.Trace(aimInfo.Surfer.SurferId)
	aimInfo.Query = &QueryPg{
		QueryId:  uf.RandomUUID(),
		SurferId: aimInfo.Surfer.SurferId,
		Count:    1,

		Noun:                 aimInfo.MakoResponse.Noun,
		Category:             aimInfo.MakoResponse.Category,
		AdjectivesSorted:     aimInfo.MakoResponse.AdjectivesSorted(),
		AntiAdjectivesSorted: make([]string, 0),
		TimeCreated:          now,
		TimeUpdated:          now,
	}

	aimInfo.Query.QueryHash = aimInfo.Query.ComputeQueryHash()

	aimInfo.ActiveQueryId = aimInfo.Query.QueryId
}
