// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// +build !integration

package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	client_id      = "O2IDtest"
	client_sercret = "secret"
)

func TestEncodeToBase64(t *testing.T) {
	a := assert.New(t)

	cl := NewSymantecClient("aaa", "", "", client_id, client_sercret)
	sign := cl.encodeToBase64()

	a.Equal("TzJJRHRlc3Q6c2VjcmV0", sign)
}
