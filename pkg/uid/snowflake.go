// Package uid provides interfaces and implementations for generating unique identifiers (UIDs).
package uid

import (
	"crypto/rand"
	"encoding/binary"

	"github.com/bwmarrin/snowflake"
)

// SnowflakeNumber is an implementation of UIDNumber that uses snowflake for generating unique numeric UIDs.
type SnowflakeNumber struct {
	node *snowflake.Node
}

// generateRandomNodeID generates a random node ID using a cryptographic random number generator.
func generateRandomNodeID() (int64, error) {
	var nodeID int64
	err := binary.Read(rand.Reader, binary.BigEndian, &nodeID)
	if err != nil {
		return 0, err
	}

	return nodeID & (1<<10 - 1), nil // Limiting to 10 bits for node ID
}

// NewSnowflakeNumber creates a new SnowflakeUIDNumber instance with a specific node ID.
func NewSnowflakeNumber() (*SnowflakeNumber, error) {
	nodeID, err := generateRandomNodeID()
	if err != nil {
		return nil, err
	}

	snowflake.Epoch = 1727197200000 // Fri Sep 25 2024 00:00:00.000 WIB

	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return nil, err
	}

	return &SnowflakeNumber{node: node}, nil
}

// Generate generates a unique identifier as a uint64 number using Snowflake.
func (s *SnowflakeNumber) Generate() uint64 {
	return uint64(s.node.Generate().Int64())
}
