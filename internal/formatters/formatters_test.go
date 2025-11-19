package formatters

import (
	"reflect"
	"testing"

	"github.com/larsks/kwctl/pkg/radio/types"
)

func TestHeadersFromStruct_WithChannel(t *testing.T) {
	expected := []string{
		"Name", "Number", "RxFreq", "RxStep", "Shift", "Reverse",
		"Tone", "CTCSS", "DCS", "ToneFreq", "CTCSSFreq", "DCSCode",
		"Offset", "Mode", "TxFreq", "TxStep", "Lockout",
	}

	result := HeadersFromStruct(types.Channel{})

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("HeadersFromStruct(Channel{}) = %v, want %v", result, expected)
	}
}

func TestHeadersFromStruct_WithPointer(t *testing.T) {
	expected := []string{
		"Name", "Number", "RxFreq", "RxStep", "Shift", "Reverse",
		"Tone", "CTCSS", "DCS", "ToneFreq", "CTCSSFreq", "DCSCode",
		"Offset", "Mode", "TxFreq", "TxStep", "Lockout",
	}

	result := HeadersFromStruct(&types.Channel{})

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("HeadersFromStruct(&Channel{}) = %v, want %v", result, expected)
	}
}

func TestHeadersFromStruct_WithCustomTags(t *testing.T) {
	type TestStruct struct {
		RxFreq int `header:"RX Frequency"`
		TxFreq int `header:"TX Frequency"`
		Name   string
	}

	expected := []string{"RX Frequency", "TX Frequency", "Name"}
	result := HeadersFromStruct(TestStruct{})

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("HeadersFromStruct(TestStruct{}) = %v, want %v", result, expected)
	}
}

func TestHeadersFromStruct_WithUnexportedFields(t *testing.T) {
	type TestStruct struct {
		PublicField   string
		privateField  string
		AnotherPublic int
	}

	expected := []string{"PublicField", "AnotherPublic"}
	result := HeadersFromStruct(TestStruct{})

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("HeadersFromStruct(TestStruct{}) = %v, want %v (unexported fields should be ignored)", result, expected)
	}
}

func TestHeadersFromStruct_WithEmptyStruct(t *testing.T) {
	type EmptyStruct struct{}

	expected := []string{}
	result := HeadersFromStruct(EmptyStruct{})

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("HeadersFromStruct(EmptyStruct{}) = %v, want %v", result, expected)
	}
}

func TestHeadersFromStruct_WithNonStruct(t *testing.T) {
	tests := []struct {
		name  string
		input any
	}{
		{"string", "test"},
		{"int", 42},
		{"slice", []string{"a", "b"}},
		{"map", map[string]int{"key": 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HeadersFromStruct(tt.input)
			if len(result) != 0 {
				t.Errorf("HeadersFromStruct(%v) = %v, want empty slice", tt.input, result)
			}
		})
	}
}

func TestHeadersFromStruct_WithMixedTags(t *testing.T) {
	type TestStruct struct {
		Field1 string `header:"Custom Header 1"`
		Field2 int
		Field3 bool `header:"Custom Header 3"`
	}

	expected := []string{"Custom Header 1", "Field2", "Custom Header 3"}
	result := HeadersFromStruct(TestStruct{})

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("HeadersFromStruct(TestStruct{}) = %v, want %v", result, expected)
	}
}
