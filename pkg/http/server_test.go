package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shandysiswandi/gostarter/pkg/errs"
	"github.com/stretchr/testify/assert"
)

// MockHandler returns a successful response.
func MockHandler(ctx context.Context, r *http.Request) (any, error) {
	return map[string]string{"message": "success"}, nil
}

// MockErrorHandler returns an error response.
func MockErrorHandler(ctx context.Context, r *http.Request) (any, error) {
	return nil, errs.New(nil, "invalid input", errs.TypeValidation, errs.CodeInvalidInput)
}

// testMiddleware is a middleware for testing.
func testMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Test-Middleware", "true")
		next.ServeHTTP(w, r)
	})
}

// TestServe_NewServe checks creation of a new Serve instance.
func TestServe_NewServe(t *testing.T) {
	serve := NewServe()

	assert.NotNil(t, serve)
}

// TestServe_Endpoint checks endpoint creation without middleware.
func TestServe_Endpoint(t *testing.T) {
	serve := NewServe()
	endpoint := serve.Endpoint(MockHandler)

	assert.NotNil(t, endpoint)
}

// TestServe_EndpointWithMiddleware checks endpoint creation with middleware.
func TestServe_EndpointWithMiddleware(t *testing.T) {
	serve := NewServe(WithMiddlewares(testMiddleware))
	endpoint := serve.Endpoint(MockHandler, testMiddleware)

	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(endpoint.ServeHTTP)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, "true", rr.Header().Get("X-Test-Middleware"))
	assert.Equal(t, http.StatusOK, rr.Code)

	var body map[string]string
	err := json.NewDecoder(rr.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Equal(t, "success", body["message"])
}

// TestEndpoint_ErrorHandling checks error response handling.
func TestEndpoint_ErrorHandling(t *testing.T) {
	serve := NewServe()
	endpoint := serve.Endpoint(MockErrorHandler)

	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(endpoint.ServeHTTP)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var body map[string]string
	err := json.NewDecoder(rr.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Equal(t, "invalid input", body["error"])
}

// TestEndpoint_ResponseEncoder checks response encoding and StatusCoder behavior.
func TestEndpoint_ResponseEncoder(t *testing.T) {
	t.Run("with_content", func(t *testing.T) {
		serve := NewServe()
		endpoint := serve.Endpoint(func(ctx context.Context, r *http.Request) (any, error) {
			return &CustomResponse{statusCode: http.StatusCreated, message: "created"}, nil
		})

		req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(endpoint.ServeHTTP)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		var body map[string]string
		err := json.NewDecoder(rr.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "created", body["message"])
	})

	t.Run("no_content", func(t *testing.T) {
		serve := NewServe()
		endpoint := serve.Endpoint(func(ctx context.Context, r *http.Request) (any, error) {
			return &CustomResponse{statusCode: http.StatusNoContent, message: "created"}, nil
		})

		req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(endpoint.ServeHTTP)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)

		var body map[string]string
		_ = json.NewDecoder(rr.Body).Decode(&body)
		assert.Equal(t, "", body["message"])
	})

	t.Run("internal", func(t *testing.T) {
		serve := NewServe()
		endpoint := serve.Endpoint(func(ctx context.Context, r *http.Request) (any, error) {
			return make(chan int), nil // Return a type that causes encoding failure
		})

		req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(endpoint.ServeHTTP)

		handler.ServeHTTP(rr, req)

		// The test should verify that the response is an internal server error if encoding fails.
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		var body map[string]string
		err := json.NewDecoder(rr.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "json: unsupported type: chan int", body["error"])
	})
}

// TestEndpoint_ErrorEncoder checks that error encoding works properly.
func TestEndpoint_ErrorEncoder(t *testing.T) {
	t.Run("not_found", func(t *testing.T) {
		serve := NewServe()
		endpoint := serve.Endpoint(func(ctx context.Context, r *http.Request) (any, error) {
			return nil, errs.New(nil, "not found", errs.TypeBusiness, errs.CodeNotFound)
		})

		req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(endpoint.ServeHTTP)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)

		var body map[string]string
		err := json.NewDecoder(rr.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "not found", body["error"])
	})

	t.Run("internal", func(t *testing.T) {
		serve := NewServe()
		endpoint := serve.Endpoint(func(ctx context.Context, r *http.Request) (any, error) {
			return nil, errs.New(nil, "internal", errs.TypeServer, errs.CodeInternal)
		})

		req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(endpoint.ServeHTTP)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		var body map[string]string
		err := json.NewDecoder(rr.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "internal", body["error"])
	})

	t.Run("internal generic", func(t *testing.T) {
		serve := NewServe()
		endpoint := serve.Endpoint(func(ctx context.Context, r *http.Request) (any, error) {
			return nil, assert.AnError
		})

		req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(endpoint.ServeHTTP)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		var body map[string]string
		err := json.NewDecoder(rr.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "assert.AnError general error for testing", body["error"])
	})
}

func TestEndpoint_errCode(t *testing.T) {

	abc := func(c errs.Code) error {
		return errs.New(nil, "", errs.TypeBusiness, c)
	}

	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{
			name:     "Non-errs.Error type",
			err:      fmt.Errorf("generic error"),
			expected: http.StatusInternalServerError,
		},
		{
			name:     "Server Error",
			err:      abc(errs.CodeUnknown),
			expected: http.StatusInternalServerError,
		},
		{
			name:     "Validation Error",
			err:      abc(errs.CodeInvalidInput),
			expected: http.StatusBadRequest,
		},
		{
			name:     "Invalid Input Error",
			err:      abc(errs.CodeInvalidInput),
			expected: http.StatusBadRequest,
		},
		{
			name:     "Not Found Error",
			err:      abc(errs.CodeNotFound),
			expected: http.StatusNotFound,
		},
		{
			name:     "Conflict Error",
			err:      abc(errs.CodeConflict),
			expected: http.StatusConflict,
		},
		{
			name:     "Unauthorized Error",
			err:      abc(errs.CodeUnauthorized),
			expected: http.StatusUnauthorized,
		},
		{
			name:     "Forbidden Error",
			err:      abc(errs.CodeForbidden),
			expected: http.StatusForbidden,
		},
		{
			name:     "Timeout Error",
			err:      abc(errs.CodeTimeout),
			expected: http.StatusRequestTimeout,
		},
		{
			name:     "Default Case",
			err:      abc(errs.CodeInternal),
			expected: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endpoint := &Endpoint{}
			actual := endpoint.errCode(tt.err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

// CustomResponse is a test implementation of StatusCoder.
type CustomResponse struct {
	statusCode int
	message    string
}

func (r *CustomResponse) StatusCode() int {
	return r.statusCode
}

func (r *CustomResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"message": r.message})
}

// TestServeOption_Apply tests the ServeOption application.
func TestServeOption_Apply(t *testing.T) {
	opt := WithMiddlewares(testMiddleware)
	serve := NewServe(opt)

	assert.NotNil(t, serve.mws)
	assert.Equal(t, 1, len(serve.mws))
}
