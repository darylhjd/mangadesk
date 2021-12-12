package utils

import "context"

// ContextWrapper : A wrapper around a context and its corresponding Cancel function.
type ContextWrapper struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

// ResetContext : Sets the new context, and return the old context,
// which the caller can use to cancel the previous context.
func (c *ContextWrapper) ResetContext() (context.Context, context.CancelFunc) {
	ctx, cancel := c.Ctx, c.Cancel
	c.Ctx, c.Cancel = context.WithCancel(context.Background())
	return ctx, cancel
}

// ToCancel : Helper function to check if the previous context should be cancelled.
// Note that this does not check the current context stored in the wrapper.
func (c *ContextWrapper) ToCancel(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
