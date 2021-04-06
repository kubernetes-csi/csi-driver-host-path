/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package proxy

import (
	"bytes"
	"context"
	"io"
	"net"
	"testing"
)

func TestProxy(t *testing.T) {
	tmpdir := t.TempDir()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	endpoint1 := tmpdir + "/a.sock"
	endpoint2 := tmpdir + "/b.sock"

	closer, err := Run(ctx, endpoint1, endpoint2)
	if err != nil {
		t.Fatalf("proxy error: %v", err)
	}
	defer closer.Close()

	t.Run("a-to-b", func(t *testing.T) {
		sendReceive(t, endpoint1, endpoint2)
	})
	t.Run("b-to-a", func(t *testing.T) {
		sendReceive(t, endpoint2, endpoint1)
	})
}

func sendReceive(t *testing.T, endpoint1, endpoint2 string) {
	conn1, err := net.Dial("unix", endpoint1)
	if err != nil {
		t.Fatalf("error connecting to first endpoint %s: %v", endpoint1, err)
	}
	defer conn1.Close()
	conn2, err := net.Dial("unix", endpoint2)
	if err != nil {
		t.Fatalf("error connecting to second endpoint %s: %v", endpoint2, err)
	}
	defer conn2.Close()

	req1 := "ping"
	if _, err := conn1.Write([]byte(req1)); err != nil {
		t.Fatalf("error writing %q: %v", req1, err)
	}
	buffer := make([]byte, 100)
	len, err := conn2.Read(buffer)
	if err != nil {
		t.Fatalf("error reading %q: %v", req1, err)
	}
	if string(buffer[:len]) != req1 {
		t.Fatalf("expected %q, got %q", req1, string(buffer[:len]))
	}

	resp1 := "pong-pong"
	if _, err := conn2.Write([]byte(resp1)); err != nil {
		t.Fatalf("error writing %q: %v", resp1, err)
	}
	buffer = make([]byte, 100)
	len, err = conn1.Read(buffer)
	if err != nil {
		t.Fatalf("error reading %q: %v", resp1, err)
	}
	if string(buffer[:len]) != resp1 {
		t.Fatalf("expected %q, got %q", resp1, string(buffer[:len]))
	}

	// Closing one side should be noticed at the other end.
	err = conn1.Close()
	if err != nil {
		t.Fatalf("error closing connection to %s: %v", endpoint1, err)
	}
	len2, err := io.Copy(&bytes.Buffer{}, conn2)
	if err != nil {
		t.Fatalf("error reading from %s: %v", endpoint2, err)
	}
	if len2 != 0 {
		t.Fatalf("unexpected data via %s: %d", endpoint2, len2)
	}
	err = conn2.Close()
	if err != nil {
		t.Fatalf("error closing connection to %s: %v", endpoint2, err)
	}
}
