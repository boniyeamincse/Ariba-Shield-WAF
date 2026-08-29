package handlers

import "testing"

func TestSimpleDiff(t *testing.T) {
	from := []byte(`{"mode":"transparent","threshold":1,"enabled":true}`)
	to := []byte(`{"mode":"blocking","threshold":5,"enabled":true,"new_field":"x"}`)

	diff := simpleDiff(from, to)
	added := diff["added"].([]string)
	changed := diff["changed"].([]string)
	removed := diff["removed"].([]string)

	if len(added) != 1 || added[0] != "new_field" {
		t.Fatalf("expected added=[new_field], got %v", added)
	}
	if len(changed) != 2 {
		t.Fatalf("expected 2 changed fields, got %v", changed)
	}
	if len(removed) != 0 {
		t.Fatalf("expected no removed fields, got %v", removed)
	}
}

func TestIdentityDiff(t *testing.T) {
	doc := []byte(`{"mode":"transparent","threshold":1}`)
	diff := simpleDiff(doc, doc)
	if len(diff["added"].([]string)) != 0 {
		t.Fatal("identity diff should have no added")
	}
	if len(diff["changed"].([]string)) != 0 {
		t.Fatal("identity diff should have no changed")
	}
}