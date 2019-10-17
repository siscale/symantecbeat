// +build !integration

package client

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
