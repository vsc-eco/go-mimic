package services

import "fmt"

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

//	{
//		"id": 1370484,
//		"name": "hiveio",
//		"owner": {
//		  "weight_threshold": 1,
//		  "account_auths": [],
//		  "key_auths": [
//			[
//			  "STM65PUAPA4yC4RgPtGgsPupxT6yJtMhmT5JHFdsT3uoCbR8WJ25s",
//			  1
//			]
//		  ]
//		},
//		"active": {
//		  "weight_threshold": 1,
//		  "account_auths": [],
//		  "key_auths": [
//			[
//			  "STM69zfrFGnZtU3gWFWpQJ6GhND1nz7TJsKBTjcWfebS1JzBEweQy",
//			  1
//			]
//		  ]
//		},
//		"posting": {
//		  "weight_threshold": 1,
//		  "account_auths": [["threespeak", 1], ["vimm.app", 1]],
//		  "key_auths": [
//			[
//			  "STM6vJmrwaX5TjgTS9dPH8KsArso5m91fVodJvv91j7G765wqcNM9",
//			  1
//			]
//		  ]
//		},
//		"memo_key": "STM7wrsg1BZogeK7X3eG4ivxmLaH69FomR8rLkBbepb3z3hm5SbXu",
//		"json_metadata": "",
//		"posting_json_metadata": "{\"profile\":{\"pinned\":\"none\",\"version\":2,\"website\":\"hive.io\",\"profile_image\":\"https://files.peakd.com/file/peakd-hive/hiveio/Jp2YHc6Q-hive-logo.png\",\"cover_image\":\"https://files.peakd.com/file/peakd-hive/hiveio/Xe1TcEBi-hive-banner.png\"}}",
//		"proxy": "",

