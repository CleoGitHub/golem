package stringbuilder

type StringBuilder struct {
	Head *Chainable
}

func NewStringBuilder() *StringBuilder {
	return &StringBuilder{}
}

func (s *StringBuilder) Append(elem Stringable) *StringBuilder {
	if s.Head.IsEmpty() {
		s.Head = &Chainable{Elem: elem}
	} else {
		current := s.Head
		for current.HasNext() {
			current = current.Next
		}
		current.Next = &Chainable{Elem: elem}
	}
	return s
}

func (s *StringBuilder) String() string {
	return s.Head.String()
}
