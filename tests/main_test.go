package tests

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

/*
== Suite lifecycle

func (*Suite) SetupSuite() {
	// Hook 1: Initialize anything in suite
}

func (*Suite) SetupTest() {
	// Hook 2: Initialize anything specific to each test case, if necessary.
}

func (*Suite) SetupSubTest() {
	// Hook 3: Prepare resources needed for sub tests.
}

func (*Suite) TearDownSubTest() {
	// Hook 4: Clean up resources used in sub tests.
}

func (*Suite) TearDownTest() {
	// Hook 5: Clean up resources after each test case, if necessary.
}

func (*Suite) TearDownSuite() {
	// Hook 6: Clean up anything in suite
}
*/

func TestMain(t *testing.T) {
	// suite.Run(t, new(RESTSuite))
	suite.Run(t, new(GQLSuite))
}
