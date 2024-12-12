package data

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"

	"github.com/google/uuid"
)

var Hasher func() hash.Hash = func() hash.Hash {
	return sha256.New()
}

const IdSize = 20

// Represents a Musician ID
type ID [IdSize]byte

var HighestId ID = [IdSize]byte{
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
}

var LowestId ID = [IdSize]byte{}

func XorIds(a, b *ID) {
	a[0] = a[0] ^ b[0]
	a[1] = a[1] ^ b[1]
	a[2] = a[2] ^ b[2]
	a[3] = a[3] ^ b[3]
	a[4] = a[4] ^ b[4]
	a[5] = a[5] ^ b[5]
	a[6] = a[6] ^ b[6]
	a[7] = a[7] ^ b[7]
	a[8] = a[8] ^ b[8]
	a[9] = a[9] ^ b[9]
	a[10] = a[10] ^ b[10]
	a[11] = a[11] ^ b[11]
	a[12] = a[12] ^ b[12]
	a[13] = a[13] ^ b[13]
	a[14] = a[14] ^ b[14]
	a[15] = a[15] ^ b[15]
	a[16] = a[16] ^ b[16]
	a[16] = a[16] ^ b[16]
	a[17] = a[17] ^ b[17]
	a[18] = a[18] ^ b[18]
	a[19] = a[19] ^ b[19]
}

func HexFromIdList(ids []ID) []string {
	res := make([]string, len(ids))
	for i := range ids {
		res[i] = ids[i].Hex()
	}
	return res
}

func IdsFromHexList(ids []string) ([]ID, error) {
	var res []ID = make([]ID, len(ids))
	for i, hex := range ids {
		id, err := IdFromHex(hex)
		if err != nil {
			return nil, err
		}
		res[i] = id
	}
	return res, nil
}

// Id from hexadecimal converts hexadecimal to a ID
func IdFromHex(id string) (ID, error) {
	if id == "" {
		return LowestId, nil
	}

	bs, err := hex.DecodeString(id)
	if len(bs) != IdSize {
		return LowestId, errors.New("wrong size")
	}
	return ID(bs), err
}

// GenId generates a random ID
func GenId() ID {
	var id ID

	_, err := rand.Read(id[:])
	if err != nil {
		panic(fmt.Errorf("unable to generate id: %w", err))
	}

	return id
}

func GenIdFromBytes(content []byte) ID {
	h := Hasher()
	h.Write(content)
	bs := h.Sum(nil)
	return ID(bs)
}

// Lower returns true when the ID is lower than cmp ID
func (id ID) Lower(cmp ID) bool {
	for i := range id {
		if id[i] < cmp[i] {
			return true
		}
		if id[i] > cmp[i] {
			return false
		}
	}

	return false
}

// Higher returns true when the ID is higher than cmp ID
func (id ID) Higher(cmp ID) bool {
	for i := range id {
		if id[i] > cmp[i] {
			return true
		}
		if id[i] < cmp[i] {
			return false
		}
	}

	return false
}

func (id ID) HigherOrEqual(cmp ID) bool {
	for i := range id {
		if id[i] >= cmp[i] {
			return true
		}
		if id[i] < cmp[i] {
			return false
		}
	}

	return false
}

// Bytes returns the representation of ID in bytes
func (id ID) Bytes() []byte {
	return []byte(id[:])
}

// Hex converts ID to hexadecimal
func (id ID) Hex() string {
	return fmt.Sprintf("%x", id)
}

type IDS []ID

// HexList converts an ID list to a hexadecimal list
func (ids IDS) Hex() []string {
	var res []string = make([]string, len(ids))
	for i, ID := range ids {
		res[i] = ID.Hex()
	}
	return res
}

// GenFileId generates a random file id
func GenFileId() string {
	return uuid.NewString()
}
