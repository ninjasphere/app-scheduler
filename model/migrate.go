package model

import "fmt"

func (m *Schedule) Migrate() *Schedule {
	for i, t := range m.Tasks {
		m.Tasks[i] = t.Migrate()
	}
	return m
}

func (m *Task) Migrate() *Task {
	for i, a := range m.Open {
		m.Open[i] = a.Migrate()
	}
	for i, a := range m.Close {
		m.Close[i] = a.Migrate()
	}
	return m
}

func (m *Action) Migrate() *Action {
	if m.ThingID != "" {
		m.SubjectID = fmt.Sprintf("thing:%s", m.ThingID)
		m.ThingID = ""
	}
	return m
}
