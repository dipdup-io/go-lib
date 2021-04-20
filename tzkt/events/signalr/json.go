package signalr

import (
	"bufio"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Separator
const (
	JSONSeparator = 0x1e
)

// JSONEncoding -
type JSONEncoding struct {
}

// NewJSONEncoding -
func NewJSONEncoding() *JSONEncoding {
	return &JSONEncoding{}
}

// Decode -
func (e *JSONEncoding) Decode(data []byte) (interface{}, error) {
	log.Trace(string(data))

	var typ Type
	if err := json.Unmarshal(data, &typ); err != nil {
		return nil, err
	}

	switch typ.Type {
	case MessageTypeInvocation:
		var msg Invocation
		err := json.Unmarshal(data, &msg)
		return msg, err
	case MessageTypeStreamItem:
		var msg StreamItem
		err := json.Unmarshal(data, &msg)
		return msg, err
	case MessageTypeCompletion:
		var msg Completion
		err := json.Unmarshal(data, &msg)
		return msg, err
	case MessageTypeStreamInvocation:
		var msg StreamInvocation
		err := json.Unmarshal(data, &msg)
		return msg, err
	case MessageTypeCancelInvocation:
		var msg CancelInvocation
		err := json.Unmarshal(data, &msg)
		return msg, err
	case MessageTypePing:
		var msg PingMessage
		err := json.Unmarshal(data, &msg)
		return msg, err
	case MessageTypeCloseMessage:
		var msg CloseMessage
		err := json.Unmarshal(data, &msg)
		return msg, err
	default:
		return nil, errors.Wrapf(ErrUnknownMessageType, "%d", typ.Type)
	}

}

// Encode -
func (e *JSONEncoding) Encode(msg interface{}) ([]byte, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return append(data, JSONSeparator), nil
}

// SplitJSON -
func SplitJSON(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i := 0; i < len(data); i++ {
		if data[i] == JSONSeparator {
			return i + 1, data[:i], nil
		}
	}
	if !atEOF {
		return 0, nil, nil
	}
	return 0, data, bufio.ErrFinalToken
}
