package i18n

import (
	"context"
	"testing"

	"golang.org/x/text/language"
)

func TestContextAwareness(t *testing.T) {
	bundle := NewBundle(language.English)
	localizer := NewLocalizer(bundle, "en")
	ctx := ContextWithLocalizer(context.Background(), localizer)
	if ctx == nil {
		t.Error("<nil> context")
	}

	recovered, ok := LocalizerFromContext(ctx)
	if !ok {
		t.Error("localizer was not recovered")
	}
	if recovered == nil {
		t.Error("recovered a <nil> localizer from context")
	}
}
