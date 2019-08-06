package bitmex

import (
	"encoding/json"
)

type BitmexSubscribe struct {
	Op   string   `json:"op"`
	Args []string `json:"args"`
}

func NewBitmexSubscribe(symbols []string, args ...string) *BitmexSubscribe {
	bs := new(BitmexSubscribe)
	bs.Op = "subscribe"
	bs.Args = make([]string, 0, len(args))

	for _, arg := range args {
		for _, symbol := range symbols {
			bs.Args = append(bs.Args, arg+":"+symbol)
		}
	}
	return bs
}

func (bs *BitmexSubscribe) GetOperateArgs() []string {
	return bs.Args
}

func (bs *BitmexSubscribe) Serialize() ([]byte, error) {
	return json.Marshal(bs)
}
