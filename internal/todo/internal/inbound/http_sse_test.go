package inbound

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mockCodec "github.com/shandysiswandi/goreng/mocker"
	"github.com/stretchr/testify/assert"
)

type sseResponseRecorder struct{}

func (sseResponseRecorder) Header() http.Header        { return make(http.Header) }
func (sseResponseRecorder) Write([]byte) (int, error)  { return 0, nil }
func (sseResponseRecorder) WriteHeader(statusCode int) {}
func (sseResponseRecorder) Result() *http.Response     { return &http.Response{Body: http.NoBody} }

func Test_sseEndpoint_TrigerEvent(t *testing.T) {
	tests := []struct {
		name   string
		mockFn func() (*sseEndpoint, *http.Request)
	}{
		{
			name: "ErrorEncode",
			mockFn: func() (*sseEndpoint, *http.Request) {
				r := httptest.NewRequest(http.MethodGet, "/", nil)
				jsonMock := mockCodec.NewMockCodec(t)

				jsonMock.EXPECT().
					Encode(Event{Name: "CREATE_TODO", Value: "TODO"}).
					Return(nil, assert.AnError)

				se := &sseEndpoint{
					codecJSON: jsonMock,
					clients:   map[chan []byte]struct{}{},
				}

				ca := make(chan []byte)
				se.addClient(ca)
				se.delClient(ca)

				return se, r
			},
		},
		{
			name: "Success",
			mockFn: func() (*sseEndpoint, *http.Request) {
				r := httptest.NewRequest(http.MethodGet, "/", nil)
				jsonMock := mockCodec.NewMockCodec(t)

				jsonMock.EXPECT().
					Encode(Event{Name: "CREATE_TODO", Value: "TODO"}).
					Return(nil, nil)

				se := &sseEndpoint{
					codecJSON: jsonMock,
					clients:   map[chan []byte]struct{}{},
				}

				ca := make(chan []byte)
				se.addClient(ca)

				return se, r
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			se, r := tt.mockFn()
			w := httptest.NewRecorder()
			se.TrigerEvent(w, r)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, http.StatusOK, res.StatusCode)
		})
	}
}

func Test_sseEndpoint_HandleEvent(t *testing.T) {
	tests := []struct {
		name   string
		want   int
		mockFn func() (*sseEndpoint, *http.Request)
	}{
		{
			name: "NotHTTPFlusher",
			want: 0,
			mockFn: func() (*sseEndpoint, *http.Request) {
				r := httptest.NewRequest(http.MethodGet, "/events", nil)
				return &sseEndpoint{}, r
			},
		},
		// next case cannot be tested, because it looping
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			se, r := tt.mockFn()
			w := sseResponseRecorder{}
			se.HandleEvent(w, r)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want, res.StatusCode)
		})
	}
}
