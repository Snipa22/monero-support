package monerocnutils

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/snipa22/monerocnutils/crypto"
	"github.com/snipa22/monerocnutils/serialization"
)

// Things to be implemented
// 1. convert_blob -> Converts a block template blob that's been modified into a hashable object
// 1.1 parse_and_validate_block_from_blob -> Created a Block struct from a blob of data provided by the Get Block Template RPC call, that is then modified to suit usages
// 1.2 get_block_hashing_blob -> Converts the blob into a block hashing blob

func ParseBlockFromTemplateBlob(blob string) (serialization.Block, error) {
	var b serialization.Block
	blobInBytes, err := hex.DecodeString(blob)
	if err != nil {
		return b, err
	}
	// Get the Major Version, uint8
	val, blobInBytes, err := serialization.ReadUint(blobInBytes)
	if err != nil {
		return b, err
	}
	b.MajorVersion = uint8(val)

	// Get the Minor Version, uint8
	val, blobInBytes, err = serialization.ReadUint(blobInBytes)
	if err != nil {
		return b, err
	}
	b.MinorVersion = uint8(val)

	// Get the Timestamp, uint64
	val, blobInBytes, err = serialization.ReadUint(blobInBytes)
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
	var t serialization.Transaction

	// Get Version, uint64
	val, blobInBytes, err = serialization.ReadUint(blobInBytes)
	if err != nil {
		return b, err
	}
	t.Version = val

	// Get UnlockTime, uint64 -- Could be a timestamp OR a block ID
	val, blobInBytes, err = serialization.ReadUint(blobInBytes)
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
	var tig serialization.TransactionInGenesis
	var ti serialization.TransactionIn

	// Get the genesis height
	val, blobInBytes, err = serialization.ReadUint(blobInBytes)
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
	var to serialization.TransactionOut

	val, blobInBytes, err = serialization.ReadUint(blobInBytes)
	if err != nil {
		return b, err
	}
	to.Amount = val

	// Outbound type is to key, or 0x02.  Increment bytes by one
	blobInBytes = blobInBytes[1:]

	// Write a new TransactionOutToKey
	var totk serialization.TransactionOutToKey
	bytesCopied = copy(totk.PublicKey[:], blobInBytes[0:32])
	blobInBytes = blobInBytes[bytesCopied:]
	totk.Used = true
	to.Key = totk

	t.TransactionsOut = append(t.TransactionsOut, to)

	// Get the number of bytes to read into "extra"
	val, blobInBytes, err = serialization.ReadUint(blobInBytes)
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
	val, blobInBytes, err = serialization.ReadUint(blobInBytes)
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

func GetBlockHashingBlob(b serialization.Block) ([]byte, error) {
	// Source: cryptonote_format_utils.cpp
	// Original: get_block_hashing_blob(const block& b, blobdata& blob) - Line: 678-ish
	var sbh []byte = b.BlockHeader.Serialize()
}

// Supporting hash systems

func getTransactionPrefixHash(t serialization.Transaction) [32]byte {
	// Given a Transaction t, extract the TransactionPrefix TP and serialize it.
	// Given the resulting serialized data, cn_fast_hash (keccak-256) it.
	hash := crypto.KeccakOneShot(t.TransactionPrefix.Serialize())
	return hash
}

func getTransactionHash(t serialization.Transaction) [32]byte {
	// Original source: cryptonote_format_utils.cpp:617-ish
	// With hashes be three, may thee get the result thoust desire.
	var hs [3][32]byte

	// Thou must take tine prefix, and hash it!
	// Original : get_transaction_prefix_hash(t (Transaction), hashes[0] (crypto::hash))
	hs[0] = getTransactionPrefixHash(t)

	// Base RingCT Transaction Hash Data - byte 0 for the main txn, due to RingCTType being null (0x0)
	// So we're gonna short-cut this...
	hs[1] = crypto.KeccakOneShot([]byte{0})

	// Null hashes are value 0, with no additional data
	hs[2] = [32]byte{}

	var ah []byte
	ah = append(ah, hs[0][:]...)
	ah = append(ah, hs[1][:]...)
	ah = append(ah, hs[2][:]...)
	var h [32]byte = crypto.KeccakOneShot(ah)
	return h
}

func getBlockMerkleTreeHash(b serialization.Block) [32]byte {
	// Original: get_tx_tree_hash(const block& b) - cryptonote_format_utils.cpp:875-ish

	// Get the transaction hash.
	var txList [][32]byte
	txList = append(txList, getTransactionHash(b.MinerTxn))

	// Shift the hashes into a new slice, first one is the txn hash, then add all other hashes to the end.
	txList = append(txList, b.TxnHashes...)

	// Get the tree hash, this is the return.  Need to abstract some of this to a support library...
	hash := [32]byte{}
	return hash
}
