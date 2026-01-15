package discord

import (
	"fmt"

	"github.com/mtzvd/ironroll/core/roll"
)

// formatResult converts a roll.Result into a Discord-friendly message.
//
// Markdown is intentionally simple to ensure compatibility
// across desktop and mobile clients.
func formatResult(r roll.Result) string {
	return fmt.Sprintf(
		"**Ironsworn Roll**\n\n"+
			"🎲 Action Die: `%d`\n"+
			"➕ Modifier: `%+d`\n"+
			"🎯 Challenge Dice: `%d`, `%d`\n\n"+
			"📊 **Total**: `%d`\n"+
			"✅ **Outcome**: **%s**",
		r.ActionDie,
		r.Modifier,
		r.ChallengeDice[0],
		r.ChallengeDice[1],
		r.Total,
		r.Outcome,
	)
}
