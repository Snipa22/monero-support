package monerocnutils

import "monerocnutils/serialization"

// Transaction Inputs
type TransactionInGenesis struct {
	Height uint64
	Used   bool
}

func (tig TransactionInGenesis) serialize() []byte {
	var s []byte
	s = serialization.WriteUint(s, tig.Height)
	return s
}

type TransactionInToScript struct {
	PreviousHash   [32]byte
	PreviousOutput uint64
	SignatureSet   []uint8
	Used           bool
}

type TransactionInToScriptHash struct {
	PreviousHash   [32]byte
	PreviousOutput uint64
	Script         TransactionOutToScript
	SignatureSet   []uint8
	Used           bool
}

type TransactionInToKey struct {
	Amount     uint64    // Amount of coin transferred
	KeyOffsets []uint64  // Key offsets are numeric ID's
	KeyImage   [32]uint8 // Key Image is a hex
	Used       bool
}

type TransactionIn struct {
	Genesis    TransactionInGenesis
	Script     TransactionInToScript
	ScriptHash TransactionInToScriptHash
	Key        TransactionInToKey
}

func (ti TransactionIn) serialize() []byte {
	var s []byte
	if ti.Genesis.Used {
		s = append(s, ti.Genesis.serialize()...)
	}
	return s
}

// Transaction Outputs
type TransactionOutToScript struct {
	PublicKey [32]byte
	Script    []uint8
	Used      bool
}

type TransactionOutToScriptHash struct {
	Hash [32]byte
	Used bool
}

type TransactionOutToKey struct {
	PublicKey [32]byte
	Used      bool
}

func (totk TransactionOutToKey) serialize() []byte {
	var s []byte
	s = append(s, totk.PublicKey[:]...)
	return s
}

type TransactionOut struct {
	Amount     uint64
	Script     TransactionOutToScript
	ScriptHash TransactionOutToScriptHash
	Key        TransactionOutToKey
}

func (to TransactionOut) serialize() []byte {
	var s []byte
	s = serialization.WriteUint(s, to.Amount)
	if to.Key.Used {
		s = append(s, to.Key.serialize()...)
	}
	return s
}

type TransactionPrefix struct {
	Version         uint64
	UnlockTime      uint64
	TransactionsIn  []TransactionIn
	TransactionsOut []TransactionOut
	Extra           []uint8
}

func (tp TransactionPrefix) serialize() []byte {
	// varint - Version
	// varint - unlocked_time
	// Vector store of the vin
	// Vector store of the vout
	// Slap on that extra data.  Mmmmm.  Extra.  Data.  Nomnom.
	var s []byte
	s = serialization.WriteUint(s, tp.Version)
	s = serialization.WriteUint(s, tp.UnlockTime)
	s = serialization.WriteUint(s, uint64(len(tp.TransactionsIn)))
	for _, e := range tp.TransactionsIn {
		s = append(s, e.serialize()...)
	}
	s = serialization.WriteUint(s, uint64(len(tp.TransactionsOut)))
	for _, e := range tp.TransactionsOut {
		s = append(s, e.serialize()...)
	}
	return s
}

type Transaction struct {
	TransactionPrefix
	Signatures [][32]byte
}
