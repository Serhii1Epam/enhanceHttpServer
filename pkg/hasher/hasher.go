/* Package implements simple HTTP server functions */
package hasher

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

type HashingData struct {
	Pass, Hash string
}

func NewHasher(p string) *HashingData {
	return &HashingData{Pass: p}
}

// Hash the pasword
func (h *HashingData) HashPassword() (string, error) {
	hash := sha256.New()
	numBytes, err := hash.Write([]byte(h.Pass))

	if err == nil && numBytes > 0 {
		h.Hash = fmt.Sprintf("%x", string(hash.Sum(nil)))
		return h.Hash, err
	}

	return "", err
}

// Check pasword hash
func (h *HashingData) CheckPasswordHash(hashed string) bool {
	ret := false
	calculatedHash, _ := h.HashPassword()
	if h.Hash != "" {
		if result := strings.Compare(calculatedHash, hashed); result == 0 {
			ret = true
		}
	}
	return ret
}
