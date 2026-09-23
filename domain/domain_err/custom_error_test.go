package domain_err

import (
	"errors"
	"testing"

	"gorm.io/gorm"
)

func TestHandleError_RecordNotFound(t *testing.T) {
	err := HandleError(gorm.ErrRecordNotFound, RecordNotFound)
	if err.Code != 404 {
		t.Fatalf("expected code 404, got %d", err.Code)
	}
	if err.DisplayMessage != RecordNotFound {
		t.Fatalf("expected display message %q, got %q", RecordNotFound, err.DisplayMessage)
	}
}

func TestHandleError_InternalError(t *testing.T) {
	err := HandleError(errors.New("db exploded"), RecordNotFound)
	if err.Code != 500 {
		t.Fatalf("expected code 500, got %d", err.Code)
	}
	if err.DisplayMessage != InternalServerError {
		t.Fatalf("expected display message %q, got %q", InternalServerError, err.DisplayMessage)
	}
	if err.OriginalMessage != "db exploded" {
		t.Fatalf("expected original message preserved, got %q", err.OriginalMessage)
	}
}

func TestNewAndError(t *testing.T) {
	err := New("original", "نمایشی", 422)
	if err.Error() != "نمایشی" {
		t.Fatalf("Error() should return display message, got %q", err.Error())
	}
	if err.Code != 422 || err.OriginalMessage != "original" {
		t.Fatalf("unexpected fields: %+v", err)
	}
}
