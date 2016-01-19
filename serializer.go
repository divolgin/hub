package hub

import (
	"encoding/json"
)

type SerializationFormat string

const (
	JSON   SerializationFormat = "json"
	BINARY SerializationFormat = "binary"
)

func (sf SerializationFormat) GetSerializer() BusSerializer {
	switch sf {
	case JSON:
		return &JSONSerializer{}
	default:
		return &JSONSerializer{}
	}
}

type BusSerializer interface {
	Serialize(obj interface{}) ([]byte, error)
	Deserialize(data []byte, obj interface{}) error
}

type JSONSerializer struct{}

func (s *JSONSerializer) Serialize(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

func (s *JSONSerializer) Deserialize(data []byte, obj interface{}) error {
	return json.Unmarshal(data, obj)
}
