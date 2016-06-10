package hub

import (
	"fmt"
	"strings"

	"github.com/pborman/uuid"
)

type Topic string
type TopicRequest string
type TopicResponse string

func (t Topic) String() string {
	return string(t)
}

// Req returns the request version of a topic.
func (t Topic) Req() Topic {
	t.IsModified()
	return Topic(fmt.Sprintf("%s.%s", t.String(), "REQ"))
}

// Res returns the response version of a topic.
func (t Topic) Res() Topic {
	t.IsModified()
	return Topic(fmt.Sprintf("%s.%s", t.String(), "RES"))
}

// ResUnique returns a unique subject where a response to
// a Topic request will be routed.
func (t Topic) ResUnique() Topic {
	t.IsModified()
	return Topic(fmt.Sprintf("%s.%s.%s", t.String(), "RES", uuid.New()))
}

// ResWildward returns a wild card subject where all
// responses to a topic request will be routed.
func (t Topic) ResWildcard() Topic {
	t.IsModified()
	return Topic(fmt.Sprintf("%s.%s.%s", t.String(), "RES", "*"))
}

func (t Topic) Stream() Topic {
	t.IsModified()
	return Topic(fmt.Sprintf("%s.%s.%s", t.String(), "STREAM", uuid.New()))
}

func (t Topic) Heartbeat() Topic {
	t.IsModified()
	return Topic(fmt.Sprintf("%s.%s.%s", t.String(), "HEARTBEAT", uuid.New()))
}

func (t Topic) IsModified() {
	if strings.Contains(string(t), ".REQ") ||
		strings.Contains(string(t), ".RES") {
		panic(fmt.Errorf("Unable to transform topic > %s", t.String()))
	}
}
