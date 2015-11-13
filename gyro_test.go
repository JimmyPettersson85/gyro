package gyro

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateLogger(t *testing.T) {
	var logger *Logger
	var err error

	logger, err = New("test")
	assert.NotNil(t, logger)
	assert.NoError(t, err)

	logger, err = New("noexist")
	assert.Nil(t, logger)
	assert.Error(t, err)
}

func TestFilenames(t *testing.T) {
	logger, err := New("test")
	require.NotNil(t, logger)
	require.NoError(t, err)

	// We need a static time to work with for the tests
	logger.SetTimeFunction(func() time.Time {
		return time.Unix(0, 0).UTC()
	})

	var tests = []struct {
		fn       func(string)
		input    string
		expected string
	}{
		{func(s string) {}, "", "1970-01-01T00.log"}, //default values test
		{logger.SetSeparator, "_", "1970-01-01T00.log"},
		{logger.SetExtension, "", "1970-01-01T00"},
		{logger.SetExtension, "txt", "1970-01-01T00.txt"},
		{logger.SetPrefix, "pre", "pre_1970-01-01T00.txt"},
		{logger.SetSuffix, "suf", "pre_1970-01-01T00_suf.txt"},
		{logger.SetLayout, "2006010215", "pre_1970010100_suf.txt"},
	}

	for _, test := range tests {
		test.fn(test.input)
		assert.Equal(t, test.expected, logger.FileName())
	}

	// Sanity check of the logger state
	assert.Equal(t, "_", logger.separator)
	assert.Equal(t, "txt", logger.extension)
	assert.Equal(t, "pre", logger.prefix)
	assert.Equal(t, "suf", logger.suffix)
	assert.Equal(t, "2006010215", logger.layout)
	assert.Equal(t, "pre_%s_suf.txt", logger.format)
	assert.True(t, time.Unix(0, 0).UTC().Equal(logger.timeFn()))
}

func TestSetTime(t *testing.T) {
	logger, err := New("test")
	require.NotNil(t, logger)
	require.NoError(t, err)

	logger.SetTimeFunction(func() time.Time {
		return time.Unix(0, 0).UTC()
	})
	assert.Equal(t, "1970-01-01T00.log", logger.FileName())

	// Switch time zone to US/Alaska
	logger.SetTimeFunction(func() time.Time {
		loc, err := time.LoadLocation("US/Alaska")
		assert.NoError(t, err)
		return time.Unix(0, 0).UTC().In(loc)
	})
	assert.Equal(t, "1969-12-31T14.log", logger.FileName())

	// Add one hour
	logger.SetTimeFunction(func() time.Time {
		loc, err := time.LoadLocation("US/Alaska")
		assert.NoError(t, err)
		return time.Unix(0, 0).UTC().In(loc).Add(1 * time.Hour)
	})
	assert.Equal(t, "1969-12-31T15.log", logger.FileName())
}

func TestWrites(t *testing.T) {
	logger, err := New("test")
	require.NotNil(t, logger)
	require.NoError(t, err)
	logger.SetTimeFunction(func() time.Time {
		return time.Unix(0, 0)
	})

	fn := filepath.Join("test", logger.FileName())

	// The test file should not be present at the start of the test
	_, err = os.Stat(fn)
	require.Error(t, err)

	// Remove the test file after the tests
	defer os.Remove(fn)

	line := []byte("test line 1\n")
	n, err := logger.Write(line)
	assert.NoError(t, err)
	assert.Equal(t, len(line), n)

	assert.NoError(t, logger.WriteString("test line 2\n"))

	// Read in the file and check that the content matches what we just wrote
	f, err := os.Open(fn)
	assert.NoError(t, err)
	contents, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	assert.Equal(t, []byte("test line 1\ntest line 2\n"), contents)
}

func BenchmarkWriteBytes(b *testing.B) {
	logger, err := New("test")
	if err != nil {
		log.Fatal(err)
	}
	logger.SetTimeFunction(func() time.Time {
		return time.Unix(0, 0)
	})

	fn := filepath.Join("test", logger.FileName())
	defer os.Remove(fn)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Write([]byte("testing\n"))
	}
}

func BenchmarkWriteStrings(b *testing.B) {
	logger, err := New("test")
	if err != nil {
		log.Fatal(err)
	}
	logger.SetTimeFunction(func() time.Time {
		return time.Unix(0, 0)
	})

	fn := filepath.Join("test", logger.FileName())
	defer os.Remove(fn)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.WriteString("testing\n")
	}
}
