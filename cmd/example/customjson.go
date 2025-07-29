package main

import (
	"mimic/lib/utils"

	"github.com/vsc-eco/hivego"
)

type customJson struct{}

func (customjson *customJson) params() hivego.HiveOperation {
	return &hivego.CustomJsonOperation{
		RequiredAuths:        []string{},
		RequiredPostingAuths: []string{utils.EnvOrPanic("TEST_USERNAME")},
		Id:                   "follow",
		Json:                 "[\"follow\",{\"follower\":\"hiveio\",\"following\":\"alice\",\"what\":[\"blog\"]}]",
	}
}
