package jira

import "strings"

// blockLevelTypes is the set of ADF node types that produce a trailing
// newline when flattened to plain text. This covers structural blocks
// defined by the Atlassian Document Format specification.
var blockLevelTypes = map[string]bool{
	"paragraph":    true,
	"heading":      true,
	"bulletList":   true,
	"orderedList":  true,
	"listItem":     true,
	"blockquote":   true,
	"codeBlock":    true,
	"rule":         true,
	"mediaSingle":  true,
	"mediaGroup":   true,
	"table":        true,
	"tableRow":     true,
	"tableCell":    true,
	"tableHeader":  true,
	"panel":        true,
	"decisionList": true,
	"decisionItem": true,
	"taskList":     true,
	"taskItem":     true,
}

// flattenADF recursively extracts plain text from an Atlassian
// Document Format node tree. The node parameter is the result of
// unmarshaling ADF JSON into any via [encoding/json]. Block-level
// nodes receive a trailing newline; text nodes yield their text value.
// Nil or non-map input returns an empty string. Trailing whitespace
// is trimmed from the final result.
func flattenADF(node any) string {
	return strings.TrimRight(flattenADFNode(node), "\n ")
}

func flattenADFNode(node any) string {
	m, ok := node.(map[string]any)
	if !ok || m == nil {
		return ""
	}

	nodeType, _ := m["type"].(string)

	if nodeType == "text" {
		text, _ := m["text"].(string)
		return text
	}

	var result string
	if content, ok := m["content"].([]any); ok {
		for _, child := range content {
			result += flattenADFNode(child)
		}
	}

	if blockLevelTypes[nodeType] {
		result += "\n"
	}

	return result
}
