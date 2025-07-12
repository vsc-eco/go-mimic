package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"slices"

	"github.com/btcsuite/btcd/btcec/v2"
)

type keyRole = string

const (
	postingKeyRole = keyRole("posting")
	activeKeyRole  = keyRole("active")
	ownerKeyRole   = keyRole("owner")
	memoKeyRole    = keyRole("memo")

	signatureCompactLen = 64
)

type HiveKeySet struct {
	ownerKey   HiveKey
	activeKey  HiveKey
	postingKey HiveKey
	memoKey    string
}

func MakeHiveKeySet(account, password string) HiveKeySet {
	key := HiveKeySet{}

	key.ownerKey = makeHiveKey(nil, account, password, ownerKeyRole)

	key.activeKey = makeHiveKey(
		key.ownerKey.PrivKey.Serialize(),
		account,
		password,
		activeKeyRole,
	)

	key.postingKey = makeHiveKey(
		key.ownerKey.PrivKey.Serialize(),
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

func (h *HiveKeySet) OwnerKey() *HiveKey   { return &h.ownerKey }
func (h *HiveKeySet) ActiveKey() *HiveKey  { return &h.activeKey }
func (h *HiveKeySet) PostingKey() *HiveKey { return &h.postingKey }
func (h *HiveKeySet) MemoKey() string      { return h.memoKey }

type HiveKey struct {
	PubKey  *btcec.PublicKey
	PrivKey *btcec.PrivateKey
}

func (h *HiveKey) PublicKeyHex() string {
	return hex.EncodeToString(h.PubKey.SerializeCompressed())
}

// the returned signature is in a compact format (64 bytes)
func (h *HiveKey) Sign(message []byte) ([]byte, error) {
	msgHash := sha256.Sum256(message)
	r, s, err := ecdsa.Sign(rand.Reader, h.PrivKey.ToECDSA(), msgHash[:])
	if err != nil {
		return nil, err
	}

	sig := make([]byte, signatureCompactLen)
	i := signatureCompactLen / 2

	copy(sig[:i], r.Bytes())
	copy(sig[i:], s.Bytes())

	return sig, nil
}

func (h *HiveKey) Verify(message, signature []byte) bool {
	if len(signature) != signatureCompactLen { // invalid signature
		return false
	}

	r, s := big.Int{}, big.Int{}
	r.SetBytes(stripZeroBytes(signature[:32]))
	s.SetBytes(stripZeroBytes(signature[32:]))

	hash := sha256.Sum256(message)

	return ecdsa.Verify(h.PubKey.ToECDSA(), hash[:], &r, &s)
}

func stripZeroBytes(buf []byte) []byte {
	i := len(buf) - 1
	for ; i >= 0; i-- {
		if buf[i] != 0 {
			break
		}
	}
	return buf[:i]
}

// Hive's implementation for key generation
// https://gitlab.syncad.com/hive/devportal/-/blob/master/tutorials/python/34_password_key_change/index.py
// https://github.com/holgern/beem/blob/2026833a836007e45f16395a9ca3b31d02e98f87/beemgraphenebase/account.py#L33
func makeHiveKey(
	keyPart []byte,
	account, password string,
	role keyRole,
) HiveKey {
	buf := sha256.Sum256(slices.Concat(
		keyPart,
		[]byte(account),
		[]byte(password),
		[]byte(role),
	))

	key := HiveKey{}
	key.PrivKey, key.PubKey = btcec.PrivKeyFromBytes(buf[:])

	return key
}
