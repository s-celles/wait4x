package ns

// Copyright 2023 The Wait4X Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"context"
	"fmt"
	"net"
	"regexp"

	"wait4x.dev/v2/checker"
)

// Option configures an CNAME DNS.
type Option func(d *NS)

// NS data structure.
type NS struct {
	nameserver          string
	address             string
	expectedNameservers []string
	resolver            *net.Resolver
}

// New creates the DNS TXT checker
func New(address string, opts ...Option) checker.Checker {
	d := &NS{
		address: address,
	}

	// apply the list of options to HTTP
	for _, opt := range opts {
		opt(d)
	}

	// Nameserver settings.
	d.resolver = net.DefaultResolver
	if d.nameserver != "" {
		d.resolver = &net.Resolver{
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				dialer := net.Dialer{}
				return dialer.DialContext(ctx, network, d.nameserver)
			},
		}
	}

	return d
}

func WithNameServer(nameserver string) Option {
	return func(d *NS) {
		d.nameserver = nameserver
	}
}

func WithExpectedNameservers(nameservers []string) Option {
	return func(d *NS) {
		d.expectedNameservers = nameservers
	}
}

// Identity returns the identity of the checker
func (d *NS) Identity() (string, error) {
	return fmt.Sprintf("NS %s %s", d.address, d.expectedNameservers), nil
}

// Check checks DNS TXT records
func (d *NS) Check(ctx context.Context) (err error) {
	values, err := d.resolver.LookupNS(ctx, d.address)
	if err != nil {
		return err
	}

	for _, ns := range values {
		if len(d.expectedNameservers) == 0 {
			return nil
		}
		for _, expectedNameserver := range d.expectedNameservers {
			matched, _ := regexp.MatchString(expectedNameserver, ns.Host)
			if matched {
				return nil
			}
		}
	}

	return checker.NewExpectedError(
		"the NS record value doesn't expect", nil,
		"actual", values, "expect", d.expectedNameservers,
	)
}
