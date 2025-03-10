package lec

import (
	"context"
	"sync"
	"time"

	"github.com/gosuit/sl"
)

type Context interface {
	Logger() sl.Logger
	SlHandler() sl.Handler
	AddValue(key string, val interface{}, forLog bool)
	GetValue(key string) *Value
	GetValues() map[string]Value
	AddErr(err error)
	GetErr() error
	HasErr() bool

	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key any) any
}

type ctx struct {
	log    sl.Logger
	data   map[string]Value
	errors []error
	mu     sync.Mutex
	base   context.Context
}

type Value struct {
	Val   interface{}
	Share bool
}

func New(log sl.Logger) Context {
	return &ctx{
		log:  log,
		data: make(map[string]Value),
		base: context.TODO(),
	}
}

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

func (c *ctx) AddValue(key string, val interface{}, share bool) {
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
