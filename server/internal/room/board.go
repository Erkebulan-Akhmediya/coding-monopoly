package room

import "server/internal/locale"

// BoardCell represents a single cell on the 32-cell perimeter board.
type BoardCell struct {
	Index       int            `json:"cell_index"`
	Name        locale.Text    `json:"name"`
	Explanation locale.Text    `json:"explanation"`
	Type        string         `json:"type"`
	Params      map[string]any `json:"params"`
}

func cell(index int, nameEn, nameRu, nameKz, explEn, explRu, explKz, cellType string, params map[string]any) BoardCell {
	if params == nil {
		params = map[string]any{}
	}
	return BoardCell{
		Index:       index,
		Name:        locale.New(nameEn, nameRu, nameKz),
		Explanation: locale.New(explEn, explRu, explKz),
		Type:        cellType,
		Params:      params,
	}
}

// DefaultBoard returns the standard 32-cell board configuration for Coding Monopoly.
// Layout: 4 corners (0: Deploy, 8: Code Freeze, 16: Coffee Break, 24: Deadline)
// and 28 perimeter effect cells.
func DefaultBoard() []BoardCell {
	cells := make([]BoardCell, 32)

	cells[0] = cell(0,
		"Deploy", "Развертывание", "Орнату",
		"Starting cell (GO). Passing this cell during a lap awards +50 XP. Landing directly on it awards a massive +100 XP Deploy bonus!",
		"Стартовая клетка (GO). Прохождение этой клетки во время круга приносит +50 XP. Попадание прямо на неё приносит огромный бонус за размещение в размере +100 XP!",
		"Бастапқы ұяшық (GO). Осы ұяшықтан бір айналым ішінде өту +50 XP береді. Тікелей осы ұяшыққа түсу үлкен +100 XP-лік орналастыру бонусын береді!",
		"deploy", map[string]any{"lap_bonus": 100},
	)
	cells[1] = cell(1,
		"Quick Bugfix", "Быстрое исправление ошибки", "Жылдам ақау түзету",
		"Fixed a minor issue quickly in development. Grants +10 XP.",
		"Быстро устранена небольшая проблема в процессе разработки. Даёт +10 XP.",
		"Дамыту кезінде шағын мәселені тез арада түзетті. +10 XP береді.",
		"xp_gain", map[string]any{"size": "S", "amount": 10},
	)
	cells[2] = cell(2,
		"Syntax Error", "Синтаксическая ошибка", "Синтаксистік қате",
		"A syntax error slipped into the codebase. Deducts 10 XP (XP cannot drop below 0).",
		"В код проникла синтаксическая ошибка. Снимается 10 XP (XP не может опускаться ниже 0).",
		"Код базасына синтаксистік қате енді. 10 XP шегеріледі (XP 0-ден төмен түсе алмайды).",
		"xp_loss", map[string]any{"size": "S", "amount": 10},
	)
	cells[3] = cell(3,
		"Mystery Box", "Таинственная коробка", "Құпия қорап",
		"Mystery Box! Triggers a random effect from the catalog (XP gain/loss, double XP, free pass, teleport, or bonus challenge).",
		"«Таинственная коробка»! Активирует случайный эффект из каталога (получение/потеря очков опыта, удвоение очков опыта, бесплатный проход, телепортация или бонусное испытание).",
		"Құпия қорап! Каталогтан кездейсоқ әсерді іске қосады (XP жинау/жоғалту, екі есе XP, тегін өту, телепортация немесе бонустық сынақ).",
		"mystery", map[string]any{},
	)
	cells[4] = cell(4,
		"Feature Merge", "Слияние функций", "Функцияны біріктіру",
		"Successfully merged a feature pull request into main. Grants +25 XP.",
		"Успешно объединил пул-реквест с функцией в основную ветку. Присвоено +25 XP.",
		"Фичерлік пул-реквестін негізгіге сәтті біріктірді. +25 XP береді.",
		"xp_gain", map[string]any{"size": "M", "amount": 25},
	)
	cells[5] = cell(5,
		"Code Review Pass", "Прохождение рецензии кода", "Кодты тексеруден өту",
		"Thorough code review completed! Doubles the XP gained from your next XP gain.",
		"Тщательный ревью кода завершено! Удваивает количество XP, получаемое при следующем наборе XP.",
		"Кодтың толық тексерісі аяқталды! Келесі XP ұпайыңыз екі есе көбейеді.",
		"double_xp", map[string]any{},
	)
	cells[6] = cell(6,
		"Merge Conflict", "Конфликт слияния", "Біріктіру қақтығысы",
		"Painful merge conflict resolved under stress. Deducts 25 XP.",
		"Конфликт слияния, вызвавший трудности, разрешён в условиях стресса. Снимается 25 XP.",
		"Ауыр біріктіру қақтығысы стресс жағдайында шешілді. 25 XP шегеріледі.",
		"xp_loss", map[string]any{"size": "M", "amount": 25},
	)
	cells[7] = cell(7,
		"CI Pass Ticket", "Прохождение CI", "CI тестілеуден өту билеті",
		"Continuous Integration pipeline passed with flying colors! Grants 1 CI Shield (Free Pass) ticket in your inventory.",
		"Конвейер непрерывной интеграции прошел на «отлично»! Вы получаете 1 билет «CI Shield» (Free Pass) в свой инвентарь.",
		"Үздіксіз интеграция құбыры тамаша нәтижемен өтті! Сіздің инвентарыңызға 1 CI Shield (Free Pass) билеті беріледі.",
		"free_pass", map[string]any{},
	)

	cells[8] = cell(8,
		"Code Freeze", "Заморозка кода", "Кодты тоқтату",
		"Emergency Code Freeze! Your deployment pipeline is locked and you enter a frozen state.",
		"Экстренное замораживание кода! Ваш конвейер развертывания заблокирован, и вы переходите в состояние «замораживания».",
		"Төтенше код мұздату! Сіздің орналастыру құбырыңыз құлыпталды және сіз мұздатылған күйге өтесіз.",
		"code_freeze", map[string]any{},
	)

	cells[9] = cell(9,
		"Refactoring", "Рефакторинг", "Рефакторинг",
		"Cleaned up technical debt and improved code quality. Grants +10 XP.",
		"Устранил технический долг и улучшил качество кода. Дает +10 XP.",
		"Техникалық қарызды жойып, код сапасын жақсартты. +10 XP береді.",
		"xp_gain", map[string]any{"size": "S", "amount": 10},
	)
	cells[10] = cell(10,
		"Meeting Overhead", "Накладные расходы на совещания", "Жиналыс презентациясы",
		"Lengthy sprint planning meeting ran over time. Your next turn is skipped.",
		"Длительное совещание по планированию спринта затянулось. Ваша следующая очередь пропущена.",
		"Ұзақ спринт жоспарлау кездесуі белгіленген уақыттан асып кетті. Сіздің келесі кезегіңіз өткізіліп кетті.",
		"skip_next", map[string]any{},
	)
	cells[11] = cell(11,
		"Major Release", "Крупный релиз", "Негізгі шығарылым",
		"Major version release launched to production successfully! Grants a huge +50 XP bonus.",
		"Выпуск новой версии успешно запущен в производственную среду! Даёт огромный бонус в размере +50 XP.",
		"Негізгі нұсқаның шығарылымы өндіріске сәтті енгізілді! Үлкен +50 XP бонусын береді.",
		"xp_gain", map[string]any{"size": "L", "amount": 50},
	)
	cells[12] = cell(12,
		"Fast-Track Pipeline", "Ускоренный конвейер", "Жылдам жол картасы",
		"Fast-Track Pipeline triggered! Instantly teleports your token straight to Cell 0 (Deploy).",
		"Запущен конвейер Fast-Track! Ваша фишка мгновенно телепортируется прямо в ячейку 0 (Развертывание).",
		"Fast-Track Pipeline іске қосылды! Токеныңызды бірден 0-ші ұяшыққа (Deploy) телепортациялайды.",
		"teleport", map[string]any{"target_position": 0},
	)
	cells[13] = cell(13,
		"Hackathon Bonus", "Бонус за хакатон", "Хакатон бонусы",
		"Built an innovative prototype during the company hackathon! Grants +30 XP.",
		"Создал инновационный прототип во время корпоративного хакатона! Приносит +30 XP.",
		"Компанияның хакерлік марафонында инновациялық прототип жасады! Гранттар +30 XP.",
		"special_challenge", map[string]any{"bonus": 30},
	)
	cells[14] = cell(14,
		"Memory Leak", "Утечка памяти", "Жады ағып кетуі",
		"A memory leak caused high resource consumption. Deducts 10 XP.",
		"Утечка памяти привела к высокому потреблению ресурсов. Снимается 10 XP.",
		"Жадтың ағып кетуі ресурстардың көп тұтынылуына әкелді. 10 XP шегеріледі.",
		"xp_loss", map[string]any{"size": "S", "amount": 10},
	)
	cells[15] = cell(15,
		"Wildcard Event", "Непредвиденное событие", "Жабайы карта оқиғасы",
		"Wildcard Event! Draws a random mystery effect from across the board.",
		"Событие «Джокер»! Выбирает случайный загадочный эффект в любой точке игрового поля.",
		"Жабайы оқиға! Тақтаның кез келген жерінен кездейсоқ жұмбақ әсерді тартады.",
		"mystery", map[string]any{},
	)

	cells[16] = cell(16,
		"Coffee Break", "Перерыв на кофе", "Кофе үзілісі",
		"Neutral rest zone. Enjoy a quick coffee break and relax while earning a +5 XP boost.",
		"Нейтральная зона отдыха. Сделайте небольшой перерыв на кофе и расслабьтесь, получив при этом бонус +5 XP.",
		"Бейтарап демалыс аймағы. Жылдам кофе үзілісін пайдаланып демалып, +5 XP қосымша ұпай жинаңыз.",
		"coffee_break", map[string]any{},
	)

	cells[17] = cell(17,
		"Performance Tuning", "Настройка производительности", "Өнімділікті баптау",
		"Optimized database queries and latency. Grants +25 XP.",
		"Оптимизированные запросы к базе данных и задержка. Дает +25 XP.",
		"Оптимизацияланған дерекқор сұраулары мен кешігу уақыты. +25 XP береді.",
		"xp_gain", map[string]any{"size": "M", "amount": 25},
	)
	cells[18] = cell(18,
		"Pair Programming", "Парное программирование", "Жұп бағдарламалау",
		"Collaborative pair programming session! Multiplies your next XP gain by 2x.",
		"Сеанс совместного парного программирования! Удваивает количество XP, которое вы получите в следующий раз.",
		"Командалық жұптық бағдарламалау сессиясы! Келесі XP ұпайыңызды 2 есе көбейтеді.",
		"double_xp", map[string]any{},
	)
	cells[19] = cell(19,
		"Failed Build", "Сбой сборки", "Сәтсіз құрастыру",
		"Broken unit test failed the staging build. Deducts 25 XP.",
		"Неисправный модульный тест привёл к сбою сборки на промежуточном этапе. Снято 25 XP.",
		"Бұзылған unit-тест staging құрастыру кезінде сәтсіз аяқталды. 25 XP шегеріледі.",
		"xp_loss", map[string]any{"size": "M", "amount": 25},
	)
	cells[20] = cell(20,
		"Express Route", "Экспресс-маршрут", "Экспресс маршрут",
		"Express Route shortcut! Instantly teleports your token to Cell 16 (Coffee Break).",
		"Ярлык «Экспресс-маршрут»! Мгновенно телепортирует вашу фишку в ячейку 16 («Кофе-брейк»).",
		"Экспресс маршрутының қысқартуы! Фишкаңызды бірден 16-шы ұяшыққа (Кофе үзілісі) телепортациялайды.",
		"teleport", map[string]any{"target_position": 16},
	)
	cells[21] = cell(21,
		"Documentation Boost", "Улучшение документации", "Құжаттаманы жетілдіру",
		"Wrote clear API documentation for the team. Grants +10 XP.",
		"Написал понятную документацию по API для команды. Даёт +10 XP.",
		"Команда үшін айқын API құжаттамасын жазды. +10 XP береді.",
		"xp_gain", map[string]any{"size": "S", "amount": 10},
	)
	cells[22] = cell(22,
		"Security Clearance", "Допуск к секретной информации", "Қауіпсіздік рұқсаты",
		"Security audit cleared without findings! Grants 1 CI Shield (Free Pass) ticket in your inventory.",
		"Аудит безопасности пройдён без замечаний! Вы получаете 1 билет «CI Shield (Free Pass)» в свой инвентарь.",
		"Қауіпсіздік аудиті ешқандай бұзушылықтарсыз өтті! Сіздің инвентарыңызға 1 CI Shield (Free Pass) билеті қосылады.",
		"free_pass", map[string]any{},
	)
	cells[23] = cell(23,
		"Prod Outage Duty", "Дежурство при сбое в производственной среде", "Prod-тағы ақау кезекшілігі",
		"Called on-call for an outage investigation. Your next turn is skipped.",
		"Вас вызвали на дежурство для расследования сбоя в работе системы. Ваш следующий ход пропускается.",
		"Ақауды тергеу үшін шақырылды. Келесі кезегіңіз өткізіліп жіберіледі.",
		"skip_next", map[string]any{},
	)

	cells[24] = cell(24,
		"Deadline", "Крайний срок", "Соңғы мерзім",
		"Critical release deadline swing event! 50% chance of +50 XP bonus, 50% chance of -20 XP penalty.",
		"Критическое событие, влияющее на срок выпуска! 50 % шанс получить бонус +50 XP, 50 % шанс получить штраф -20 XP.",
		"Критикалық шығарылым мерзімінің ауытқу оқиғасы! +50 XP бонусын алу ықтималдығы 50%, -20 XP айыппұл алу ықтималдығы 50%.",
		"deadline", map[string]any{},
	)

	cells[25] = cell(25,
		"Architecture Upgrade", "Модернизация архитектуры", "Архитектураны жаңарту",
		"Migrated legacy monolith to a modern resilient architecture. Grants +50 XP.",
		"Перенес устаревшую монолитную архитектуру на современную отказоустойчивую архитектуру. Даёт +50 XP.",
		"Ескі мұрагерлік монолитті заманауи, төзімді архитектураға көшірді. +50 XP береді.",
		"xp_gain", map[string]any{"size": "L", "amount": 50},
	)
	cells[26] = cell(26,
		"Bug Bounty", "Программа поощрения за обнаружение ошибок", "Қателіктерге сыйақы",
		"Discovered and patched a critical vulnerability! Grants +35 XP.",
		"Обнаружена и устранена критическая уязвимость! Начисляется +35 XP.",
		"Маңызды қауіпсіздік олқылығын анықтап, жөндедім! +35 тәжірибе ұпайы беріледі.",
		"special_challenge", map[string]any{"bonus": 35},
	)
	cells[27] = cell(27,
		"Dependency Hell", "«Ад зависимостей»", "Тәуелділік тозағы",
		"Incompatible package versions caused build failure. Deducts 10 XP.",
		"Несовместимые версии пакетов привели к сбою сборки. Снимается 10 XP.",
		"Сәйкес келмейтін пакет нұсқалары құрастырудың сәтсіздігіне әкелді. 10 XP шегеріледі.",
		"xp_loss", map[string]any{"size": "S", "amount": 10},
	)
	cells[28] = cell(28,
		"Surprise Audit", "Неожиданный аудит", "Күтпеген аудит",
		"Surprise compliance audit! Triggers a random mystery event effect.",
		"Неожиданная проверка на соответствие требованиям! Вызывает случайный эффект загадочного события.",
		"Күтпеген сәйкестік аудиті! Кездейсоқ жұмбақ оқиға әсерін іске қосады.",
		"mystery", map[string]any{},
	)
	cells[29] = cell(29,
		"Test Coverage 100%", "100% тестовое покрытие", "Тест қамтуы 100%",
		"Achieved 100% test coverage with comprehensive test suites. Grants +25 XP.",
		"Достигнут 100% охват тестированием с помощью комплексных наборов тестов. Начисляется +25 XP.",
		"Кешенді тест жиынтықтарымен 100% тест қамтуына қол жеткізілді. +25 XP береді.",
		"xp_gain", map[string]any{"size": "M", "amount": 25},
	)
	cells[30] = cell(30,
		"Hotfix Shift", "Смена при выпуске исправления", "Hotfix ауысымы",
		"Emergency hotfix required! Instantly teleports your token to Cell 8 (Code Freeze).",
		"Требуется экстренное исправление! Мгновенно телепортирует вашу фишку в ячейку 8 («Заморозка кода»).",
		"Төтенше хотфикс қажет! Токеныңызды бірден 8-ші ұяшыққа (Кодты мұздату) телепортациялайды.",
		"teleport", map[string]any{"target_position": 8},
	)
	cells[31] = cell(31,
		"Linter Pass", "Успешная проверка линтером", "Линтер өтуі",
		"Passed strict automated linter and static analysis checks. Grants +10 XP.",
		"Прошло строгие проверки с помощью автоматизированного линтера и статического анализа. Дает +10 XP.",
		"Қатаң автоматтандырылған линтер мен статикалық талдау тексерістерінен сәтті өтті. +10 XP береді.",
		"xp_gain", map[string]any{"size": "S", "amount": 10},
	)

	return cells
}
