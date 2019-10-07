package crypto

import (
	"encoding/hex"
	"testing"
)

func TestTreeHash(t *testing.T) {
	var testVals [][32]byte
	var tempBuf [32]byte

	hexBuf, _ := hex.DecodeString("676567f8b1b470207c20d8efbaacfa64b2753301b46139562111636f36304bb8")
	copy(tempBuf[:], hexBuf[0:32])
	testVals = append(testVals, tempBuf)
	result := TreeHash(testVals)

	if hex.EncodeToString(result[:]) != "676567f8b1b470207c20d8efbaacfa64b2753301b46139562111636f36304bb8" {
		t.Fatal("Unable to hash a single length object properly")
	}

	hexBuf, _ = hex.DecodeString("3124758667bc8e76e25403eee75a1044175d58fcd3b984e0745d0ab18f473984")
	copy(tempBuf[:], hexBuf[0:32])
	testVals[0] = tempBuf
	hexBuf, _ = hex.DecodeString("975ce54240407d80eedba2b395bcad5be99b5c920abc2423865e3066edd4847a")
	copy(tempBuf[:], hexBuf[0:32])
	testVals = append(testVals, tempBuf)
	result = TreeHash(testVals)

	if hex.EncodeToString(result[:]) != "5077570fed2363a14fa978218185b914059e23517faf366f08a87cf3c47fd58e" {
		t.Fatalf("Unable to hash a double length object properly, received %v", hex.EncodeToString(result[:]))
	}

	hexBuf, _ = hex.DecodeString("decc1e0aa505d7d5fbe8ed823d7f5da55307c4cc7008e306da82dbce492a0576")
	copy(tempBuf[:], hexBuf[0:32])
	testVals[0] = tempBuf
	hexBuf, _ = hex.DecodeString("dbcf0c26646d36b36a92408941f5f2539f7715bcb1e2b1309cedb86ae4211554")
	copy(tempBuf[:], hexBuf[0:32])
	testVals[1] = tempBuf
	hexBuf, _ = hex.DecodeString("f56f5e6b2fce16536e44c851d473d1f994793873996ba448dd59b3b4b922b183")
	copy(tempBuf[:], hexBuf[0:32])
	testVals = append(testVals, tempBuf)

	result = TreeHash(testVals)

	if hex.EncodeToString(result[:]) != "f8e26aaa7c36523cea4c5202f2df159c62bf70d10670c96aed516dbfd5cb5227" {
		t.Fatalf("Unable to hash a triple length object properly, received %v", hex.EncodeToString(result[:]))
	}

	hexBuf, _ = hex.DecodeString("53edbbf98d3fa50a85fd2d46c42502aafad3fea30bc25ba4f16ec8bf4a475c4d")
	copy(tempBuf[:], hexBuf[0:32])
	testVals[0] = tempBuf
	hexBuf, _ = hex.DecodeString("87da8ad3e5c90aae0b10a559a77a0985608eaa3cc3dd338239be52572c3bdf4b")
	copy(tempBuf[:], hexBuf[0:32])
	testVals[1] = tempBuf
	hexBuf, _ = hex.DecodeString("a403d27466991997b3cf4e8d238d002a1451ccc9c4790269d0f0085d9382d60f")
	copy(tempBuf[:], hexBuf[0:32])
	testVals[2] = tempBuf
	hexBuf, _ = hex.DecodeString("ef37717f59726e4cc8787d5d2d75238ba9adb9627a8f4aeeec8d80465ed3f5fb")
	copy(tempBuf[:], hexBuf[0:32])
	testVals = append(testVals, tempBuf)

	result = TreeHash(testVals)

	if hex.EncodeToString(result[:]) != "45f6e06fc0263e667caddd8fba84c9fb723a961a01a5b115f7cab7fe8f2c7e44" {
		t.Fatalf("Unable to hash a quad length object properly, received %v", hex.EncodeToString(result[:]))
	}
}
