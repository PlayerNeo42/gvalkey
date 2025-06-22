package resp

var (
	OK   = []byte("+OK\r\n")
	NULL = []byte("$-1\r\n")
	ERR  = []byte("-ERR\r\n")
)
