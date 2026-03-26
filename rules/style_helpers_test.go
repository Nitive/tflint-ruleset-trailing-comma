package rules

import (
	"errors"
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type runnerWithOverrides struct {
	*helper.Runner
	getFilesErr error
	files       map[string]*hcl.File
	diags       hcl.Diagnostics
}

func (r *runnerWithOverrides) GetFiles() (map[string]*hcl.File, error) {
	if r.getFilesErr != nil {
		return nil, r.getFilesErr
	}
	if r.files != nil {
		return r.files, nil
	}
	return r.Runner.GetFiles()
}

func (r *runnerWithOverrides) WalkExpressions(walker tflint.ExprWalker) hcl.Diagnostics {
	if r.diags != nil {
		return r.diags
	}
	return r.Runner.WalkExpressions(walker)
}

func TestExprCheckWalkerEnterSkipsAfterError(t *testing.T) {
	expectedErr := errors.New("boom")
	calls := 0
	walker := &exprCheckWalker{
		check: func(hcl.Expression) error {
			calls++
			return expectedErr
		},
	}

	if diags := walker.Enter(nil); diags != nil {
		t.Fatalf("expected nil diagnostics, got %#v", diags)
	}
	if walker.err != expectedErr {
		t.Fatalf("unexpected walker error: %v", walker.err)
	}

	if diags := walker.Enter(nil); diags != nil {
		t.Fatalf("expected nil diagnostics on repeated enter, got %#v", diags)
	}
	if calls != 1 {
		t.Fatalf("expected callback to be invoked once, got %d", calls)
	}

	if diags := walker.Exit(nil); diags != nil {
		t.Fatalf("expected nil diagnostics from exit, got %#v", diags)
	}
}

func TestWalkExpressionsReturnsWalkerError(t *testing.T) {
	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
locals {
  value = [1, 2]
}`,
	})

	expectedErr := errors.New("walk failed")
	err := walkExpressions(runner, func(hcl.Expression) error {
		return expectedErr
	})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWalkExpressionsReturnsDiagnostics(t *testing.T) {
	runner := &runnerWithOverrides{
		Runner: helper.TestRunner(t, map[string]string{
			"main.tf": `
locals {
  value = [1, 2]
}`,
		}),
		diags: hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "walk diagnostics",
			},
		},
	}

	err := walkExpressions(runner, func(hcl.Expression) error {
		t.Fatal("check callback should not be called when walk diagnostics are returned")
		return nil
	})
	var diags hcl.Diagnostics
	if !errors.As(err, &diags) {
		t.Fatalf("expected diagnostics error, got %T", err)
	}
}

func TestAnalyzeDelimitedExpression(t *testing.T) {
	file := mustParseTestFile(t, "main.tf", `
locals {
  tuple = [
    [1, 2],
    3,
  ]
}`)

	attr := mustFindLocalAttribute(t, file, "tuple")
	expr, ok := attr.Expr.(*hclsyntax.TupleConsExpr)
	if !ok {
		t.Fatalf("unexpected expression type: %T", attr.Expr)
	}

	layout, err := analyzeDelimitedExpression(file, expr, expr.OpenRange, hclsyntax.TokenOBrack, hclsyntax.TokenCBrack)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !layout.multiline {
		t.Fatal("expected multiline layout")
	}
	if !layout.hasTrailingComma {
		t.Fatal("expected trailing comma to be detected")
	}
	if len(layout.topLevelCommas) != 2 {
		t.Fatalf("unexpected top-level comma count: %d", len(layout.topLevelCommas))
	}
	if layout.closeRange.Start.Line != 6 {
		t.Fatalf("unexpected close range line: %d", layout.closeRange.Start.Line)
	}
}

func TestAnalyzeDelimitedExpressionWrongCloser(t *testing.T) {
	file := mustParseTestFile(t, "main.tf", `
locals {
  tuple = [1, 2]
}`)

	attr := mustFindLocalAttribute(t, file, "tuple")
	expr := attr.Expr.(*hclsyntax.TupleConsExpr)

	_, err := analyzeDelimitedExpression(file, expr, expr.OpenRange, hclsyntax.TokenOBrack, hclsyntax.TokenCBrace)
	if err == nil {
		t.Fatal("expected error for mismatched closing token")
	}
}

func TestAnalyzeDelimitedExpressionOutOfBounds(t *testing.T) {
	file := mustParseTestFile(t, "main.tf", `
locals {
  tuple = [1, 2]
}`)

	attr := mustFindLocalAttribute(t, file, "tuple")
	expr := attr.Expr.(*hclsyntax.TupleConsExpr)

	file.Bytes = nil
	_, err := analyzeDelimitedExpression(file, expr, expr.OpenRange, hclsyntax.TokenOBrack, hclsyntax.TokenCBrack)
	if err == nil {
		t.Fatal("expected slicing error")
	}
}

func TestRuleCheckGetFilesError(t *testing.T) {
	expectedErr := errors.New("get files failed")
	runner := &runnerWithOverrides{
		Runner:      helper.TestRunner(t, map[string]string{"main.tf": "locals {}"}),
		getFilesErr: expectedErr,
	}

	if err := NewMultilineTrailingCommaRule().Check(runner); !errors.Is(err, expectedErr) {
		t.Fatalf("unexpected trailing comma error: %v", err)
	}
	if err := NewMultilineMapNoCommaRule().Check(runner); !errors.Is(err, expectedErr) {
		t.Fatalf("unexpected map comma error: %v", err)
	}
}

func TestRuleCheckSkipsExpressionsWhenSourceFileIsUnavailable(t *testing.T) {
	runner := &runnerWithOverrides{
		Runner: helper.TestRunner(t, map[string]string{
			"main.tf": `
locals {
  tuple = [1, 2,]
  call = format(
    "%s",
    "a"
  )
  object = {
    a = 1,
  }
}`,
		}),
		files: map[string]*hcl.File{},
	}

	if err := NewMultilineTrailingCommaRule().Check(runner); err != nil {
		t.Fatalf("unexpected trailing comma error: %v", err)
	}
	if err := NewMultilineMapNoCommaRule().Check(runner); err != nil {
		t.Fatalf("unexpected map comma error: %v", err)
	}
}

func mustParseTestFile(t *testing.T, filename, src string) *hcl.File {
	t.Helper()

	file, diags := hclsyntax.ParseConfig([]byte(src), filename, hcl.InitialPos)
	if diags.HasErrors() {
		t.Fatalf("unexpected parse diagnostics: %s", diags.Error())
	}
	return file
}

func mustFindLocalAttribute(t *testing.T, file *hcl.File, name string) *hclsyntax.Attribute {
	t.Helper()

	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		t.Fatalf("unexpected body type: %T", file.Body)
	}

	for _, block := range body.Blocks {
		if block.Type != "locals" {
			continue
		}

		attr, ok := block.Body.Attributes[name]
		if !ok {
			t.Fatalf("attribute %q not found in locals block", name)
		}
		return attr
	}

	t.Fatal(`block "locals" not found`)
	return nil
}
