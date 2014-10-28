package model

import (
	"encoding/json"
	"testing"
)

func assert(t *testing.T, description string, assertion func() bool) {
	if !assertion() {
		t.Fatalf("assertion failed: %s\n", description)
	}
}

func TestJSONRoundTrip(t *testing.T) {
	item := &Task{
		Uuid:        "u",
		Description: "d",
		Window: &Window{
			From:  "10:00",
			Until: "12:00",
		},
		Open: []*ThingAction{
			&ThingAction{
				ActionType: "thing-action",
				ThingId:    "thing-id",
				Action:     "on",
			},
		},
		Close: []*ThingAction{
			&ThingAction{
				ActionType: "thing-action",
				ThingId:    "thing-id",
				Action:     "off",
			},
		},
	}

	serialized, err := json.Marshal(item)
	if err != nil {
	   t.Fatalf("marhsalling - %s", err)
	}
	deserialized := &Task{}
	err =json.Unmarshal(serialized, deserialized)
	if err != nil {
	   t.Fatalf("unmarhsalling - %s", err)
	}

	assert(t, "uuid", func() bool { return deserialized.Uuid == item.Uuid })
	assert(t, "Open[0].Action", func() bool { return deserialized.Open[0].Action == item.Open[0].Action })
}
