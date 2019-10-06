package crypto

import (
	cnutil "github.com/snipa22/monerocnutils"
	"github.com/snipa22/monerocnutils/crypto"
)

func TreeHash(b cnutil.Block) [32]byte {
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
