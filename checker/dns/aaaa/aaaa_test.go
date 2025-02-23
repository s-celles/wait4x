// Copyright 2019-2025 The Wait4X Authors
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

package aaaa

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"wait4x.dev/v2/checker"
)

const server = "wait4x.dev"

func TestCheckExistenceAAAA(t *testing.T) {
	d := New(server)
	assert.Nil(t, d.Check(context.Background()))
}

func TestCorrectAAAA(t *testing.T) {
	d := New(server, WithExpectedIPV6s([]string{"2606:4700:3034::6815:591"}))
	assert.Nil(t, d.Check(context.Background()))
}

func TestIncorrectAAAA(t *testing.T) {
	var expectedError *checker.ExpectedError
	d := New(server, WithExpectedIPV6s([]string{"127.0.0.1"}))
	assert.ErrorAs(t, d.Check(context.Background()), &expectedError)
}

func TestCustomNSCorrectAAAA(t *testing.T) {
	d := New(server, WithNameServer("8.8.8.8:53"), WithExpectedIPV6s([]string{"2606:4700:3034::6815:591"}))
	assert.Nil(t, d.Check(context.Background()))
}
