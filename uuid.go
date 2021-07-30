package uuid

import (
	"crypto/rand"
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/zeebo/errs"
)

// UUID is 128bits = 16bytes
type UUID [16]byte

// String stringify UUID to 'xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx' format
func (uuid UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
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

// Equal compares uuid by bytes
func Equal(uuid1 UUID, uuid2 UUID) bool {
	for i, v := range uuid1 {
		if v != uuid2[i] {
			return false
		}
	}

	return true
}

// ParseOpt is equal to Parse, but in error case generate zeroed UUID
func ParseOpt(s string) *UUID {
	id, err := Parse(s)
	if err != nil {
		return &UUID{}
	}
	return id
}

// Parse generate UUID from string.
// Returns error in in case of invalid input format.
func Parse(s string) (*UUID, error) {
	// the string format should be either in
	// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx (or)
	// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

	// If the uuid is marshaled by us we add " " around the uuid.
	// while parsing this, we have to remove the " " around the
	// uuid. So we check if uuid has " " around it, if yes we remove
	// it.

	if strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") {
		s = s[1 : len(s)-1]
	}

	uuid := new(UUID)
	switch len(s) {
	case 36:
		if ba, err := hex.DecodeString(s[0:8]); err == nil {
			copy(uuid[:4], ba)
		} else {
			return nil, err
		}
		if ba, err := hex.DecodeString(s[9:13]); err == nil {
			copy(uuid[4:], ba)
		} else {
			return nil, err
		}
		if ba, err := hex.DecodeString(s[14:18]); err == nil {
			copy(uuid[6:], ba)
		} else {
			return nil, err
		}
		if ba, err := hex.DecodeString(s[19:23]); err == nil {
			copy(uuid[8:], ba)
		} else {
			return nil, err
		}
		if ba, err := hex.DecodeString(s[24:]); err == nil {
			copy(uuid[10:], ba)
		} else {
			return nil, err
		}
	case 32:
		if ba, err := hex.DecodeString(s); err == nil {
			copy(uuid[:], ba)
		} else {
			return nil, err
		}
	default:
		return nil, errors.New("unknown UUID string " + s)
	}

	return uuid, nil
}

// FromBytes is used to convert byte slice to UUID
func FromBytes(data []byte) (UUID, error) {
	var id UUID

	copy(id[:], data)
	if len(id) != len(data) {
		return UUID{}, errs.New("Invalid uuid")
	}

	return id, nil
}

// FromBytesOpt is equal to FromBytes, but in error case generate zeroed UUID
func FromBytesOpt(data []byte) UUID {
	var id UUID

	copy(id[:], data)
	if len(id) != len(data) {
		return UUID{}
	}

	return id
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

	toUUID, err := FromBytes(b)
	if err != nil {
		return err
	}

	*uuid = toUUID
	return nil
}
