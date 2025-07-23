package opvalidator

import "github.com/vsc-eco/hivego"

type customJsonValidator struct{}

func (c *customJsonValidator) Validate(
	_ hivego.HiveOperation,
) error {
	panic("not implemented") // TODO: Implement
}
