package handler

import (
	"fmt"
	"sync"

	"github.com/PlayerNeo42/gvalkey/resp"
)

type Command struct {
	// name of the command
	Name resp.BulkString

	// number of arguments
	// positive value means fixed number of arguments
	// negative value means at least that number of arguments
	Args int

	// handler of the command
	Handler func(command resp.Array) (resp.Marshaler, error)
}

type CommandTable struct {
	m sync.Map
}

func NewCommandTable() *CommandTable {
	return &CommandTable{}
}

func (c *CommandTable) MustRegister(command *Command) {
	_, ok := c.m.LoadOrStore(command.Name.Upper(), command)
	if ok {
		panic(fmt.Sprintf("command %s already registered", command.Name.Upper()))
	}
}

func (c *CommandTable) Get(name resp.BulkString) (*Command, bool) {
	val, ok := c.m.Load(name.Upper())
	if !ok {
		return nil, false
	}
	return val.(*Command), true
}
