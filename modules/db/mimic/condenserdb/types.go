package condenserdb

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
