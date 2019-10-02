package monero

import (
	"bytes"
	"errors"
	"github.com/snipa22/monero-support/base58"
	"github.com/snipa22/monero-support/crypto"
)

const ChecksumSize = 4

var Tag byte
var Tags []byte

type Coin int
type Network int
type AddressType int

const (
	Monero Coin = 0
)

const (
	Testnet  Network = 0
	Mainnet  Network = 1
	Stagenet Network = 2
)

const (
	Normal     AddressType = 0
	Integrated AddressType = 1
	Subaddress AddressType = 2
)

func SetValidTags(c Coin, n Network) {
	switch c {
	case Monero:
		switch n {
		case Mainnet:
			Tags = []byte{byte(0x12), byte(0x13), byte(0x2A)}
			break
		case Stagenet:
			Tags = []byte{byte(0x18), byte(0x19), byte(0x24)}
			break
		case Testnet:
			Tags = []byte{byte(0x35), byte(0x36), byte(0x3F)}
			break
		}
		break
	}
}

func SetActiveTag(a AddressType) error {
	if int(a) < len(Tags) {
		Tag = Tags[a]
		return nil
	}
	return InvalidTagSelected
}

var (
	InvalidAddressLength = errors.New("invalid address length")
	CorruptAddress       = errors.New("address has invalid checksum")
	InvalidAddressTag    = errors.New("address has invalid prefix")
	InvalidAddress       = errors.New("address contains invalid keys")
	InvalidTagSelected   = errors.New("tag not available")
	InvalidPaymentID     = errors.New("invalid payment id for integrated address")
)

// DecodeAddress decodes an address from the standard textual representation.
func DecodeAddress(s string) (*Address, error) {
	pa := new(Address)
	err := pa.UnmarshalText([]byte(s))
	if err != nil {
		return nil, err
	}
	return pa, nil
}

// Address contains public keys for the spend and view aspects of a Monero account.
type Address struct {
	spend, view [32]byte
	paymentID   [8]byte
}

func (a *Address) MarshalBinary() (data []byte, err error) {
	// make this long enough to hold a full hash on the end
	data = make([]byte, 112)
	// copy tag
	n := 1
	data[0] = Tag

	//copy keys
	copy(data[n:], a.spend[:])
	if Tag == Tags[Integrated] {
		copy(data[n+32:n+64], a.view[:])
		copy(data[n+64:], a.paymentID[:])
	} else {
		copy(data[n+32:], a.view[:])
	}

	// checksum
	hash := crypto.NewHash()
	if Tag == Tags[Integrated] {
		hash.Write(data[:n+72])
		// hash straight to the slice
		hash.Sum(data[:n+72])
		return data[:n+76], nil
	} else {
		hash.Write(data[:n+64])
		// hash straight to the slice
		hash.Sum(data[:n+64])
		return data[:n+68], nil
	}

}

func (a *Address) UnmarshalBinary(data []byte) error {
	if len(data) < ChecksumSize {
		return InvalidAddressLength
	}

	// Verify checksum
	checksum := data[len(data)-ChecksumSize:]
	data = data[:len(data)-ChecksumSize]
	hash := crypto.NewHash()
	hash.Write(data)
	digest := hash.Sum(nil)
	if !bytes.Equal(checksum, digest[:ChecksumSize]) {
		return CorruptAddress
	}

	// check address prefix
	if data[0] != Tag {
		return InvalidAddressTag
	}

	data = data[1:]

	if len(data) == 64 {
		copy(a.spend[:], data[0:32])
		copy(a.view[:], data[32:64])
	} else if len(data) == 72 {
		copy(a.spend[:], data[0:32])
		copy(a.view[:], data[32:64])
		copy(a.paymentID[:], data[64:72])
	} else {
		return InvalidAddressLength
	}

	// don't check the keys yet
	return nil
}

func (a *Address) String() string {
	text, _ := a.MarshalText()
	return string(text)
}

func (a *Address) MarshalText() (text []byte, err error) {
	data, _ := a.MarshalBinary()
	text = make([]byte, base58.EncodedLen(len(data)))
	base58.Encode(text, data)
	return text, nil
}

func (a *Address) UnmarshalText(text []byte) error {
	// Decode from base58
	b := make([]byte, base58.DecodedLen(len(text)))
	_, err := base58.Decode(b, text)
	if err != nil {
		return err
	}
	return a.UnmarshalBinary(b)
}
