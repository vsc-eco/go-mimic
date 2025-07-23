package main

import "github.com/vsc-eco/hivego"

type customJson struct{}

func (customjson *customJson) params() any {
	op := &hivego.CustomJsonOperation{
		RequiredAuths:        []string{},
		RequiredPostingAuths: []string{"hiveio"},
		Id:                   "follow",
		Json:                 "[\"follow\",{\"follower\":\"hiveio\",\"following\":\"alice\",\"what\":[\"blog\"]}]",
	}

	return op
}

func (customjson *customJson) OpName() string {
	return "custom_json"
}

func (c *customJson) SerializeOp() ([]byte, error) {
	panic("Not implemented")
}

/*
func customJson() error {
	op := hivego.CustomJsonOperation{
		RequiredAuths:        []string{},
		RequiredPostingAuths: []string{"hiveio"},
		Id:                   "follow",
		Json:                 "[\"follow\",{\"follower\":\"hiveio\",\"following\":\"alice\",\"what\":[\"blog\"]}]",
	}
	return nil
}
*/
