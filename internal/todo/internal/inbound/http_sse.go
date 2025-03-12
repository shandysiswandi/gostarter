//nolint:errcheck,revive // it will be ignored
package inbound

// Known Issue:
//  - WriteTimeout Limitation: The current setup shares a router with REST endpoints,
//    which may be subject to a WriteTimeout limit on the server, potentially causing
//    the connection to close prematurely. To address this, consider separating
//    the SSE endpoints onto a different server or
//    adjusting the WriteTimeout setting for this particular route.

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/shandysiswandi/goreng/codec"
)

type Event struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

type sseEndpoint struct {
	codecJSON codec.Codec
	clients   map[chan []byte]struct{}
	mu        sync.RWMutex
}

func (s *sseEndpoint) TrigerEvent(w http.ResponseWriter, r *http.Request) {
	s.doBackground()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(time.Now().Format(time.RFC3339)))
}

func (s *sseEndpoint) addClient(ch chan []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[ch] = struct{}{}
}

func (s *sseEndpoint) delClient(ch chan []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, ch)
	close(ch)
}

func (s *sseEndpoint) doBackground() {
	event := Event{Name: "CREATE_TODO", Value: "TODO"}
	data, err := s.codecJSON.Encode(event)
	if err != nil {
		log.Println("Failed to encode event:", err)

		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	for client := range s.clients {
		// Send data asynchronously to avoid blocking on slow clients
		go func(ch chan []byte) {
			ch <- data
		}(client)
	}
}

func (s *sseEndpoint) HandleEvent(w http.ResponseWriter, r *http.Request) {
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
