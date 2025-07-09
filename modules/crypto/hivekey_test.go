package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"io"
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
	account := "hive-io-account"
	password := "hive-io-password"

	key1 := makeHiveKey(nil, account, password, ownerKeyRole)
	key2 := makeHiveKey(nil, account, password, ownerKeyRole)

	t.Run("generates key pairs deterministically.", func(t *testing.T) {
		assert.Equal(
			t,
			key1.pubKey.SerializeCompressed(),
			key2.pubKey.SerializeCompressed(),
			"public keys deterministically generated.",
		)

		assert.Equal(
			t,
			key1.privKey.Serialize(),
			key2.privKey.Serialize(),
			"private keys deterministically generated.",
		)
	})

	t.Run("parses the private key correctly.", func(t *testing.T) {
		keyBytes := key1.privKey.Serialize()
		privKey := sha256.Sum256(
			slices.Concat(
				nil,
				[]byte(account),
				[]byte(password),
				[]byte(ownerKeyRole),
			),
		)
		assert.Equal(t, privKey[:], keyBytes)
	})

	t.Run("accepts valid signature", func(t *testing.T) {
		message := makeMessage(1024)
		sig, err := key1.Sign(message)
		assert.Nil(t, err)
		assert.True(t, key1.Verify(message, sig))
	})

	t.Run("rejects invalid signature", func(t *testing.T) {
		msg1 := makeMessage(1024)
		sig1, err := key1.Sign(msg1)
		assert.Nil(t, err)

		msg2 := makeMessage(1024)
		sig2, err := key2.Sign(msg2)
		assert.Nil(t, err)

		assert.False(t, key1.Verify(msg1, sig2))
		assert.False(t, key1.Verify(msg2, sig1))
		assert.False(t, key2.Verify(msg1, sig2))
		assert.False(t, key2.Verify(msg2, sig1))

		sig2[0] ^= 1 // making sig invalid
		assert.False(t, key2.Verify(msg2, sig2))
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
