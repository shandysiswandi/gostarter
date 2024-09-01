// Package http provides the HTTP handlers and routes for handling
// Server-Sent Events (SSE) and triggering events in a concurrent-safe manner.
//
//nolint:errcheck,revive,contextcheck,contextcheck,nlreturn // it will be ignored
package http

// Known Issue:
//  - WriteTimeout Limitation: The current setup shares a router with REST endpoints,
//    which may be subject to a WriteTimeout limit on the server, potentially causing
//    the connection to close prematurely. To address this, consider separating
//    the SSE endpoints onto a different server or
//    adjusting the WriteTimeout setting for this particular route.

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/shandysiswandi/gostarter/pkg/codec"
	"github.com/shandysiswandi/gostarter/pkg/http/middleware"
	"github.com/shandysiswandi/gostarter/pkg/logger"
)

// RegisterSSEEndpoint registers the SSE endpoints for handling event streams
// and triggering events using the provided httprouter.Router.
func RegisterSSEEndpoint(router *httprouter.Router, h *SSE) {
	router.Handler(http.MethodGet, "/events", middleware.Recovery(http.HandlerFunc(h.HandleEvent)))
	router.Handler(http.MethodGet, "/triger-event", middleware.Recovery(http.HandlerFunc(h.TrigerEvent)))
}

// Event represents a simple event with a name and value to be sent to clients.
type Event struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

// SSE provides the functionality for handling Server-Sent Events (SSE) and
// managing connected clients in a thread-safe manner.
type SSE struct {
	CodecJSON codec.Codec              // Codec for encoding event data into JSON
	Logger    logger.Logger            // Logger for logging information and errors
	Clients   map[chan []byte]struct{} // Map of clients connected via SSE
	mu        sync.RWMutex             // Mutex to ensure thread-safe access to Clients map
}

// NewSSE creates a new SSE handler with the provided codec for encoding JSON
// and a logger for logging events.
func NewSSE(codecJSON codec.Codec, logger logger.Logger) *SSE {
	return &SSE{
		CodecJSON: codecJSON,
		Logger:    logger,
		Clients:   make(map[chan []byte]struct{}),
	}
}

// TrigerEvent is an HTTP handler that triggers an event to all connected clients.
// It simulates a background task that sends a predefined event to the clients.
func (s *SSE) TrigerEvent(w http.ResponseWriter, r *http.Request) {
	s.doBackground()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(time.Now().Format(time.RFC3339)))
}

// addClient adds a new client to the Clients map in a thread-safe manner.
func (s *SSE) addClient(ch chan []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Clients[ch] = struct{}{}
}

// delClient removes a client from the Clients map in a thread-safe manner.
func (s *SSE) delClient(ch chan []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Clients, ch)
}

// doBackground sends a predefined event to all connected clients asynchronously.
// It is used to simulate the event being triggered and sent to all clients.
func (s *SSE) doBackground() {
	event := Event{Name: "CREATE_TODO", Value: "TODO"}
	data, err := s.CodecJSON.Encode(event)
	if err != nil {
		s.Logger.Error(context.Background(), "Failed to encode event:", err)
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	for client := range s.Clients {
		// Send data asynchronously to avoid blocking on slow clients
		go func(ch chan []byte) {
			ch <- data
		}(client)
	}
}

// HandleEvent is an HTTP handler that manages a Server-Sent Events (SSE) connection
// with a client. It sends periodic keepalive messages and listens for triggered events
// to forward them to the connected client.
func (s *SSE) HandleEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	keepAliveTickler := time.NewTicker(10 * time.Second)
	messageChan := make(chan []byte)
	s.addClient(messageChan)

	defer func() {
		s.delClient(messageChan)
	}()

	go func() {
		<-r.Context().Done()
		s.delClient(messageChan)
		keepAliveTickler.Stop()
		log.Println("client connection has been done")
	}()

	for {
		select {
		case <-keepAliveTickler.C:
			fmt.Fprintf(w, ":keepalive\n\n")
			flusher.Flush()
		case data := <-messageChan:
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}
