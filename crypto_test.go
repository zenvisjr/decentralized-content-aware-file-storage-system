package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyEncrypt(t *testing.T) {
	payload := "large file mimic"
	src := bytes.NewReader([]byte(payload))
	dist := new(bytes.Buffer)
	key := newEncryptionKey()

	_, err := copyEncrypt(dist, src, key)
	assert.NoError(t, err)

	fmt.Println(dist.Bytes())

	out := new(bytes.Buffer)
	n, err := copyDecrypt(out, dist, key)
	assert.NoError(t, err)

	assert.Equal(t, n, 16+len(payload))
	assert.Equal(t, payload, out.String())

}
