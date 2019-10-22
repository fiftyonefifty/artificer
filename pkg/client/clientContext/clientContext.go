package clientContext

import (
	"artificer/pkg/client/models"
	"context"
)

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type key int

// userIPkey is the context key for the user IP address.  Its value of zero is
// arbitrary.  If this package defined other context keys, they would have
// different integer values.
const clientKey key = 0

// NewContext returns a new Context carrying userIP.
func NewContext(ctx context.Context, client models.Client) context.Context {
	return context.WithValue(ctx, clientKey, client)
}

// FromContext extracts the user IP address from ctx, if present.
func FromContext(ctx context.Context) (models.Client, bool) {
	// ctx.Value returns nil if ctx has no value for the key;
	// the net.IP type assertion returns ok=false for nil.
	client, ok := ctx.Value(clientKey).(models.Client)
	return client, ok
}
