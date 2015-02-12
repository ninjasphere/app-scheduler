package model

import (
	"strings"
)

func (m *Action) GetThingID() string {
	if x := strings.Index(m.SubjectID, "thing:"); x >= 0 {
		return m.SubjectID[x+len("thing:"):]
	} else {
		return ""
	}
}

func (m *Action) GetSceneID() string {
	if x := strings.Index(m.SubjectID, "scene:"); x >= 0 {
		return m.SubjectID[x+len("scene:"):]
	} else {
		return ""
	}
}
