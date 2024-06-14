package util

import (
	"crypto/sha1"
	"fmt"
)

func Sha1(s string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(s)))
}
