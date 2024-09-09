// Package codec provides a MsgpackCodec implementation of the Codec interface.
//
// MsgpackCodec is used for encoding and decoding data in MessagePack format.
package codec

import (
	"github.com/vmihailenco/msgpack/v5"
)

// MsgPackCodec is a Codec implementation for MessagePack encoding/decoding.
type MsgPackCodec struct{}

// NewMsgPackCodec creates a new instance of MsgpackCodec.
func NewMsgPackCodec() *MsgPackCodec {
	return &MsgPackCodec{}
}

// Encode encodes a value into MessagePack format.
//
// v: The value to encode. It can be of any type that is supported by MessagePack encoding.
// Returns: A byte slice containing the MessagePack-encoded data, and an error if encoding fails.
func (*MsgPackCodec) Encode(v any) ([]byte, error) {
	return msgpack.Marshal(v)
}

// Decode decodes MessagePack data from a byte slice into a value.
//
// data: The byte slice containing MessagePack-encoded data.
// v: A pointer to the value where the decoded data will be stored.
// Returns: An error if decoding fails.
func (*MsgPackCodec) Decode(data []byte, v any) error {
	return msgpack.Unmarshal(data, v)
}
