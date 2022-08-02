package signalr

import (
	"bufio"
	"bytes"
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
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

const (
	chunkSize = 64 * 1024
)

// JSONReader -
type JSONReader struct {
	reader io.Reader
	buffer *bytes.Buffer
}

// NewJSONReader -
func NewJSONReader(r io.Reader) *JSONReader {
	return &JSONReader{
		reader: r,
		buffer: bytes.NewBuffer(make([]byte, 0)),
	}
}

// Scan -
func (r *JSONReader) Scan() error {
	data := make([]byte, chunkSize)
	for {
		count, err := r.reader.Read(data)
		switch err {
		case io.EOF:
			return nil
		case nil:
			if _, err := r.buffer.Write(data[:count]); err != nil {
				return err
			}
		default:
			return err
		}
	}
}

// Bytes -
func (r *JSONReader) Bytes() ([]byte, error) {
	if r.buffer.Len() == 0 {
		return nil, nil
	}

	data, err := r.buffer.ReadBytes(JSONSeparator)
	if err != nil {
		return nil, err
	}
	return data[:len(data)-1], nil
}
