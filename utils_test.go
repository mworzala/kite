package main

import (
	"bytes"
	"testing"
)

func TestReadUUID(t *testing.T) {
	// Create a test case with a sample UUID
	uuid := "123e4567-e89b-12d3-a456-426614174000"
	buf := bytes.NewReader([]byte{
		0x12, 0x3e, 0x45, 0x67, 0xe8, 0x9b, 0x12, 0xd3,
		0xa4, 0x56, 0x42, 0x66, 0x14, 0x17, 0x40, 0x00,
	})

	// Call the readUUID function
	result, err := readUUID(buf)

	// Check if the result matches the expected UUID
	if result != uuid {
		t.Errorf("Expected UUID: %s, Got: %s", uuid, result)
	}

	// Check if there was any error during decoding
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestReadUUID_InvalidInput(t *testing.T) {
	// Create a test case with invalid input
	buf := bytes.NewReader([]byte{
		0x12, 0x3e, 0x45, 0x67, 0xe8, 0x9b, 0x12, 0xd3,
		0xa4, 0x56, 0x42, 0x66, 0x14, 0x17, 0x40, // Missing last byte
	})

	// Call the readUUID function
	result, err := readUUID(buf)

	// Check if the result is empty
	if result != "" {
		t.Errorf("Expected empty result, Got: %s", result)
	}

	// Check if there was an error during decoding
	if err == nil {
		t.Error("Expected error, Got nil")
	}
}

func TestWriteVarInt(t *testing.T) {
	testCases := []struct {
		value    int
		expected []byte
	}{
		{0, []byte{0x00}},
		{127, []byte{0x7f}},
		{128, []byte{0x80, 0x01}},
		{300, []byte{0xac, 0x02}},
		{16383, []byte{0xff, 0x7f}},
		{16384, []byte{0x80, 0x80, 0x01}},
	}

	for _, tc := range testCases {
		buf := new(bytes.Buffer)
		err := writeVarInt(buf, tc.value)
		if err != nil {
			t.Errorf("Unexpected error for value %d: %v", tc.value, err)
		}
		if !bytes.Equal(buf.Bytes(), tc.expected) {
			t.Errorf("For value %d, Expected: %v, Got: %v", tc.value, tc.expected, buf.Bytes())
		}
	}
}
func TestReadVarInt(t *testing.T) {
	testCases := []struct {
		value    int
		encoding []byte
	}{
		{0, []byte{0x00}},
		{127, []byte{0x7f}},
		{128, []byte{0x80, 0x01}},
		{300, []byte{0xac, 0x02}},
		{16383, []byte{0xff, 0x7f}},
		{16384, []byte{0x80, 0x80, 0x01}},
	}

	for _, tc := range testCases {
		buf := bytes.NewBuffer(tc.encoding)
		reader := bytes.NewReader(buf.Bytes())
		result, err := readVarInt(reader)
		if err != nil {
			t.Errorf("Unexpected error for value %d: %v", tc.value, err)
		}
		if result != tc.value {
			t.Errorf("Expected: %d, Got: %d for value %d", tc.value, result, tc.value)
		}
	}
}
