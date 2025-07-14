package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"

	"github.com/decred/dcrd/dcrec/secp256k1/v2"
	"github.com/vsc-eco/hivego"
)

type keyRole = string

type HiveKey struct {
	*hivego.KeyPair
}

const (
	postingKeyRole = keyRole("posting")
	activeKeyRole  = keyRole("active")
	ownerKeyRole   = keyRole("owner")
	memoKeyRole    = keyRole("memo")

	signatureCompactLen = 65
)

type HiveKeySet struct {
	ownerKey   HiveKey
	activeKey  HiveKey
	postingKey HiveKey
	memoKey    string
}

func MakeHiveKeySet(account, password string) HiveKeySet {
	key := HiveKeySet{}

	var (
		accountBytes  = []byte(account)
		passwordBytes = []byte(password)
	)

	key.ownerKey = makeHiveKey(ownerKeyRole, accountBytes, passwordBytes)

	key.activeKey = makeHiveKey(
		activeKeyRole,
		key.ownerKey.PrivateKey.Serialize(),
		accountBytes,
		passwordBytes,
	)

	key.postingKey = makeHiveKey(
		postingKeyRole,
		key.ownerKey.PrivateKey.Serialize(),
		accountBytes,
		passwordBytes,
	)

	memoKeyParts := sha256.Sum256(slices.Concat(
		[]byte(memoKeyRole),
		accountBytes,
		passwordBytes,
	))

	key.memoKey = hex.EncodeToString(memoKeyParts[:])

	return key
}

func (h *HiveKeySet) OwnerKey() *HiveKey   { return &h.ownerKey }
func (h *HiveKeySet) ActiveKey() *HiveKey  { return &h.activeKey }
func (h *HiveKeySet) PostingKey() *HiveKey { return &h.postingKey }
func (h *HiveKeySet) MemoKey() string      { return h.memoKey }

func (h *HiveKey) Sign(message []byte) ([]byte, error) {
	if h.PrivateKey == nil {
		return nil, errors.New("nil private key.")
	}

	digest := sha256.Sum256(message)
	return secp256k1.SignCompact(h.PrivateKey, digest[:], true)
}

func Verify(pubKeyWif string, message, signature []byte) (bool, error) {
	pubKey, err := hivego.DecodePublicKey(pubKeyWif)
	if err != nil {
		return false, fmt.Errorf("failed to decode public key: %v", err)
	}

	digest := sha256.Sum256(message)

	recoveredPubKey, compacted, err := secp256k1.RecoverCompact(signature, digest[:])
	if err != nil {
		return false, fmt.Errorf("failed to decode signature key: %v", err)
	}

	if !compacted {
		return false, errors.New("invalid signature format")
	}

	return recoveredPubKey.IsEqual(pubKey), nil
}

// Hive's implementation for key generation
// https://gitlab.syncad.com/hive/devportal/-/blob/master/tutorials/python/34_password_key_change/index.py
// https://github.com/holgern/beem/blob/2026833a836007e45f16395a9ca3b31d02e98f87/beemgraphenebase/account.py#L33
func makeHiveKey(
	role keyRole,
	keyParts ...[]byte,
) HiveKey {
	buf := slices.Concat(
		keyParts...,
	)

	digest := sha256.Sum256(append(buf, []byte(role)...))

	return HiveKey{hivego.KeyPairFromBytes(digest[:])}
}
