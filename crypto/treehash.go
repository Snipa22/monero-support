package crypto

func TreeHash(hs [][32]byte) [32]byte {
	hash := [32]byte{}
	if len(hs) == 1 {
		return hs[0]
	} else if len(hs) == 2 {
		h := NewHash() // Hash Object
		h.Write(hs[0][:])
		h.Write(hs[1][:])
		var tempHash []byte
		h.Sum(tempHash)
		copy(hash[:], tempHash[:])
	} else {
		// Borrowed from the tree-hash.c inmplenm
		var i, j uint64
		var cnt uint64 = uint64(len(hs) - 1)
		i = j
		j = i
		i = cnt
	}
	return hash
}
