package controller

import (
	"time"
)

type mockclock struct {
	now       time.Time
	callbacks map[time.Time][]callback
}

func initMockClock(now time.Time) *mockclock {
	mock := &mockclock{
		now:       now,
		callbacks: make(map[time.Time][]callback),
	}
	clock = mock
	return mock
}

func (m *mockclock) Now() time.Time {
	return m.now
}

func (m *mockclock) SetNow(now time.Time) {
	m.now = now
	saved := make(map[time.Time][]callback)
	for t, cbs := range m.callbacks {
		if t.Sub(now) >= 0 {
			for _, cb := range cbs {
				cb()
			}
		} else {
			saved[t] = cbs
		}
	}
	m.callbacks = saved
}

func (m *mockclock) AfterFunc(delay time.Duration, then func()) {
	if delay <= 0 {
		then()
	} else {
		when := m.now.Add(delay)
		var list []callback
		if list, ok := m.callbacks[when]; !ok {
			list = make([]callback, 0)
			m.callbacks[when] = list
		}
		m.callbacks[when] = append(list, then)
	}
}

func (m *mockclock) Location() *time.Location {
	return m.now.Location()
}
