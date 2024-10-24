// Package http provides the HTTP handlers and routes for handling
// Server-Sent Events (SSE) and triggering events in a concurrent-safe manner.
//
//nolint:errcheck,revive,contextcheck,contextcheck,nlreturn // it will be ignored
package http

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/shandysiswandi/gostarter/pkg/codec"
)

func TestRegisterSSEEndpoint(t *testing.T) {
	type args struct {
		router *httprouter.Router
		h      *SSE
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterSSEEndpoint(tt.args.router, tt.args.h)
		})
	}
}

func TestNewSSE(t *testing.T) {
	type args struct {
		codecJSON codec.Codec
	}
	tests := []struct {
		name string
		args args
		want *SSE
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSSE(tt.args.codecJSON); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSSE() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSSE_TrigerEvent(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		s    *SSE
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.TrigerEvent(tt.args.w, tt.args.r)
		})
	}
}

func TestSSE_addClient(t *testing.T) {
	type args struct {
		ch chan []byte
	}
	tests := []struct {
		name string
		s    *SSE
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.addClient(tt.args.ch)
		})
	}
}

func TestSSE_delClient(t *testing.T) {
	type args struct {
		ch chan []byte
	}
	tests := []struct {
		name string
		s    *SSE
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.delClient(tt.args.ch)
		})
	}
}

func TestSSE_doBackground(t *testing.T) {
	tests := []struct {
		name string
		s    *SSE
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.doBackground()
		})
	}
}

func TestSSE_HandleEvent(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		s    *SSE
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.HandleEvent(tt.args.w, tt.args.r)
		})
	}
}
