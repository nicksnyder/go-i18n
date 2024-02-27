package i18n

import "context"

type contextKeyType struct{}

var contextKey = contextKeyType{}

// LocalizerFromContext raises request-scoped localizer from context.
// Returns `<nil>, false` if there is no localizer in context.
func LocalizerFromContext(ctx context.Context) (l *Localizer, ok bool) {
	l, ok = ctx.Value(contextKey).(*Localizer)
	return
}

// ContextWithLocalizer adds localizer into context as a value.
// Use [LocalizerFromContext] to recover it later.
func ContextWithLocalizer(parent context.Context, l *Localizer) context.Context {
	return context.WithValue(parent, contextKey, l)
}
