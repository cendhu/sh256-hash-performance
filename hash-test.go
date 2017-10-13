package main

import (
	"fmt"
	"time"

	"crypto/rand"
	"crypto/sha256"
)

// getRandomBytes returns a random bytes of a given size
func getRandomBytes(size int) []byte {
	bytes := make([]byte, size)
	rand.Read(bytes)
	return bytes
}

func main() {

	// size of the key and number of iterations
	minKeySize := 1
	maxKeySize := 2000
	keySizeIncr := 1
	maxIteration := 1000

	// for each key size, compute hash and elapsed time.
	for keySize := minKeySize; keySize <= maxKeySize; keySize = keySize + keySizeIncr {

		var totalElapsed time.Duration
		var data []byte
		// measure time taken to compute the hash for a given key size
		for iter := 1; iter <= maxIteration; iter++ {
			// get random bytes of size equal to keySize
			data = getRandomBytes(keySize)

			start := time.Now()
			hasher := sha256.New()
			hasher.Write(data)
			hash := hasher.Sum(nil)
			_ = hash
			totalElapsed += time.Since(start)
		}
		elapsed := totalElapsed.Nanoseconds() / int64(maxIteration)

		fmt.Printf("HASH: keySize %d hashComputeTime %d ns\n", len(data), elapsed)
	}
}
