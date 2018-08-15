package todo

import (
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randStringBytesMaskImpr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func fillStorageWithMockData(s *storage) (err error) {
	var id string
	var added = time.Now().Add(time.Duration(-30) * time.Second)
	var done = time.Now()
	id, _ = NewItemID()
	s.put(&Item{
		ID:    id,
		Note:  "Note1",
		Added: &added,
		Done:  &done,
	})
	id, _ = NewItemID()
	added = time.Now().Add(time.Duration(-30) * time.Second)
	s.put(&Item{
		ID:    id,
		Note:  "Note2",
		Added: &added,
	})
	return
}