//		"last_owner_update": "1970-01-01T00:00:00",
//		"last_account_update": "2020-11-12T01:20:48",
//		"created": "2020-03-06T12:22:48",
//		"mined": false,
//		"recovery_account": "steempeak",
//		"last_account_recovery": "1970-01-01T00:00:00",
//		"reset_account": "null",
//		"comment_count": 0,
//		"lifetime_vote_count": 0,
//		"post_count": 31,
//		"can_vote": true,
//		"voting_manabar": {
//		  "current_mana": "598442432741",
//		  "last_update_time": 1591297380
//		},
//		"downvote_manabar": {
//		  "current_mana": "149610608184",
//		  "last_update_time": 1591297380
//		},
//		"voting_power": 0,
//		"balance": "11.682 HIVE",
//		"savings_balance": "0.000 HIVE",
//		"hbd_balance": "43.575 HBD",
//		"hbd_seconds": "0",
//		"hbd_seconds_last_update": "2020-10-21T02:45:12",
//		"hbd_last_interest_payment": "2020-10-21T02:45:12",
//		"savings_hbd_balance": "0.000 HBD",
//		"savings_hbd_seconds": "0",
//		"savings_hbd_seconds_last_update": "1970-01-01T00:00:00",
//		"savings_hbd_last_interest_payment": "1970-01-01T00:00:00",
//		"savings_withdraw_requests": 0,
//		"reward_hbd_balance": "0.000 HBD",
//		"reward_hive_balance": "0.000 HIVE",
//		"reward_vesting_balance": "0.000000 VESTS",
//		"reward_vesting_hive": "0.000 HIVE",
//		"vesting_shares": "598442.432741 VESTS",
//		"delegated_vesting_shares": "0.000000 VESTS",
//		"received_vesting_shares": "0.000000 VESTS",
//		"vesting_withdraw_rate": "0.000000 VESTS",
//		"post_voting_power": "598442.432741 VESTS",
//		"next_vesting_withdrawal": "1969-12-31T23:59:59",
//		"withdrawn": 0,
//		"to_withdraw": 0,
//		"withdraw_routes": 0,
//		"pending_transfers": 0,
//		"curation_rewards": 0,
//		"posting_rewards": 604589,
//		"proxied_vsf_votes": [0, 0, 0, 0],
//		"witnesses_voted_for": 0,
//		"last_post": "2021-03-23T18:05:48",
//		"last_root_post": "2021-03-23T18:05:48",
//		"last_vote_time": "1970-01-01T00:00:00",
//		"post_bandwidth": 0,
//		"pending_claimed_accounts": 0,
//		"delayed_votes": [
//		  {
//			"time": "2021-02-24T05:08:21",
//			"val": "11550765516955"
//		  },
//		  {
//			"time": "2021-02-26T15:46:06",
//			"val": "633465684569"
//		  },
//		  {
//			"time": "2021-03-07T17:54:39",
//			"val": "1000000037683"
//		  },
//		  {
//			"time": "2021-03-16T05:54:33",
//			"val": "999978763511"
//		  },
//		  {
//			"time": "2021-03-18T06:06:00",
//			"val": "1000000171317"
//		  }
//		],
//		"vesting_balance": "0.000 HIVE",
//		"reputation": "88826789432105",
//		"transfer_history": [],
//		"market_history": [],
//		"post_history": [],
//		"vote_history": [],
//		"other_history": [],
//		"witness_votes": [],
//		"tags_usage": [],
//		"guest_bloggers": []
//	  }
type AccountAuthority struct {
	WeightThreshold int           `json:"weight_threshold"`
	AccountAuths    []interface{} `json:"account_auths"`
	KeyAuths        []interface{} `json:"key_auths"`
}
type GetAccountsReply struct {
	Id                  int              `json:"id"`
	Name                string           `json:"name"`
	Owner               AccountAuthority `json:"owner"`
	Active              AccountAuthority `json:"active"`
	Posting             AccountAuthority `json:"posting"`
	MemoKey             string           `json:"memo_key"`
	JsonMeta            string           `json:"json_metadata"`
	JsonPostingMetadata string           `json:"posting_json_metadata"`
	Proxy               string           `json:"proxy"`
	LastOwnerUpdate     string           `json:"last_owner_update"`
	LastAccountUpdate   string           `json:"last_account_update"`
	Created             string           `json:"created"`
	Mined               bool             `json:"mined"`
	RecoveryAccount     string           `json:"recovery_account"`
	LastAccountRecovery string           `json:"last_account_recovery"`
	ResetAccount        string           `json:"reset_account"`
	CommentCount        int              `json:"comment_count"`
	LifetimeVoteCount   int              `json:"lifetime_vote_count"`
	PostCount           int              `json:"post_count"`
	CanVote             bool             `json:"can_vote"`

	VotingManabar   VotingManabar   `json:"voting_manabar"`
	DownvoteManabar DownvoteManabar `json:"downvote_manabar"`
	VotingPower     int             `json:"voting_power"`

	Balance                       string `json:"balance"`
	SavingsBalance                string `json:"savings_balance"`
	HbdBalance                    string `json:"hbd_balance"`
	HbdSeconds                    string `json:"hbd_seconds"`
	HbdSecondsLastUpdate          string `json:"hbd_seconds_last_update"`
	HbdLastInterestPayment        string `json:"hbd_last_interest_payment"`
	SavingsHbdBalance             string `json:"savings_hbd_balance"`
	SavingsHbdSeconds             string `json:"savings_hbd_seconds"`
	SavingsHbdSecondsLastUpdate   string `json:"savings_hbd_seconds_last_update"`
	SavingsHbdLastInterestPayment string `json:"savings_hbd_last_interest_payment"`
	SavingsWithdrawRequests       int    `json:"savings_withdraw_requests"`
	RewardHbdBalance              string `json:"reward_hbd_balance"`
	RewardHiveBalance             string `json:"reward_hive_balance"`
	RewardVestingBalance          string `json:"reward_vesting_balance"`
	RewardVestingHive             string `json:"reward_vesting_hive"`
	VestingShares                 string `json:"vesting_shares"`

	DelegatedVestingShares string `json:"delegated_vesting_shares"`
	ReceivedVestingShares  string `json:"received_vesting_shares"`
	VestingWithdrawRate    string `json:"vesting_withdraw_rate"`
	PostVotingPower        string `json:"post_voting_power"`
	NextVestingWithdrawal  string `json:"next_vesting_withdrawal"`
	Withdrawn              int    `json:"withdrawn"`
	ToWithdraw             int    `json:"to_withdraw"`
	WithdrawRoutes         int    `json:"withdraw_routes"`
	PendingTransfers       int    `json:"pending_transfers"`
	CurationRewards        int    `json:"curation_rewards"`
	PostingRewards         int    `json:"posting_rewards"`
	ProxiedVsfVotes        []int  `json:"proxied_vsf_votes"`
	WitnessesVotedFor      int    `json:"witnesses_voted_for"`
	LastPost               string `json:"last_post"`
	LastRootPost           string `json:"last_root_post"`
	LastVoteTime           string `json:"last_vote_time"`
	PostBandwidth          int    `json:"post_bandwidth"`
	PendingClaimedAccounts int    `json:"pending_claimed_accounts"`
	DelayedVotes           []any  `json:"delayed_votes"`
	VestingBalance         string `json:"vesting_balance"`
	Reputation             string `json:"reputation"`
	TransferHistory        []any  `json:"transfer_history"`
	MarketHistory          []any  `json:"market_history"`
	PostHistory            []any  `json:"post_history"`
	VoteHistory            []any  `json:"vote_history"`
	OtherHistory           []any  `json:"other_history"`
	WitnessVotes           []any  `json:"witness_votes"`
	TagsUsage              []any  `json:"tags_usage"`
	GuestBloggers          []any  `json:"guest_bloggers"`
}

