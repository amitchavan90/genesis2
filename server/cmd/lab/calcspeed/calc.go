package main

/*
 * calculate and see how fast go can calc hashes
 */

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

func uuid2byte(uid string) []byte {
	return uuid.FromStringOrNil(uid).Bytes()
}

func byte2uuid(b []byte) string {
	return uuid.FromBytesOrNil(b).String()
}

func main() {
	id2s := []string{}
	for i := 0; i < 1000000; i++ {
		u, _ := uuid.NewV4()
		id2s = append(id2s, u.String())
	}

	t := time.Now()
	for _, id := range id2s {
		_ = sha256.Sum256(uuid2byte(id))
	}

	t2 := time.Now().Sub(t)
	fmt.Printf("took %s time to calc %d\n", t2, len(id2s))
}
