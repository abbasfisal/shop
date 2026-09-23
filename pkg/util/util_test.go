package util

import "testing"

func TestStringToUint(t *testing.T) {
	if got := StringToUint("123"); got != 123 {
		t.Fatalf("expected 123, got %d", got)
	}
	if got := StringToUint("abc"); got != 0 {
		t.Fatalf("expected 0 for invalid input, got %d", got)
	}
}

func TestHasSuffix(t *testing.T) {
	if !HasSuffix("main.go", ".go") {
		t.Fatal("expected suffix match")
	}
	if HasSuffix("main", ".go") {
		t.Fatal("expected no match")
	}
}

func TestGeneratePageNumbers_Window(t *testing.T) {
	// current=5, total=10 → start=3, show min(10,4)=4 → [3,4,5,6]
	got := GeneratePageNumbers(5, 10)
	want := []int{3, 4, 5, 6}
	if len(got) != len(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected %v, got %v", want, got)
		}
	}
}

func TestGeneratePageNumbers_SmallTotal(t *testing.T) {
	got := GeneratePageNumbers(1, 2)
	if len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Fatalf("unexpected pages: %v", got)
	}
}

func TestValidateIRMobile(t *testing.T) {
	if !ValidateIRMobile("09121234567") {
		t.Fatal("valid mobile rejected")
	}
	if ValidateIRMobile("08121234567") {
		t.Fatal("invalid mobile accepted")
	}
	if ValidateIRMobile("0912123456") {
		t.Fatal("short mobile accepted")
	}
}

func TestAllowImageExtensions(t *testing.T) {
	exts := AllowImageExtensions()
	found := map[string]bool{}
	for _, e := range exts {
		found[e] = true
	}
	for _, want := range []string{".jpg", ".png", ".jpeg", ".webp"} {
		if !found[want] {
			t.Fatalf("missing extension %s", want)
		}
	}
}

func TestGetProductStoragePath_Local(t *testing.T) {
	t.Setenv("STORAGE_STATUS", "deactive")
	if got := GetProductStoragePath(); got != "/uploads/media/products/" {
		t.Fatalf("unexpected local storage path: %q", got)
	}
}
