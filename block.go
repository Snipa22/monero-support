package moneroCryptoNoteUtils

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"github.com/snipa22/monero-support/crypto"
)

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

type BlockHeader struct {
	MajorVersion uint8
	MinorVersion uint8
	Timestamp    uint64
	PreviousID   [32]byte
	Nonce        uint32
}

type Block struct {
	BlockHeader
	MinerTxn  Transaction
	TxnHashes [][32]byte
}

// Things to be implemented
// 1. convert_blob -> Converts a block template blob that's been modified into a hashable object
// 1.1 parse_and_validate_block_from_blob -> Created a Block struct from a blob of data provided by the Get Block Template RPC call, that is then modified to suit usages
// 1.2 get_block_hashing_blob -> Converts the blob into a block hashing blob

var (
	errorBadBlob     error = errors.New("bad blob provided, unable to continue")
	errorUintTooLong error = errors.New("bad blob provided, uint is unable to be decoded, unable to continue")
)

func readUint(b []byte) (uint64, []byte, error) {
	val, byteCount := binary.Uvarint(b)
	if byteCount == 0 {
		return 0, b, errorBadBlob
	}
	if byteCount <= 0 {
		return 0, b, errorUintTooLong
	}
	return val, b[byteCount:], nil
}

func ParseBlockFromTemplateBlob(blob string) (Block, error) {
	var b Block
	blobInBytes, err := hex.DecodeString(blob)
	if err != nil {
		return b, err
	}
	// Get the Major Version, uint8
	val, blobInBytes, err := readUint(blobInBytes)
	if err != nil {
		return b, err
	}
	b.MajorVersion = uint8(val)

	// Get the Minor Version, uint8
	val, blobInBytes, err = readUint(blobInBytes)
	if err != nil {
		return b, err
	}
	b.MinorVersion = uint8(val)

	// Get the Timestamp, uint64
	val, blobInBytes, err = readUint(blobInBytes)
	if err != nil {
		return b, err
	}
	b.Timestamp = val

	// Get the previous hash, which is an array of 32 bytes in uint8 form, stored as 32 bytes in the array
	bytesCopied := copy(b.PreviousID[:], blobInBytes[0:32])
	blobInBytes = blobInBytes[bytesCopied:]

	// Get the nonce, uint32, but is stored as a block of 4 bytes...  Jackassery.
	b.Nonce = binary.BigEndian.Uint32(blobInBytes[0:4])
	blobInBytes = blobInBytes[4:]

	// Start Transaction Processing (Miner Transaction)
	var t Transaction

	// Get Version, uint64
	val, blobInBytes, err = readUint(blobInBytes)
	if err != nil {
		return b, err
	}
	t.Version = val

	// Get UnlockTime, uint64 -- Could be a timestamp OR a block ID
	val, blobInBytes, err = readUint(blobInBytes)
	if err != nil {
		return b, err
	}
	t.UnlockTime = val

	// Start processing the t.vin fields
	// These are the Variant In fields.

	// Move forwards by 1 as the array is one object in length
	blobInBytes = blobInBytes[1:]

	// Move forwards by 1 as the key type for the array is going to be the TransactionInGenesis type
	blobInBytes = blobInBytes[1:]

	// Load the resulting blob data into the correct portion of Transaction/TransactionsIn
	var tig TransactionInGenesis
	var ti TransactionIn

	// Get the genesis height
	val, blobInBytes, err = readUint(blobInBytes)
	if err != nil {
		return b, err
	}
	tig.Height = val
	tig.Used = true
	ti.Genesis = tig

	t.TransactionsIn = append(t.TransactionsIn, ti)

	// Move forwards by 1 as the array is one object in length
	blobInBytes = blobInBytes[1:]

	// Load the resulting blob data into the correct portion of Transaction/TransactionsOut
	var to TransactionOut

	val, blobInBytes, err = readUint(blobInBytes)
	if err != nil {
		return b, err
	}
	to.Amount = val

	// Outbound type is to key, or 0x02.  Increment bytes by one
	blobInBytes = blobInBytes[1:]

	// Write a new TransactionOutToKey
	var totk TransactionOutToKey
	bytesCopied = copy(totk.PublicKey[:], blobInBytes[0:32])
	blobInBytes = blobInBytes[bytesCopied:]
	totk.Used = true
	to.Key = totk

	t.TransactionsOut = append(t.TransactionsOut, to)

	// Get the number of bytes to read into "extra"
	val, blobInBytes, err = readUint(blobInBytes)
	if err != nil {
		return b, err
	}

	// With val set to the # of bytes to read, slice and go
	t.Extra = blobInBytes[0:val]
	blobInBytes = blobInBytes[val:]

	// RingCT Type is 0.  Advance one byte
	blobInBytes = blobInBytes[1:]

	// Miner Transaction Complete!  Load to the main storage.  AYE AYE CAPTAIN!
	b.MinerTxn = t

	// Get the number of hashes in the tx_hashes field
	val, blobInBytes, err = readUint(blobInBytes)
	if err != nil {
		return b, err
	}

	// Attempt to get <val> hashes and append to the main store
	for ; val > 0; val-- {
		var iSlice [32]byte
		bytesCopied = copy(iSlice[:], blobInBytes[0:32])
		b.TxnHashes = append(b.TxnHashes, iSlice)
		blobInBytes = blobInBytes[bytesCopied:]
	}

	return b, nil
}

