package argparser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testNestedArgs struct {
	Nested string
}

type testArgs struct {
	Value      string
	Number     int
	unexported int
	Bool       bool
	Sub        testNestedArgs
	Slice      []string
}

func TestSimpleParsing(t *testing.T) {
	var obj testArgs
	err := LoadArgs(&obj, strings.Split("--value=abc --number=123 --bool=true --sub-nested=qef", " ")...)
	assert.NoError(t, err)
	assert.Equal(t, obj.Value, "abc")
	assert.Equal(t, obj.Number, 123)
	assert.Equal(t, obj.Bool, true)
	assert.Equal(t, obj.Sub.Nested, "qef")
	assert.Equal(t, obj.unexported, 0)
}

func TestSettingUnexportedField(t *testing.T) {
	var obj testArgs
	err := LoadArgs(&obj, "--unexported=5")
	assert.Error(t, err)
}

func TestBadArgFormat(t *testing.T) {
	var obj testArgs
	err := LoadArgs(&obj, "--unexported 5")
	assert.Error(t, err)
}

func TestParseSlice(t *testing.T) {
	var obj testArgs
	err := LoadArgs(&obj, strings.Split("--slice=a --slice=bc", " ")...)
	assert.NoError(t, err)

	assert.Equal(t, []string{"a", "bc"}, obj.Slice)
}
