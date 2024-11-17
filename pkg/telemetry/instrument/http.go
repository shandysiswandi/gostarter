package instrument

import (
	"bufio"
	"net"
	"net/http"

	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/requestid"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	xRequestID = "X-Request-ID"
)

func UseTelemetryServer(tel *telemetry.Telemetry, sid func() string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ihs := &instarumentHTTPServer{tel: tel, uuid: sid, next: h}

			switch tel.TracerCollector() {
			case telemetry.OPENTELEMETRY:
				hand := otelhttp.NewHandler(ihs, r.URL.Path,
					otelhttp.WithTracerProvider(tel.TracerProvider()))
				hand.ServeHTTP(w, r)

			case telemetry.NOOP:
				h.ServeHTTP(w, r)

			default:
				h.ServeHTTP(w, r)
			}
		})
	}
}

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code and writes it.
func (srw *statusResponseWriter) WriteHeader(statusCode int) {
	srw.statusCode = statusCode
	srw.ResponseWriter.WriteHeader(statusCode)
}

// Hijack implements the http.Hijacker interface if the underlying ResponseWriter supports it.
func (srw *statusResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := srw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}

	return hijacker.Hijack()
}

// Flush implements the http.Flusher interface if the underlying ResponseWriter supports it.
func (srw *statusResponseWriter) Flush() {
	if flusher, ok := srw.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

type instarumentHTTPServer struct {
	next http.Handler
	uuid func() string
	tel  *telemetry.Telemetry
}

func (ihs *instarumentHTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if rid := r.Header.Get(xRequestID); rid == "" {
		r.Header.Set(xRequestID, ihs.uuid())
	}

	ctx := r.Context()
	ctx = requestid.Set(ctx, r.Header.Get(xRequestID))
	r = r.WithContext(ctx)

	srw := &statusResponseWriter{ResponseWriter: w}

	ihs.next.ServeHTTP(srw, r)

	ihs.tel.Logger().Info(ctx, "http request response",
		logger.KeyVal("http.method", r.Method),
		logger.KeyVal("http.path", r.URL.Path),
		logger.KeyVal("http.status", srw.statusCode),
		logger.KeyVal("http.header", map[string]string{
			"user-agent":   r.Header.Get("User-Agent"),
			"content-type": r.Header.Get("Content-Type"),
			"host":         r.Header.Get("Host"),
			"client-ip":    r.Header.Get("X-Forwarded-For"),
		}),
	)
}
