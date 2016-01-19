package hub

import (
	"fmt"

	"github.com/pborman/uuid"
)

type Topic string

func (t Topic) String() string {
	return string(t)
}

// Req returns a subject were all requests of a
// specific topic will be routed.
func (t Topic) Req() string {
	return fmt.Sprintf("%s.%s", t.String(), "REQ")
}

// ResUnique returns a unique subject where a response to
// a Topic request will be routed.
func (t Topic) ResUnique() string {
	return fmt.Sprintf("%s.%s.%s", t.String(), "RES", uuid.New())
}

// ResWildward returns a wild card subject where all
// responses to a topic request will be routed.
func (t Topic) ResWildcard() string {
	return fmt.Sprintf("%s.%s.%s", t.String(), "RES", "*")
}
