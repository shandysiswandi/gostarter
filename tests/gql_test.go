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
	token      string
}

func (gs *GQLSuite) SetupSuite() {
	// Hook 1: Initialize anything in suite

	gs.baseURL = "http://localhost:8081"
	gs.httpClient = resty.New()
	gs.idTodo = new(string)
	gs.token = gs.getAccessToken()
}

func (gs *GQLSuite) getAccessToken() string {
	responseBody := struct {
		AccessToken string `json:"access_token"`
	}{}

	resp, err := gs.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]any{
			"email":    "admin@admin.com",
			"password": "admin123",
		}).
		SetResult(&responseBody).
		Post(gs.baseURL + "/auth/login")

	gs.Assert().NoError(err)
	gs.Assert().NotNil(resp)

	if responseBody.AccessToken == "" {
		gs.T().Fatal("access token is empty", resp, err)
	}

	return responseBody.AccessToken
}

func (gs *GQLSuite) TearDownSuite() {
	// Hook 6: Clean up anything in suite

	gs.httpClient = nil
	gs.baseURL = ""
	gs.idTodo = nil
	gs.token = ""
}

// -:::::::::::::-  -:::::::::::::-  -:::::::::::::-  -:::::::::::::-  -:::::::::::::-

func (gs *GQLSuite) TestGQLTodos() {
	tests := []struct {
		name  string
		query func() string
		cb    func(data any)
	}{
		{
			name: "CREATE",
			query: func() string {
				return `mutation{
					create(in: {title: "some title gql", description: "some description gql" })
				}`
			},
			cb: func(v any) {
				data, ok := v.(map[string]any)
				if !ok {
					return
				}

				obj, ok := data["create"].(string)
				if !ok {
					return
				}

				*gs.idTodo = obj
			},
		},
		{
			name: "UPDATE_STATUS",
			query: func() string {
				return `mutation{
					updateStatus(in: { id: "` + *gs.idTodo + `", status: IN_PROGRESS }) {
						id status
					}
				}`
			},
		},
		{
			name: "UPDATE",
			query: func() string {
				return `mutation{
					update(in: {
						id: "` + *gs.idTodo + `",
						title: "update title gql",
						description: "update description gql",
						status: DROP
					}) { id title description status }
				}`
			},
		},
		{
			name: "FIND",
			query: func() string {
				return `query{
					find(id: "` + *gs.idTodo + `") { id title description status }
				}`
			},
		},
		{
			name: "FETCH",
			query: func() string {
				return `query{ fetch{ id title description status } }`
			},
		},
		{
			name: "DELETE",
			query: func() string {
				return `mutation{ delete(id: "` + *gs.idTodo + `") }`
			},
		},
	}

	for _, tt := range tests {
		gs.Run(tt.name, func() {
			// Arrange
			var result map[string]any
			query := map[string]any{"query": tt.query()}
			requestBody, err := json.Marshal(query)
			if err != nil {
				gs.Assert().NoError(err, "json.Marshal cannot be error")
			}

			// Action
			resp, err := gs.httpClient.R().
				SetHeader("Content-Type", "application/json").
				SetHeader("Accept", "application/json").
				SetAuthToken(gs.token).
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
