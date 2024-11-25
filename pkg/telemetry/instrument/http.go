package instrument

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"net/http"
	"slices"
	"strings"

	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/requestid"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	xRequestID = "X-Request-ID"
)

func UseTelemetryServer(tel *telemetry.Telemetry, sid func() string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ihs := &instarumentHTTPServer{tel: tel, uuid: sid, next: h}

			switch tel.Collector() {
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
	statusCode   int
	responseBody []byte
}

func (srw *statusResponseWriter) Write(bytes []byte) (int, error) {
	srw.responseBody = bytes

	return srw.ResponseWriter.Write(bytes)
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

	r = r.WithContext(requestid.Set(r.Context(), r.Header.Get(xRequestID)))

	srw := &statusResponseWriter{ResponseWriter: w}
	filter := ihs.tel.Filter()
	fields := []logger.Field{
		logger.KeyVal("http.method", r.Method),
		logger.KeyVal("http.path", r.URL.Path),
	}

	var okBody bool
	var rwBody []byte
	if slices.Contains([]string{http.MethodPost, http.MethodPut, http.MethodPatch}, r.Method) {
		rawBody, err := io.ReadAll(r.Body)
		if err == nil {
			okBody = true
			rwBody = rawBody
			r.Body = io.NopCloser(bytes.NewBuffer(rawBody))
		}
	}

	ihs.next.ServeHTTP(srw, r)
	ctx := r.Context()
	fields = append(fields, logger.KeyVal("http.status", srw.statusCode))

	if okBody {
		fields = append(fields, logger.KeyVal("http.body", filter.Body(rwBody)))
	}

	if r.Method == http.MethodGet {
		fields = append(fields, logger.KeyVal("http.query", filter.Query(r.URL.RawQuery)))
	}

	fields = append(fields, logger.KeyVal("http.response", filter.Body(srw.responseBody)))
	fields = append(fields, logger.KeyVal("http.header", filter.Header(r.Header)))

	log.Println("middle", ctx)

	ic, err := ihs.tel.Meter().Int64Counter("http_requests")
	if err == nil {
		ic.Add(ctx, 1, metric.WithAttributes(
			attribute.Int("status", srw.statusCode),
			attribute.String("path", r.URL.Path),
			attribute.String("method", strings.ToLower(r.Method)),
		))
	} else {
		log.Println("sudahlah", err)
	}

	ihs.tel.Logger().Info(ctx, "http request response", fields...)
}
