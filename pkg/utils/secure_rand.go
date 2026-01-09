package utils

import (
	"crypto/rand"
	"math/big"
)

func SecureRandInt(max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil || max <= 0 {
		return 0
	}
	return int(n.Int64())
}

func SecureRandFloat32() float32 {
	n, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		return 0
	}
	return float32(n.Int64()) / 10000.0
}
