package hub

import (
	"fmt"

	"github.com/pborman/uuid"
)

type Topic string
type TopicRequest string
type TopicResponse string

func (t Topic) String() string {
	return string(t)
}

func (req TopicRequest) String() string {
	return string(req)
}

func (res TopicResponse) String() string {
	return string(res)
}

// Req returns the request version of a topic.
func (t Topic) Req() TopicRequest {
	return TopicRequest(fmt.Sprintf("%s.%s", t.String(), "REQ"))
}

// Res returns the response version of a topic.
func (t Topic) Res() TopicResponse {
	return TopicResponse(fmt.Sprintf("%s.%s", t.String(), "RES"))
}

// ResUnique returns a unique subject where a response to
// a Topic request will be routed.
func (t Topic) ResUnique() TopicResponse {
	return TopicResponse(fmt.Sprintf("%s.%s.%s", t.String(), "RES", uuid.New()))
}

// ResWildward returns a wild card subject where all
// responses to a topic request will be routed.
func (t Topic) ResWildcard() TopicResponse {
	return TopicResponse(fmt.Sprintf("%s.%s.%s", t.String(), "RES", "*"))
}
