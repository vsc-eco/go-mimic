package hiveop

import (
	"mimic/lib/hive"

	"github.com/vsc-eco/hivego"
)

// CustomJson implements Operation
type CustomJson struct {
	*hivego.CustomJsonOperation `json:",inline"`
}

func (c *CustomJson) SigningAuthorities() []hive.SigningAuthorities {
	cap := len(c.RequiredAuths) + len(c.RequiredPostingAuths)

	s := make([]hive.SigningAuthorities, 0, cap)

	for i := range len(c.RequiredAuths) {
		s = append(s, hive.SigningAuthorities{
			Account: c.RequiredAuths[i],
			KeyType: hive.ActiveKeyRole,
		})
	}

	for i := range len(c.RequiredPostingAuths) {
		s = append(s, hive.SigningAuthorities{
			Account: c.RequiredPostingAuths[i],
			KeyType: hive.PostingKeyRole,
		})
	}

	return s
}
