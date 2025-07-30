package generatePackID

import (
	"crypto/rand"
	"encoding/binary"
)

func GeneratePackID() (int64, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return 0, err
	}

	id := int64(binary.BigEndian.Uint64(buf[:8]))
	if id < 0 {
		id = -id
	}
	return id, nil
}
