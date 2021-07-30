package uuid

import "github.com/zeebo/errs"

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
