package commands

import (
	"fmt"
	"strings"
)

type Command struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	Args            string `json:"-"`
	ArgsPlaceholder string `json:"args_placeholder"`
	Flags           []Flag `json:"flags"`
}

type Flag struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Default     string `json:"default"`
	Value       string `json:"-"`
}

var (
	list []Command
)

func Import(l []Command) {
	list = l
}

// List returns a list of commands.
func List() []Command {
	return list
}

// Map returns commands as a map.
func Map() map[string]Command {
	m := make(map[string]Command, len(list))
	for _, cmd := range list {
		m[cmd.Name] = cmd
	}
	return m
}

func (c Command) CombinedArgs() []string {
	args := []string{c.Name}
	args = append(args, c.ArgsSlice()...)
	args = append(args, c.FlagsSlice()...)
	return args
}

func (c Command) ArgsSlice() []string {
	if c.Args == "" {
		return []string{}
	}
	return strings.Split(c.Args, " ")
}

func (c Command) FlagsSlice() []string {
	flags := []string{}
	for _, f := range c.Flags {
		flags = append(flags, f.encode())
	}
	return flags
}

func (c Command) FlagsString() string {
	return strings.Join(c.FlagsSlice(), " ")
}

func (c Command) String() string {
	return strings.Join(c.CombinedArgs(), " ")
}

func (f Flag) encode() string {
	return fmt.Sprintf("--%s=%s", f.Name, f.Value)
}
