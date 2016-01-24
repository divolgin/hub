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
	if t.IsAlreadyReqRes() {
		panic(fmt.Errorf("Unable to transform topic to .REQ  > %s", t.String()))
	}
	return Topic(fmt.Sprintf("%s.%s", t.String(), "REQ"))
}

// Res returns the response version of a topic.
func (t Topic) Res() Topic {
	if t.IsAlreadyReqRes() {
		panic(fmt.Errorf("Unable to transform topic to .RES  > %s", t.String()))
	}
	return Topic(fmt.Sprintf("%s.%s", t.String(), "RES"))
}

// ResUnique returns a unique subject where a response to
// a Topic request will be routed.
func (t Topic) ResUnique() Topic {
	if t.IsAlreadyReqRes() {
		panic(fmt.Errorf("Unable to transform topic to .RES.UUID > %s", t.String()))
	}
	return Topic(fmt.Sprintf("%s.%s.%s", t.String(), "RES", uuid.New()))
}

// ResWildward returns a wild card subject where all
// responses to a topic request will be routed.
func (t Topic) ResWildcard() Topic {
	if t.IsAlreadyReqRes() {
		panic(fmt.Errorf("Unable to transform topic to .RES.* > %s", t.String()))
	}
	return Topic(fmt.Sprintf("%s.%s.%s", t.String(), "RES", "*"))
}

func (t Topic) IsAlreadyReqRes() bool {
	if strings.Contains(string(t), ".REQ") || strings.Contains(string(t), ".RES") {
		return true
	}
	return false
}
