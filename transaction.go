package moneroCryptoNoteUtils

// Transaction Inputs
type TransactionInGenesis struct {
	Height uint64
	Used   bool
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

type TransactionOut struct {
	Amount     uint64
	Script     TransactionOutToScript
	ScriptHash TransactionOutToScriptHash
	Key        TransactionOutToKey
}

type TransactionPrefix struct {
	Version         uint64
	UnlockTime      uint64
	TransactionsIn  []TransactionIn
	TransactionsOut []TransactionOut
	Extra           []uint8
}

type Transaction struct {
	TransactionPrefix
	Signatures [][32]byte
}
