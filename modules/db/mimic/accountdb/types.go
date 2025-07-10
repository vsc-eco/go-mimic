package accountdb

import (
	"encoding/json"
	"mimic/lib/encoder"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Account struct {
	ObjectId primitive.ObjectID `json:"-"  bson:"_id,omitempty"`
	Id       int                `json:"id" bson:"-"`

	Name                string           `json:"name"`
	Owner               AccountAuthority `json:"owner"`
	Active              AccountAuthority `json:"active"`
	Posting             AccountAuthority `json:"posting"`
	MemoKey             string           `json:"memo_key"`
	JsonMeta            string           `json:"json_metadata"`
	JsonPostingMetadata string           `json:"posting_json_metadata"`
	LastOwnerUpdate     string           `json:"last_owner_update"`
	LastAccountUpdate   string           `json:"last_account_update"`
	Created             string           `json:"created"`
	Balance             string           `json:"balance"`
	HbdBalance          string           `json:"hbd_balance"`
	SavingsHbdBalance   string           `json:"savings_hbd_balance"`
	VestingShares       string           `json:"vesting_shares"`
	Reputation          int              `json:"reputation"`
}

type AccountAuthority struct {
	WeightThreshold int           `json:"weight_threshold"`
	AccountAuths    []AccountAuth `json:"account_auths"`
	KeyAuths        []KeyAuth     `json:"key_auths"`
}

// Format: [account_name, weight]
type AccountAuth struct {
	Account string `json:"account"`
	Weight  int    `json:"weight"`
}

func (a *AccountAuth) MarshalJSON() ([]byte, error) {
	buf := [2]any{a.Account, a.Weight}
	return json.Marshal(buf)
}

func (a *AccountAuth) UnmarshalJSON(raw []byte) error {
	return encoder.JsonArrayDeserialize(a, raw)
}

// Format: [public_key, weight]
type KeyAuth struct {
	PublicKey string `json:"public_key"`
	Weight    int    `json:"weight"`
}

func (k *KeyAuth) MarshalJSON() ([]byte, error) {
	buf := [2]any{k.PublicKey, k.Weight}
	return json.Marshal(buf)
}

func (a *KeyAuth) UnmarshalJSON(raw []byte) error {
	return encoder.JsonArrayDeserialize(a, raw)
}

type VotingManabar struct {
	CurrentMana    int `json:"current_mana"`
	LastUpdateTime int `json:"last_update_time"`
}

type DownvoteManabar struct {
	CurrentMana    int `json:"current_mana"`
	LastUpdateTime int `json:"last_update_time"`
}
