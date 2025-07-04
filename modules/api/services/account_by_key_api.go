package services

type AccountByKeyAPI struct{}

func (a *AccountByKeyAPI) AccountUpdate(args *any, reply *any) {}

func (a *AccountByKeyAPI) Expose(rm RegisterMethod) {
	rm("account_update", "AccountUpdate")
}
