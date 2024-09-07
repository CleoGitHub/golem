package stringbuilder

type Chainable struct {
	Elem Stringable
	Next *Chainable
}

func (c *Chainable) IsEmpty() bool {
	if c == nil {
		return true
	}
	return c.Elem == nil
}

func (c *Chainable) HasNext() bool {
	return c.Next != nil
}

func (c *Chainable) String() string {
	if !c.HasNext() {
		return c.Elem.String()
	}
	return c.Elem.String() + c.Next.String()
}
