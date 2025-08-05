package analytics

import (
	"github.com/dchest/siphash"
)

type UserID uint64

func NewUserID(salt, domain, ip, userAgent string) UserID {
	return UserID(0)
}
