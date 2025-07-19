package main

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCryptoPathTransformFunc(t *testing.T) {
	path := "watashinosoulsociety"
	pathname := CryptoPathTransformFunc(path)
	exprectedFileName := "cbd29d6f71beae1f2b5737d26a0f498a05a0b09f"
	exprectedPathName := "cbd29/d6f71/beae1/f2b57/37d26/a0f49/8a05a/0b09f"
	assert.Equal(t, exprectedPathName, pathname.FilePath)
	assert.Equal(t, exprectedFileName, pathname.FileName)
}
func TestStoreFile(t *testing.T) {
	s := newStore()
	id := generateID()
	defer teardown(t, s)

	for i := 0; i < 1; i++ {
		key := fmt.Sprintf("watashinosoulsociety_%d", i)
		data := "My name is zenvis"
		_, err, _ := s.Write(id, key, strings.NewReader(data))
		assert.NoError(t, err)

		exist := s.Has(id, key)
		assert.True(t, exist)

		_, buf, err, _ := s.Read(id, key)
		assert.NoError(t, err)

		b, err := io.ReadAll(buf)
		fmt.Println(string(b))
		assert.NoError(t, err)
		assert.Equal(t, string(b), data)

		s.Delete(id, key)

		exist = s.Has(id, key)
		assert.False(t, !exist)
	}

}

// func TestDeletePathAndFile(t *testing.T) {
// 	s := NewStore(StoreOps{
// 		PathTransformFunc: CryptoPathTransformFunc,
// 	})
// 	key := "watashinosoulsociety"
// 	data := "My name is zenvis"
// 	err := s.WriteStream(key, strings.NewReader(data))
// 	assert.NoError(t, err)

// 	err = s.DeletePathAndFile(key)
// 	assert.NoError(t, err)

// }

func newStore() *Store {
	ops := StoreOps{
		PathTransformFunc: CryptoPathTransformFunc,
	}
	s := NewStore(ops)
	return s
}

func teardown(t *testing.T, s *Store) {
	err := s.ClearRoot()
	assert.NoError(t, err)
}
