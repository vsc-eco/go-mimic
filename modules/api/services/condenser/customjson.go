package condenser

import (
	"encoding/json"
	"fmt"

	"github.com/vsc-eco/hivego"
)

func (c *Condenser) CustomJSON(
	args *CondenserParam[hivego.CustomJsonOperation],
	reply *any,
) {
	j, _ := json.MarshalIndent(args, "", "  ")
	fmt.Println(string(j), reply)
	*reply = make(map[string]any)
}
