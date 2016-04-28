package human_size

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstBytes(t *testing.T) {
	assert.Equal(t, KB, 1024)
	assert.Equal(t, MB, KB*1024)
	assert.Equal(t, GB, MB*1024)
	assert.Equal(t, TB, GB*1024)
}

func TestParse(t *testing.T) {
	var (
		value int
		ok    bool
	)

	value, ok = Parse("")
	assert.False(t, ok)

	value, ok = Parse("90")
	assert.True(t, ok)
	assert.Equal(t, value, 90)

	value, ok = Parse("90 B")
	assert.True(t, ok)
	assert.Equal(t, value, 90)

	value, ok = Parse("90 KiB")
	assert.True(t, ok)
	assert.Equal(t, value, 90*KB)

	value, ok = Parse("90 kB")
	assert.True(t, ok)
	assert.Equal(t, value, 90*KB)

	value, ok = Parse("90 kiB")
	assert.True(t, ok)
	assert.Equal(t, value, 90*KB)

	value, ok = Parse("90 MB")
	assert.True(t, ok)
	assert.Equal(t, value, 90*MB)

	value, ok = Parse("90 GiB")
	assert.True(t, ok)
	assert.Equal(t, value, 90*GB)

	value, ok = Parse("90 TB")
	assert.True(t, ok)
	assert.Equal(t, value, 90*TB)

	value, ok = Parse("90x")
	assert.False(t, ok)

	value, ok = Parse("90 Ki")
	assert.False(t, ok)

	value, ok = Parse("90 k")
	assert.False(t, ok)

	value, ok = Parse("90 GB")
	assert.True(t, ok)
	assert.Equal(t, value, 90*GB)
}

func TestFormat(t *testing.T) {
	var out string

	out = Format(1023)
	assert.Equal(t, out, "1023 B")

	out = Format(KB)
	assert.Equal(t, out, "1.00 kB")

	out = Format(1536)
	assert.Equal(t, out, "1.50 kB")

	out = Format(1537)
	assert.Equal(t, out, "1.50 kB")

	out = Format(MB)
	assert.Equal(t, out, "1.00 MB")

	out = Format(GB)
	assert.Equal(t, out, "1.00 GB")

	out = Format(TB)
	assert.Equal(t, out, "1.00 TB")

	out = Format(10866 * TB)
	assert.Equal(t, out, "10866.00 TB")

	out = Format(478841856)
	assert.Equal(t, out, "456.66 MB")
}
