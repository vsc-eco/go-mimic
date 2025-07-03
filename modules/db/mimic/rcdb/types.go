package rcdb

type RCAccount struct {
	Account                 string                  `json:"account"`
	DelegatedRC             int64                   `json:"delegated_rc"`
	MaxRC                   int64                   `json:"max_rc"`
	MaxRCCreationAdjustment MaxRCCreationAdjustment `json:"max_rc_creation_adjustment"`
	RCManabar               RCManabar               `json:"rc_manabar"`
	ReceivedDelegatedRC     int64                   `json:"received_delegated_rc"`
}

type MaxRCCreationAdjustment struct {
	Amount    string `json:"amount"`
	Nai       string `json:"nai"`
	Precision int64  `json:"precision"`
}

type RCManabar struct {
	CurrentMana    int64 `json:"current_mana"`
	LastUpdateTime int64 `json:"last_update_time"`
}
