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

// Package proxy makes it possible to forward a listening socket in
// situations where the proxy cannot connect to some other address.
// Instead, it creates two listening sockets, pairs two incoming
// connections and then moves data back and forth. This matches
// the behavior of the following socat command:
// socat -d -d -d UNIX-LISTEN:/tmp/socat,fork TCP-LISTEN:9000,reuseport
//
// The advantage over that command is that both listening
// sockets are always open, in contrast to the socat solution
// where the TCP port is only open when there actually is a connection
// available.
//
// To establish a connection, someone has to poll the proxy with a dialer.
package proxy

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/golang/glog"

	"github.com/kubernetes-csi/csi-driver-host-path/internal/endpoint"
)

// New listens on both endpoints and starts accepting connections
// until closed or the context is done.
func Run(ctx context.Context, endpoint1, endpoint2 string) (io.Closer, error) {
	proxy := &proxy{}
	failedProxy := proxy
	defer func() {
		if failedProxy != nil {
			failedProxy.Close()
		}
	}()

	proxy.ctx, proxy.cancel = context.WithCancel(ctx)

	var err error
	proxy.s1, proxy.cleanup1, err = endpoint.Listen(endpoint1)
	if err != nil {
		return nil, fmt.Errorf("listen %s: %v", endpoint1, err)
	}
	proxy.s2, proxy.cleanup2, err = endpoint.Listen(endpoint2)
	if err != nil {
		return nil, fmt.Errorf("listen %s: %v", endpoint2, err)
	}

	glog.V(3).Infof("proxy listening on %s and %s", endpoint1, endpoint2)

	go func() {
		for {
			// We block on the first listening socket.
			// The Linux kernel proactively accepts connections
			// on the second one which we will take over below.
			conn1 := accept(proxy.ctx, proxy.s1, endpoint1)
			if conn1 == nil {
				// Done, shut down.
				glog.V(5).Infof("proxy endpoint %s closed, shutting down", endpoint1)
				return
			}
			conn2 := accept(proxy.ctx, proxy.s2, endpoint2)
			if conn2 == nil {
				// Done, shut down. The already accepted
				// connection gets closed.
				glog.V(5).Infof("proxy endpoint %s closed, shutting down and close established connection", endpoint2)
				conn1.Close()
				return
			}

			glog.V(3).Infof("proxy established a new connection between %s and %s", endpoint1, endpoint2)
			go copy(conn1, conn2, endpoint1, endpoint2)
			go copy(conn2, conn1, endpoint2, endpoint1)
		}
	}()

	failedProxy = nil
	return proxy, nil
}

type proxy struct {
	ctx                context.Context
	cancel             func()
	s1, s2             net.Listener
	cleanup1, cleanup2 func()
}

func (p *proxy) Close() error {
	if p.cancel != nil {
		p.cancel()
	}
	if p.s1 != nil {
		p.s1.Close()
	}
	if p.s2 != nil {
		p.s2.Close()
	}
	if p.cleanup1 != nil {
		p.cleanup1()
	}
	if p.cleanup2 != nil {
		p.cleanup2()
	}
	return nil
}

func copy(from, to net.Conn, fromEndpoint, toEndpoint string) {
	glog.V(5).Infof("starting to copy %s -> %s", fromEndpoint, toEndpoint)
	// Signal recipient that no more data is going to come.
	// This also stops reading from it.
	defer to.Close()
	// Copy data until EOF.
	cnt, err := io.Copy(to, from)
	glog.V(5).Infof("done copying %s -> %s: %d bytes, %v", fromEndpoint, toEndpoint, cnt, err)
}

func accept(ctx context.Context, s net.Listener, endpoint string) net.Conn {
	for {
		c, err := s.Accept()
		if err == nil {
			return c
		}
		// Ignore error if we are shutting down.
		if ctx.Err() != nil {
			return nil
		}
		glog.V(3).Infof("accept on %s failed: %v", endpoint, err)
	}
}
