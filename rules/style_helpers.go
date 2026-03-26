package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type delimitedExprLayout struct {
	multiline        bool
	hasTrailingComma bool
	closeRange       hcl.Range
	lastItemRange    hcl.Range
	topLevelCommas   []hcl.Range
}

type exprCheckWalker struct {
	check func(hcl.Expression) error
	err   error
}

func (w *exprCheckWalker) Enter(expr hcl.Expression) hcl.Diagnostics {
	if w.err != nil {
		return nil
	}

	w.err = w.check(expr)
	return nil
}

func (w *exprCheckWalker) Exit(hcl.Expression) hcl.Diagnostics {
	return nil
}

func walkExpressions(runner tflint.Runner, check func(hcl.Expression) error) error {
	walker := &exprCheckWalker{check: check}
	diags := runner.WalkExpressions(walker)
	if diags.HasErrors() {
		return diags
	}
	return walker.err
}

func analyzeDelimitedExpression(file *hcl.File, expr hcl.Expression, openRange hcl.Range, openType, closeType hclsyntax.TokenType) (delimitedExprLayout, error) {
	layout := delimitedExprLayout{}
	srcRange := expr.Range()
	if !srcRange.CanSliceBytes(file.Bytes) {
		return layout, fmt.Errorf("unable to slice %s from %s", srcRange.String(), srcRange.Filename)
	}

	src := srcRange.SliceBytes(file.Bytes)
	tokens, diags := hclsyntax.LexExpression(src, srcRange.Filename, srcRange.Start)
	if diags.HasErrors() {
		return layout, diags
	}

	started := false
	depth := 0
	lastSignificant := hclsyntax.TokenType(0)
	lastSignificantRange := hcl.Range{}
	hasLastSignificant := false

	for _, token := range tokens {
		if !started {
			if token.Type == openType && sameRange(token.Range, openRange) {
				started = true
				depth = 1
			}
			continue
		}

		if token.Type == closeType && depth == 1 {
			layout.hasTrailingComma = hasLastSignificant && lastSignificant == hclsyntax.TokenComma
			layout.closeRange = token.Range
			layout.lastItemRange = lastSignificantRange
			return layout, nil
		}

		if token.Type == hclsyntax.TokenNewline && depth == 1 {
			layout.multiline = true
		}
		if token.Type == hclsyntax.TokenComma && depth == 1 {
			layout.topLevelCommas = append(layout.topLevelCommas, token.Range)
		}

		switch token.Type {
		case hclsyntax.TokenOBrace, hclsyntax.TokenOBrack, hclsyntax.TokenOParen:
			depth++
		case hclsyntax.TokenCBrace, hclsyntax.TokenCBrack, hclsyntax.TokenCParen:
			depth--
		}

		if isIgnoredToken(token.Type) {
			continue
		}
		lastSignificant = token.Type
		lastSignificantRange = token.Range
		hasLastSignificant = true
	}

	return layout, fmt.Errorf("unable to locate closing %q for %s", closeType, srcRange.String())
}

func sameRange(left, right hcl.Range) bool {
	return left.Filename == right.Filename &&
		left.Start.Byte == right.Start.Byte &&
		left.End.Byte == right.End.Byte
}

func isIgnoredToken(tokenType hclsyntax.TokenType) bool {
	switch tokenType {
	case hclsyntax.TokenComment, hclsyntax.TokenNewline, hclsyntax.TokenTabs, hclsyntax.TokenEOF:
		return true
	default:
		return false
	}
}
