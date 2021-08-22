package ui

import (
	"testing"
)

func TestUI(t *testing.T) {
	ui := New()

	if ui.App == nil {
		t.Errorf("failed New() App is Nil")
	}
	if ui.Pages == nil {
		t.Errorf("failed New() Pages is Nil")
	}
	if ui.UrlField == nil {
		t.Errorf("failed New() UrlField is Nil")
	}
	if ui.ParamsTable == nil {
		t.Errorf("failed New() ParamsTable is Nil")
	}
	if ui.BodyTable == nil {
		t.Errorf("failed New() BodyTable is Nil")
	}
	if ui.ResTextView == nil {
		t.Errorf("failed New() ResTextView is Nil")
	}
	if ui.HTTPTextView == nil {
		t.Errorf("failed New() HTTPTextView is Nil")
	}
}
