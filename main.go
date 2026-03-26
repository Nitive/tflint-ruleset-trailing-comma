package main

import (
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/Nitive/tflint-ruleset-trailing-comma/rules"
)

func newRuleSet() *tflint.BuiltinRuleSet {
	return &tflint.BuiltinRuleSet{
		Name:    "tflint-ruleset-trailing-comma",
		Version: "0.1.0",
		Rules: []tflint.Rule{
			rules.NewMultilineTrailingCommaRule(),
			rules.NewMultilineMapNoCommaRule(),
		},
	}
}

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: newRuleSet(),
	})
}
