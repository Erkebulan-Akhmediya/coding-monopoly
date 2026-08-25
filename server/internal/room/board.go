package room

// BoardCell represents a single cell on the 32-cell perimeter board.
type BoardCell struct {
	Index  int            `json:"cell_index"`
	NameEn string         `json:"name_en"`
	NameRu string         `json:"name_ru"`
	NameKz string         `json:"name_kz"`
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
		NameEn: "Deploy",
		NameRu: "Развернуть",
		NameKz: "Орнату",
		Type:   "deploy",
		Params: map[string]any{"lap_bonus": 100},
	}

	// Side 1: Cells 1 to 7
	cells[1] = BoardCell{Index: 1, NameEn: "Quick Bugfix", NameRu: "Быстрое исправление ошибки", NameKz: "Жылдам ақау түзету", Type: "xp_gain", Params: map[string]any{"size": "S", "amount": 10}}
	cells[2] = BoardCell{Index: 2, NameEn: "Syntax Error", NameRu: "Ошибка синтаксиса", NameKz: "Синтаксистік қате", Type: "xp_loss", Params: map[string]any{"size": "S", "amount": 10}}
	cells[3] = BoardCell{Index: 3, NameEn: "Mystery Box", NameRu: "«Загадочная коробка»", NameKz: "Құпия қорап", Type: "mystery", Params: map[string]any{}}
	cells[4] = BoardCell{Index: 4, NameEn: "Feature Merge", NameRu: "Объединение функций", NameKz: "Функцияны біріктіру", Type: "xp_gain", Params: map[string]any{"size": "M", "amount": 25}}
	cells[5] = BoardCell{Index: 5, NameEn: "Code Review Pass", NameRu: "Прохождение проверки кода", NameKz: "Кодты тексеру сәтті өтті", Type: "double_xp", Params: map[string]any{}}
	cells[6] = BoardCell{Index: 6, NameEn: "Merge Conflict", NameRu: "Конфликт маршрутов", NameKz: "Қиылысу қақтығысын біріктіру", Type: "xp_loss", Params: map[string]any{"size": "M", "amount": 25}}
	cells[7] = BoardCell{Index: 7, NameEn: "CI Pass Ticket", NameRu: "Билет CI Pass", NameKz: "CI өткізу билеті", Type: "free_pass", Params: map[string]any{}}

	// Corner 1: Code Freeze (Jail equivalent)
	cells[8] = BoardCell{
		Index:  8,
		NameEn: "Code Freeze",
		NameRu: "Заморозка кода",
		NameKz: "Кодты тоңдыру",
		Type:   "code_freeze",
		Params: map[string]any{},
	}

	// Side 2: Cells 9 to 15
	cells[9] = BoardCell{Index: 9, NameEn: "Refactoring", NameRu: "Рефакторинг", NameKz: "Қайта құру", Type: "xp_gain", Params: map[string]any{"size": "S", "amount": 10}}
	cells[10] = BoardCell{Index: 10, NameEn: "Meeting Overhead", NameRu: "Расходы на проведение встречи", NameKz: "Жыйынның тақта жазбасы", Type: "skip_next", Params: map[string]any{}}
	cells[11] = BoardCell{Index: 11, NameEn: "Major Release", NameRu: "Крупное обновление", NameKz: "Негізгі шығарылым", Type: "xp_gain", Params: map[string]any{"size": "L", "amount": 50}}
	cells[12] = BoardCell{Index: 12, NameEn: "Fast-Track Pipeline", NameRu: "Ускоренная процедура рассмотрения", NameKz: "Жылдам жолмен құбыр желісі", Type: "teleport", Params: map[string]any{"target_position": 0}}
	cells[13] = BoardCell{Index: 13, NameEn: "Hackathon Bonus", NameRu: "Бонус хакатона", NameKz: "Хакатон бонусы", Type: "special_challenge", Params: map[string]any{"bonus": 30}}
	cells[14] = BoardCell{Index: 14, NameEn: "Memory Leak", NameRu: "Утечка памяти", NameKz: "Естеліктегі ақау", Type: "xp_loss", Params: map[string]any{"size": "S", "amount": 10}}
	cells[15] = BoardCell{Index: 15, NameEn: "Wildcard Event", NameRu: "Мероприятие с участием «Wildcard»", NameKz: "Жабайы карта оқиғасы", Type: "mystery", Params: map[string]any{}}

	// Corner 2: Coffee Break (Rest cell)
	cells[16] = BoardCell{
		Index:  16,
		NameEn: "Coffee Break",
		NameRu: "Кофе-брейк",
		NameKz: "Кофе үзілісі",
		Type:   "coffee_break",
		Params: map[string]any{},
	}

	// Side 3: Cells 17 to 23
	cells[17] = BoardCell{Index: 17, NameEn: "Performance Tuning", NameRu: "Оптимизация производительности", NameKz: "Өнімділікті баптау", Type: "xp_gain", Params: map[string]any{"size": "M", "amount": 25}}
	cells[18] = BoardCell{Index: 18, NameEn: "Pair Programming", NameRu: "Парное программирование", NameKz: "Жұп бағдарламалау", Type: "double_xp", Params: map[string]any{}}
	cells[19] = BoardCell{Index: 19, NameEn: "Failed Build", NameRu: "Сбой сборки", NameKz: "Құрастыру сәтсіз аяқталды", Type: "xp_loss", Params: map[string]any{"size": "M", "amount": 25}}
	cells[20] = BoardCell{Index: 20, NameEn: "Express Route", NameRu: "Экспресс-маршрут", NameKz: "Экспресс маршруты", Type: "teleport", Params: map[string]any{"target_position": 16}}
	cells[21] = BoardCell{Index: 21, NameEn: "Documentation Boost", NameRu: "Расширение документации", NameKz: "Құжаттаманы жақсарту", Type: "xp_gain", Params: map[string]any{"size": "S", "amount": 10}}
	cells[22] = BoardCell{Index: 22, NameEn: "Security Clearance", NameRu: "Допуск к секретной информации", NameKz: "Қауіпсіздік рұқсаты", Type: "free_pass", Params: map[string]any{}}
	cells[23] = BoardCell{Index: 23, NameEn: "Prod Outage Duty", NameRu: "Дежурство при сбое в производстве", NameKz: "Өнімділіктің үзіліс кезеңіндегі міндет", Type: "skip_next", Params: map[string]any{}}

	// Corner 3: Deadline (Swing event)
	cells[24] = BoardCell{
		Index:  24,
		NameEn: "Deadline",
		NameRu: "Крайний срок",
		NameKz: "Соңғы мерзім",
		Type:   "deadline",
		Params: map[string]any{},
	}

	// Side 4: Cells 25 to 31
	cells[25] = BoardCell{Index: 25, NameEn: "Architecture Upgrade", NameRu: "Модернизация архитектуры", NameKz: "Архитектураны жаңарту", Type: "xp_gain", Params: map[string]any{"size": "L", "amount": 50}}
	cells[26] = BoardCell{Index: 26, NameEn: "Bug Bounty", NameRu: "Программа поощрения за обнаружение ошибок", NameKz: "Қателіктерге сыйақы", Type: "special_challenge", Params: map[string]any{"bonus": 35}}
	cells[27] = BoardCell{Index: 27, NameEn: "Dependency Hell", NameRu: "«Ад зависимостей»", NameKz: "Тәуелділік тозағы", Type: "xp_loss", Params: map[string]any{"size": "S", "amount": 10}}
	cells[28] = BoardCell{Index: 28, NameEn: "Surprise Audit", NameRu: "Незапланированная проверка", NameKz: "Күтпеген тексеріс", Type: "mystery", Params: map[string]any{}}
	cells[29] = BoardCell{Index: 29, NameEn: "Test Coverage 100%", NameRu: "Покрытие тестами — 100 %", NameKz: "Сынақ қамтуы 100%", Type: "xp_gain", Params: map[string]any{"size": "M", "amount": 25}}
	cells[30] = BoardCell{Index: 30, NameEn: "Hotfix Shift", NameRu: "Исправление «Shift»", NameKz: "Хотфикс Шифт", Type: "teleport", Params: map[string]any{"target_position": 8}}
	cells[31] = BoardCell{Index: 31, NameEn: "Linter Pass", NameRu: "Перевал Линтер", NameKz: "Линтер асуы", Type: "xp_gain", Params: map[string]any{"size": "S", "amount": 10}}

	return cells
}
