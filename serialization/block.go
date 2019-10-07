package serialization

import (
	"encoding/binary"
)

type BlockHeader struct {
	MajorVersion uint8
	MinorVersion uint8
	Timestamp    uint64
	PreviousID   [32]byte
	Nonce        uint32
}

func (bh BlockHeader) Serialize() []byte {
	var s []byte

	s = WriteUint(s, uint64(bh.MajorVersion))
	s = WriteUint(s, uint64(bh.MinorVersion))
	s = WriteUint(s, bh.Timestamp)

	// Previous Block ID
	var tempBlob []byte = make([]byte, 32)
	copy(tempBlob[:], bh.PreviousID[:])
	s = append(s, tempBlob...)

	// Nonce
	binary.BigEndian.PutUint32(tempBlob, bh.Nonce)
	s = append(s, tempBlob[0:4]...)

	return s
}

type Block struct {
	BlockHeader
	MinerTxn  Transaction
	TxnHashes [][32]byte
}