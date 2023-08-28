package uf

import (
	"github.com/gocql/gocql"
)

var ZerosUuidString = "00000000-0000-0000-0000-000000000000"

var FridayyBytes = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
var FridayyUuid, _ = gocql.UUIDFromBytes(FridayyBytes)

var ZerosBytes = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
var ZerosBytes16 = [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
var ZerosUuid, _ = gocql.UUIDFromBytes(ZerosBytes)

var TwosBytes = []byte{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
var TwosBytes16 = [16]byte{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
var TwosUuid, _ = gocql.UUIDFromBytes(TwosBytes)

var SevensBytes = []byte{119, 119, 119, 119, 119, 119, 119, 119, 119, 119, 119, 119, 119, 119, 119, 119}
var SevensUuid, _ = gocql.UUIDFromBytes(SevensBytes)

var NinesBytes = []byte{153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153, 153}
var NinesBytes16 = [16]byte{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9}
var NinesUuid, _ = gocql.UUIDFromBytes(NinesBytes)

var SendBlueBytes = []byte{5, 8, 5, 8, 5, 8, 5, 8, 5, 8, 5, 8, 5, 8, 5, 8}
var SendBlueBytes16 = [16]byte{5, 8, 5, 8, 5, 8, 5, 8, 5, 8, 5, 8, 5, 8, 5, 8}

type EmptyStruct struct {
}

func (out *EmptyStruct) Sanitize() error {
	return nil
}
