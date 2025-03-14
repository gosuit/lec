package lec

import (
	"context"
	"sync"
	"time"

	"github.com/gosuit/sl"
)

// Context is an interface that defines methods for logging, error handling,
// and storing key-value pairs.
type Context interface {
	// Logger returns the logger associated with this context.
	Logger() sl.Logger

	// SlHandler returns the handler associated with the logger.
	SlHandler() sl.Handler

	// AddValue adds a key-value pair to the context. The 'forLog' parameter indicates
	// whether the value should be added to logger attributes.
	AddValue(key string, val any, forLog bool)

	// GetValue retrieves the value associated with the given key. Returns nil if the key does not exist.
	GetValue(key string) *Value

	// GetValues returns all key-value pairs stored in the context.
	GetValues() map[string]Value

	// AddErr adds an error to the context's error list.
	AddErr(err error)

	// GetErr retrieves and removes the first error from the context's error list.
	// Returns nil if there are no errors.
	GetErr() error

	// HasErr checks if there are any errors stored in the context.
	HasErr() bool

	// Deadline returns the deadline associated with the context, if any.
	// Implementation of context.Context.
	Deadline() (deadline time.Time, ok bool)

	// Done returns a channel that is closed when the context is done.
	// Implementation of context.Context.
	Done() <-chan struct{}

	// Err returns an error that was encountered during the context's lifetime, if any.
	// Implementation of context.Context.
	Err() error

	// Value retrieves a value from the base context associated with the given key.
	// Implementation of context.Context.
	Value(key any) any
}

type ctx struct {
	log    sl.Logger
	data   map[string]Value
	errors []error
	mu     sync.Mutex
	base   context.Context
}

// Value represents a key-value pair that can optionally be shared (logged).
type Value struct {
	Val   any  // The actual value
	Share bool // Indicates if the value should be shared (logged).
}

// New creates a new Context with a logger.
func New(log sl.Logger) Context {
	return &ctx{
		log:  log,
		data: make(map[string]Value),
		base: context.TODO(),
	}
}

// NewWithCtx creates a new Context with a specified base context and logger.
func NewWithCtx(base context.Context, log sl.Logger) Context {
	return &ctx{
		log:  log,
		data: make(map[string]Value),
		base: base,
	}
}

func (c *ctx) GetValue(key string) *Value {
	c.mu.Lock()
	defer c.mu.Unlock()

	value := c.data[key]
	if value.Val == nil {
		return nil
	}

	return &value
}

func (c *ctx) GetValues() map[string]Value {
	return c.data
}

func (c *ctx) AddValue(key string, val any, share bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = Value{
		val, share,
	}

	if share {
		c.log = c.log.With(key, val)
	}
}

func (c *ctx) Logger() sl.Logger {
	return c.log
}

func (c *ctx) SlHandler() sl.Handler {
	return c.log.Handler()
}

func (c *ctx) AddErr(err error) {
	c.mu.Lock()
	c.errors = append(c.errors, err)
	c.mu.Unlock()
}

func (c *ctx) GetErr() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.errors) == 0 {
		return nil
	}

	err := c.errors[0]

	if len(c.errors) > 1 {
		c.errors = c.errors[1:]
	} else {
		c.errors = make([]error, 0)
	}

	return err
}

func (c *ctx) HasErr() bool {
	c.mu.Lock()
	errCount := len(c.errors)
	c.mu.Unlock()

	return errCount > 0
}

func (c *ctx) Deadline() (deadline time.Time, ok bool) {
	return c.base.Deadline()
}

func (c *ctx) Done() <-chan struct{} {
	return c.base.Done()
}

func (c *ctx) Err() error {
	return c.base.Err()
}

func (c *ctx) Value(key any) any {
	return c.base.Value(key)
}
