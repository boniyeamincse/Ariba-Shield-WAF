package ratelimit

import (
	"testing"
	"time"
)

func TestAllowWithinLimit(t *testing.T) {
	r := New(3, time.Minute)
	for i := 0; i < 3; i++ {
		ok, remaining := r.Allow("1.2.3.4")
		if !ok {
			t.Fatalf("request %d should be allowed", i+1)
		}
		if remaining != 3-i-1 {
			t.Fatalf("request %d: expected remaining %d, got %d", i+1, 3-i-1, remaining)
		}
	}
	ok, _ := r.Allow("1.2.3.4")
	if ok {
		t.Fatal("4th request should be rate-limited")
	}
}

func TestWindowReset(t *testing.T) {
	r := NewWithMax(1, 50*time.Millisecond, 100)
	if ok, _ := r.Allow("k"); !ok {
		t.Fatal("first request should pass")
	}
	if ok, _ := r.Allow("k"); ok {
		t.Fatal("second request within window should be limited")
	}
	time.Sleep(60 * time.Millisecond)
	if ok, _ := r.Allow("k"); !ok {
		t.Fatal("request after window should pass again")
	}
}

func TestEvictionBounded(t *testing.T) {
	// Tiny maxKeys so eviction must kick in.
	r := NewWithMax(5, time.Minute, 10)
	// Touch many distinct keys.
	for i := 0; i < 1000; i++ {
		key := string(rune('a' + i%26)) + itoa(i)
		r.Allow(key)
	}
	if n := r.Len(); n > 100 {
		t.Fatalf("buckets map should be bounded, got %d entries", n)
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [16]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}