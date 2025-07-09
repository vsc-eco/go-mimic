package encoder_test

import (
	"mimic/lib/encoder"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	field1 string
	Field1 string

	field2 int
	Field2 int

	Field3 any
}

func TestJsonArrayMarshaler(t *testing.T) {
	raw := []byte(`["Field1 exported", 12, {}]`)
	buf := testStruct{}
	err := encoder.JsonArrayDeserialize(&buf, raw)
	assert.Nil(t, err)
	assert.Empty(t, buf.field1)
	assert.Equal(t, "Field1 exported", buf.Field1)
	assert.Empty(t, buf.field2)
	assert.Equal(t, 12, buf.Field2)
	assert.Empty(t, buf.Field3)
}
