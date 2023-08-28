package uf_border

import "strings"

func (in *SendBlueIn) Clean() {
	in.Content = strings.TrimSpace(in.Content)
}
