package room

import "server/internal/locale"

// BoardCell represents a single cell on the 32-cell perimeter board.
type BoardCell struct {
	Index  int            `json:"cell_index"`
	Name   locale.Text    `json:"name"`
	Type   string         `json:"type"`
	Params map[string]any `json:"params"`
}

func cell(index int, nameEn, nameRu, nameKz, cellType string, params map[string]any) BoardCell {
	if params == nil {
		params = map[string]any{}
	}
	return BoardCell{
		Index:  index,
		Name:   locale.New(nameEn, nameRu, nameKz),
		Type:   cellType,
		Params: params,
	}
}

// DefaultBoard returns the standard 32-cell board configuration for Coding Monopoly.
// Layout: 4 corners (0: Deploy, 8: Code Freeze, 16: Coffee Break, 24: Deadline)
// and 28 perimeter effect cells.
func DefaultBoard() []BoardCell {
	cells := make([]BoardCell, 32)

	cells[0] = cell(0, "Deploy", "Развертывание", "Орнату", "deploy", map[string]any{"lap_bonus": 100})
	cells[1] = cell(1, "Quick Bugfix", "Быстрое исправление ошибки", "Жылдам ақау түзету", "xp_gain", map[string]any{"size": "S", "amount": 10})
	cells[2] = cell(2, "Syntax Error", "Синтаксическая ошибка", "Синтаксистік қате", "xp_loss", map[string]any{"size": "S", "amount": 10})
	cells[3] = cell(3, "Mystery Box", "Таинственная коробка", "Құпия қорап", "mystery", map[string]any{})
	cells[4] = cell(4, "Feature Merge", "Слияние функций", "Функцияны біріктіру", "xp_gain", map[string]any{"size": "M", "amount": 25})
	cells[5] = cell(5, "Code Review Pass", "Прохождение рецензии кода", "Кодты тексеруден өту", "double_xp", map[string]any{})
	cells[6] = cell(6, "Merge Conflict", "Конфликт слияния", "Біріктіру қақтығысы", "xp_loss", map[string]any{"size": "M", "amount": 25})
	cells[7] = cell(7, "CI Pass Ticket", "Прохождение CI", "CI тестілеуден өту билеті", "free_pass", map[string]any{})

	cells[8] = cell(8, "Code Freeze", "Заморозка кода", "Кодты тоқтату", "code_freeze", map[string]any{})

	cells[9] = cell(9, "Refactoring", "Рефакторинг", "Рефакторинг", "xp_gain", map[string]any{"size": "S", "amount": 10})
	cells[10] = cell(10, "Meeting Overhead", "Накладные расходы на совещания", "Жиналыс презентациясы", "skip_next", map[string]any{})
	cells[11] = cell(11, "Major Release", "Крупный релиз", "Негізгі шығарылым", "xp_gain", map[string]any{"size": "L", "amount": 50})
	cells[12] = cell(12, "Fast-Track Pipeline", "Ускоренный конвейер", "Жылдам жол картасы", "teleport", map[string]any{"target_position": 0})
	cells[13] = cell(13, "Hackathon Bonus", "Бонус за хакатон", "Хакатон бонусы", "special_challenge", map[string]any{"bonus": 30})
	cells[14] = cell(14, "Memory Leak", "Утечка памяти", "Жады ағып кетуі", "xp_loss", map[string]any{"size": "S", "amount": 10})
	cells[15] = cell(15, "Wildcard Event", "Непредвиденное событие", "Жабайы карта оқиғасы", "mystery", map[string]any{})

	cells[16] = cell(16, "Coffee Break", "Перерыв на кофе", "Кофе үзілісі", "coffee_break", map[string]any{})

	cells[17] = cell(17, "Performance Tuning", "Настройка производительности", "Өнімділікті баптау", "xp_gain", map[string]any{"size": "M", "amount": 25})
	cells[18] = cell(18, "Pair Programming", "Парное программирование", "Жұп бағдарламалау", "double_xp", map[string]any{})
	cells[19] = cell(19, "Failed Build", "Сбой сборки", "Сәтсіз құрастыру", "xp_loss", map[string]any{"size": "M", "amount": 25})
	cells[20] = cell(20, "Express Route", "Экспресс-маршрут", "Экспресс маршрут", "teleport", map[string]any{"target_position": 16})
	cells[21] = cell(21, "Documentation Boost", "Улучшение документации", "Құжаттаманы жетілдіру", "xp_gain", map[string]any{"size": "S", "amount": 10})
	cells[22] = cell(22, "Security Clearance", "Допуск к секретной информации", "Қауіпсіздік рұқсаты", "free_pass", map[string]any{})
	cells[23] = cell(23, "Prod Outage Duty", "Дежурство при сбое в производственной среде", "Prod-тағы ақау кезекшілігі", "skip_next", map[string]any{})

	cells[24] = cell(24, "Deadline", "Крайний срок", "Соңғы мерзім", "deadline", map[string]any{})

	cells[25] = cell(25, "Architecture Upgrade", "Модернизация архитектуры", "Архитектураны жаңарту", "xp_gain", map[string]any{"size": "L", "amount": 50})
	cells[26] = cell(26, "Bug Bounty", "Программа поощрения за обнаружение ошибок", "Қателіктерге сыйақы", "special_challenge", map[string]any{"bonus": 35})
	cells[27] = cell(27, "Dependency Hell", "«Ад зависимостей»", "Тәуелділік тозағы", "xp_loss", map[string]any{"size": "S", "amount": 10})
	cells[28] = cell(28, "Surprise Audit", "Неожиданный аудит", "Күтпеген аудит", "mystery", map[string]any{})
	cells[29] = cell(29, "Test Coverage 100%", "100% тестовое покрытие", "Тест қамтуы 100%", "xp_gain", map[string]any{"size": "M", "amount": 25})
	cells[30] = cell(30, "Hotfix Shift", "Смена при выпуске исправления", "Hotfix ауысымы", "teleport", map[string]any{"target_position": 8})
	cells[31] = cell(31, "Linter Pass", "Успешная проверка линтером", "Линтер өтуі", "xp_gain", map[string]any{"size": "S", "amount": 10})

	return cells
}
