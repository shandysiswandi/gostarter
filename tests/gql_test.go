package tests

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
)

type GQLSuite struct {
	suite.Suite

	httpClient *resty.Client
	baseURL    string
	idTodo     *string
}

func (gs *GQLSuite) SetupSuite() {
	// Hook 1: Initialize anything in suite

	gs.baseURL = "http://localhost:8081"
	gs.httpClient = resty.New()
	gs.idTodo = new(string)
}

func (gs *GQLSuite) TearDownSuite() {
	// Hook 6: Clean up anything in suite

	gs.httpClient = nil
	gs.baseURL = ""
}

// -:::::::::::::-  -:::::::::::::-  -:::::::::::::-  -:::::::::::::-  -:::::::::::::-

func (gs *GQLSuite) TestTodos() {
	tests := []struct {
		name  string
		query func(id string) string
		cb    func(data any)
	}{
		{
			name: "CREATE",
			query: func(id string) string {
				return `mutation{
					create(in: {title: "some title gql", description: "some description gql" }) { 
						id title description status 
					}
				}`
			},
			cb: func(v any) {
				data, ok := v.(map[string]any)
				if !ok {
					return
				}

				obj, ok := data["create"].(map[string]any)
				if !ok {
					return
				}

				*gs.idTodo = obj["id"].(string)
			},
		},
		{
			name: "UPDATE_STATUS",
			query: func(id string) string {
				return `mutation{
					updateStatus(id: "` + id + `", status: DONE) {
						id status
					}
				}`
			},
		},
		{
			name: "UPDATE",
			query: func(id string) string {
				return `mutation{
					update(
						id: "` + id + `", in: {
						title: "update title gql", 
						description: "update description gql", 
						status: DROP
					}) { id title description status }
				}`
			},
		},
		{
			name: "FIND",
			query: func(id string) string {
				return `query{
					find(id: "` + id + `") { id title description status }
				}`
			},
		},
		{
			name: "FETCH",
			query: func(id string) string {
				return `query{ fetch{ id title description status } }`
			},
		},
		{
			name: "DELETE",
			query: func(id string) string {
				return `mutation{ delete(id: "` + id + `") }`
			},
		},
	}

	for _, tt := range tests {
		gs.Run(tt.name, func() {
			// Arrange
			var result map[string]any
			query := map[string]any{"query": tt.query(*gs.idTodo)}
			requestBody, err := json.Marshal(query)
			if err != nil {
				gs.Assert().NoError(err, "json.Marshal cannot be error")
			}

			// Action
			resp, err := gs.httpClient.R().
				SetHeader("Content-Type", "application/json").
				SetHeader("Accept", "application/json").
				SetBody(requestBody).
				SetResult(&result).
				SetError(&result).
				Post(gs.baseURL + "/graphql")

			// Assert
			gs.Assert().NoError(err)
			gs.Assert().NotNil(resp)
			gs.Assert().NotNil(result)
			gs.Assert().NotEmpty(result)

			if !gs.Assert().Empty(result["errors"]) {
				gs.T().Log("query", query)
			}
			if gs.Assert().NotEmpty(result["data"]) && tt.cb != nil {
				tt.cb(result["data"])
			}
		})
	}
}
