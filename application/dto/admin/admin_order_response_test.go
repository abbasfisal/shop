package responses

import (
	"testing"

	"shop/domain/entities"
)

func TestToOrderItemAttributes_SkipsUnlinkedValues(t *testing.T) {
	vavs := []*entities.VariantAttributeValue{
		{AttributeValue: &entities.AttributeValue{AttributeTitle: "سایز", Value: "L"}},
		{AttributeValue: nil}, // not preloaded — must be skipped, not panic
		{AttributeValue: &entities.AttributeValue{AttributeTitle: "رنگ", Value: "آبی"}},
	}

	out := ToOrderItemAttributes(vavs)
	if len(out.Data) != 2 {
		t.Fatalf("expected 2 attributes, got %d", len(out.Data))
	}
	if out.Data[0].Title != "سایز" || out.Data[0].Value != "L" {
		t.Fatalf("unexpected first attribute: %+v", out.Data[0])
	}
	if out.Data[1].Title != "رنگ" || out.Data[1].Value != "آبی" {
		t.Fatalf("unexpected second attribute: %+v", out.Data[1])
	}
}

func TestToOrderItemAttributes_Empty(t *testing.T) {
	out := ToOrderItemAttributes(nil)
	if out == nil || len(out.Data) != 0 {
		t.Fatalf("expected empty result, got %+v", out)
	}
}
