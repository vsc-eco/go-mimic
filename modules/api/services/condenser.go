package services

import (
	"log/slog"
	"mimic/mock"
	cdb "mimic/modules/db/mimic/condenserdb"
	"slices"
	"strings"
)

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

// get_accounts
func (t *Condenser) GetAccounts(args *GetAccountsArgs, reply *[]cdb.Account) {
	nameMatched := (*args)[0]
	db := cdb.Collection()

	if err := db.QueryGetAccounts(reply, nameMatched); err != nil {
		slog.Error("Failed to query for accounts.", "err", err)
		return
	}
}

type GlobalProps struct {
	HeadBlockNumber              int    `json:"head_block_number"`
	HeadBlockId                  string `json:"head_block_id"`
	Time                         string `json:"time"`
	CurrentWitness               string `json:"current_witness"`
	TotalPow                     string `json:"total_pow"`
	NumPowWitnesses              int    `json:"num_pow_witnesses"`
	VirtualSupply                string `json:"virtual_supply"`
	CurrentSupply                string `json:"current_supply"`
	ConfidentialSupply           string `json:"confidential_supply"`
	CurrentHbdSupply             string `json:"current_hbd_supply"`
	ConfidentialHbdSupply        string `json:"confidential_hbd_supply"`
	TotalVestingFundHive         string `json:"total_vesting_fund_hive"`
	TotalVestingShares           string `json:"total_vesting_shares"`
	TotalRewardFundHive          string `json:"total_reward_fund_hive"`
	TotalRewardShares2           string `json:"total_reward_shares2"`
	PendingRewardedVestingShares string `json:"pending_rewarded_vesting_shares"`
	PendingRewardedVestingHive   string `json:"pending_rewarded_vesting_hive"`
	HbdInterestRate              int    `json:"hbd_interest_rate"`
	HbdPrintRate                 int    `json:"hbd_print_rate"`
	MaximumBlockSize             int    `json:"maximum_block_size"`
	CurrentAslot                 int    `json:"current_aslot"`
	RecentSlotsFilled            string `json:"recent_slots_filled"`
	ParticipationCount           int    `json:"participation_count"`
	LastIrreversibleBlockNum     int    `json:"last_irreversible_block_num"`
	VotePowerReserveRate         int    `json:"vote_power_reserve_rate"`
}

// get_dynamic_global_properties
func (t *Condenser) GetDynamicGlobalProperties(args *[]string, reply *GlobalProps) {
	//Fake data for now until it gets hooked up with the rest of the mock context
	reply.HeadBlockNumber = 100
	reply.HeadBlockId = "1234567890"
	reply.Time = "2023-10-01T00:00:00"
	reply.CurrentWitness = "test"
	reply.TotalPow = "0"
	reply.NumPowWitnesses = 0
	reply.VirtualSupply = "100.000 HIVE"
	reply.CurrentSupply = "100.000 HIVE"
	reply.ConfidentialSupply = "0.000 HIVE"
	reply.CurrentHbdSupply = "100.000 HBD"
	reply.ConfidentialHbdSupply = "0.000 HBD"
	reply.TotalVestingFundHive = "100.000 HIVE"
	reply.TotalVestingShares = "100.000 HIVE"
	reply.TotalRewardFundHive = "100.000 HIVE"
	reply.TotalRewardShares2 = "0"
	reply.PendingRewardedVestingShares = "0.000 HIVE"
	reply.PendingRewardedVestingHive = "0.000 HIVE"
	reply.HbdInterestRate = 0
	reply.HbdPrintRate = 10000
	reply.MaximumBlockSize = 1000000
	reply.CurrentAslot = 0
	reply.RecentSlotsFilled = "0"
	reply.ParticipationCount = 0
	reply.LastIrreversibleBlockNum = 100
	reply.VotePowerReserveRate = 40

}

type MediumPrice struct {
	Base  string `json:"base"`
	Quote string `json:"quote"`
}

