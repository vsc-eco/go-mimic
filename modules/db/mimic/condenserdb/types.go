package condenserdb

import "go.mongodb.org/mongo-driver/bson/primitive"

type Account struct {
	ObjectId            primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Id                  int                `json:"id" bson:"-"`
	Name                string             `json:"name"`
	Owner               AccountAuthority   `json:"owner"`
	Active              AccountAuthority   `json:"active"`
	Posting             AccountAuthority   `json:"posting"`
	MemoKey             string             `json:"memo_key"`
	JsonMeta            string             `json:"json_metadata"`
	JsonPostingMetadata string             `json:"posting_json_metadata"`
	Proxy               string             `json:"proxy"`
	LastOwnerUpdate     string             `json:"last_owner_update"`
	LastAccountUpdate   string             `json:"last_account_update"`
	Created             string             `json:"created"`
	Mined               bool               `json:"mined"`
	RecoveryAccount     string             `json:"recovery_account"`
	LastAccountRecovery string             `json:"last_account_recovery"`
	ResetAccount        string             `json:"reset_account"`
	CommentCount        int                `json:"comment_count"`
	LifetimeVoteCount   int                `json:"lifetime_vote_count"`
	PostCount           int                `json:"post_count"`
	CanVote             bool               `json:"can_vote"`

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
	Reputation             int    `json:"reputation"`
	TransferHistory        []any  `json:"transfer_history"`
	MarketHistory          []any  `json:"market_history"`
	PostHistory            []any  `json:"post_history"`
	VoteHistory            []any  `json:"vote_history"`
	OtherHistory           []any  `json:"other_history"`
	WitnessVotes           []any  `json:"witness_votes"`
	TagsUsage              []any  `json:"tags_usage"`
	GuestBloggers          []any  `json:"guest_bloggers"`
}

type AccountAuthority struct {
	WeightThreshold int   `json:"weight_threshold"`
	AccountAuths    []any `json:"account_auths"`
	KeyAuths        []any `json:"key_auths"`
}

type VotingManabar struct {
	CurrentMana    int `json:"current_mana"`
	LastUpdateTime int `json:"last_update_time"`
}

type DownvoteManabar struct {
	CurrentMana    int `json:"current_mana"`
	LastUpdateTime int `json:"last_update_time"`
}
