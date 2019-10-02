package monero

import "testing"

const normal = "4AMGENEQLdPGSqhGSgTdzH8dWxWoVwiTfgf2oTjPjxsgbUJS7kkK7euAhm94snzXVhHtZLwAXLiZQ6nDaWmqWHeSTafpXVw"
const integrated = "4L3wFB3twtuGSqhGSgTdzH8dWxWoVwiTfgf2oTjPjxsgbUJS7kkK7euAhm94snzXVhHtZLwAXLiZQ6nDaWmqWHeSgTmh7tVYUx65eb5iE7"
const subaddress = "8AiHrLaxEACUgytKhaVVEN4JELJ8m9uc5DbXWxTvavKqFYPbMGmPE75N7RFUVHhgxABW7y7tqih6r8CVUWzcc42DBTMihBd"

func TestDecodeEncodeAddressNormal(t *testing.T) {
	SetValidTags(Monero, Mainnet)
	err := SetActiveTag(Normal)
	if err != nil {
		t.Fatal("Error setting tags,", err)
	}
	addr, err := DecodeAddress(normal)
	if err != nil {
		t.Fatal("Error decoding address,", err)
	}
	if normal != addr.String() {
		t.Errorf("Decoding and encoding failed,\nwanted %s,\ngot    %s", normal, addr)
	}
}

func TestDecodeEncodeAddressIntegrated(t *testing.T) {
	SetValidTags(Monero, Mainnet)
	err := SetActiveTag(Integrated)
	if err != nil {
		t.Fatal("Error setting tags,", err)
	}
	addr, err := DecodeAddress(integrated)
	if err != nil {
		t.Fatal("Error decoding address,", err)
	}
	if integrated != addr.String() {
		t.Errorf("Decoding and encoding failed,\nwanted %s,\ngot    %s", integrated, addr)
	}
}

func TestDecodeEncodeAddressSubaddress(t *testing.T) {
	SetValidTags(Monero, Mainnet)
	err := SetActiveTag(Subaddress)
	if err != nil {
		t.Fatal("Error setting tags,", err)
	}
	addr, err := DecodeAddress(subaddress)
	if err != nil {
		t.Fatal("Error decoding address,", err)
	}
	if subaddress != addr.String() {
		t.Errorf("Decoding and encoding failed,\nwanted %s,\ngot    %s", subaddress, addr)
	}
}
