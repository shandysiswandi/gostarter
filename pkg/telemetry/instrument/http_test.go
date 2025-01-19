package instrument

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/stretchr/testify/assert"
)

func TestUseTelemetryServer(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if wf, ok := w.(http.Flusher); ok {
			wf.Flush()
		}

		if wh, ok := w.(http.Hijacker); ok {
			wh.Hijack()
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	tests := []struct {
		name       string
		telemetry  *telemetry.Telemetry
		handler    http.Handler
		wantStatus int
	}{
		{
			name:       "CollectorNoop",
			telemetry:  telemetry.NewTelemetry(),
			handler:    handler,
			wantStatus: http.StatusOK,
		},
		{
			name: "CollectorOpenTelemetry",
			telemetry: telemetry.NewTelemetry(
				telemetry.WithOTLP("http://"),
				telemetry.WithVerbose(),
				telemetry.WithLogFilter(),
			),
			handler:    handler,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(UseTelemetryServer(tt.telemetry)(tt.handler))
			defer ts.Close()

			resp, err := http.Get(ts.URL)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			resp, err = http.Post(ts.URL, "application/json", bytes.NewReader([]byte(`{"name":"ok"}`)))
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

// func TestStatusResponseWriter(t *testing.T) {
// 	rec := httptest.NewRecorder()
// 	srw := &statusResponseWriter{ResponseWriter: rec}

// 	srw.WriteHeader(http.StatusAccepted)
// 	assert.Equal(t, http.StatusAccepted, srw.statusCode)

// 	body := []byte("response body")
// 	srw.Write(body)
// 	assert.Equal(t, body, srw.responseBody)
// }

// func TestInstarumentHTTPServer_ServeHTTP(t *testing.T) {
// 	tel := &telemetry.Telemetry{
// 		// Initialize telemetry fields as needed
// 	}

// 	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 	})

// 	ihs := &instarumentHTTPServer{tel: tel, next: handler}

// 	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
// 	rec := httptest.NewRecorder()

// 	ihs.ServeHTTP(rec, req)

// 	assert.Equal(t, http.StatusOK, rec.Code)
// }
