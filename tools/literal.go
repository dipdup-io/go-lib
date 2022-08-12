package tools

import (
	"regexp"
	"strings"

	"github.com/dipdup-net/go-lib/tools/consts"
)

// IsLiteral -
func IsLiteral(prim string) bool {
	for _, s := range []string{
		consts.CONTRACT, consts.BYTES, consts.ADDRESS, consts.KEYHASH,
		consts.KEY, consts.TIMESTAMP, consts.BOOL, consts.MUTEZ,
		consts.NAT, consts.STRING, consts.INT, consts.SIGNATURE,
	} {
		if prim == s {
			return true
		}
	}
	return false
}

// IsContract -
func IsContract(address string) bool {
	return len(address) == 36 && strings.HasPrefix(address, "KT")
}

// IsAddressLazy -
func IsAddressLazy(address string) bool {
	return len(address) == 36 && (strings.HasPrefix(address, "KT") || strings.HasPrefix(address, "tz"))
}

var (
	addressRegex       = regexp.MustCompile("(tz|KT)[0-9A-Za-z]{34}")
	operationHashRegex = regexp.MustCompile("(o)[0-9A-Za-z]{50}")
	bigMapKeyHashRegex = regexp.MustCompile("(expr)[0-9A-Za-z]{50}")
	bakerHashRegex     = regexp.MustCompile("(SG1)[0-9A-Za-z]{33}")
)

// IsAddress -
func IsAddress(str string) bool {
	return addressRegex.MatchString(str)
}

// IsOperationHash -
func IsOperationHash(str string) bool {
	return operationHashRegex.MatchString(str)
}

// IsBigMapKeyHash -
func IsBigMapKeyHash(str string) bool {
	return bigMapKeyHashRegex.MatchString(str)
}

// IsBakerHash -
func IsBakerHash(str string) bool {
	return bakerHashRegex.MatchString(str)
}
