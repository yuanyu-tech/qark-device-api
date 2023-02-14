package device

import (
	"fmt"
	"testing"
)

func TestSign(t *testing.T) {
	appKey := "7344ca0736ebb67a60f85b7e9765ef61ad838e7500000001602f9ca9"
	noid := "6973886222418313494"
	seed := "1625737328"
	suffix := "__qark_env_sign"
	sign, err := BuildSign(appKey, noid, seed, suffix)
	if err != nil {
		fmt.Println("build sign err %w", err)
		return
	}

	fmt.Printf("sign(%s, %s, %s, %s)=%s\n", appKey, noid, seed, suffix, sign)
}
