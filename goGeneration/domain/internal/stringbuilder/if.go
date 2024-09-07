package stringbuilder

import "fmt"

type If struct {
	Condition string
	Then      *Chainable
	ElseIf    map[string]*Chainable // map of condition => Stringable
	Else      *Chainable
}

func (i *If) String() string {
	str := fmt.Sprintf("if %s {\n%s\n}\n", i.Condition, i.Then.String())
	for condition, then := range i.ElseIf {
		str += fmt.Sprintf("else if %s {\n%s\n}\n", condition, then.String())
	}
	if !i.Else.IsEmpty() {
		str += fmt.Sprintf("else {\n%s\n}\n", i.Else.String())
	}
	return str
}
