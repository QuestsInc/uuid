package uuid

import (
	"encoding/hex"
	"errors"
	"strings"
)

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
