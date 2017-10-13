package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"crypto/rand"
	"crypto/sha256"
	"github.com/syndtr/goleveldb/leveldb"
	mrand "math/rand"
)

// keyHash holds the key, length of the key (in terms of bytes) and hash of the key
type keyHash struct {
	Key    string
	KeyLen int
	Hash   []byte
}

// dbRead holds time taken to read a key of specified length
type dbRead struct {
	KeyLen  int
	ValLen  int
	elapsed uint64
}

// getRandomBytes returns a random bytes of a given size
func getRandomBytes(size int) []byte {
	bytes := make([]byte, size)
	rand.Read(bytes)
	return bytes
}

func main() {

	// size of the key and number of iterations
	// total number of keys written to db = 2000 * 1000 = 2 million entries
	// db size = 2.2 GB
	minKeySize := 1
	maxKeySize := 2000
	keySizeIncr := 1
	maxIteration := 1000

	// total number of keys to be read = maxKeySize * 3 = 2000 * 3 = 6000 keys
	maxDBRead := maxKeySize * 3 // for each key size, read 3 items from db
	// read 6000 keys 100 times and compute the average
	maxDBReadIteration := 100

	var hash []byte
	var keyToHash []*keyHash

	// open database to store key and hash
	path := "/tmp/test.db"
	os.Remove(path)
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		fmt.Println("ERROR: Unable to open leveldb database")
		return
	}
	defer db.Close()

	for keySize := minKeySize; keySize <= maxKeySize; keySize = keySize + keySizeIncr {

		for iter := 1; iter <= maxIteration; iter++ {
			// get random bytes of size equal to keySize
			data := getRandomBytes(keySize)

			hasher := sha256.New()
			hasher.Write(data)
			hash = hasher.Sum(nil)

			// add key and hash to goleveldb as well
			err = db.Put(data, hash, nil)
			if err != nil {
				fmt.Println("ERROR: Unable to write to db")
				return
			}

			// choose 3 keys per size to be read during read phase
			if iter == 1 || iter == maxIteration/2 || iter == maxIteration {
				keyToHash = append(keyToHash, &keyHash{string(data), len(data), hash})
			}
		}

	}

	// TEST: Read from goleveldb
	mrand.Seed(time.Now().Unix())
	// generate random number sequence in the range [0, maxDBRead].
	// as a result, we would read all keys.
	randomlyOrderedKeyIndices := mrand.Perm(maxDBRead)

	dbReadTime := make(map[int]uint64)

	// read the whole database multiple times
	for iter := 1; iter <= maxDBReadIteration; iter++ {
		// clear buffer cache during every iteration so that read always goes to disk
		cmd := exec.Command("sh", "-c", "echo 1 > /proc/sys/vm/drop_caches")
		cmd.Run()
		// read keys stored on the db
		for nread := 0; nread < maxDBRead; nread++ {
			keyIndex := randomlyOrderedKeyIndices[nread]
			start := time.Now()
			value, err := db.Get([]byte(keyToHash[keyIndex].Key), nil)
			_ = value
			if err != nil {
				continue
			}
			dbReadTime[keyIndex] += uint64(time.Since(start).Nanoseconds())
		}
	}

	for nread := 1; nread <= maxDBRead; nread++ {
		index := nread - 1
		fmt.Printf("DB: keySize %d readTime %d ns \n", keyToHash[index].KeyLen, dbReadTime[index]/uint64(maxDBReadIteration))
	}

}
