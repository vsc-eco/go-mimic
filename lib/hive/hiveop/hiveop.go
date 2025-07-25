package hiveop

import (
	"mimic/lib/hive"

	"github.com/vsc-eco/hivego"
)

type Operation interface {
	hivego.HiveOperation
	SigningAuthorities() []hive.SigningAuthorities
}
