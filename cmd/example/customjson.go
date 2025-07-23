package main

import "github.com/vsc-eco/hivego"

type customJson struct{}

func (customjson *customJson) params() hivego.HiveOperation {
	return &hivego.CustomJsonOperation{
		RequiredAuths:        []string{},
		RequiredPostingAuths: []string{"hiveio"},
		Id:                   "follow",
		Json:                 "[\"follow\",{\"follower\":\"hiveio\",\"following\":\"alice\",\"what\":[\"blog\"]}]",
	}
}
