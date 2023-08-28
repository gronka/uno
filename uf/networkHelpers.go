package uf

import (
	"bytes"
	"encoding/json"
)

func PrintResponseBody(target []byte) error {
	amap := make(map[string]interface{})
	reader := bytes.NewReader(target)
	err := json.NewDecoder(reader).Decode(&amap)
	Debug(amap)
	return err
}
