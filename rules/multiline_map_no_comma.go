package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// MultilineMapNoCommaRule enforces comma style for maps.
type MultilineMapNoCommaRule struct {
	tflint.DefaultRule
}

// NewMultilineMapNoCommaRule returns a new rule.
func NewMultilineMapNoCommaRule() *MultilineMapNoCommaRule {
	return &MultilineMapNoCommaRule{}
}

// Name returns the rule name.
func (r *MultilineMapNoCommaRule) Name() string {
	return "multiline_map_no_comma"
}

// Enabled returns whether the rule is enabled by default.
func (r *MultilineMapNoCommaRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *MultilineMapNoCommaRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *MultilineMapNoCommaRule) Link() string {
	return ""
}

// Check checks whether maps follow the configured comma style.
func (r *MultilineMapNoCommaRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	return walkExpressions(runner, func(expr hcl.Expression) error {
		objectExpr, ok := expr.(*hclsyntax.ObjectConsExpr)
		if !ok || len(objectExpr.Items) == 0 {
			return nil
		}

		file, ok := files[expr.Range().Filename]
		if !ok {
			return nil
		}

		layout, err := analyzeDelimitedExpression(file, objectExpr, objectExpr.OpenRange, hclsyntax.TokenOBrace, hclsyntax.TokenCBrace)
		if err != nil {
			return err
		}

		if layout.multiline {
			for _, commaRange := range layout.topLevelCommas {
				if err := runner.EmitIssueWithFix(
					r,
					"multiline maps must not contain commas",
					commaRange,
					func(fixer tflint.Fixer) error {
						return fixer.Remove(commaRange)
					},
				); err != nil {
					return err
				}
			}
			return nil
		}

		if !layout.hasTrailingComma {
			return nil
		}

		commaRange := layout.topLevelCommas[len(layout.topLevelCommas)-1]
		return runner.EmitIssueWithFix(
			r,
			"single-line maps must not end with a trailing comma",
			commaRange,
			func(fixer tflint.Fixer) error {
				return fixer.Remove(commaRange)
			},
		)
	})
}
