package analytics

import (
	"bytes"

	"github.com/dchest/siphash"
)

type UserID uint64

func NewUserID(salt uint64, domain, ip, userAgent string) UserID {
	var buf bytes.Buffer
	total := len(domain) + len(ip) + len(userAgent)
	buf.Grow(total)

	buf.WriteString(domain)
	buf.WriteString(ip)
	buf.WriteString(userAgent)

	return UserID(siphash.Hash(salt, 0, buf.Bytes()))
}
