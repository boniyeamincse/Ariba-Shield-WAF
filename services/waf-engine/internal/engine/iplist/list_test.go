package iplist

import "testing"

func TestSetAndCheck(t *testing.T) {
	l := New()
	if err := l.SetBlocked([]string{"203.0.113.0/24", "10.0.0.1"}); err != nil {
		t.Fatalf("SetBlocked: %v", err)
	}
	if err := l.SetAllowed([]string{"198.51.100.5"}); err != nil {
		t.Fatalf("SetAllowed: %v", err)
	}

	if !l.IsBlocked("203.0.113.9") {
		t.Error("203.0.113.9 should be blocked by /24")
	}
	if !l.IsBlocked("10.0.0.1") {
		t.Error("10.0.0.1 should be blocked by /32")
	}
	if l.IsBlocked("203.0.114.1") {
		t.Error("203.0.114.1 should not be blocked")
	}
	if !l.IsAllowed("198.51.100.5") {
		t.Error("198.51.100.5 should be allowed")
	}
	if l.IsAllowed("198.51.100.6") {
		t.Error("198.51.100.6 should not be allowed (only .5)")
	}
}

func TestInvalidInputs(t *testing.T) {
	l := New()
	if err := l.SetBlocked([]string{"not-an-ip"}); err == nil {
		t.Error("expected error for invalid CIDR")
	}
	// Invalid addresses are simply not matched, never panic.
	if l.IsBlocked("garbage") {
		t.Error("garbage should not match")
	}
	if l.IsAllowed("") {
		t.Error("empty should not match")
	}
}

func TestAtomicReplace(t *testing.T) {
	l := New()
	_ = l.SetBlocked([]string{"203.0.113.0/24"})
	_ = l.SetBlocked([]string{"198.51.100.0/24"}) // replaces, not appends
	if l.IsBlocked("203.0.113.1") {
		t.Error("old list should be replaced, not appended")
	}
	if !l.IsBlocked("198.51.100.1") {
		t.Error("new list should apply")
	}
}