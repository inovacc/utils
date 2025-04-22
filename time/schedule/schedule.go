package schedule

import (
	"context"
	"errors"
	"log/slog"

	cron2 "github.com/inovacc/utils/v2/time/schedule/cron"
)

/*
Entry                  | Description                                  | Equivalent To
-----                  | -----------                                  | -------------
@yearly (or @annually) | Run once a year, midnight, Jan. 1st          | 0 0 0 1 1 *
@monthly               | Run once a month, midnight, first of month   | 0 0 0 1 * *
@weekly                | Run once a week, midnight between Sat/Sun    | 0 0 0 * * 0
@daily (or @midnight)  | Run once a day, midnight                     | 0 0 0 * * *
@hourly                | Run once an hour, beginning of hour          | 0 0 * * * *
@every <duration>      | Run every <duration>, starting now           | 0 0 0 0 * *
@weekday               | Run every weekday, midnight between Mon/Fri  | 0 0 0 * * 1-5
*/

const (
	Yearly    = "@yearly"
	Annually  = "@annually"
	Monthly   = "@monthly"
	Weekly    = "@weekly"
	Daily     = "@daily"
	Midnight  = "@midnight"
	Hourly    = "@hourly"
	Minute    = "@minute"
	Every     = "@every"
	Weekday   = "@weekday"
	Monday    = "1"
	Tuesday   = "2"
	Wednesday = "3"
	Thursday  = "4"
	Friday    = "5"
	Saturday  = "6"
	Sunday    = "0"
)

type printfLogger struct{}

func (pl printfLogger) Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func (pl printfLogger) Error(err error, msg string, args ...any) {
	slog.Error(msg, append([]any{"err", err}, args...)...)
}

type Cron struct {
	cron *cron2.Cron
	ctx  context.Context
}

func NewCronScheduler(ctx context.Context) (*Cron, error) {
	if ctx == nil {
		return nil, errors.New("context cannot be nil")
	}

	c := cron2.New(cron2.WithSeconds(), cron2.WithLogger(printfLogger{}))
	c.Start()

	go func() {
		<-ctx.Done()
		c.Stop()
	}()

	return &Cron{
		cron: c,
		ctx:  ctx,
	}, nil
}

func (c *Cron) AddFunc(spec string, cmd func()) (int, error) {
	id, err := c.cron.AddFunc(c.fixWeekday(spec), cmd)
	return int(id), err
}

func (c *Cron) fixWeekday(spec string) string {
	switch spec {
	case Weekday:
		return "0 0 0 * * 1-5"
	case Monday:
		return "0 0 0 * * 1"
	case Tuesday:
		return "0 0 0 * * 2"
	case Wednesday:
		return "0 0 0 * * 3"
	case Thursday:
		return "0 0 0 * * 4"
	case Friday:
		return "0 0 0 * * 5"
	case Saturday:
		return "0 0 0 * * 6"
	case Sunday:
		return "0 0 0 * * 0"
	case Minute:
		return "0 * * * * *"
	default:
		return spec
	}
}
