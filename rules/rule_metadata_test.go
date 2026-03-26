package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func TestRuleMetadata(t *testing.T) {
	tests := []struct {
		name        string
		rule        tflint.Rule
		wantName    string
		wantEnabled bool
		wantLevel   tflint.Severity
		wantLink    string
	}{
		{
			name:        "multiline_trailing_comma",
			rule:        NewMultilineTrailingCommaRule(),
			wantName:    "multiline_trailing_comma",
			wantEnabled: true,
			wantLevel:   tflint.ERROR,
			wantLink:    "",
		},
		{
			name:        "multiline_map_no_comma",
			rule:        NewMultilineMapNoCommaRule(),
			wantName:    "multiline_map_no_comma",
			wantEnabled: true,
			wantLevel:   tflint.ERROR,
			wantLink:    "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.rule.Name(); got != test.wantName {
				t.Fatalf("unexpected rule name: got %s want %s", got, test.wantName)
			}

			if got := test.rule.Enabled(); got != test.wantEnabled {
				t.Fatalf("unexpected enabled flag: got %t want %t", got, test.wantEnabled)
			}

			if got := test.rule.Severity(); got != test.wantLevel {
				t.Fatalf("unexpected severity: got %s want %s", got, test.wantLevel)
			}

			if got := test.rule.Link(); got != test.wantLink {
				t.Fatalf("unexpected link: got %s want %s", got, test.wantLink)
			}
		})
	}
}
