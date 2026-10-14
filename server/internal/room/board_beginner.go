package room

import "maps"

// BeginnerBoard returns a gentler 32-cell board for newcomers: mostly XP
// gains and bonuses, only two small XP losses, a single forward-only
// teleport, and no skip-turn, Code Freeze, or Deadline cells.
//
// Layout: 4 corners (0: Deploy, 8: Coffee Break, 16: Hackathon, 24: Demo Day)
// and 28 perimeter cells.
func BeginnerBoard() []BoardCell {
	standard := DefaultBoard()
	// reuse copies a default-board cell (same effect and texts) to index.
	reuse := func(index int, from BoardCell) BoardCell {
		from.Index = index
		from.Params = maps.Clone(from.Params)
		return from
	}

	cells := make([]BoardCell, BoardSize)

	cells[0] = reuse(0, standard[0])
	cells[1] = cell(1,
		"Hello World", "Hello World", "Hello World",
		"Your first program runs! Grants +10 XP.",
		"Ваша первая программа заработала! Даёт +10 XP.",
		"Алғашқы бағдарламаңыз іске қосылды! +10 XP береді.",
		"xp_gain", map[string]any{"amount": 10},
	)
	cells[2] = cell(2,
		"First Commit", "Первый коммит", "Алғашқы коммит",
		"Saved your work with a first commit. Grants +10 XP.",
		"Вы сохранили работу первым коммитом. Даёт +10 XP.",
		"Жұмысыңызды алғашқы коммитпен сақтадыңыз. +10 XP береді.",
		"xp_gain", map[string]any{"amount": 10},
	)
	cells[3] = cell(3,
		"Pair Programming", "Парное программирование", "Жұп бағдарламалау",
		"Coding together with a friend! Your next XP gain is doubled.",
		"Программируете вместе с другом! Следующее начисление XP удваивается.",
		"Досыңызбен бірге код жазып жатырсыз! Келесі XP ұпайыңыз екі есе көбейеді.",
		"double_xp", map[string]any{},
	)
	cells[4] = cell(4,
		"Unit Test Passed", "Модульный тест пройден", "Unit-тест өтті",
		"Your first unit test is green. Grants +15 XP.",
		"Ваш первый модульный тест прошёл успешно. Даёт +15 XP.",
		"Алғашқы unit-тестіңіз сәтті өтті. +15 XP береді.",
		"xp_gain", map[string]any{"amount": 15},
	)
	cells[5] = reuse(5, standard[7]) // CI Pass Ticket (free pass)
	cells[6] = cell(6,
		"Typo", "Опечатка", "Қате теру",
		"A small typo broke the build. Deducts 10 XP (XP cannot drop below 0).",
		"Небольшая опечатка сломала сборку. Снимается 10 XP (XP не может опускаться ниже 0).",
		"Кішкене қате теру құрастыруды бұзды. 10 XP шегеріледі (XP 0-ден төмен түсе алмайды).",
		"xp_loss", map[string]any{"amount": 10},
	)
	cells[7] = cell(7,
		"Code Review Approved", "Код-ревью одобрено", "Код тексеруі мақұлданды",
		"A teammate approved your pull request. Grants +20 XP.",
		"Коллега одобрил ваш пул-реквест. Даёт +20 XP.",
		"Әріптесіңіз пул-реквестіңізді мақұлдады. +20 XP береді.",
		"xp_gain", map[string]any{"amount": 20},
	)

	cells[8] = cell(8,
		"Coffee Break", "Перерыв на кофе", "Кофе үзілісі",
		"Rest corner. Enjoy a coffee and earn +10 XP.",
		"Угол отдыха. Выпейте кофе и получите +10 XP.",
		"Демалыс бұрышы. Кофе ішіп, +10 XP алыңыз.",
		"coffee_break", map[string]any{"amount": 10},
	)

	cells[9] = cell(9,
		"Clean Code", "Чистый код", "Таза код",
		"Kept your code tidy and simple. Grants +10 XP.",
		"Ваш код аккуратный и простой. Даёт +10 XP.",
		"Кодыңызды ұқыпты әрі қарапайым ұстадыңыз. +10 XP береді.",
		"xp_gain", map[string]any{"amount": 10},
	)
	cells[10] = cell(10,
		"Coding Kata", "Ката по программированию", "Бағдарламалау катасы",
		"Solved a short practice exercise! Grants +25 bonus XP (doubled with Double XP).",
		"Решили короткое тренировочное упражнение! Даёт +25 бонусных XP (вдвое больше с Двойным XP).",
		"Қысқа жаттығу есебін шештіңіз! +25 бонус XP береді (Қос XP болса, екі есе көп).",
		"special_challenge", map[string]any{"bonus": 25},
	)
	cells[11] = cell(11,
		"Good Variable Names", "Понятные имена", "Түсінікті атаулар",
		"Chose clear, descriptive names. Grants +15 XP.",
		"Выбрали понятные, говорящие имена. Даёт +15 XP.",
		"Анық әрі мағыналы атаулар таңдадыңыз. +15 XP береді.",
		"xp_gain", map[string]any{"amount": 15},
	)
	cells[12] = reuse(12, standard[3]) // Mystery Box
	cells[13] = cell(13,
		"Feature Shipped", "Функция выпущена", "Функция шығарылды",
		"Your first feature reached users. Grants +20 XP.",
		"Ваша первая функция дошла до пользователей. Даёт +20 XP.",
		"Алғашқы функцияңыз пайдаланушыларға жетті. +20 XP береді.",
		"xp_gain", map[string]any{"amount": 20},
	)
	cells[14] = cell(14,
		"Mentor Session", "Встреча с наставником", "Тәлімгермен кездесу",
		"A mentor shared great tips! Your next XP gain is doubled.",
		"Наставник поделился полезными советами! Следующее начисление XP удваивается.",
		"Тәлімгер пайдалы кеңестер берді! Келесі XP ұпайыңыз екі есе көбейеді.",
		"double_xp", map[string]any{},
	)
	cells[15] = cell(15,
		"Docs Updated", "Документация обновлена", "Құжаттама жаңартылды",
		"Wrote a helpful README. Grants +10 XP.",
		"Написали полезный README. Даёт +10 XP.",
		"Пайдалы README жаздыңыз. +10 XP береді.",
		"xp_gain", map[string]any{"amount": 10},
	)

	cells[16] = cell(16,
		"Hackathon", "Хакатон", "Хакатон",
		"Built a cool prototype at the hackathon! Grants +40 bonus XP (doubled with Double XP).",
		"Создали классный прототип на хакатоне! Даёт +40 бонусных XP (вдвое больше с Двойным XP).",
		"Хакатонда керемет прототип жасадыңыз! +40 бонус XP береді (Қос XP болса, екі есе көп).",
		"special_challenge", map[string]any{"bonus": 40},
	)

	cells[17] = cell(17,
		"Bug Squashed", "Баг исправлен", "Қате түзетілді",
		"Found and fixed a bug. Grants +15 XP.",
		"Нашли и исправили баг. Даёт +15 XP.",
		"Қатені тауып, түзеттіңіз. +15 XP береді.",
		"xp_gain", map[string]any{"amount": 15},
	)
	cells[18] = cell(18,
		"Green Build", "Зелёная сборка", "Жасыл құрастыру",
		"All checks passed! Grants 1 CI Shield (Free Pass) ticket in your inventory.",
		"Все проверки пройдены! Вы получаете 1 билет «CI Shield» (Free Pass) в свой инвентарь.",
		"Барлық тексерулер өтті! Инвентарыңызға 1 CI Shield (Free Pass) билеті беріледі.",
		"free_pass", map[string]any{},
	)
	cells[19] = cell(19,
		"Refactoring Done", "Рефакторинг завершён", "Рефакторинг аяқталды",
		"Made the code easier to read and change. Grants +25 XP.",
		"Сделали код проще для чтения и изменений. Даёт +25 XP.",
		"Кодты оқуға және өзгертуге жеңілдеттіңіз. +25 XP береді.",
		"xp_gain", map[string]any{"amount": 25},
	)
	cells[20] = cell(20,
		"Shortcut", "Короткий путь", "Қысқа жол",
		"Found a shortcut! Instantly moves your token forward to Cell 28 (Release Candidate).",
		"Нашли короткий путь! Ваша фишка мгновенно перемещается вперёд на клетку 28 («Релиз-кандидат»).",
		"Қысқа жол таптыңыз! Фишкаңыз бірден алға, 28-ші ұяшыққа (Релиз-кандидат) жылжиды.",
		"teleport", map[string]any{"target_position": 28},
	)
	cells[21] = cell(21,
		"Small Fix", "Небольшое исправление", "Шағын түзету",
		"Fixed a small issue quickly. Grants +10 XP.",
		"Быстро исправили небольшую проблему. Даёт +10 XP.",
		"Шағын мәселені тез түзеттіңіз. +10 XP береді.",
		"xp_gain", map[string]any{"amount": 10},
	)
	cells[22] = cell(22,
		"Flaky Test", "Нестабильный тест", "Тұрақсыз тест",
		"A flaky test failed at random. Deducts 10 XP (XP cannot drop below 0).",
		"Нестабильный тест упал случайным образом. Снимается 10 XP (XP не может опускаться ниже 0).",
		"Тұрақсыз тест кездейсоқ құлады. 10 XP шегеріледі (XP 0-ден төмен түсе алмайды).",
		"xp_loss", map[string]any{"amount": 10},
	)
	cells[23] = cell(23,
		"Pull Request Merged", "Пул-реквест принят", "Пул-реквест біріктірілді",
		"Your changes are merged into main. Grants +20 XP.",
		"Ваши изменения влиты в основную ветку. Даёт +20 XP.",
		"Өзгерістеріңіз негізгі тармаққа біріктірілді. +20 XP береді.",
		"xp_gain", map[string]any{"amount": 20},
	)

	cells[24] = cell(24,
		"Demo Day", "День демо", "Демо күні",
		"Presented your project to the class! Grants +50 bonus XP (doubled with Double XP).",
		"Представили свой проект группе! Даёт +50 бонусных XP (вдвое больше с Двойным XP).",
		"Жобаңызды топқа таныстырдыңыз! +50 бонус XP береді (Қос XP болса, екі есе көп).",
		"special_challenge", map[string]any{"bonus": 50},
	)

	cells[25] = cell(25,
		"Clear Commit Message", "Понятное сообщение коммита", "Түсінікті коммит хабарламасы",
		"Wrote a commit message everyone understands. Grants +15 XP.",
		"Написали сообщение коммита, понятное всем. Даёт +15 XP.",
		"Барлығына түсінікті коммит хабарламасын жаздыңыз. +15 XP береді.",
		"xp_gain", map[string]any{"amount": 15},
	)
	cells[26] = cell(26,
		"Focus Time", "Время фокуса", "Шоғырлану уақыты",
		"No distractions, deep focus! Your next XP gain is doubled.",
		"Никаких отвлечений, полная концентрация! Следующее начисление XP удваивается.",
		"Еш алаңдаусыз, толық шоғырлану! Келесі XP ұпайыңыз екі есе көбейеді.",
		"double_xp", map[string]any{},
	)
	cells[27] = cell(27,
		"Performance Win", "Ускорение программы", "Өнімділік жеңісі",
		"Made the program run faster. Grants +20 XP.",
		"Программа стала работать быстрее. Даёт +20 XP.",
		"Бағдарламаны жылдамырақ жұмыс істететін еттіңіз. +20 XP береді.",
		"xp_gain", map[string]any{"amount": 20},
	)
	cells[28] = cell(28,
		"Release Candidate", "Релиз-кандидат", "Релиз-кандидат",
		"The release is almost ready. Grants +25 XP.",
		"Релиз почти готов. Даёт +25 XP.",
		"Шығарылым дерлік дайын. +25 XP береді.",
		"xp_gain", map[string]any{"amount": 25},
	)
	cells[29] = cell(29,
		"Bonus Quiz", "Бонусная викторина", "Бонус викторина",
		"Answered a quick bonus quiz! Grants +30 bonus XP (doubled with Double XP).",
		"Ответили на бонусную викторину! Даёт +30 бонусных XP (вдвое больше с Двойным XP).",
		"Бонус викторинаға жауап бердіңіз! +30 бонус XP береді (Қос XP болса, екі есе көп).",
		"special_challenge", map[string]any{"bonus": 30},
	)
	surprise := reuse(30, standard[3]) // same random effect as Mystery Box
	surprise.Name.En, surprise.Name.Ru, surprise.Name.Kz = "Surprise Box", "Коробка-сюрприз", "Тосын қорап"
	cells[30] = surprise
	cells[31] = cell(31,
		"Ready to Deploy", "Готово к развёртыванию", "Орнатуға дайын",
		"Everything is tested and ready to ship. Grants +30 XP.",
		"Всё протестировано и готово к выпуску. Даёт +30 XP.",
		"Барлығы тексерілді және шығаруға дайын. +30 XP береді.",
		"xp_gain", map[string]any{"amount": 30},
	)

	return cells
}
