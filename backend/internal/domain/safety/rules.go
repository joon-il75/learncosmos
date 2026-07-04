package safety

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"
)

var builtinRules = []Rule{
	{RuleType: "keyword", Pattern: "죽이는 법", RiskType: "unsafe_instruction", Action: ActionBlock, Description: "Harmful how-to phrasing"},
	{RuleType: "keyword", Pattern: "폭탄 만들기", RiskType: "unsafe_instruction", Action: ActionBlock, Description: "Weapon construction phrasing"},
	{RuleType: "keyword", Pattern: "자살 방법", RiskType: "self_harm", Action: ActionBlock, Description: "Self-harm instruction phrasing"},
	{RuleType: "keyword", Pattern: "how to kill", RiskType: "unsafe_instruction", Action: ActionBlock, Description: "Harmful how-to phrasing"},
	{RuleType: "keyword", Pattern: "make a bomb", RiskType: "unsafe_instruction", Action: ActionBlock, Description: "Weapon construction phrasing"},
	{RuleType: "keyword", Pattern: "suicide method", RiskType: "self_harm", Action: ActionBlock, Description: "Self-harm instruction phrasing"},
}

func matchRule(normalized string, rules []Rule) *Rule {
	for _, rule := range rules {
		pattern := normalizeText(rule.Pattern)
		if pattern == "" {
			continue
		}
		switch strings.TrimSpace(rule.RuleType) {
		case "regex":
			re, err := regexp.Compile(pattern)
			if err == nil && re.MatchString(normalized) {
				matched := rule
				return &matched
			}
		default:
			if strings.Contains(normalized, pattern) {
				matched := rule
				return &matched
			}
		}
	}
	return nil
}

func sha256Hex(input string) string {
	sum := sha256.Sum256([]byte(input))
	return hex.EncodeToString(sum[:])
}
