package template

// IdentityParser is an Parser that does no parsing and returns template string unchanged.
type IdentityParser struct{}

func (IdentityParser) Cacheable() bool {
	// Caching is not necessary because Parse is cheap.
	return false
}

func (IdentityParser) Parse(src, leftDelim, rightDelim string) (ParsedTemplate, error) {
	return &identityParsedTemplate{src: src}, nil
}

type identityParsedTemplate struct {
	src string
}

func (t *identityParsedTemplate) Execute(data any) (string, error) {
	return t.src, nil
}
