// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

// package monkit is not a real monkit package. it's a reimplementation that
// does nothing, to avoid apache v2 vs gpl v2 licensing incompatibility. it also
// means there are no metrics or monkit collection
package monkit

import "context"

// Scope
type Scope struct{}

// Package
func Package() *Scope { return &Scope{} }

// ScopeNamed
func ScopeNamed(string) *Scope { return &Scope{} }

// Event
func (s *Scope) Event(string) {}

// Task
func (s *Scope) Task() func(*context.Context, ...interface{}) func(*error) {
	return func(*context.Context, ...interface{}) func(*error) {
		return func(*error) {}
	}
}

// TaskNamed
func (s *Scope) TaskNamed(string) func(*context.Context) func(*error) {
	return func(*context.Context) func(*error) {
		return func(*error) {}
	}
}

// Meter
type Meter struct{}

// Mark
func (m *Meter) Mark(int64) {}

// Meter
func (s *Scope) Meter(string) *Meter { return &Meter{} }

// IntVal
type IntVal struct{}

// Observe
func (v *IntVal) Observe(int64) {}

// IntVal
func (s *Scope) IntVal(string) *IntVal { return &IntVal{} }

// Func
type Func struct{}

// FuncNamed
func (s *Scope) FuncNamed(string) *Func { return &Func{} }

// Func
func (s *Scope) Func() *Func { return &Func{} }

// RestartTrace
func (f *Func) RestartTrace(*context.Context) func(*error) {
	return func(*error) {}
}

// RemoteTrace
func (f *Func) RemoteTrace(*context.Context, int64, *Trace) func(*error) {
	return func(*error) {}
}

// Trace
type Trace struct{}

// NewTrace
func NewTrace(int64) *Trace { return &Trace{} }

// Get
func (t *Trace) Get(interface{}) interface{} { return nil }

// Id
func (t *Trace) Id() int64 { return 0 }

// NewId
func NewId() int64 { return 0 }

// Span
type Span struct{}

// SpanFromCtx
func SpanFromCtx(context.Context) *Span {
	return &Span{}
}

// Parent
func (s *Span) Parent() *Span { return &Span{} }

// Trace
func (s *Span) Trace() *Trace { return &Trace{} }

// Id
func (s *Span) Id() int64 { return 0 }

// Registry
type Registry struct{}

// Default
var Default = &Registry{}

// Stats
func (r *Registry) Stats(cb func(SeriesKey, string, float64)) {}

// SeriesKey
type SeriesKey struct{}

// WithField
func (k *SeriesKey) WithField(string) string { return "" }

// ResetContextSpan
func ResetContextSpan(ctx context.Context) context.Context { return ctx }
