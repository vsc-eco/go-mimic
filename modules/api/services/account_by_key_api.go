package services

import (
	"encoding/json"
	"fmt"
)

type AccountByKeyArg[T any] struct {
	Command string
	Param   T
}

func (a *AccountByKeyArg[T]) UnmarshalJSON(data []byte) error {
	buf := [2]any{}
	if err := json.Unmarshal(data, &buf); err != nil {
		return err
	}

	var ok bool
	a.Command, ok = buf[0].(string)
	if !ok {
		return fmt.Errorf("not a valid string: `%s`", buf[0])
	}

	paramJson, err := json.Marshal(buf[1])
	if err != nil {
		return fmt.Errorf("failed to decode param type: `%T`, `%T`",
			a.Param, buf[1])
	}

	return json.Unmarshal(paramJson, &a.Param)
}

type AccountByKeyAPI struct{}

type AccountUpdate struct {
	Account      string  `json:"account"`
	Posting      Posting `json:"posting"`
	MemoKey      string  `json:"memo_key"`
	JSONMetadata string  `json:"json_metadata"`
}

type Posting struct {
	WeightThreshold int       `json:"weight_threshold"`
	AccountAuths    []any     `json:"account_auths"`
	KeyAuths        []KeyAuth `json:"-"` // json key `key_auths`, need manual deserialization
}

type KeyAuth struct {
	Value  string
	Weight int
}

func (p *Posting) UnmarshalJSON(data []byte) error {
	mapBuf := make(map[string]any)
	if err := json.Unmarshal(data, &mapBuf); err != nil {
		return err
	}

	p.WeightThreshold = int(mapBuf["weight_threshold"].(float64))
	p.AccountAuths = mapBuf["account_auths"].([]any)

	keyAuthRaw := mapBuf["key_auths"].([]any)
	p.KeyAuths = make([]KeyAuth, len(keyAuthRaw))
	for i, entry := range keyAuthRaw {
		v := entry.([]any)
		p.KeyAuths[i].Value = v[0].(string)
		p.KeyAuths[i].Weight = int(v[1].(float64))
	}

	return nil
}

func (a *AccountByKeyAPI) AccountUpdate(args *AccountByKeyArg[AccountUpdate], reply *any) {
	fmt.Println()
	fmt.Println(args)
	fmt.Println()

}

func (a *AccountByKeyAPI) Expose(rm RegisterMethod) {
	rm("account_update", "AccountUpdate")
}
