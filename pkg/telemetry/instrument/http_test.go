package instrument

import (
	"bufio"
	"net"
	"net/http"
	"reflect"
	"testing"

	"github.com/shandysiswandi/gostarter/pkg/telemetry"
)

func TestUseTelemetryServer(t *testing.T) {
	type args struct {
		tel *telemetry.Telemetry
		sid func() string
	}
	tests := []struct {
		name string
		args args
		want func(http.Handler) http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := UseTelemetryServer(tt.args.tel, tt.args.sid); !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("UseTelemetryServer() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func Test_statusResponseWriter_WriteHeader(t *testing.T) {
	type args struct {
		statusCode int
	}
	tests := []struct {
		name string
		srw  *statusResponseWriter
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.srw.WriteHeader(tt.args.statusCode)
		})
	}
}

func Test_statusResponseWriter_Hijack(t *testing.T) {
	tests := []struct {
		name    string
		srw     *statusResponseWriter
		want    net.Conn
		want1   *bufio.ReadWriter
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.srw.Hijack()
			if (err != nil) != tt.wantErr {
				t.Errorf("statusResponseWriter.Hijack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("statusResponseWriter.Hijack() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("statusResponseWriter.Hijack() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_statusResponseWriter_Flush(t *testing.T) {
	tests := []struct {
		name string
		srw  *statusResponseWriter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.srw.Flush()
		})
	}
}

func Test_instarumentHTTPServer_ServeHTTP(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		ihs  *instarumentHTTPServer
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ihs.ServeHTTP(tt.args.w, tt.args.r)
		})
	}
}
