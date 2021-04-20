package signalr

// Encoding -
type Encoding interface {
	Decode(data []byte) (interface{}, error)
	Encode(msg interface{}) ([]byte, error)
}
