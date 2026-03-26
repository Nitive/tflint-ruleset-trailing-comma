package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestMultilineMapNoCommaRule(t *testing.T) {
	tests := []struct {
		name            string
		content         string
		expected        helper.Issues
		expectedChanges map[string]string
	}{
		{
			name: "issues found",
			content: `
locals {
  bad_single_line_map = { a = 1, b = 2, }
  bad_multiline_map = {
    a = [1, 2, 3],
    b = {
      c = 3
    },
  }
}`,
			expected: helper.Issues{
				{
					Rule:    NewMultilineMapNoCommaRule(),
					Message: "single-line maps must not end with a trailing comma",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 39},
						End:      hcl.Pos{Line: 3, Column: 40},
					},
				},
				{
					Rule:    NewMultilineMapNoCommaRule(),
					Message: "multiline maps must not contain commas",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 5, Column: 18},
						End:      hcl.Pos{Line: 5, Column: 19},
					},
				},
				{
					Rule:    NewMultilineMapNoCommaRule(),
					Message: "multiline maps must not contain commas",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 8, Column: 6},
						End:      hcl.Pos{Line: 8, Column: 7},
					},
				},
			},
			expectedChanges: map[string]string{
				"main.tf": `
locals {
  bad_single_line_map = { a = 1, b = 2 }
  bad_multiline_map = {
    a = [1, 2, 3]
    b = {
      c = 3
    }
  }
}`,
			},
		},
		{
			name: "issues not found",
			content: `
locals {
  empty_map = {}
  good_single_line_map = { a = 1, b = 2 }
  good_multiline_map = {
    a = [1, 2, 3]
    b = {
      c = 3
    }
  }
}`,
			expected:        helper.Issues{},
			expectedChanges: map[string]string{},
		},
	}

	rule := NewMultilineMapNoCommaRule()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"main.tf": test.content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, test.expected, runner.Issues)
			helper.AssertChanges(t, test.expectedChanges, runner.Changes())
		})
	}
}
