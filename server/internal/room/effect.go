package room

import (
	"fmt"
)

// EffectResult details the outcome of applying a cell effect to a player.
type EffectResult struct {
	EffectType  string                 `json:"effect_type"`
	Description string                 `json:"description"`
	XPDelta     int                    `json:"xp_delta"`
	NewXP       int                    `json:"new_xp"`
	NewPosition int                    `json:"new_position"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

// EffectFunc defines the function signature for cell effect handlers.
type EffectFunc func(r *Room, p *Player, cell BoardCell) EffectResult

// EffectDispatchTable maps cell effect type strings to their respective handler functions.
var EffectDispatchTable = map[string]EffectFunc{
	"xp_gain":           handleXPGain,
	"xp_loss":           handleXPLoss,
	"teleport":          handleTeleport,
	"skip_next":         handleSkipNext,
	"double_xp":         handleDoubleXP,
	"free_pass":         handleFreePass,
	"special_challenge": handleSpecialChallenge,
	"mystery":           handleMystery,
	"deploy":            handleDeploy,
	"code_freeze":       handleCodeFreeze,
	"coffee_break":      handleCoffeeBreak,
	"deadline":          handleDeadline,
}

// ApplyCellEffect executes the effect handler for a cell's type via the dispatch table.
func (r *Room) ApplyCellEffect(p *Player, cell BoardCell) EffectResult {
	handler, found := EffectDispatchTable[cell.Type]
	if !found {
		return EffectResult{
			EffectType:  cell.Type,
			Description: fmt.Sprintf("No effect defined for cell type '%s'", cell.Type),
			XPDelta:     0,
			NewXP:       p.XP,
			NewPosition: p.Position,
		}
	}
	return handler(r, p, cell)
}

func parseAmount(params map[string]interface{}, defaultVal int) int {
	if val, ok := params["amount"]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		}
	}
	if size, ok := params["size"].(string); ok {
		switch size {
		case "S":
			return 10
		case "M":
			return 25
		case "L":
			return 50
		}
	}
	return defaultVal
}

func handleXPGain(r *Room, p *Player, cell BoardCell) EffectResult {
	amount := parseAmount(cell.Params, 10)
	xpMultiplier := 1
	if p.DoubleXP {
		xpMultiplier = 2
		p.DoubleXP = false
	}
	gained := amount * xpMultiplier
	p.XP += gained

	desc := fmt.Sprintf("Gained +%d XP", gained)
	if xpMultiplier > 1 {
		desc += " (Double XP multiplier active!)"
	}

	return EffectResult{
		EffectType:  "xp_gain",
		Description: desc,
		XPDelta:     gained,
		NewXP:       p.XP,
		NewPosition: p.Position,
	}
}

func handleXPLoss(r *Room, p *Player, cell BoardCell) EffectResult {
	amount := parseAmount(cell.Params, 10)
	oldXP := p.XP
	p.XP -= amount
	if p.XP < 0 {
		p.XP = 0
	}
	actualLoss := oldXP - p.XP

	return EffectResult{
		EffectType:  "xp_loss",
		Description: fmt.Sprintf("Lost -%d XP", actualLoss),
		XPDelta:     -actualLoss,
		NewXP:       p.XP,
		NewPosition: p.Position,
	}
}

func handleTeleport(r *Room, p *Player, cell BoardCell) EffectResult {
	targetPos := 0
	if val, ok := cell.Params["target_position"]; ok {
		switch v := val.(type) {
		case int:
			targetPos = v
		case float64:
			targetPos = int(v)
		}
	}
	p.Position = targetPos % 32

	return EffectResult{
		EffectType:  "teleport",
		Description: fmt.Sprintf("Teleported to cell %d (%s)", p.Position, r.GetCell(p.Position).Name),
		XPDelta:     0,
		NewXP:       p.XP,
		NewPosition: p.Position,
	}
}

func handleSkipNext(r *Room, p *Player, cell BoardCell) EffectResult {
	p.SkipNextTurn = true
	return EffectResult{
		EffectType:  "skip_next",
		Description: "Your next turn will be skipped",
		XPDelta:     0,
		NewXP:       p.XP,
		NewPosition: p.Position,
	}
}

func handleDoubleXP(r *Room, p *Player, cell BoardCell) EffectResult {
	p.DoubleXP = true
	return EffectResult{
		EffectType:  "double_xp",
		Description: "Double XP modifier granted for your next XP gain",
		XPDelta:     0,
		NewXP:       p.XP,
		NewPosition: p.Position,
	}
}

func handleFreePass(r *Room, p *Player, cell BoardCell) EffectResult {
	p.FreePasses++
	return EffectResult{
		EffectType:  "free_pass",
		Description: fmt.Sprintf("Received a Free Pass (total: %d)", p.FreePasses),
		XPDelta:     0,
		NewXP:       p.XP,
		NewPosition: p.Position,
	}
}

func handleSpecialChallenge(r *Room, p *Player, cell BoardCell) EffectResult {
	bonus := 30
	if val, ok := cell.Params["bonus"]; ok {
		switch v := val.(type) {
		case int:
			bonus = v
		case float64:
			bonus = int(v)
		}
	}
	multiplier := 1
	if p.DoubleXP {
		multiplier = 2
		p.DoubleXP = false
	}
	gained := bonus * multiplier
	p.XP += gained

	return EffectResult{
		EffectType:  "special_challenge",
		Description: fmt.Sprintf("Completed special bonus challenge! Gained +%d XP", gained),
		XPDelta:     gained,
		NewXP:       p.XP,
		NewPosition: p.Position,
	}
}

func handleMystery(r *Room, p *Player, cell BoardCell) EffectResult {
	// Mystery box picks a random effect from the catalog
	mysteryOptions := []string{
		"xp_gain_large",
		"xp_loss_small",
		"double_xp",
		"free_pass",
		"teleport_deploy",
		"special_challenge",
	}

	choiceIdx := r.RNGIntn(len(mysteryOptions))
	choice := mysteryOptions[choiceIdx]

	switch choice {
	case "xp_gain_large":
		p.XP += 40
		return EffectResult{
			EffectType:  "mystery",
			Description: "Mystery Box! Gained +40 XP",
			XPDelta:     40,
			NewXP:       p.XP,
			NewPosition: p.Position,
		}
	case "xp_loss_small":
		oldXP := p.XP
		p.XP -= 15
		if p.XP < 0 {
			p.XP = 0
		}
		loss := oldXP - p.XP
		return EffectResult{
			EffectType:  "mystery",
			Description: fmt.Sprintf("Mystery Box! Minor setback, lost -%d XP", loss),
			XPDelta:     -loss,
			NewXP:       p.XP,
			NewPosition: p.Position,
		}
	case "double_xp":
		p.DoubleXP = true
		return EffectResult{
			EffectType:  "mystery",
			Description: "Mystery Box! Granted Double XP modifier",
			XPDelta:     0,
			NewXP:       p.XP,
			NewPosition: p.Position,
		}
	case "free_pass":
		p.FreePasses++
		return EffectResult{
			EffectType:  "mystery",
			Description: "Mystery Box! Granted a Free Pass",
			XPDelta:     0,
			NewXP:       p.XP,
			NewPosition: p.Position,
		}
	case "teleport_deploy":
		p.Position = 0
		return EffectResult{
			EffectType:  "mystery",
			Description: "Mystery Box! Teleported back to Deploy (GO)",
			XPDelta:     0,
			NewXP:       p.XP,
			NewPosition: 0,
		}
	case "special_challenge":
		p.XP += 30
		return EffectResult{
			EffectType:  "mystery",
			Description: "Mystery Box! Instant Challenge Bonus +30 XP",
			XPDelta:     30,
			NewXP:       p.XP,
			NewPosition: p.Position,
		}
	default:
		p.XP += 10
		return EffectResult{
			EffectType:  "mystery",
			Description: "Mystery Box! Gained +10 XP",
			XPDelta:     10,
			NewXP:       p.XP,
			NewPosition: p.Position,
		}
	}
}

func handleDeploy(r *Room, p *Player, cell BoardCell) EffectResult {
	lapBonus := 100
	p.XP += lapBonus
	return EffectResult{
		EffectType:  "deploy",
		Description: fmt.Sprintf("Landed on Deploy! Collected +%d XP bonus", lapBonus),
		XPDelta:     lapBonus,
		NewXP:       p.XP,
		NewPosition: p.Position,
	}
}

func handleCodeFreeze(r *Room, p *Player, cell BoardCell) EffectResult {
	p.InCodeFreeze = true
	return EffectResult{
		EffectType:  "code_freeze",
		Description: "Landed in Code Freeze! Solve an easy question or wait to leave",
		XPDelta:     0,
		NewXP:       p.XP,
		NewPosition: p.Position,
	}
}

func handleCoffeeBreak(r *Room, p *Player, cell BoardCell) EffectResult {
	p.XP += 5
	return EffectResult{
		EffectType:  "coffee_break",
		Description: "Coffee Break! Rested and gained +5 XP",
		XPDelta:     5,
		NewXP:       p.XP,
		NewPosition: p.Position,
	}
}

func handleDeadline(r *Room, p *Player, cell BoardCell) EffectResult {
	// Swing event: 50% chance +50 XP, 50% chance -20 XP
	if r.RNGIntn(2) == 0 {
		p.XP += 50
		return EffectResult{
			EffectType:  "deadline",
			Description: "Met the Deadline! Big swing bonus +50 XP",
			XPDelta:     50,
			NewXP:       p.XP,
			NewPosition: p.Position,
		}
	}
	oldXP := p.XP
	p.XP -= 20
	if p.XP < 0 {
		p.XP = 0
	}
	loss := oldXP - p.XP
	return EffectResult{
		EffectType:  "deadline",
		Description: fmt.Sprintf("Missed the Deadline! Penalty -%d XP", loss),
		XPDelta:     -loss,
		NewXP:       p.XP,
		NewPosition: p.Position,
	}
}
