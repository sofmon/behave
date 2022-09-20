package behave

// Then_do simple do function
func Eval(do func() Action) *eval {
	return &eval{do}
}

// DoSomething action
type eval struct {
	f func() Action
}

/* Action implementation */

func (x *eval) String(res any) string {
	a := x.f()
	if a == nil {
		panic("can not evaluate statement as action is null")
	}
	return a.String(res)
}

// Do the action
func (x *eval) Do(res any) any {
	a := x.f()
	if a == nil {
		panic("can not evaluate statement as action is null")
	}
	return a.Do(res)
}
