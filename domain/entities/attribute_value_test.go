package entities

import (
	"testing"

	"gorm.io/datatypes"
)

func TestColorHex_ParsesValidHex(t *testing.T) {
	v := &AttributeValue{Meta: datatypes.JSON([]byte(`{"hex": "#a020f0"}`))}
	if got := v.ColorHex(); got != "#a020f0" {
		t.Fatalf("want #a020f0 got %q", got)
	}
}

func TestColorHex_RejectsGarbage(t *testing.T) {
	cases := []string{"", "{}", `{"hex": ""}`, `{"hex": "red"}`, `{"hex": "#fff"}`,
		`{"hex": "#gggggg"}`, `{"hex": "#1234567"}`, `not-json`, `{"hex": "#ABCDEF"}`}
	// note: #ABCDEF (uppercase) IS valid and must pass
	for _, raw := range cases[:len(cases)-1] {
		v := &AttributeValue{Meta: datatypes.JSON([]byte(raw))}
		if got := v.ColorHex(); got != "" {
			t.Fatalf("raw %q: want empty got %q", raw, got)
		}
	}
	v := &AttributeValue{Meta: datatypes.JSON([]byte(cases[len(cases)-1]))}
	if got := v.ColorHex(); got != "#ABCDEF" {
		t.Fatalf("uppercase hex: want #ABCDEF got %q", got)
	}
	if got := (&AttributeValue{}).ColorHex(); got != "" {
		t.Fatalf("nil meta: want empty got %q", got)
	}
}

func TestAttributeInputTypeHelpers(t *testing.T) {
	if !(&Attribute{InputType: AttributeInputColor}).IsColor() {
		t.Fatal("color attribute must report IsColor")
	}
	if (&Attribute{InputType: AttributeInputText}).IsColor() {
		t.Fatal("text attribute must not report IsColor")
	}
	if (&Attribute{}).IsColor() {
		t.Fatal("nil/empty attribute must not report IsColor")
	}
	if !ValidAttributeInputType("text") || !ValidAttributeInputType("color") || ValidAttributeInputType("image") {
		t.Fatal("ValidAttributeInputType mismatch")
	}
}
