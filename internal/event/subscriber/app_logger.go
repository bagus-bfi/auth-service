package subscriber

import (
	"bravo-go-template/internal/event"

	"github.com/bfi-finance/bfi-go-pkg/eventbus"
	"github.com/rs/zerolog"
)

type Bar struct {
	zlogger zerolog.Logger
}

// NewAppLogger is a sample of internal events handler to log app events.
func NewAppLogger(zlogger zerolog.Logger) *Bar {
	return &Bar{
		zlogger: zlogger,
	}
}

func (s *Bar) Name() string {
	return "internal/event/subscriber/app_logger" // or name it accordingly.
}

func (s *Bar) SubscribedEvents() []string {
	return []string{
		event.AppInternalError,
	}
}

func (s *Bar) Handle(e *eventbus.Event) {
	switch e.Name {
	case event.AppInternalError:
		if err, ok := e.Data.(*event.InternalError); ok {
			s.zlogger.Error().
				Err(err.Err).
				Str("sub", s.Name()).
				Str("caller", err.Caller).
				Str("request_id", err.RequestID).
				Interface("request_data", err.ReqData).
				Msg(err.Err.Error()) // original error as message
			return
		}
		if err, ok := e.Data.(error); ok {
			s.zlogger.Error().Str("sub", s.Name()).Msg(err.Error())
		}
	default:
	}
}
