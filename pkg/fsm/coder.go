package fsm

import (
	"encoding/base64"
	"encoding/json"
)

// Encodes and Decodes a Match to a byte array:
// a state that can resume a single match

// Encodes a Match to a byte array that the frontend can save
func (m *Match) Encode() []byte {
	m.CStateId = m.CState.stateId()
	data, err := json.Marshal(m)
	if err != nil {
		return []byte{}
	}
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(encoded, data)
	return encoded
}

// Decodes a byte array match from the frontend
func Decode(encoded []byte) *Match {
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(encoded)))
	n, err := base64.StdEncoding.Decode(decoded, encoded)
	if err != nil {
		return NewMatch()
	}

	m := &Match{}
	err = json.Unmarshal(decoded[:n], m)
	if err != nil {
		return NewMatch()
	}

	m.bindStates()
	switch m.CStateId {
	case 1:
		m.CState = m.Playing
	case 2:
		m.CState = m.Announcing
	case 3:
		m.CState = m.Responding
	case 0:
		m.CState = m.End
	default:
		m.CState = m.Playing
	}

	return m
}
