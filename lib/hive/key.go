package hive

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"slices"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/decred/dcrd/dcrec/secp256k1/v2"
	"github.com/vsc-eco/hivego"
)

type (
	KeyRole = string
)

type SigningAuthorities struct {
	Account string
	KeyType KeyRole
}

type HiveKey struct {
	*hivego.KeyPair
}

const (
	PostingKeyRole = KeyRole("posting")
	ActiveKeyRole  = KeyRole("active")
	OwnerKeyRole   = KeyRole("owner")
	MemoKeyRole    = KeyRole("memo")

	signatureLen     = 65
	signatureCompact = true

	hiveNetworkID = 0x80
)

const version = byte(0x00)

type HiveKeySet struct {
	ownerKey   HiveKey
	activeKey  HiveKey
	postingKey HiveKey
	memoKey    HiveKey
}

func MakeHiveKeySet(account, password string) HiveKeySet {
	key := HiveKeySet{}

	// make owner key
	key.ownerKey = makeHiveKey(nil, OwnerKeyRole, account, password)

	// make active key
	key.activeKey = makeHiveKey(
		key.ownerKey.PrivateKey.Serialize(),
		ActiveKeyRole,
		account,
		password,
	)

	// make posting key
	key.postingKey = makeHiveKey(
		key.ownerKey.PrivateKey.Serialize(),
		PostingKeyRole,
		account,
		password,
	)

	// make memo key
	key.memoKey = makeHiveKey(
		nil,
		MemoKeyRole,
		account,
		password,
	)

	return key
}

func (h *HiveKeySet) OwnerKey() *HiveKey   { return &h.ownerKey }
func (h *HiveKeySet) ActiveKey() *HiveKey  { return &h.activeKey }
func (h *HiveKeySet) PostingKey() *HiveKey { return &h.postingKey }
func (h *HiveKeySet) MemoKey() *HiveKey    { return &h.memoKey }

func (h *HiveKey) PrivateKeyWif() string {
	privKeyRaw := h.PrivateKey.Serialize()

	buf := make([]byte, len(privKeyRaw)+1)
	buf[0] = hiveNetworkID

	copy(buf[1:], privKeyRaw)
	return "5" + base58.CheckEncode(buf, version)
}

func NewHiveKeyFromPrivateWif(encodedKey string) (*HiveKey, error) {
	prefix, encodedKey := encodedKey[:1], encodedKey[1:]
	if prefix != "5" {
		return nil, errors.New("invalid private key WIF prefix")
	}

	key, _, err := base58.CheckDecode(encodedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %v", err)
	}

	networkID, key := key[0], key[1:]
	if networkID != hiveNetworkID {
		return nil, errors.New("invalid network id")
	}

	return &HiveKey{hivego.KeyPairFromBytes(key)}, nil
}

func (h *HiveKey) Sign(message []byte) ([]byte, error) {
	if h.PrivateKey == nil {
		return nil, errors.New("nil private key")
	}

	digest := sha256.Sum256(message)
	return secp256k1.SignCompact(h.PrivateKey, digest[:], signatureCompact)
}

func Verify(pubKeyWif string, message, signature []byte) (bool, error) {
	pubKey, err := hivego.DecodePublicKey(pubKeyWif)
	if err != nil {
		return false, fmt.Errorf("failed to decode public key: %v", err)
	}

	digest := sha256.Sum256(message)

	recoveredPubKey, compacted, err := secp256k1.RecoverCompact(
		signature,
		digest[:],
	)
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
	keyPart []byte,
	role KeyRole,
	username, password string,
) HiveKey {
	buf := slices.Concat(
		keyPart,
		[]byte(username),
		[]byte(role),
		[]byte(password),
	)
	digest := sha256.Sum256(buf)
	return HiveKey{hivego.KeyPairFromBytes(digest[:])}
}
