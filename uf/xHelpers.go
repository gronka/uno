package uf

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/gocql/gocql"
	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
	"github.com/twmb/murmur3"
)

func RandomUUID() gocql.UUID {
	uuid, err := gocql.RandomUUID()
	if err != nil {
		Error(err)
	}
	return uuid
}

func StringToUUID(str string) gocql.UUID {
	strMurmur := murmur3.StringSum64(str)
	arr := make([]byte, 16)
	binary.LittleEndian.PutUint64(arr, strMurmur)
	uuid, err := gocql.UUIDFromBytes(arr)
	if err != nil {
		Error(err)
	}
	return uuid
}

func UuidSliceToStringSlice(uuids []gocql.UUID) []string {
	arr := make([]string, 0)
	for _, uuid := range uuids {
		arr = append(arr, uuid.String())
	}
	return arr
}

func StringSliceToUuidSlice(strs []string) []gocql.UUID {
	arr := make([]gocql.UUID, 0)
	for _, str := range strs {
		//arr = append(arr, gocql.UUID.FromBytes([]byte(str)))
		uuid, err := gocql.UUIDFromBytes([]byte(str))
		if err != nil {
			Error(err)
		}
		arr = append(arr, uuid)
	}
	return arr
}

func UuidArrayToUuidSlice(arr pgtype.UUIDArray) []gocql.UUID {
	uuids := make([]gocql.UUID, 0)
	for _, pgUuid := range arr.Elements {
		//var tmp [16]byte
		tmp := make([]byte, 16)
		pgUuid.AssignTo(tmp)
		copy(tmp[:], pgUuid.Bytes[:])
		uuid, err := gocql.UUIDFromBytes(tmp)
		Trace("trying to convert")
		if err != nil {
			Error(err)
		}
		uuids = append(uuids, uuid)
	}
	return uuids
}

func Check(err error, msg string) {
	if err != nil {
		panic(errors.Wrap(err, msg))
	}
}

func PrintBytes(bytes []byte) {
	for _, n := range bytes {
		fmt.Printf("%8b", n) // prints 1111111111111101
	}
	fmt.Printf("\n")
}

func IntToPriceString(price int) string {
	priceAsFloat := float32(price) / 100.0
	return fmt.Sprintf("$%.2f", priceAsFloat)
}

func UuidEqual(a, b gocql.UUID) bool {
	if bytes.Equal(a.Bytes(), b.Bytes()) {
		return true
	}
	return false
}

func UuidIsZero(a gocql.UUID) bool {
	if bytes.Equal(a.Bytes(), ZerosBytes) {
		return true
	}
	return false
}
