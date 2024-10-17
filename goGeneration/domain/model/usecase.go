package model

type Usecase struct {
	Function *Function
	Request  *Struct
	Result   *Struct
	Roles    []string
}

func (u *Usecase) Copy() *Usecase {
	return &Usecase{
		Function: u.Function.Copy(),
		Request:  u.Request.Copy(),
		Result:   u.Result.Copy(),
		Roles:    u.Roles,
	}
}
