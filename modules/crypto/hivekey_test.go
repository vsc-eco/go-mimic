package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"io"
	mrand "math/rand"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeMessage(msgLen int) []byte {
	msg := make([]byte, msgLen)
	if _, err := io.ReadFull(rand.Reader, msg); err != nil {
		panic(err)
	}
	return msg
}

func TestHiveKey(t *testing.T) {
	account := []byte("hive-io-account")
	password := []byte("hive-io-password")

	key1 := makeHiveKey(ownerKeyRole, account, password)
	key2 := makeHiveKey(ownerKeyRole, account, password)

	t.Run("generates key pairs deterministically.", func(t *testing.T) {
		assert.True(
			t,
			key1.PublicKey.IsEqual(key2.PublicKey),
			"public keys deterministically generated.",
		)

		assert.Equal(
			t,
			key1.PrivateKey.Serialize(),
			key2.PrivateKey.Serialize(),
			"private keys deterministically generated.",
		)
	})

	t.Run("parses the private key correctly.", func(t *testing.T) {
		keyBytes := key1.PrivateKey.Serialize()
		privKey := sha256.Sum256(
			slices.Concat(
				[]byte(account),
				[]byte(password),
				[]byte(ownerKeyRole),
			),
		)
		assert.Equal(t, privKey[:], keyBytes)
	})

	t.Run("signs/verifies valid signatures for messages of random length.", func(t *testing.T) {
		for range 0xff {
			msg := makeMessage(mrand.Intn(0xffff))
			sig, err := key1.Sign(msg)
			assert.Nil(t, err)
			assert.Equal(t, signatureLen, len(sig))

			pubKeyWif := key1.GetPublicKeyString()
			sigOk, err := Verify(*pubKeyWif, msg, sig)
			assert.Nil(t, err)
			assert.True(t, sigOk)
		}
	})

	t.Run("rejects on invalid signature", func(t *testing.T) {
		msg := makeMessage(1024)
		sig, err := key1.Sign(msg)
		assert.Nil(t, err)

		sigCpy := make([]byte, len(sig))
		copy(sigCpy, sig)

		sigCpy[2] ^= 1

		pubKeyWif := key1.GetPublicKeyString()

		sigOk, err := Verify(*pubKeyWif, msg, sigCpy)
		assert.Nil(t, err)
		assert.False(t, sigOk)
	})

	t.Run("rejects on invalid message", func(t *testing.T) {
		msg := makeMessage(1024)
		sig, err := key1.Sign(msg)
		assert.Nil(t, err)

		pubKeyWif := key1.GetPublicKeyString()

		sigOk, err := Verify(*pubKeyWif, makeMessage(1024), sig)
		assert.Nil(t, err)
		assert.False(t, sigOk)
	})

	t.Run(
		"produces different signatures on different messages.",
		func(t *testing.T) {
			msg1 := makeMessage(1024)

			msg2 := make([]byte, len(msg1))
			copy(msg2, msg1)
			msg2[0] ^= 1 // flip a bit of msg2, make msg1 != msg2

			sig1, err := key1.Sign(msg1)
			assert.Nil(t, err)

			sig2, err := key1.Sign(msg2)
			assert.Nil(t, err)

			assert.NotEqual(t, sig1, sig2)
		},
	)
}