type VotingManabar struct {
	CurrentMana    string `json:"current_mana"`
	LastUpdateTime int    `json:"last_update_time"`
}

type DownvoteManabar struct {
	CurrentMana    string `json:"current_mana"`
	LastUpdateTime int    `json:"last_update_time"`
}

// get_accounts
func (t *Condenser) GetAccounts(args *GetAccountsArgs, reply *[]GetAccountsReply) {

	for _, account := range *args {
		*reply = append(*reply, GetAccountsReply{
			Id:   1,
			Name: account[0],
			Owner: AccountAuthority{
				WeightThreshold: 1,
				AccountAuths: []interface{}{
					1, "",
				},
				KeyAuths: []interface{}{},
			},
			Active: AccountAuthority{
				WeightThreshold: 1,
				AccountAuths:    []interface{}{},
				KeyAuths:        []interface{}{},
			},
			Posting: AccountAuthority{
				WeightThreshold: 1,
				AccountAuths:    []interface{}{},
				KeyAuths:        []interface{}{},
			},
			VotingManabar: struct {
				CurrentMana    string `json:"current_mana"`
				LastUpdateTime int    `json:"last_update_time"`
			}{
				CurrentMana:    "1000000000",
				LastUpdateTime: 1000000000,
			},
			MemoKey:                       "STM7wrsg1BZogeK7X3eG4ivxmLaH69FomR8rLkBbepb3z3hm5SbXu",
			JsonMeta:                      "",
			Proxy:                         "",
			Created:                       "2023-10-01T00:00:00",
			Mined:                         false,
			RecoveryAccount:               "test",
			LastAccountRecovery:           "2023-10-01T00:00:00",
			ResetAccount:                  "null",
			CommentCount:                  0,
			LifetimeVoteCount:             0,
			PostCount:                     0,
			CanVote:                       true,
			VotingPower:                   0,
			Balance:                       "1.100 HIVE",
			SavingsBalance:                "1.200 HIVE",
			HbdBalance:                    "1.300 HBD",
			HbdSeconds:                    "0",
			HbdSecondsLastUpdate:          "2023-10-01T00:00:00",
			HbdLastInterestPayment:        "2023-10-01T00:00:00",
			SavingsHbdBalance:             "10.000 HBD",
			SavingsHbdSeconds:             "0",
			SavingsHbdSecondsLastUpdate:   "2023-10-01T00:00:00",
			SavingsHbdLastInterestPayment: "2023-10-01T00:00:00",
			SavingsWithdrawRequests:       0,
			RewardHbdBalance:              "0.000 HBD",
			RewardHiveBalance:             "0.000 HIVE",
			RewardVestingBalance:          "0.000000 VESTS",
			RewardVestingHive:             "0.000 HIVE",
			VestingShares:                 "10000000.000000 VESTS",
			DelegatedVestingShares:        "0.000000 VESTS",
			ReceivedVestingShares:         "0.000000 VESTS",

			VestingWithdrawRate:    "0.000000 VESTS",
			PostVotingPower:        "0.000000 VESTS",
			NextVestingWithdrawal:  "2023-10-01T00:00:00",
			Withdrawn:              0,
			ToWithdraw:             0,
			WithdrawRoutes:         0,
			PendingTransfers:       0,
			CurationRewards:        0,
			PostingRewards:         0,
			ProxiedVsfVotes:        []int{0, 0, 0, 0},
			WitnessesVotedFor:      0,
			LastPost:               "2023-10-01T00:00:00",
			LastRootPost:           "2023-10-01T00:00:00",
			LastVoteTime:           "2023-10-01T00:00:00",
			PostBandwidth:          0,
			PendingClaimedAccounts: 0,
			DelayedVotes:           []any{},
			VestingBalance:         "10000000.000000 VESTS",
			Reputation:             "0",
			TransferHistory:        []any{},
			MarketHistory:          []any{},
			PostHistory:            []any{},
			VoteHistory:            []any{},
			OtherHistory:           []any{},
			WitnessVotes:           []any{},
			TagsUsage:              []any{},
			GuestBloggers:          []any{},
		})
		fmt.Println("act", account[0])
	}
	fmt.Println("GetAccounts", args)
}

