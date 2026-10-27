package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestOverlayHTTP(t *testing.T) {
	store := &overlayStore{}
	want := overlayData{
		ThreeStarCount: 103, FourStarCount: 64, FiveStarCount: 18, SixStarCount: 5,
		PityCount: 51, TotalPullCount: 190,
		ThreeStarNew: 4, FourStarNew: 9, FiveStarNew: 6, SixStarNew: 3,
		TotalNew: 22, SixStarChance: 6,
	}
	store.publish(want)
	handler := newOverlayHandler(store)

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/state", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("GET /state: status %d", response.Code)
	}
	if response.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("counter responses must disable caching")
	}
	if !strings.HasPrefix(response.Header().Get("Content-Type"), "application/json") {
		t.Fatal("counter response must be JSON")
	}
	var got overlayData
	if err := json.Unmarshal(response.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("snapshot mismatch: got %+v, want %+v", got, want)
	}
	var properties map[string]int
	if err := json.Unmarshal(response.Body.Bytes(), &properties); err != nil {
		t.Fatal(err)
	}
	if len(properties) != 12 || properties["sixStarCount"] != 5 || properties["sixStarChance"] != 6 {
		t.Fatalf("unexpected JSON property names: %v", properties)
	}

	// A second request must see an updated snapshot, including decreases such
	// as undo and a reset to zero; the server must not cache the first state.
	store.publish(overlayData{SixStarChance: 2})
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/state", nil))
	if err := json.Unmarshal(response.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got != (overlayData{SixStarChance: 2}) {
		t.Fatalf("reset was not reflected: %+v", got)
	}

	for _, tc := range []struct {
		method string
		path   string
		status int
	}{
		{http.MethodGet, "/", http.StatusOK},
		{http.MethodHead, "/", http.StatusOK},
		{http.MethodHead, "/state", http.StatusOK},
		{http.MethodPost, "/", http.StatusMethodNotAllowed},
		{http.MethodPost, "/state", http.StatusMethodNotAllowed},
		{http.MethodGet, "/missing", http.StatusNotFound},
		{http.MethodGet, "/main.go", http.StatusNotFound},
	} {
		t.Run(tc.method+tc.path, func(t *testing.T) {
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, httptest.NewRequest(tc.method, tc.path, nil))
			if rr.Code != tc.status {
				t.Fatalf("status %d, want %d", rr.Code, tc.status)
			}
			if tc.method == http.MethodHead && rr.Body.Len() != 0 {
				t.Fatal("HEAD must not return a response body")
			}
			if tc.method == http.MethodGet && tc.path == "/" && rr.Body.String() != overlayHTML {
				t.Fatal("root did not serve the embedded page")
			}
		})
	}
}

func TestOverlayConcurrentSnapshots(t *testing.T) {
	store := &overlayStore{}
	var workers sync.WaitGroup
	workers.Add(5)
	go func() {
		defer workers.Done()
		for n := 0; n < 2000; n++ {
			store.publish(overlayData{TotalPullCount: n, SixStarCount: n, TotalNew: n})
		}
	}()
	for n := 0; n < 4; n++ {
		go func() {
			defer workers.Done()
			for j := 0; j < 2000; j++ {
				data := store.snapshot()
				if data.TotalPullCount != data.SixStarCount || data.TotalPullCount != data.TotalNew {
					t.Errorf("read part of two different snapshots: %+v", data)
					return
				}
			}
		}()
	}
	workers.Wait()
}

func TestOverlayServerLifecycle(t *testing.T) {
	store := &overlayStore{}
	store.publish(overlayData{TotalPullCount: 24, SixStarChance: 2})
	server, err := startOverlayServer(store, "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = server.Close() })
	client := &http.Client{Timeout: 2 * time.Second}
	response, err := client.Get("http://" + server.Addr + "/state")
	if err != nil {
		t.Fatal(err)
	}
	body, err := io.ReadAll(response.Body)
	_ = response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK || !strings.Contains(string(body), `"totalPullCount":24`) {
		t.Fatalf("unexpected real HTTP response: %s", body)
	}

	other, err := startOverlayServer(store, server.Addr)
	if err == nil {
		_ = other.Close()
		t.Fatal("a second server must report an occupied port")
	}
	if err := server.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Get("http://" + server.Addr + "/state"); err == nil {
		t.Fatal("closed server still accepted requests")
	}
	// Closing the server must free its port for a new tracker process.
	restarted, err := startOverlayServer(store, server.Addr)
	if err != nil {
		t.Fatalf("could not restart the overlay: %v", err)
	}
	_ = restarted.Close()
}
