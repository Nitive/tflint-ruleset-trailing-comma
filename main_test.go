package main

import "testing"

func TestNewRuleSet(t *testing.T) {
	ruleSet := newRuleSet()

	if got := ruleSet.RuleSetName(); got != "tflint-ruleset-trailing-comma" {
		t.Fatalf("unexpected ruleset name: %s", got)
	}

	if got := ruleSet.RuleSetVersion(); got != "0.1.0" {
		t.Fatalf("unexpected ruleset version: %s", got)
	}

	wantRules := []string{"multiline_trailing_comma", "multiline_map_no_comma"}
	gotRules := ruleSet.RuleNames()
	if len(gotRules) != len(wantRules) {
		t.Fatalf("unexpected rule count: got %d want %d", len(gotRules), len(wantRules))
	}

	for i, want := range wantRules {
		if gotRules[i] != want {
			t.Fatalf("unexpected rule at index %d: got %s want %s", i, gotRules[i], want)
		}
	}
}
