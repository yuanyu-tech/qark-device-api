package device

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash"
)

func BuildSign(args ...string) (string, error) {
	return BuildSignBase(sha1.New(), args...)
}

func BuildSignBase(hash hash.Hash, args ...string) (string, error) {
	signSeeds := "=~!@#$%^&*()_+{}|\\;:',./<>\"%%`~?~"
	seedsLen := len(signSeeds)
	random := make([]byte, 4)

	if _, err := hash.Write([]byte("^")); err != nil {
		return "", fmt.Errorf("sign write error for prefix")
	}

	for _, seed := range args {
		v := seed
		z := len(v)

		// mix the current seed first
		if _, err := hash.Write([]byte(v)); err != nil {
			return "", fmt.Errorf("sign write error for seed %s", v)
		}

		// generate the random seed
		h := uint32(v[0] & 0xFF)
		h = h*uint32(131) + uint32(v[z/2]&0xFF)
		h = h*uint32(1331) + uint32(v[z-1]&0xFF)
		h = h & 0x7FFFFFFF
		i := int(h % uint32(seedsLen))
		for z := 0; z < 3; z++ {
			random[z] = signSeeds[i]
			i++
			if i >= seedsLen {
				i = 0
			}
		}
		random[3] = '|'

		// mix the random salt values
		if _, err := hash.Write(random); err != nil {
			return "", fmt.Errorf("sign write error for random seed %s", v)
		}
	}

	if _, err := hash.Write([]byte("$")); err != nil {
		return "", fmt.Errorf("sign write error for suffix")
	}

	return hex.EncodeToString(hash.Sum([]byte(""))), nil
}
