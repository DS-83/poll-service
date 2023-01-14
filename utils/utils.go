package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"
)

func CreateId() string {
	n := 1000000
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(n) + rand.Intn(n*2)
	h := sha1.Sum([]byte(strconv.Itoa(i)))
	s := hex.EncodeToString(h[:])
	return s
}
