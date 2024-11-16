package codec

import (
	"github.com/vmihailenco/msgpack/v5"
)

// MsgPackCodec is a Codec implementation for MessagePack encoding/decoding.
type MsgPackCodec struct{}

// NewMsgPackCodec creates a new instance of MsgPackCodec.
func NewMsgPackCodec() *MsgPackCodec {
	return &MsgPackCodec{}
}

// Encode encodes a value into MessagePack format.
func (*MsgPackCodec) Encode(v any) ([]byte, error) {
	return msgpack.Marshal(v)
}

// Decode decodes MessagePack data from a byte slice into a value.
func (*MsgPackCodec) Decode(data []byte, v any) error {
	return msgpack.Unmarshal(data, v)
}
