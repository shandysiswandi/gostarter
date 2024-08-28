package tests

import (
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
)

type RESTSuite struct {
	suite.Suite

	httpClient *resty.Client
	baseURL    string
}

func (rs *RESTSuite) SetupSuite() {
	// Hook 1: Initialize anything in suite

	rs.baseURL = "http://localhost:8081"
	rs.httpClient = resty.New()
}

func (rs *RESTSuite) SetupTest() {
	// Hook 2: Initialize anything specific to each test case, if necessary.
}

func (rs *RESTSuite) SetupSubTest() {
	// Hook 3: Prepare resources needed for subtests.
}

func (rs *RESTSuite) TearDownSubTest() {
	// Hook 4: Clean up resources used in subtests.
}

func (rs *RESTSuite) TearDownTest() {
	// Hook 5: Clean up resources after each test case, if necessary.
}

func (rs *RESTSuite) TearDownSuite() {
	// Hook 6: Clean up anything in suite

	rs.httpClient = nil
	rs.baseURL = ""
}

// -:::::::::::::-  -:::::::::::::-  -:::::::::::::-  -:::::::::::::-  -:::::::::::::-

func (rs *RESTSuite) TestTodos() {
	id := uint64(0)

	rs.Run("Create", func() {
		// Arrange
		requestBody := map[string]any{
			"title":       "tests title 200",
			"description": "test description 200",
		}
		responseBody := struct {
			ID uint64 `json:"id"`
		}{}

		// Action
		resp, err := rs.httpClient.R().
			SetHeader("Content-Type", "application/json").
			SetBody(requestBody).
			SetResult(&responseBody).
			Post(rs.baseURL + "/todos")

		// Assert
		rs.Assert().NoError(err)
		rs.Assert().NotNil(resp)
		if ok := rs.Assert().Equal(http.StatusOK, resp.StatusCode()); ok {
			rs.Assert().NotEmpty(responseBody.ID)
		} else {
			rs.T().Log(rs.T().Name(), "response:", resp)
		}

		id = responseBody.ID
	})

	rs.Run("UpdateStatus", func() {
		// Arrange
		status := "DROP" // valid values: "UNKNOWN" "INITIATE" "IN_PROGRESS" "DROP" "DONE"
		requestBody := map[string]any{
			"status": status,
		}
		responseBody := struct {
			ID     uint64 `json:"id"`
			Status string `json:"status"`
		}{}

		// Action
		resp, err := rs.httpClient.R().
			SetHeader("Content-Type", "application/json").
			SetPathParam("todoId", strconv.FormatUint(id, 10)).
			SetBody(requestBody).
			SetResult(&responseBody).
			Patch(rs.baseURL + "/todos/{todoId}/status")

		// Assert
		rs.Assert().NoError(err)
		rs.Assert().NotNil(resp)
		rs.Assert().Equal(http.StatusOK, resp.StatusCode())
		rs.Assert().Equal(id, responseBody.ID)
		rs.Assert().Equal(status, responseBody.Status)
	})

	rs.Run("Update", func() {
		// Arrange
		requestBody := map[string]any{
			"title":       "tests title update 200",
			"description": "test description update 200",
			"status":      "DONE",
		}
		responseBody := struct {
			ID          uint64 `json:"id"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Status      string `json:"status"`
		}{}

		// Action
		resp, err := rs.httpClient.R().
			SetHeader("Content-Type", "application/json").
			SetPathParam("todoId", strconv.FormatUint(id, 10)).
			SetBody(requestBody).
			SetResult(&responseBody).
			Put(rs.baseURL + "/todos/{todoId}")

		// Assert
		rs.Assert().NoError(err)
		rs.Assert().NotNil(resp)
		rs.Assert().Equal(http.StatusOK, resp.StatusCode())
		rs.Assert().Equal(id, responseBody.ID)
		rs.Assert().Equal(requestBody["title"], responseBody.Title)
		rs.Assert().Equal(requestBody["description"], responseBody.Description)
		rs.Assert().Equal(requestBody["status"], responseBody.Status)
	})

	rs.Run("Find", func() {
		// Arrange
		responseBody := struct {
			ID          uint64 `json:"id"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Status      string `json:"status"`
		}{}

		// Action
		resp, err := rs.httpClient.R().
			SetPathParam("todoId", strconv.FormatUint(id, 10)).
			SetResult(&responseBody).
			Get(rs.baseURL + "/todos/{todoId}")

		// Assert
		rs.Assert().NoError(err)
		rs.Assert().NotNil(resp)
		rs.Assert().Equal(http.StatusOK, resp.StatusCode())
		rs.Assert().Equal("tests title update 200", responseBody.Title)
		rs.Assert().Equal("test description update 200", responseBody.Description)
		rs.Assert().Equal("DONE", responseBody.Status)
	})

	rs.Run("Fetch", func() {
		// Arrange
		var responseBody struct {
			Todos []struct {
				ID          uint64 `json:"id"`
				Title       string `json:"title"`
				Description string `json:"description"`
				Status      string `json:"status"`
			} `json:"todos"`
		}
		var todo struct {
			ID          uint64 `json:"id"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Status      string `json:"status"`
		}

		// Action
		resp, err := rs.httpClient.R().
			SetResult(&responseBody).
			Get(rs.baseURL + "/todos")

		// Assert
		rs.Assert().NoError(err)
		rs.Assert().NotNil(resp)
		rs.Assert().Equal(http.StatusOK, resp.StatusCode())

		for _, todoItem := range responseBody.Todos {
			if todoItem.ID == id {
				todo = todoItem
				break
			}
		}

		rs.Assert().Equal("tests title update 200", todo.Title)
		rs.Assert().Equal("test description update 200", todo.Description)
		rs.Assert().Equal("DONE", todo.Status)
	})

	rs.Run("Delete", func() {
		// Arrange
		responseBody := struct {
			ID uint64 `json:"id"`
		}{}

		// Action
		resp, err := rs.httpClient.R().
			SetPathParam("todoId", strconv.FormatUint(id, 10)).
			SetResult(&responseBody).
			Delete(rs.baseURL + "/todos/{todoId}")

		// Assert
		rs.Assert().NoError(err)
		rs.Assert().NotNil(resp)
		rs.Assert().Equal(http.StatusOK, resp.StatusCode())
		rs.Assert().Equal(id, responseBody.ID)
	})
}
