package event

import (
	"context"
	"fmt"
	"runtime"

	"github.com/bfi-finance/bfi-go-pkg/eventbus"
	"github.com/bfi-finance/bfi-go-pkg/tracer"
)

// List of event name constants that is happens inside the application.
const (
	// AppInternalError is an event for internal error that the details is not returned to the client.
	AppInternalError = "app:internal_error"
)

type option struct {
	requestData any
}

type Option func(*option)

func WithRequestData(data any) Option {
	return func(o *option) { o.requestData = data }
}

type InternalError struct {
	Err       error
	ReqData   any
	Caller    string
	RequestID string // for error tracking
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("[%s]: %s", e.Caller, e.Err.Error())
}

// Following methods are event construct helpers.

// NewAppInternalError is event to publish an internal/unknown error.
// It also include caller info (full method name and line) in the event.
func NewAppInternalError(ctx context.Context, err error, opts ...Option) *eventbus.Event {
	opt := &option{}
	for _, o := range opts {
		o(opt)
	}
	pc, _, line, ok := runtime.Caller(1)
	caller := runtime.FuncForPC(pc)
	if ok && caller != nil {
		info := &InternalError{
			Caller: fmt.Sprintf("%s:%d", caller.Name(), line),
			Err:    err,
		}
		if reqID := tracer.RequestIDFromContext(ctx); reqID != "" {
			info.RequestID = reqID
		}
		if opt.requestData != nil {
			info.ReqData = opt.requestData
		}
		err = info
	}
	return &eventbus.Event{Name: AppInternalError, Data: err}
}
