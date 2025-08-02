package contract

import (
	"time"
)

type (
	Cache interface {
		Get(key string) (any, bool)
		Set(key string, value any, ttl time.Duration) error
		Delete(key string) error
		Flush() error
	}
)
