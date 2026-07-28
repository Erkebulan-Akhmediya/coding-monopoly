package room

// BoardCell represents a single cell on the 32-cell perimeter board.
type BoardCell struct {
	Index  int            `json:"cell_index"`
	Name   string         `json:"name"`
	Type   string         `json:"type"`
	Params map[string]any `json:"params"`
}

// DefaultBoard returns the standard 32-cell board configuration for Coding Monopoly.
// Layout: 4 corners (0: Deploy, 8: Code Freeze, 16: Coffee Break, 24: Deadline)
// and 28 perimeter effect cells.
func DefaultBoard() []BoardCell {
	cells := make([]BoardCell, 32)

	// Corner 0: Deploy (Start / GO)
	cells[0] = BoardCell{
		Index:  0,
		Name:   "Deploy",
		Type:   "deploy",
		Params: map[string]any{"lap_bonus": 100},
	}

	// Side 1: Cells 1 to 7
	cells[1] = BoardCell{Index: 1, Name: "Quick Bugfix", Type: "xp_gain", Params: map[string]any{"size": "S", "amount": 10}}
	cells[2] = BoardCell{Index: 2, Name: "Syntax Error", Type: "xp_loss", Params: map[string]any{"size": "S", "amount": 10}}
	cells[3] = BoardCell{Index: 3, Name: "Mystery Box", Type: "mystery", Params: map[string]any{}}
	cells[4] = BoardCell{Index: 4, Name: "Feature Merge", Type: "xp_gain", Params: map[string]any{"size": "M", "amount": 25}}
	cells[5] = BoardCell{Index: 5, Name: "Code Review Pass", Type: "double_xp", Params: map[string]any{}}
	cells[6] = BoardCell{Index: 6, Name: "Merge Conflict", Type: "xp_loss", Params: map[string]any{"size": "M", "amount": 25}}
	cells[7] = BoardCell{Index: 7, Name: "CI Pass Ticket", Type: "free_pass", Params: map[string]any{}}

	// Corner 1: Code Freeze (Jail equivalent)
	cells[8] = BoardCell{
		Index:  8,
		Name:   "Code Freeze",
		Type:   "code_freeze",
		Params: map[string]any{},
	}

	// Side 2: Cells 9 to 15
	cells[9] = BoardCell{Index: 9, Name: "Refactoring", Type: "xp_gain", Params: map[string]any{"size": "S", "amount": 10}}
	cells[10] = BoardCell{Index: 10, Name: "Meeting Overhead", Type: "skip_next", Params: map[string]any{}}
	cells[11] = BoardCell{Index: 11, Name: "Major Release", Type: "xp_gain", Params: map[string]any{"size": "L", "amount": 50}}
	cells[12] = BoardCell{Index: 12, Name: "Fast-Track Pipeline", Type: "teleport", Params: map[string]any{"target_position": 0}}
	cells[13] = BoardCell{Index: 13, Name: "Hackathon Bonus", Type: "special_challenge", Params: map[string]any{"bonus": 30}}
	cells[14] = BoardCell{Index: 14, Name: "Memory Leak", Type: "xp_loss", Params: map[string]any{"size": "S", "amount": 10}}
	cells[15] = BoardCell{Index: 15, Name: "Wildcard Event", Type: "mystery", Params: map[string]any{}}

	// Corner 2: Coffee Break (Rest cell)
	cells[16] = BoardCell{
		Index:  16,
		Name:   "Coffee Break",
		Type:   "coffee_break",
		Params: map[string]any{},
	}

	// Side 3: Cells 17 to 23
	cells[17] = BoardCell{Index: 17, Name: "Performance Tuning", Type: "xp_gain", Params: map[string]any{"size": "M", "amount": 25}}
	cells[18] = BoardCell{Index: 18, Name: "Pair Programming", Type: "double_xp", Params: map[string]any{}}
	cells[19] = BoardCell{Index: 19, Name: "Failed Build", Type: "xp_loss", Params: map[string]any{"size": "M", "amount": 25}}
	cells[20] = BoardCell{Index: 20, Name: "Express Route", Type: "teleport", Params: map[string]any{"target_position": 16}}
	cells[21] = BoardCell{Index: 21, Name: "Documentation Boost", Type: "xp_gain", Params: map[string]any{"size": "S", "amount": 10}}
	cells[22] = BoardCell{Index: 22, Name: "Security Clearance", Type: "free_pass", Params: map[string]any{}}
	cells[23] = BoardCell{Index: 23, Name: "Prod Outage Duty", Type: "skip_next", Params: map[string]any{}}

	// Corner 3: Deadline (Swing event)
	cells[24] = BoardCell{
		Index:  24,
		Name:   "Deadline",
		Type:   "deadline",
		Params: map[string]any{},
	}

	// Side 4: Cells 25 to 31
	cells[25] = BoardCell{Index: 25, Name: "Architecture Upgrade", Type: "xp_gain", Params: map[string]any{"size": "L", "amount": 50}}
	cells[26] = BoardCell{Index: 26, Name: "Bug Bounty", Type: "special_challenge", Params: map[string]any{"bonus": 35}}
	cells[27] = BoardCell{Index: 27, Name: "Dependency Hell", Type: "xp_loss", Params: map[string]any{"size": "S", "amount": 10}}
	cells[28] = BoardCell{Index: 28, Name: "Surprise Audit", Type: "mystery", Params: map[string]any{}}
	cells[29] = BoardCell{Index: 29, Name: "Test Coverage 100%", Type: "xp_gain", Params: map[string]any{"size": "M", "amount": 25}}
	cells[30] = BoardCell{Index: 30, Name: "Hotfix Shift", Type: "teleport", Params: map[string]any{"target_position": 8}}
	cells[31] = BoardCell{Index: 31, Name: "Linter Pass", Type: "xp_gain", Params: map[string]any{"size": "S", "amount": 10}}

	return cells
}
