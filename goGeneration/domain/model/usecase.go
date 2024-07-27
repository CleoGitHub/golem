package model

type Usecase struct {
	Function *Function
	Request  *Struct
	Result   *Struct
	Roles    []string
}

func (u *Usecase) Copy() *Usecase {
	return &Usecase{
		Function: u.Function.Copy().(*Function),
		Request:  u.Request.Copy().(*Struct),
		Result:   u.Result.Copy().(*Struct),
		Roles:    u.Roles,
	}
}
