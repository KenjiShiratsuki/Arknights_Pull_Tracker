package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

const overlayAddress = "127.0.0.1:8080"
const overlayURL = "http://" + overlayAddress + "/"

// The compiler copies this file into the executable at build time.
//
//go:embed overlay.html
var overlayHTML string

// These fields are exported so encoding/json can read them. The tags choose
// the property names that JavaScript receives; State itself stays unchanged.
type overlayData struct {
	ThreeStarCount int `json:"threeStarCount"`
	FourStarCount  int `json:"fourStarCount"`
	FiveStarCount  int `json:"fiveStarCount"`
	SixStarCount   int `json:"sixStarCount"`
	PityCount      int `json:"pityCount"`
	TotalPullCount int `json:"totalPullCount"`
	ThreeStarNew   int `json:"threeStarNew"`
	FourStarNew    int `json:"fourStarNew"`
	FiveStarNew    int `json:"fiveStarNew"`
	SixStarNew     int `json:"sixStarNew"`
	TotalNew       int `json:"totalNew"`
	SixStarChance  int `json:"sixStarChance"`
}

// Only this copy is shared between Fyne and HTTP handlers. The history and
// widgets remain on the UI goroutine; the server never reads them directly.
type overlayStore struct {
	mu   sync.RWMutex
	data overlayData
}

func (s *overlayStore) publish(data overlayData) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = data
}

func (s *overlayStore) snapshot() overlayData {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data
}

func newOverlayHandler(store *overlayStore) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		if r.Method == http.MethodHead {
			return
		}
		// snapshot returns a value copy and releases its lock before any I/O.
		if err := json.NewEncoder(w).Encode(store.snapshot()); err != nil {
			log.Printf("OBS overlay: writing state: %v", err)
		}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		if r.Method == http.MethodHead {
			return
		}
		if _, err := io.WriteString(w, overlayHTML); err != nil {
			log.Printf("OBS overlay: writing page: %v", err)
		}
	})

	return mux
}

func startOverlayServer(store *overlayStore, address string) (*http.Server, error) {
	// Listen before starting a goroutine so a port conflict is reported to the
	// caller immediately, where it can be shown in the Fyne window.
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("listen on %s: %w", address, err)
	}

	server := &http.Server{
		Addr:              listener.Addr().String(),
		Handler:           newOverlayHandler(store),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Serve blocks until shutdown, so it must run alongside the Fyne event loop.
	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("OBS overlay: server stopped: %v", err)
		}
	}()

	return server, nil
}
