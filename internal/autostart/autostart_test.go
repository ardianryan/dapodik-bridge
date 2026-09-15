package autostart

import "testing"

func TestNewManager(t *testing.T) {
	mgr := NewManager()
	if mgr == nil {
		t.Fatal("expected non-nil autostart manager")
	}
}
