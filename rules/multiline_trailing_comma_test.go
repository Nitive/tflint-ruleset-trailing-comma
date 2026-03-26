package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestMultilineTrailingCommaRule(t *testing.T) {
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
  bad_single_line_list = [1, 2,]
  bad_multiline_list = [
    "a",
    "b"
  ]
  bad_single_line_call = format("%s %s", "a", "b",)
  bad_multiline_call = format(
    "%s %s",
    "a",
    "b"
  )
}`,
			expected: helper.Issues{
				{
					Rule:    NewMultilineTrailingCommaRule(),
					Message: "single-line lists must not end with a trailing comma",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 31},
						End:      hcl.Pos{Line: 3, Column: 32},
					},
				},
				{
					Rule:    NewMultilineTrailingCommaRule(),
					Message: "multiline lists must end with a trailing comma",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 7, Column: 3},
						End:      hcl.Pos{Line: 7, Column: 4},
					},
				},
				{
					Rule:    NewMultilineTrailingCommaRule(),
					Message: "single-line function calls must not end with a trailing comma",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 8, Column: 50},
						End:      hcl.Pos{Line: 8, Column: 51},
					},
				},
				{
					Rule:    NewMultilineTrailingCommaRule(),
					Message: "multiline function calls must end with a trailing comma",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 13, Column: 3},
						End:      hcl.Pos{Line: 13, Column: 4},
					},
				},
			},
			expectedChanges: map[string]string{
				"main.tf": `
locals {
  bad_single_line_list = [1, 2]
  bad_multiline_list = [
    "a",
    "b",
  ]
  bad_single_line_call = format("%s %s", "a", "b")
  bad_multiline_call = format(
    "%s %s",
    "a",
    "b",
  )
}`,
			},
		},
		{
			name: "issues not found",
			content: `
locals {
  empty_list = []
  empty_call = timestamp()
  good_single_line_list = [1, 2]
  good_multiline_list = [
    "a",
    format(
      "%s %s",
      "a",
      "b",
    ),
  ]
  good_single_line_call = format("%s %s", "a", "b")
  good_multiline_call = format(
    "%s %s",
    "a",
    "b",
  )
}`,
			expected:        helper.Issues{},
			expectedChanges: map[string]string{},
		},
	}

	rule := NewMultilineTrailingCommaRule()

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
