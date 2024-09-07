package stringbuilder

type Call struct {
	Name string
	Args []Stringable
}

func (c *Call) String() string {
	str := c.Name + "("
	for _, arg := range c.Args {
		str += arg.String()
	}
	str += ")"
	return str
}
