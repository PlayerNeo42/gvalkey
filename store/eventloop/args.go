package eventloop

const (
	CmdGet = iota
	CmdSet
	CmdDel
)

type cmd struct {
	typ     int
	payload any
	resp    any
}

type operationResult struct {
	Value any
	OK    bool
}
