package services

type TestMethodArgs struct {
	A int `json:"a"`
	B int `json:"b"`
}

// TestMethodReply is the output from exampleservice.test_method.
type TestMethodReply struct {
	Sum     int `json:"sum"`
	Product int `json:"product"`
}
type Condenser struct {
}

func (t *Condenser) GetBlock(args *TestMethodArgs, reply *TestMethodReply) error {
	// Fill reply pointer to send the data back
	reply.Sum = args.A + args.B + 1
	reply.Product = args.A * args.B
	return nil
}

type GetAccountsArgs [][]string

type GetAccountsReply struct {
}

func (t *Condenser) GetAccounts(args *GetAccountsArgs, reply *GetAccountsReply) {

}

type GlobalProps struct {
	HeadBlockNumber int    `json:"head_block_number"`
	HeadBlockId     string `json:"head_block_id"`
	Time            string `json:"time"`
}

// get_dynamic_global_properties
func (t *Condenser) GetDynamicGlobalProperties(args any, reply *GlobalProps) {
	//Fake data for now until it gets hooked up with the rest of the mock context
	reply.HeadBlockNumber = 100
	reply.HeadBlockId = "1234567890"
	reply.Time = "2023-10-01T00:00:00"
}

type MediumPrice struct {
	Base  string `json:"base"`
	Quote string `json:"quote"`
}

// get_current_median_history_price
func (t *Condenser) GetCurrentMedianHistoryPrice(args any, reply *MediumPrice) {
	//Fake data for now until it gets hooked up with the rest of the mock context
	reply.Base = "100.000 SBD"
	reply.Quote = "100.000 HIVE"
}

//	{
//	  "id": 0,
//	  "name": "",
//	  "reward_balance": "0.000 HIVE",
//	  "recent_claims": "0",
//	  "last_update": "1970-01-01T00:00:00",
//	  "content_constant": "0",
//	  "percent_curation_rewards": 0,
//	  "percent_content_rewards": 0,
//	  "author_reward_curve": "quadratic",
//	  "curation_reward_curve": "34723648"
//	}
type RewardFund struct {
	Id                  int    `json:"id"`
	Name                string `json:"name"`
	RewardBalance       string `json:"reward_balance"`
	RecentClaims        string `json:"recent_claims"`
	LastUpdate          string `json:"last_update"`
	ContentConstant     string `json:"content_constant"`
	PercentCuration     int    `json:"percent_curation_rewards"`
	PercentContent      int    `json:"percent_content_rewards"`
	AuthorRewardCurve   string `json:"author_reward_curve"`
	CurationRewardCurve string `json:"curation_reward_curve"`
}

// get_reward_fund
func (t *Condenser) GetRewardFund(args []string, reply *RewardFund) {
	//Fake data for now until it gets hooked up with the rest of the mock context
	reply.Id = 1
	reply.Name = "test"
	reply.RewardBalance = "100.000 HIVE"
	reply.RecentClaims = "1000"
	reply.LastUpdate = "2023-10-01T00:00:00"
	reply.ContentConstant = "1000"
	reply.PercentCuration = 50
	reply.PercentContent = 50
	reply.AuthorRewardCurve = "linear"
	reply.CurationRewardCurve = "quadratic"
}

// [
//
//	{
//	  "id": 0,
//	  "from_account": "",
//	  "to_account": "",
//	  "percent": 0,
//	  "auto_vest": false
//	}
//
// ]
type WithdrawRoute struct {
	Id          int    `json:"id"`
	FromAccount string `json:"from_account"`
	ToAccount   string `json:"to_account"`
	Percent     int    `json:"percent"`
	AutoVest    bool   `json:"auto_vest"`
}

// get_withdraw_routes
func (t *Condenser) GetWithdrawRoutes(args [2]string, reply *[]WithdrawRoute) {
	//Fake data for now until it gets hooked up with the rest of the mock context
	reply = &[]WithdrawRoute{
		{
			Id:          1,
			FromAccount: "test",
			ToAccount:   "test2",
			Percent:     50,
			AutoVest:    true,
		},
	}

}

// get_open_orders
func (t *Condenser) GetOpenOrders(args any, reply *[]string) {
	//Fake data for now until it gets hooked up with the rest of the mock context
	reply = &[]string{}

}

// get_conversion_requests
func (t *Condenser) GetConversionRequests(args any) {

}

// get_collateralized_conversion_requests
func (t *Condenser) GetCollateralizedConversionRequests() {}

func (t *Condenser) Expose(rm RegisterMethod) {
	rm("get_block", "GetBlock")
	rm("get_dynamic_global_properties", "GetDynamicGlobalProperties")
	rm("get_current_median_history_price", "GetCurrentMedianHistoryPrice")
	rm("get_reward_fund", "GetRewardFund")
	rm("get_withdraw_routes", "GetWithdrawRoutes")
	rm("get_open_orders", "GetOpenOrders")
	rm("get_conversion_requests", "GetConversionRequests")
	rm("get_collateralized_conversion_requests", "GetCollateralizedConversionRequests")
	rm("get_accounts", "GetAccounts")
}