//	{
//		"head_block_number": 0,
//		"head_block_id": "0000000000000000000000000000000000000000",
//		"time": "1970-01-01T00:00:00",
//		"current_witness": "",
//		"total_pow": "18446744073709551615",
//		"num_pow_witnesses": 0,
//		"virtual_supply": "0.000 HIVE",
//		"current_supply": "0.000 HIVE",
//		"confidential_supply": "0.000 HIVE",
//		"current_hbd_supply": "0.000 HIVE",
//		"confidential_hbd_supply": "0.000 HIVE",
//		"total_vesting_fund_hive": "0.000 HIVE",
//		"total_vesting_shares": "0.000 HIVE",
//		"total_reward_fund_hive": "0.000 HIVE",
//		"total_reward_shares2": "0",
//		"pending_rewarded_vesting_shares": "0.000 HIVE",
//		"pending_rewarded_vesting_hive": "0.000 HIVE",
//		"hbd_interest_rate": 0,
//		"hbd_print_rate": 10000,
//		"maximum_block_size": 0,
//		"current_aslot": 0,
//		"recent_slots_filled": "0",
//		"participation_count": 0,
//		"last_irreversible_block_num": 0,
//		"vote_power_reserve_rate": 40
//	  }
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

type OpenOrder struct {
	Created    string `json:"created"`
	Expiration string `json:"expiration"`
	ForSale    int    `json:"for_sale"`
	Id         int    `json:"id"`
	OrderId    int    `json:"orderid"`
	RealPrice  string `json:"real_price"`
	Rewarded   bool   `json:"rewarded"`
	SellPrice  struct {
		Base  string `json:"base"`
		Quote string `json:"quote"`
	} `json:"sell_price"`
	Seller string `json:"seller"`
}

// get_open_orders
func (t *Condenser) GetOpenOrders(args *[]string, reply *[]OpenOrder) {
	//Fake data for now until it gets hooked up with the rest of the mock context
	reply = &[]OpenOrder{}
}

type ConversionRequest struct {
}

// get_conversion_requests
// aka hbd -> hive conversion
func (t *Condenser) GetConversionRequests(args *[]string, reply *[]ConversionRequest) {
	//For now send empty response until decided as necessary and implemented
	*reply = []ConversionRequest{}
}

// get_collateralized_conversion_requests
// aka hive -> hbd conversion
func (t *Condenser) GetCollateralizedConversionRequests(args *[]string, reply *[]ConversionRequest) {
	//For now send empty response until decided as necessary and implemented
	*reply = []ConversionRequest{}
}

// list_proposals
func (t *Condenser) ListProposals(args *[]interface{}, reply *[]string) {
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
