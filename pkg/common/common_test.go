package common //nolint:revive

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadZlibFromBuffer(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name          string
		data          []byte
		expectedData  []byte
		expectedError error
	}{
		{
			name:          "empty data",
			data:          []byte("\x78\xda\x01\x00\x00\xff\xff\x00\x00\x00\x01"),
			expectedData:  []byte(""),
			expectedError: nil,
		},
		{
			name:          "some data",
			data:          []byte("\x78\xda\x73\x4f\xcc\x4d\x2d\x4f\x2c\x4b\x05\x00\x0d\xbe\x03\x2e"),
			expectedData:  []byte("Gamewave"),
			expectedError: nil,
		},
	}

	for _, tt := range cases {
		name := tt.name
		data := tt.data
		expectedData := tt.expectedData
		expectedError := tt.expectedError
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			data, err := ReadZlibFromBuffer(data)
			require.Equal(t, expectedError, err)
			require.Equal(t, expectedData, data)
		})
	}
}

func TestWriteZlibToBuffer(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name          string
		data          []byte
		expectedData  []byte
		expectedError error
	}{
		{
			name:          "empty data",
			data:          []byte{},
			expectedData:  []byte("\x78\xda\x01\x00\x00\xff\xff\x00\x00\x00\x01"),
			expectedError: nil,
		},
		{
			name:          "some data",
			data:          []byte("Gamewave"),
			expectedData:  []byte("\x78\xda\x72\x4f\xcc\x4d\x2d\x4f\x2c\x4b\x05\x04\x00\x00\xff\xff\x0d\xbe\x03\x2e"),
			expectedError: nil,
		},
	}

	for _, tt := range cases {
		name := tt.name
		data := tt.data
		expectedData := tt.expectedData
		expectedError := tt.expectedError
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			data, err := WriteZlibToBuffer(data)
			require.Equal(t, expectedError, err)
			require.Equal(t, expectedData, data)
		})
	}
}
