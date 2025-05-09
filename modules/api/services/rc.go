package services

type RcApi struct {
}

func (api RcApi) FindRcAccounts() {

}

func (api RcApi) Expose(mr RegisterMethod) {
	mr("find_rc_accounts", "FindRcAccounts")
}
