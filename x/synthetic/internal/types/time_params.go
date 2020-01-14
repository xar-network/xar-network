package types

import "time"

type IntervalTimer struct {
	TimeStart      time.Time     `json:"time_start" yaml:"time_start"`
	Interval       time.Duration `json:"time_interval" yaml:"time_interval"`
	DeadLine       time.Time     `json:"dead_line" yaml:"dead_line"` // we can either Recalculate deadline with interval every time we need, or we can save a deadline once
	IsScheduling   bool          `json:"is_scheduling" yaml:"is_scheduling"`
	StopScheduling bool          `json:"stop_scheduling" yaml:"stop_scheduling"`
}

func TimerFromInterval(interval time.Duration) IntervalTimer {
	return IntervalTimer{
		time.Now(),
		interval,
		time.Now().Add(interval),
		false,
		false,
	}
}

func NewTimeParams(start time.Time, interval time.Duration, deadLine time.Time) IntervalTimer {
	return IntervalTimer{
		start,
		interval,
		deadLine,
		false,
		false,
	}
}

func (t *IntervalTimer) Schedule(f func()) {
	go t.executeAtDeadline(f)
}

func (t *IntervalTimer) StopSchedule() {
	t.StopScheduling = true
}

func (t *IntervalTimer) executeAtDeadline(f func()) {
	t.IsScheduling = true
	for {
		if t.StopScheduling {
			t.StopScheduling = false
			t.IsScheduling = false
			return
		}

		waitDuration := t.DeadLine.Sub(time.Now())
		time.Sleep(waitDuration)
		f()
		t.Reset()
	}
}

func (t *IntervalTimer) Reset() {
	t.TimeStart = time.Now()
	t.DeadLine = time.Now().Add(t.Interval)
}

func (t *IntervalTimer) IsExpired(now time.Time) bool {
	if !now.After(t.TimeStart) {
		panic("a time passed is out of lower bound")
	}

	return now.Before(t.DeadLine)
}

func (t *IntervalTimer) IntervalIsZero() bool {
	return t.Interval == time.Duration(0)
}
