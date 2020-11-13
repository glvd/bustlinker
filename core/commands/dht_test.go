package commands

import (
	"testing"

	"github.com/ipfs/go-ipfs/namesys"

	blns "github.com/ipfs/go-ipns"
	"github.com/libp2p/go-libp2p-core/test"
)

func TestKeyTranslation(t *testing.T) {
	pid := test.RandPeerIDFatal(t)
	pkname := namesys.PkKeyForID(pid)
	blnsname := blns.RecordKey(pid)

	pkk, err := escapeDhtKey("/pk/" + pid.Pretty())
	if err != nil {
		t.Fatal(err)
	}

	blnsk, err := escapeDhtKey("/blns/" + pid.Pretty())
	if err != nil {
		t.Fatal(err)
	}

	if pkk != pkname {
		t.Fatal("keys didnt match!")
	}

	if blnsk != blnsname {
		t.Fatal("keys didnt match!")
	}
}
