package entities

import "testing"

func uintRef(v uint) *uint { return &v }

// TestQuoteOrderFees covers the checkout fee split: schedule-less tariffs,
// the free-shipping threshold, inactive/missing tariffs and the grand total.
func TestQuoteOrderFees(t *testing.T) {
	shipActive := &FeeRate{Kind: FeeKindShipping, Title: "s", Amount: 104_000, Status: true}
	packActive := &FeeRate{Kind: FeeKindPackaging, Title: "p", Amount: 23_000, Status: true}

	cases := []struct {
		name                        string
		items                       uint
		ship, pack                  *FeeRate
		wantShip, wantPack, wantGrd uint
		wantFree                    bool
	}{
		{"plain quote", 676_875, shipActive, packActive, 104_000, 23_000, 803_875, false},
		{"no tariffs", 100, nil, nil, 0, 0, 100, false},
		{"inactive ignored", 100,
			&FeeRate{Kind: FeeKindShipping, Amount: 50, Status: false},
			&FeeRate{Kind: FeeKindPackaging, Amount: 60, Status: false},
			0, 0, 100, false},
		{"zero items still charges fees", 0, shipActive, packActive, 104_000, 23_000, 127_000, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ship, pack, grand, free := QuoteOrderFees(tc.items, tc.ship, tc.pack)
			if ship != tc.wantShip || pack != tc.wantPack || grand != tc.wantGrd || free != tc.wantFree {
				t.Fatalf("got (%d,%d,%d,%v) want (%d,%d,%d,%v)",
					ship, pack, grand, free, tc.wantShip, tc.wantPack, tc.wantGrd, tc.wantFree)
			}
		})
	}
}

func TestQuoteOrderFees_FreeShippingThreshold(t *testing.T) {
	ship := &FeeRate{Kind: FeeKindShipping, Amount: 104_000, Status: true, FreeThreshold: uintRef(1_000_000)}

	// below threshold → shipping charged
	shipAmt, _, grand, free := QuoteOrderFees(676_875, ship, nil)
	if shipAmt != 104_000 || grand != 780_875 || free {
		t.Fatalf("below threshold: got (%d,%d,%v)", shipAmt, grand, free)
	}

	// at/above threshold → free shipping
	for _, total := range []uint{1_000_000, 100_000_000} {
		shipAmt, _, grand, free = QuoteOrderFees(total, ship, nil)
		if shipAmt != 0 || grand != total || !free {
			t.Fatalf("above threshold (%d): got (%d,%d,%v)", total, shipAmt, grand, free)
		}
	}

	// a zero threshold disables the rule (never free)
	noRule := &FeeRate{Kind: FeeKindShipping, Amount: 104_000, Status: true, FreeThreshold: uintRef(0)}
	if shipAmt, _, _, free := QuoteOrderFees(50_000_000, noRule, nil); shipAmt != 104_000 || free {
		t.Fatalf("zero threshold must not trigger free shipping: got (%d,%v)", shipAmt, free)
	}
}

func TestFeeKindLabel(t *testing.T) {
	if FeeKindLabel(FeeKindShipping) == "" || FeeKindLabel(FeeKindPackaging) == "" {
		t.Fatal("fee kind labels must not be empty")
	}
	if !ValidFeeKind(FeeKindShipping) || !ValidFeeKind(FeeKindPackaging) || ValidFeeKind("other") {
		t.Fatal("ValidFeeKind mismatch")
	}
}
