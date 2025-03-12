package sqlkit

import (
	"database/sql"

	"github.com/shandysiswandi/goreng/telemetry/logger"
)

type MockDB struct {
}

func NewMockDB(driver string, db *sql.DB, log logger.Logger) *MockDB {
	return &MockDB{}
}
