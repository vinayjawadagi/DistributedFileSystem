package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
)

const defaultRootFolderName = "SomethingDefault"

type PathKey struct {
	PathName string
	FileName string
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.FileName)
}

// FirstPathName returns the 'root' directory of the file
func (p PathKey) FirstPathName() string {
	paths := strings.Split(p.PathName, "/")

	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

// Content addressable path transform function, transforms the given key to a path on the disk
func CASPathTransformFunc(key string) PathKey {
	// generate the hash for the key
	hash := sha1.Sum([]byte(key))
	// 1. [170 244 198 29 220 197 232 162 218 190 222 15 59 72 44 217 174 169 67 77]

	// Encodes the size 20 byte array to 40 digit hexadecimal string
	hashStr := hex.EncodeToString(hash[:]) // to convert fixed array to slice: [20]byte => []byte => [:]

	// 2. "a8657b2b508f11451d28ac94285a85e44ce65362"
	blockSize := 5

	sliceLen := len(hashStr) / blockSize

	paths := make([]string, sliceLen)

	// 3. [a8657 b2b50 8f114 51d28 ac942 85a85 e44ce 65362]
	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}

	// 4. a8657/b2b50/8f114/51d28/ac942/85a85/e44ce/65362
	//return strings.Join(paths, "/")
	return PathKey{
		PathName: strings.Join(paths, "/"),
		FileName: hashStr,
	}
}

type PathTransformFunc func(string) PathKey

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		PathName: key,
		FileName: key,
	}
}

type StoreOpts struct {
	// Root is the folder name of the root, containing all the folder/files of the system
	Root              string
	PathTransformFunc PathTransformFunc
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}
	if len(opts.Root) == 0 {
		opts.Root = defaultRootFolderName
	}
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransformFunc(key)

	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())
	_, err := os.Stat(fullPathWithRoot)
	if fs.ErrNotExist == err {
		return false
	}
	return true
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)

	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.FileName)
	}()

	pathKeyWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FirstPathName())
	// remove the parent directory
	return os.RemoveAll(pathKeyWithRoot)
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())
	return os.Open(fullPathWithRoot)
}

func (s *Store) writeStream(key string, r io.Reader) error {

	// creates the path of the file based on the key
	pathKey := s.PathTransformFunc(key)

	// path without the filename
	pathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.PathName)
	// creates file within the path structure (MkDir creates just the file)
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return err
	}

	// path with the file name
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.FullPath())

	// creates a new file, if file already exists it truncates the file(clears the contents)
	f, err := os.Create(fullPathWithRoot)
	if err != nil {
		return err
	}

	// copys the contents of the file from the buf
	n, err := io.Copy(f, r)

	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s", n, fullPathWithRoot)

	return nil
}
