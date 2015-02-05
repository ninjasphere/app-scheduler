package controller

import (
	"time"
)

type mockclock struct {
	systemclock
	now       time.Time
	jitter    time.Duration
	callbacks map[time.Time][]callback
}

func initMockClock(now time.Time, jitter time.Duration) *mockclock {
	mock := &mockclock{
		now:       now.Add(jitter),
		jitter:    jitter,
		callbacks: make(map[time.Time][]callback),
	}
	clock = mock
	return mock
}

func (m *mockclock) Now() time.Time {
	return m.now.Truncate(defaultRounding)
}

func (m *mockclock) SetNow(now time.Time) {
	m.now = now.Add(m.jitter)
	saved := make(map[time.Time][]callback)
	for t, cbs := range m.callbacks {
		if t.Sub(now) <= 0 {
			for _, cb := range cbs {
				log.Debugf("firing event for %v because time is now %v", t, now)
				cb()
			}
		} else {
			log.Debugf("keeping event for %v because time is now %v", t, now)
			saved[t] = cbs
		}
	}
	m.callbacks = saved
}

func (m *mockclock) AfterFunc(delay time.Duration, then func()) {
	if delay <= 0 {
		log.Debugf("event for %d s fired now %v", delay/time.Second, m.now)
		then()
	} else {
		when := m.now.Add(delay)
		log.Debugf("at %s event for %d s scheduled for %v", m.now, delay/time.Second, when)
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
