package redis

import _ "embed"

var (
	//go:embed scripts/unlock.lua
	unlockScriptSrc string

	//go:embed scripts/extend.lua
	extendScriptSrc string

	//go:embed scripts/deduct_stock.lua
	deductStockSrc string

	//go:embed scripts/check_limit.lua
	checkLimitSrc string
)
