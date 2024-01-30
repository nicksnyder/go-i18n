// Package template defines a generic interface for template parsers and implementations of that interface.
package template

// Parser parses strings into executable templates.
type Parser interface {
	// Parse parses src and returns a ParsedTemplate.
	Parse(src, leftDelim, rightDelim string) (ParsedTemplate, error)

	// Cacheable returns true if Parse returns ParsedTemplates that are always safe to cache.
	Cacheable() bool
}

// ParsedTemplate is an executable template.
type ParsedTemplate interface {
	// Execute applies a parsed template to the specified data.
	Execute(data any) (string, error)
}
