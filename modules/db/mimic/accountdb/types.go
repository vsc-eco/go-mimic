package accountdb

import (
	"github.com/vsc-eco/hivego"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Account struct {
	ObjectId primitive.ObjectID `json:"-"  bson:"_id,omitempty"`
	Id       int                `json:"id" bson:"-"`

	Name                string       `json:"name"`
	Owner               hivego.Auths `json:"owner"`
	Active              hivego.Auths `json:"active"`
	Posting             hivego.Auths `json:"posting"`
	MemoKey             string       `json:"memo_key"`
	JsonMeta            string       `json:"json_metadata"`
	JsonPostingMetadata string       `json:"posting_json_metadata"`
	LastOwnerUpdate     string       `json:"last_owner_update"`
	LastAccountUpdate   string       `json:"last_account_update"`
	Created             string       `json:"created"`
	Balance             string       `json:"balance"`
	HbdBalance          string       `json:"hbd_balance"`
	SavingsHbdBalance   string       `json:"savings_hbd_balance"`
	VestingShares       string       `json:"vesting_shares"`
	Reputation          int          `json:"reputation"`
}

type VotingManabar struct {
	CurrentMana    int `json:"current_mana"`
	LastUpdateTime int `json:"last_update_time"`
}

type DownvoteManabar struct {
	CurrentMana    int `json:"current_mana"`
	LastUpdateTime int `json:"last_update_time"`
}
