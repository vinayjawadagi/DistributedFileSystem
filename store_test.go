package main

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "hi how are you bruh"

	pathKey := CASPathTransformFunc(key)

	expectedPathName := "a8657/b2b50/8f114/51d28/ac942/85a85/e44ce/65362"
	expectedFileName := "a8657b2b508f11451d28ac94285a85e44ce65362"

	if pathKey.PathName != expectedPathName {
		t.Errorf("wrong path name expected: %s got : %s", pathKey.PathName, expectedPathName)
	}

	if pathKey.FileName != expectedFileName {
		t.Errorf("wrong file name expected: %s got : %s", pathKey.FileName, expectedFileName)
	}
}

func TestStoreDeleteKey(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	store := NewStore(opts)

	key := "hellohibyebye"
	data := []byte("why you still here? i said bye")

	err := store.writeStream(key, bytes.NewReader(data))

	if err != nil {
		t.Error(err)
	}

	err = store.Delete(key)

	if err != nil {
		t.Error(err)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	s := NewStore(opts)

	key := "sumn special"
	data := []byte("some jpg bytes")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if ok := s.Has(key); !ok {
		t.Errorf("Expected to have key: %s", key)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := ioutil.ReadAll(r)

	if string(b) != string(data) {
		t.Error()
	}

	s.Delete(key)
}
