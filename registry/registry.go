package registry

import "asd/handlers"

// Registry holds the registered file handlers.
type Registry struct {
	handlers []handlers.Handler
	fallback handlers.Handler
}

// New returns a new Registry with the provided fallback handler.
func New(fallback handlers.Handler) *Registry {
	return &Registry{
		handlers: make([]handlers.Handler, 0),
		fallback: fallback,
	}
}

// Register adds a handler to the registry.
func (r *Registry) Register(h handlers.Handler) {
	r.handlers = append(r.handlers, h)
}

// Dispatch finds the first handler that can handle the given mime and extension.
func (r *Registry) Dispatch(mime, ext string) handlers.Handler {
	for _, h := range r.handlers {
		if h.CanHandle(mime, ext) {
			return h
		}
	}
	return r.fallback
}
