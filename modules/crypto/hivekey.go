package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"slices"

	"github.com/btcsuite/btcd/btcec/v2"
)

type keyRole = string

const (
	postingKeyRole = keyRole("posting")
	activeKeyRole  = keyRole("active")
	ownerKeyRole   = keyRole("owner")
	memoKeyRole    = keyRole("memo")
)

type HiveKeySet struct {
	ownerKey   hiveKey
	activeKey  hiveKey
	postingKey hiveKey
	memoKey    string
}

func MakeHiveKeySet(account, password string) HiveKeySet {
	key := HiveKeySet{}

	key.ownerKey = makeHiveKey(nil, account, password, ownerKeyRole)

	key.activeKey = makeHiveKey(
		key.ownerKey.privKey.Serialize(),
		account,
		password,
		activeKeyRole,
	)

	key.postingKey = makeHiveKey(
		key.ownerKey.privKey.Serialize(),
		account,
		password,
		postingKeyRole,
	)

	memoKeyParts := sha256.Sum256(slices.Concat(
		[]byte(account),
		[]byte(password),
		[]byte(memoKeyRole),
	))

	key.memoKey = hex.EncodeToString(memoKeyParts[:])

	return key
}

func (h *HiveKeySet) OwnerKey() *hiveKey   { return &h.ownerKey }
func (h *HiveKeySet) ActiveKey() *hiveKey  { return &h.activeKey }
func (h *HiveKeySet) PostingKey() *hiveKey { return &h.postingKey }
func (h *HiveKeySet) MemoKey() string      { return h.memoKey }

type hiveKey struct {
	pubKey  *btcec.PublicKey
	privKey *btcec.PrivateKey
}

func (h *hiveKey) PublicKeyHex() string {
	return hex.EncodeToString(h.pubKey.SerializeCompressed())
}

func (h *hiveKey) Sign(message []byte) ([]byte, error) {
	msgHash := sha256.Sum256(message)
	return ecdsa.SignASN1(rand.Reader, h.privKey.ToECDSA(), msgHash[:])
}

func (h *hiveKey) Verify(message, signature []byte) bool {
	msgHash := sha256.Sum256(message)
	return ecdsa.VerifyASN1(h.pubKey.ToECDSA(), msgHash[:], signature)
}

// Hive's implementation for key generation
// https://gitlab.syncad.com/hive/devportal/-/blob/master/tutorials/python/34_password_key_change/index.py
// https://github.com/holgern/beem/blob/2026833a836007e45f16395a9ca3b31d02e98f87/beemgraphenebase/account.py#L33
func makeHiveKey(
	keyPart []byte,
	account, password string,
	role keyRole,
) hiveKey {
	buf := sha256.Sum256(slices.Concat(
		keyPart,
		[]byte(account),
		[]byte(password),
		[]byte(role),
	))

	key := hiveKey{}
	key.privKey, key.pubKey = btcec.PrivKeyFromBytes(buf[:])

	return key
}
