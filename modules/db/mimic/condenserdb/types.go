package condenserdb

type GlobalProperties struct {
	HeadBlockNumber              uint32 `json:"head_block_number"`
	HeadBlockID                  string `json:"head_block_id"`
	Time                         string `json:"time"`
	CurrentWitness               string `json:"current_witness"`
	TotalPow                     string `json:"total_pow"`
	NumPowWitnesses              int64  `json:"num_pow_witnesses"`
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
	HbdInterestRate              int64  `json:"hbd_interest_rate"`
	HbdPrintRate                 int64  `json:"hbd_print_rate"`
	MaximumBlockSize             int64  `json:"maximum_block_size"`
	CurrentAslot                 int64  `json:"current_aslot"`
	RecentSlotsFilled            string `json:"recent_slots_filled"`
	ParticipationCount           int64  `json:"participation_count"`
	LastIrreversibleBlockNum     int64  `json:"last_irreversible_block_num"`
	VotePowerReserveRate         int64  `json:"vote_power_reserve_rate"`
}

type OpenOrder struct {
	ID         int64     `json:"id"`
	Created    string    `json:"created"`
	Expiration string    `json:"expiration"`
	Seller     string    `json:"seller"`
	Orderid    int64     `json:"orderid"`
	ForSale    int64     `json:"for_sale"`
	SellPrice  SellPrice `json:"sell_price"`
}

type SellPrice struct {
	Base  BasePrice `json:"base"`
	Quote BasePrice `json:"quote"`
}

type BasePrice struct {
	Amount    string `json:"amount"`
	Precision int64  `json:"precision"`
	Nai       string `json:"nai"`
}

type ConversionRequest struct {
	ID             int64  `json:"id"`
	Owner          string `json:"owner"`
	Requestid      int64  `json:"requestid"`
	Amount         Amount `json:"amount"`
	ConversionDate string `json:"conversion_date"`
}

type Amount struct {
	Amount    string `json:"amount"`
	Precision int64  `json:"precision"`
	Nai       string `json:"nai"`
}

type WithdrawRoute struct {
	ID          int64  `json:"id"`
	FromAccount string `json:"from_account"`
	ToAccount   string `json:"to_account"`
	Percent     int64  `json:"percent"`
	AutoVest    bool   `json:"auto_vest"`
}

type RewardFund struct {
	ID                     int64  `json:"id"`
	Name                   string `json:"name"`
	RewardBalance          string `json:"reward_balance"`
	RecentClaims           string `json:"recent_claims"`
	LastUpdate             string `json:"last_update"`
	ContentConstant        string `json:"content_constant"`
	PercentCurationRewards int64  `json:"percent_curation_rewards"`
	PercentContentRewards  int64  `json:"percent_content_rewards"`
	AuthorRewardCurve      string `json:"author_reward_curve"`
	CurationRewardCurve    string `json:"curation_reward_curve"`
}

type MedianPrice struct {
	Base  string `json:"base"`
	Quote string `json:"quote"`
}
