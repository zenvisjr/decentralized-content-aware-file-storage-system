package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyEncrypt(t *testing.T) {
	// payload := "large file mimic"
	// src := bytes.NewReader([]byte(payload))
	// dist := new(bytes.Buffer)
	// key := newEncryptionKey()

	// _, err := copyEncrypt(dist, src, key)
	// assert.NoError(t, err)

	// fmt.Println(dist.Bytes())

	// out := new(bytes.Buffer)
	// n, err := copyDecrypt(out, dist, key)
	// assert.NoError(t, err)

	// assert.Equal(t, n, 16+len(payload))
	// assert.Equal(t, payload, out.String())

	key := "bill.pdf"
	keywithoutExtension := removeExtension(key)
	hashkey1 := hashKey(key)
	hashkey2 := hashKey(key)
	
	fmt.Println(hashkey1)
	fmt.Println(hashkey2)
	assert.Equal(t, hashkey1, hashkey2)

	hashkey3 := hashKey(keywithoutExtension)
	fmt.Println(hashkey3)
	hashkey4 := hashKey(keywithoutExtension)
	fmt.Println(hashkey4)
	assert.Equal(t, hashkey3, hashkey4)

	hashkey5 := hashKey(hashKey(key))
	hashkey6 := hashKey(hashKey(key))
	
	fmt.Println(hashkey5)
	fmt.Println(hashkey6)
	assert.Equal(t, hashkey5, hashkey6)

	hashkey7 := hashKey(hashKey(keywithoutExtension))
	fmt.Println(hashkey7)
	hashkey8 := hashKey(hashKey(keywithoutExtension))
	fmt.Println(hashkey8)
	assert.Equal(t, hashkey7, hashkey8)

	

}
