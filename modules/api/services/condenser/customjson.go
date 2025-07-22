package condenser

import (
	"fmt"
	"log"
	"mimic/lib/validator"

	"github.com/vsc-eco/hivego"
)

func (c *Condenser) CustomJSON(
	args *jsonRpcParam[hivego.CustomJsonOperation],
	reply *any,
) {
	if err := validator.New().Struct(args.Op); err != nil {
		log.Println(err) // TODO: log error with slog here
		return
	}

	fmt.Println(args.Op, reply)
	*reply = make(map[string]any)
}
