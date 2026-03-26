package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// MultilineTrailingCommaRule enforces trailing comma style for lists and function calls.
type MultilineTrailingCommaRule struct {
	tflint.DefaultRule
}

// NewMultilineTrailingCommaRule returns a new rule.
func NewMultilineTrailingCommaRule() *MultilineTrailingCommaRule {
	return &MultilineTrailingCommaRule{}
}

// Name returns the rule name.
func (r *MultilineTrailingCommaRule) Name() string {
	return "multiline_trailing_comma"
}

// Enabled returns whether the rule is enabled by default.
func (r *MultilineTrailingCommaRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *MultilineTrailingCommaRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *MultilineTrailingCommaRule) Link() string {
	return ""
}

// Check checks whether lists and function calls follow the configured comma style.
func (r *MultilineTrailingCommaRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	return walkExpressions(runner, func(expr hcl.Expression) error {
		filename := expr.Range().Filename
		file, ok := files[filename]
		if !ok {
			return nil
		}

		switch expr := expr.(type) {
		case *hclsyntax.TupleConsExpr:
			if len(expr.Exprs) == 0 {
				return nil
			}
			layout, err := analyzeDelimitedExpression(file, expr, expr.OpenRange, hclsyntax.TokenOBrack, hclsyntax.TokenCBrack)
			if err != nil {
				return err
			}
			return emitTrailingCommaIssue(runner, r, "list", layout)
		case *hclsyntax.FunctionCallExpr:
			if len(expr.Args) == 0 {
				return nil
			}
			layout, err := analyzeDelimitedExpression(file, expr, expr.OpenParenRange, hclsyntax.TokenOParen, hclsyntax.TokenCParen)
			if err != nil {
				return err
			}
			return emitTrailingCommaIssue(runner, r, "function call", layout)
		default:
			return nil
		}
	})
}

func emitTrailingCommaIssue(runner tflint.Runner, rule tflint.Rule, subject string, layout delimitedExprLayout) error {
	if layout.multiline {
		if layout.hasTrailingComma {
			return nil
		}

		return runner.EmitIssueWithFix(
			rule,
			"multiline "+subject+"s must end with a trailing comma",
			layout.closeRange,
			func(fixer tflint.Fixer) error {
				return fixer.InsertTextAfter(layout.lastItemRange, ",")
			},
		)
	}

	if !layout.hasTrailingComma {
		return nil
	}

	commaRange := layout.topLevelCommas[len(layout.topLevelCommas)-1]
	return runner.EmitIssueWithFix(
		rule,
		"single-line "+subject+"s must not end with a trailing comma",
		commaRange,
		func(fixer tflint.Fixer) error {
			return fixer.Remove(commaRange)
		},
	)
}
