package uuid

import (
	"errors"
	"testing"
)

func TestNewMust(t *testing.T) {
	id := NewMust()

	err := validateID(id)
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkNewMust(b *testing.B) {
	for i := 0; i < b.N; i++ {
		id := NewMust()

		err := validateID(id)
		if err != nil {
			b.Error(err)
		}
	}
}

func validateID(id *UUID) error {
	if id == nil {
		return errors.New("uuid should not be nil")
	}

	if len(id) != 16 {
		return errors.New("uuid should be 16 bytes long")
	}

	if id.String() == "" {
		return errors.New("uuid should not be empty")
	}

	if id.String() == "00000000-0000-0000-0000-000000000000" {
		return errors.New("uuid should not be zero")
	}

	return nil
}
