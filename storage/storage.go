package storage

import (
	"fmt"

	"github.com/peterbourgon/diskv"
)

var storage *diskv.Diskv

func Start() {
	// Simplest transform function: put all the data files into the base dir.
	flatTransform := func(s string) []string { return []string{} }

	// Initialize a new diskv store, rooted at "my-data-dir", with a 1MB cache.
	storage = diskv.New(diskv.Options{
		BasePath:     "alexa-storage",
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})

	key := "alpha"
	storage.Write(key, []byte{'1', '2', '3'})

	// Read the value back out of the store.
	value, _ := storage.Read(key)
	fmt.Printf("%v\n", value)
}

func Set(key string, value []byte) {
	storage.Write(key, value)
}

func Get(key string) ([]byte, error) {
	return storage.Read(key)
}
