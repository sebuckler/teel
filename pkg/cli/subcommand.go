package cli

import "fmt"

type subcommand struct {
	args        []string
	name        string
	parent      *subcommand
	subcommands []*subcommand
}

var parsed []*subcommand

func New(n string, s []*subcommand) *subcommand {
	return &subcommand{
		args:        []string{},
		name:        n,
		subcommands: s,
	}
}

func (s *subcommand) Parse(a []string) {
	if len(a) == 0 {
		return
	}

	for _, arg := range a {
		found := false

		for _, sub := range s.subcommands {
			found = sub.parse(arg)

			if found {
				break
			}
		}

		if found {
			continue
		}

		if len(parsed) == 0 {
			s.args = append(s.args, arg)
			parsed = append(parsed, s)

			continue
		}

		parsed[len(parsed)-1].args = append(parsed[len(parsed)-1].args, arg)
	}

	for _, cmd := range parsed {
		fmt.Println(cmd)
	}
}

func (s *subcommand) parse(a string) bool {
	if a == s.name {
		parsed = append(parsed, s)

		return true
	}

	for _, sub := range s.subcommands {
		if sub.parse(a) {
			return true
		}
	}

	return false
}