func GetBlockHashingBlob(b Block) ([]byte, error) {
	var blob []byte
	var tempBlob []byte = make([]byte, 32)

	// Time to serialize the headers.  Lets go.

	// Major Version
	bytesWritten := binary.PutUvarint(tempBlob, uint64(b.MajorVersion))
	blob = append(blob, tempBlob[0:bytesWritten]...)

	// Minor Version
	bytesWritten = binary.PutUvarint(tempBlob, uint64(b.MinorVersion))
	blob = append(blob, tempBlob[0:bytesWritten]...)

	// Timestamp
	bytesWritten = binary.PutUvarint(tempBlob, b.Timestamp)
	blob = append(blob, tempBlob[0:bytesWritten]...)

	// Previous Hash
	copy(tempBlob[:], b.PreviousID[:])
	blob = append(blob, tempBlob...)

	// Nonce
	binary.BigEndian.PutUint32(tempBlob, b.Nonce)
	blob = append(blob, tempBlob[0:4]...)

	return blob, nil
}

func getTransactionPrefixSerialized(tp TransactionPrefix) []byte {
	// Serialize a prefix
	// varint - Version
	// varint - unlocked_time
	// Vector store of the vin
	// Vector store of the vout
	// Slap on that extra data.  Mmmmm.  Extra.  Data.  Nomnom.
}

func getTransactionPrefixHash(t Transaction) [32]byte {
	// Given a Transaction t, extract the TransactionPrefix TP and serialize it.
	// Given the resulting serialized data, cn_fast_hash (keccak-256) it.
}

func getTransactionHash(t Transaction) [32]byte {
	// Original source: cryptonote_format_utils.cpp:617-ish
	// With hashes be three, may thee get the result thoust desire.
	var hs [3][32]byte

	// Thou must take tine prefix, and hash it!
	// Original : get_transaction_prefix_hash(t (Transaction), hashes[0] (crypto::hash))
}

func getBlockMerkleTreeHash(b Block) [32]byte {
	// Get the transaction hash?  Wtfh.
	// Original: get_transaction_hash(b.miner_tx (Transaction), h (crypto::hash), bl_sz (size_t (uint64 for us!)));

	// Shift the hashes into a new slice, first one is the txn hash, then add all other hashes to the end.

	// Get the tree hash, this is the return.  Need to abstract some of this to a support library...
}

func treeHash(b Block) [32]byte {
	hash := [32]byte{}
	if len(b.TxnHashes) == 1 {
		return b.TxnHashes[0]
	} else if len(b.TxnHashes) == 2 {
		h := crypto.NewHash() // Hash Object
		h.Write(b.TxnHashes[0][:])
		h.Write(b.TxnHashes[1][:])
		var tempHash []byte
		h.Sum(tempHash)
		copy(hash[:], tempHash[:])
	} else {
		// Borrowed from the tree-hash.c inmplenm
		var i, j uint64
		var cnt uint64 = uint64(len(b.TxnHashes) - 1)
	}
	return hash
}
