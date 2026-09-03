package mock

import (
	"crypto/md5"
	"uuid"
)

func newUUID(seed string) uuid.UUID {
	b := md5.Sum([]byte(seed))
	// Cast the dereferenced 16-byte array pointer
	return uuid.UUID(b)
}
