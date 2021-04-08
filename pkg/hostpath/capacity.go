/*
Copyright 2021 The Kubernetes Authors.

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

package hostpath

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/api/resource"
)

// Capacity simulates linear storage of certain types ("fast",
// "slow"). When volumes of those types get created, they must
// allocate storage (which can fail!) and that storage must
// be freed again when volumes get destroyed.
//
// Available capacity is configurable with a command line flag
// -capacity <type>=<size> where <type> is a string and <size>
// is a quantity (1T, 1Gi). More than one of those
// flags can be used.
//
// The underlying map will be initialized if needed by Set,
// which makes it possible to define and use a Capacity instance
// without explicit initialization (`var capacity Capacity` or as
// member in a struct).
type Capacity map[string]resource.Quantity

// Set is an implementation of flag.Value.Set.
func (c *Capacity) Set(arg string) error {
	parts := strings.SplitN(arg, "=", 2)
	if len(parts) != 2 {
		return errors.New("must be of format <type>=<size>")
	}
	quantity, err := resource.ParseQuantity(parts[1])
	if err != nil {
		return err
	}

	// We overwrite any previous value.
	if *c == nil {
		*c = Capacity{}
	}
	(*c)[parts[0]] = quantity
	return nil
}

func (c *Capacity) String() string {
	return fmt.Sprintf("%v", map[string]resource.Quantity(*c))
}

var _ flag.Value = &Capacity{}

// Alloc reserves a certain amount of bytes. Errors are
// usable as result of gRPC calls. Empty kind means
// that any large enough one is fine.
func (c *Capacity) Alloc(kind string, size int64) (actualKind string, err error) {
	requested := *resource.NewQuantity(size, resource.BinarySI)

	if kind == "" {
		for k, quantity := range *c {
			if quantity.Value() >= size {
				kind = k
				break
			}
		}
		// Still nothing?
		if kind == "" {
			available := c.Check("")
			return "", status.Error(codes.ResourceExhausted,
				fmt.Sprintf("not enough capacity: have %s, need %s", available.String(), requested.String()))
		}
	}

	available, ok := (*c)[kind]
	if !ok {
		return "", status.Error(codes.InvalidArgument, fmt.Sprintf("unknown capacity kind: %q", kind))
	}
	if available.Cmp(requested) < 0 {
		return "", status.Error(codes.ResourceExhausted,
			fmt.Sprintf("not enough capacity of kind %q: have %s, need %s", kind, available.String(), requested.String()))
	}
	available.Sub(requested)
	(*c)[kind] = available
	return kind, nil
}

// Free returns capacity reserved earlier with Alloc.
func (c *Capacity) Free(kind string, size int64) {
	available := (*c)[kind]
	available.Add(*resource.NewQuantity(size, resource.BinarySI))
	(*c)[kind] = available
}

// Check reports available capacity for a certain kind.
// If empty, it reports the maximum capacity.
func (c *Capacity) Check(kind string) resource.Quantity {
	if kind != "" {
		quantity := (*c)[kind]
		return quantity
	}
	available := resource.Quantity{}
	for _, q := range *c {
		if q.Cmp(available) >= 0 {
			available = q
		}
	}
	return available
}

// Enabled returns true if capacities are configured.
func (c *Capacity) Enabled() bool {
	return len(*c) > 0
}
