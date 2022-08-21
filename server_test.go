package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

func TestServer_Run(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)

	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen port %v", err)
	}
	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	})

	eg.Go(func() error {
		s := NewServer(l, mux)
		return s.Run(ctx)
	})
	time.Sleep(1 * time.Second)
	in := "message"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)
	t.Logf("try request to %q", url)

	rsp, err := http.Get(url)
	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}
	defer rsp.Body.Close()
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %+v", err)
	}

	want := fmt.Sprintf("Hello, %s!", in)
	if string(got) != want {
		t.Errorf("want: %q, but got: %q", want, got)
	}

	cancel()

	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
