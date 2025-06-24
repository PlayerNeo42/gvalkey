package resp

type Null struct{}

func (n Null) MarshalRESP() []byte {
	return []byte("$-1\r\n")
}
