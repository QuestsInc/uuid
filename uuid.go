package uuid

import (
	"crypto/rand"
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// UUID is 128bits = 16bytes
type UUID [16]byte

func NewMust() *UUID {
	id, err := New()
	if err != nil {
		return NewMust()
	}

	return id
}

// New generates new unique UUID v4
func New() (*UUID, error) {
	uuid := new(UUID)

	n, err := io.ReadFull(rand.Reader, uuid[:])
	if err != nil {
		return nil, err
	}

	if n != len(uuid) {
		return nil, errors.New(fmt.Sprintf("insufficient random data (expected: %d, read: %d)", len(uuid), n))
	}

	// variant bits; for more info
	// see https://www.ietf.org/rfc/rfc4122.txt section 4.1.1
	uuid[8] = uuid[8]&0x3f | 0x80

	// version 4 (pseudo-random); for more info
	// see https://www.ietf.org/rfc/rfc4122.txt section 4.1.3
	uuid[6] = uuid[6]&0x0f | 0x40

	return uuid, nil
}

// String stringify UUID to 'xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx' format
func (uuid UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func (uuid UUID) RawString() string {
	return fmt.Sprintf("%x", uuid[:])
}

func (uuid UUID) IsZero() bool {
	var zeroUuid UUID
	return Equal(zeroUuid, uuid)
}

func (uuid UUID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + uuid.String() + `"`), nil
}

func (uuid *UUID) UnmarshalJSON(b []byte) error {
	if u, err := Parse(string(b)); err != nil {
		return err
	} else {
		copy(uuid[:], u[:])
		return nil
	}
}

// Equal compares uuid by bytes
func Equal(uuid1 UUID, uuid2 UUID) bool {
	for i, v := range uuid1 {
		if v != uuid2[i] {
			return false
		}
	}

	return true
}

func (uuid UUID) Value() (driver.Value, error) {
	return uuid[:], nil
}

func (uuid *UUID) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		// return Error.New("unexpected type %T for uuid", src)
		*uuid = UUID{}
		return nil
	}

	switch len(b) {
	case 32:
		dst := make([]byte, 16)
		_, err := hex.Decode(dst, b)
		if err != nil {
			return err
		}

		b = dst
		fallthrough

	case 16:
		toUUID, err := FromBytes(b)
		if err != nil {
			return err
		}

		*uuid = toUUID
		return nil

	default:
		id, err := Parse(string(b))
		if err != nil {
			return err
		}

		*uuid = *id
		return nil
	}
}
