package i18n

import "context"

type contextKeyType struct{}

var contextKey = contextKeyType{}

// LocalizerFromContext raises request-scoped localizer.
// Warning: localizer will be <nil> if it was not set
// using [ContextWithLocalizer].
func LocalizerFromContext(ctx context.Context) *Localizer {
	return ctx.Value(contextKey).(*Localizer)
}

// ContextWithLocalizer adds localizer into context as a value.
// Use [LocalizerFromContext] to recover it later.
func ContextWithLocalizer(parent context.Context, l *Localizer) context.Context {
	return context.WithValue(parent, contextKey, l)
}
