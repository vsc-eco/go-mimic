package producers

import (
	"mimic/lib/hive"
	"testing"
)

func TestProducerPubKeyQuery(t *testing.T) {
	keyBuf := make(map[string]map[hive.KeyRole]string)
	keyBuf["@testvscsodf"] = map[hive.KeyRole]string{
		hive.ActiveKeyRole: "",
	}
	t.Log("Finish this test")
}