// get_current_median_history_price
func (t *Condenser) GetCurrentMedianHistoryPrice(args *[]string, reply *MediumPrice) {
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
func (t *Condenser) GetRewardFund(args *[]string, reply *RewardFund) {
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

// get_withdraw_routes
func (t *Condenser) GetWithdrawRoutes(args *[]string, reply *[]cdb.WithdrawRoute) {
	var (
		routes      []cdb.WithdrawRoute
		mockApiData = "condenser_api.get_withdraw_routes"
	)

	if err := mock.GetMockData(&routes, mockApiData); err != nil {
		slog.Error("Failed to read mock data",
			"mock-json", mockApiData,
			"err", err)
		return
	}

	*reply = make([]cdb.WithdrawRoute, 0, len(routes))

	user, transferDirection := (*args)[0], (*args)[1]

	allowedDirection := []string{"all", "incoming", "outgoing"}
	if !slices.Contains(allowedDirection, transferDirection) {
		slog.Warn("Invalid transfer direction query, allowed values: incoming, outgoing, all")
		return
	}

	filterMap(&routes, reply, func(r *cdb.WithdrawRoute) bool {
		switch transferDirection {

		case "incoming":
			return strings.EqualFold(user, r.ToAccount)

		case "outgoing":
			return strings.EqualFold(user, r.FromAccount)

		case "all":
			return strings.EqualFold(user, r.FromAccount) ||
				strings.EqualFold(user, r.ToAccount)

		default:
			panic("invalid transfer direction")
		}
	})
}

// get_open_orders
func (t *Condenser) GetOpenOrders(args *[]string, reply *[]cdb.OpenOrder) {
	var (
		orders       []cdb.OpenOrder
		mockFilePath = "condenser_api.get_open_orders"
	)

	if err := mock.GetMockData(&orders, mockFilePath); err != nil {
		slog.Error("Failed to read mock data",
			"mock-json", mockFilePath,
			"err", err)
		return
	}

	*reply = make([]cdb.OpenOrder, 0, len(orders))

	filterMap(&orders, reply, func(o *cdb.OpenOrder) bool {
		return slices.Contains(*args, o.Seller)
	})
}

// get_conversion_requests
// aka hbd -> hive conversion
func (t *Condenser) GetConversionRequests(args *[]int, reply *[]cdb.ConversionRequest) {
	var (
		conversionRequests []cdb.ConversionRequest
		mockFilePath       = "condenser_api.get_conversion_requests"
	)

	if err := mock.GetMockData(&conversionRequests, mockFilePath); err != nil {
		slog.Error("Failed to read mock data",
			"mock-json", mockFilePath,
			"err", err)
		return
	}

	*reply = make([]cdb.ConversionRequest, 0, len(conversionRequests))

	filterMap(
		&conversionRequests,
		reply,
		func(e *cdb.ConversionRequest) bool {
			return slices.Contains(*args, int(e.ID))
		},
	)
}

// Filters elements from `data` that matches the predicate `filterFunc`, then
// writes to `buf`
func filterMap[T any](data, buf *[]T, filterFunc func(*T) bool) {
	for _, d := range *data {
		if filterFunc(&d) {
			*buf = append(*buf, d)
		}
	}
}

// get_collateralized_conversion_requests
// aka hive -> hbd conversion
// NOTE: docs is empty right now...
// https://developers.hive.io/apidefinitions/#condenser_api.get_collateralized_conversion_requests
func (t *Condenser) GetCollateralizedConversionRequests(
	args *[]string,
	reply *[]cdb.ConversionRequest,
) {
	//For now send empty response until decided as necessary and implemented
	*reply = []cdb.ConversionRequest{}
}

// list_proposals
func (t *Condenser) ListProposals(args *[]any, reply *[]string) {
	//For now send empty response until decided as necessary and implemented
	*reply = []string{}
}

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
	rm("list_proposals", "ListProposals")
}
