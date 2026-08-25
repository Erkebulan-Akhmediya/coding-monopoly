package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

type BoardCellData struct {
	NameEn string
	NameRu string
	NameKz string
	Type   string
	Params map[string]interface{}
}

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/monopoly?sslmode=disable"
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(ctx)

	log.Println("Starting database seeding...")

	// Truncate tables to ensure idempotency
	_, err = conn.Exec(ctx, "TRUNCATE games, board_cells, problems, problem_options, problem_accepted_answers, submissions, game_events CASCADE")
	if err != nil {
		log.Fatalf("Failed to truncate tables: %v\n", err)
	}
	log.Println("Database truncated successfully.")

	// 1. Seed a default game
	defaultGameID := "00000000-0000-0000-0000-000000000000"
	_, err = conn.Exec(ctx, "INSERT INTO games (id, status) VALUES ($1, $2)", defaultGameID, "lobby")
	if err != nil {
		log.Fatalf("Failed to insert default game: %v\n", err)
	}
	log.Println("Seeded default game.")

	// 2. Define the 32 board cells
	// Index 0: Deploy, Index 8: Code Freeze, Index 16: Coffee Break, Index 24: Deadline.
	// The rest are effects.
	boardCells := make(map[int]BoardCellData)
	boardCells[0] = BoardCellData{NameEn: "Deploy", NameRu: "Развернуть", NameKz: "Орнату", Type: "deploy", Params: map[string]interface{}{"bonus": 200}}
	boardCells[8] = BoardCellData{NameEn: "Code Freeze", NameRu: "Заморозка кода", NameKz: "Кодты тоңдыру", Type: "code_freeze", Params: map[string]interface{}{}}
	boardCells[16] = BoardCellData{NameEn: "Coffee Break", NameRu: "Кофе-брейк", NameKz: "Кофе үзілісі", Type: "coffee_break", Params: map[string]interface{}{}}
	boardCells[24] = BoardCellData{NameEn: "Deadline", NameRu: "Крайний срок", NameKz: "Соңғы мерзім", Type: "deadline", Params: map[string]interface{}{}}

	// Distribute the other 28 cells carrying effects:
	// XP gain (S/M/L), XP loss (S/M), mystery/random event, teleport, skip-next-turn, double-XP, free pass, special bonus challenge.
	effects := []BoardCellData{
		{NameEn: "XP Gain (S)", NameRu: "Прирост очков опыта (S)", NameKz: "XP өсімі (S)", Type: "xp_gain", Params: map[string]interface{}{"size": "S", "amount": 50}},
		{NameEn: "Mystery Event", NameRu: "Загадочное событие", NameKz: "Сиқырлы іс-шара", Type: "mystery", Params: map[string]interface{}{}},
		{NameEn: "XP Loss (S)", NameRu: "Потеря XP (S)", NameKz: "XP жоғалту (S)", Type: "xp_loss", Params: map[string]interface{}{"size": "S", "amount": 30}},
		{NameEn: "Double XP", NameRu: "Двойной XP", NameKz: "Екі есе XP", Type: "double_xp", Params: map[string]interface{}{}},
		{NameEn: "XP Gain (M)", NameRu: "Прирост XP (млн)", NameKz: "XP өсімі (M)", Type: "xp_gain", Params: map[string]interface{}{"size": "M", "amount": 100}},
		{NameEn: "Skip Next Turn", NameRu: "Пропустить следующий ход", NameKz: "Келесі кезекті өткізіп жіберу", Type: "skip_turn", Params: map[string]interface{}{}},
		{NameEn: "XP Gain (S)", NameRu: "Прирост очков опыта (S)", NameKz: "XP өсімі (S)", Type: "xp_gain", Params: map[string]interface{}{"size": "S", "amount": 50}},
		// index 8 is code_freeze
		{NameEn: "XP Gain (M)", NameRu: "Прирост XP (млн)", NameKz: "XP өсімі (M)", Type: "xp_gain", Params: map[string]interface{}{"size": "M", "amount": 100}},
		{NameEn: "Teleport", NameRu: "Телепорт", NameKz: "Телепортация", Type: "teleport", Params: map[string]interface{}{"target": 18}},
		{NameEn: "XP Loss (M)", NameRu: "Потеря XP (млн)", NameKz: "XP жоғалту (M)", Type: "xp_loss", Params: map[string]interface{}{"size": "M", "amount": 60}},
		{NameEn: "Free Pass", NameRu: "Бесплатный пропуск", NameKz: "Тегін өту", Type: "free_pass", Params: map[string]interface{}{}},
		{NameEn: "XP Gain (L)", NameRu: "Прирост опыта (L)", NameKz: "XP ұпайы (L)", Type: "xp_gain", Params: map[string]interface{}{"size": "L", "amount": 200}},
		{NameEn: "Mystery Event", NameRu: "Загадочное событие", NameKz: "Сиқырлы іс-шара", Type: "mystery", Params: map[string]interface{}{}},
		{NameEn: "Special Bonus Challenge", NameRu: "Специальное бонусное испытание", NameKz: "Арнайы бонус сынағы", Type: "bonus_challenge", Params: map[string]interface{}{}},
		// index 16 is coffee_break
		{NameEn: "XP Gain (S)", NameRu: "Прирост очков опыта (S)", NameKz: "XP өсімі (S)", Type: "xp_gain", Params: map[string]interface{}{"size": "S", "amount": 50}},
		{NameEn: "Mystery Event", NameRu: "Загадочное событие", NameKz: "Сиқырлы іс-шара", Type: "mystery", Params: map[string]interface{}{}},
		{NameEn: "XP Loss (S)", NameRu: "Потеря XP (S)", NameKz: "XP жоғалту (S)", Type: "xp_loss", Params: map[string]interface{}{"size": "S", "amount": 30}},
		{NameEn: "Double XP", NameRu: "Двойной XP", NameKz: "Екі есе XP", Type: "double_xp", Params: map[string]interface{}{}},
		{NameEn: "XP Gain (M)", NameRu: "Прирост XP (млн)", NameKz: "XP өсімі (M)", Type: "xp_gain", Params: map[string]interface{}{"size": "M", "amount": 100}},
		{NameEn: "Skip Next Turn", NameRu: "Пропустить следующий ход", NameKz: "Келесі кезекті өткізіп жіберу", Type: "skip_turn", Params: map[string]interface{}{}},
		{NameEn: "XP Gain (S)", NameRu: "Прирост очков опыта (S)", NameKz: "XP өсімі (S)", Type: "xp_gain", Params: map[string]interface{}{"size": "S", "amount": 50}},
		// index 24 is deadline
		{NameEn: "XP Gain (M)", NameRu: "Прирост XP (млн)", NameKz: "XP өсімі (M)", Type: "xp_gain", Params: map[string]interface{}{"size": "M", "amount": 100}},
		{NameEn: "Teleport", NameRu: "Телепорт", NameKz: "Телепортация", Type: "teleport", Params: map[string]interface{}{"target": 2}},
		{NameEn: "XP Loss (M)", NameRu: "Потеря XP (млн)", NameKz: "XP жоғалту (M)", Type: "xp_loss", Params: map[string]interface{}{"size": "M", "amount": 60}},
		{NameEn: "Free Pass", NameRu: "Бесплатный пропуск", NameKz: "Тегін өту", Type: "free_pass", Params: map[string]interface{}{}},
		{NameEn: "XP Gain (L)", NameRu: "Прирост опыта (L)", NameKz: "XP ұпайы (L)", Type: "xp_gain", Params: map[string]interface{}{"size": "L", "amount": 200}},
		{NameEn: "Mystery Event", NameRu: "Загадочное событие", NameKz: "Сиқырлы іс-шара", Type: "mystery", Params: map[string]interface{}{}},
		{NameEn: "Special Bonus Challenge", NameRu: "Специальное бонусное испытание", NameKz: "Арнайы бонус сынағы", Type: "bonus_challenge", Params: map[string]interface{}{}},
	}

	effectIdx := 0
	for i := 0; i < 32; i++ {
		if i == 0 || i == 8 || i == 16 || i == 24 {
			continue
		}
		boardCells[i] = effects[effectIdx]
		effectIdx++
	}

	// Insert board cells
	for i := 0; i < 32; i++ {
		cell := boardCells[i]
		if cell.Params == nil {
			cell.Params = map[string]interface{}{}
		}
		paramsJSON, err := json.Marshal(cell.Params)
		if err != nil {
			log.Fatalf("Failed to marshal cell params: %v\n", err)
		}
		_, err = conn.Exec(ctx, "INSERT INTO board_cells (game_id, cell_index, name_en, name_ru, name_kz, type, params) VALUES ($1, $2, $3, $4, $5, $6, $7)", defaultGameID, i, cell.NameEn, cell.NameRu, cell.NameKz, cell.Type, paramsJSON)
		if err != nil {
			log.Fatalf("Failed to insert board cell %d: %v\n", i, err)
		}
	}
	log.Printf("Seeded 32 board cells for game %s.\n", defaultGameID)

	// 3. Seed problems (18 minimum, 3 per difficulty per type)
	type SeedOption struct {
		TextEn    string
		TextRu    string
		TextKz    string
		IsCorrect bool
	}
	type SeedProblem struct {
		Type            string
		Difficulty      string
		TitleEn         string
		TitleRu         string
		TitleKz         string
		PromptEn        string
		PromptRu        string
		PromptKz        string
		Options         []SeedOption
		AcceptedAnswers []string
	}

	problemsToSeed := []SeedProblem{
		// EASY MCQ
		{
			Type:       "mcq",
			Difficulty: "easy",
			TitleEn:    "Go Variable Declaration",
			TitleRu:    "Объявление переменных в Go",
			TitleKz:    "Go тіліндегі айнымалы жариялау",
			PromptEn:   "Which of the following is the correct short variable declaration in Go?",
			PromptRu:   "Какой из приведённых ниже вариантов является правильным объявлением короткой переменной в языке Go?",
			PromptKz:   "Төмендегілердің қайсысы Go тілінде дұрыс қысқа айнымалы жариялау?",
			Options: []SeedOption{
				{TextEn: "x := 10", TextRu: "x := 10", TextKz: "x := 10", IsCorrect: true},
				{TextEn: "x = 10", TextRu: "x = 10", TextKz: "x = 10", IsCorrect: false},
				{TextEn: "var x = 10", TextRu: "var x = 10", TextKz: "var x = 10", IsCorrect: false},
				{TextEn: "let x = 10", TextRu: "let x = 10", TextKz: "let x = 10", IsCorrect: false},
			},
		},
		{
			Type:       "mcq",
			Difficulty: "easy",
			TitleEn:    "Slices vs Arrays in Go",
			TitleRu:    "Срезы и массивы в Go",
			TitleKz:    "Go тіліндегі слайстар мен массивтер",
			PromptEn:   "What is the key difference in declaration between a slice and an array in Go?",
			PromptRu:   "В чём заключается основное различие в объявлении среза и массива в языке Go?",
			PromptKz:   "Go тілінде slice пен array-ды жариялаудың негізгі айырмашылығы неде?",
			Options: []SeedOption{
				{TextEn: "An array has a fixed size specified in its type, while a slice has a dynamic size.", TextRu: "Массив имеет фиксированный размер, заданный его типом, тогда как срез имеет динамический размер.", TextKz: "Массивтің өлшемі оның типімен анықталған тұрақты өлшемге ие, ал слайстың өлшемі динамикалық.", IsCorrect: true},
				{TextEn: "A slice is declared with parentheses `()` and an array with brackets `[]`.", TextRu: "Срез объявляется с помощью круглых скобок `()`, а массив — с помощью квадратных скобок `[]`.", TextKz: "Бір бөлік скобкалармен `()` жарияланады, ал массив тік жақшалармен `[]`.", IsCorrect: false},
				{TextEn: "Slices are value types, while arrays are reference types.", TextRu: "Срезки являются типом по значению, тогда как массивы — типом по ссылке.", TextKz: "Слайстер – мән типтері, ал массивтер – сілтеме типтері.", IsCorrect: false},
			},
		},
		{
			Type:       "mcq",
			Difficulty: "easy",
			TitleEn:    "Go Package Main",
			TitleRu:    "Основной пакет Go",
			TitleKz:    "Go негізгі пакеті",
			PromptEn:   "Which package must be declared for a Go file to contain the entry point of an executable program?",
			PromptRu:   "Какой пакет необходимо указать в файле Go, чтобы он содержал точку входа исполняемой программы?",
			PromptKz:   "Орындалатын бағдарламаның кіру нүктесін қамтитын Go файлы үшін қай пакет жариялануы керек?",
			Options: []SeedOption{
				{TextEn: "main", TextRu: "main", TextKz: "main", IsCorrect: true},
				{TextEn: "root", TextRu: "root", TextKz: "root", IsCorrect: false},
				{TextEn: "start", TextRu: "start", TextKz: "start", IsCorrect: false},
				{TextEn: "global", TextRu: "global", TextKz: "global", IsCorrect: false},
			},
		},

		// MEDIUM MCQ
		{
			Type:       "mcq",
			Difficulty: "medium",
			TitleEn:    "Go Channel Directional Constraint",
			TitleRu:    "Ограничение направления канала в Go",
			TitleKz:    "Go тіліндегі арна бағытының шектеуі",
			PromptEn:   "How do you define a write-only channel of type int in a Go function signature?",
			PromptRu:   "Как определить канал типа `int` только для записи в сигнатуре функции на языке Go?",
			PromptKz:   "Go функциясының қолтаңбасында int типіндегі тек жазуға арналған арнаны қалай анықтайсыз?",
			Options: []SeedOption{
				{TextEn: "chan<- int", TextRu: "chan<- int", TextKz: "chan<- int", IsCorrect: true},
				{TextEn: "<-chan int", TextRu: "<-chan int", TextKz: "<-chan int", IsCorrect: false},
				{TextEn: "chan int", TextRu: "chan int", TextKz: "chan int", IsCorrect: false},
				{TextEn: "write chan int", TextRu: "write chan int", TextKz: "write chan int", IsCorrect: false},
			},
		},
		{
			Type:       "mcq",
			Difficulty: "medium",
			TitleEn:    "Defer Statement Execution Order",
			TitleRu:    "Порядок выполнения отложенных вызовов",
			TitleKz:    "Кейінге қалдырылған шақырулардың орындалу тәртібі",
			PromptEn:   "In what order are deferred function calls executed in Go?",
			PromptRu:   "В каком порядке выполняются отложенные вызовы функций в языке Go?",
			PromptKz:   "Go тілінде кейінге қалдырылған функция шақырулары қандай ретпен орындалады?",
			Options: []SeedOption{
				{TextEn: "LIFO (Last-In, First-Out)", TextRu: "LIFO (последний поступил — первый выбыл)", TextKz: "LIFO (соңғы кірген, бірінші шыққан)", IsCorrect: true},
				{TextEn: "FIFO (First-In, First-Out)", TextRu: "FIFO (первым поступил — первым вышел)", TextKz: "FIFO (алғаш кірген, алғаш шыққан)", IsCorrect: false},
				{TextEn: "Random order", TextRu: "Случайный порядок", TextKz: "Кездейсоқ тәртіп", IsCorrect: false},
				{TextEn: "Concurrent order", TextRu: "Конкурентный порядок", TextKz: "Бір мезгілдегі тәртіп", IsCorrect: false},
			},
		},
		{
			Type:       "mcq",
			Difficulty: "medium",
			TitleEn:    "Interface Implementation in Go",
			TitleRu:    "Реализация интерфейса в Go",
			TitleKz:    "Go тілінде интерфейсті іске асыру",
			PromptEn:   "How does a custom type in Go implement an interface?",
			PromptRu:   "Как пользовательский тип в Go реализует интерфейс?",
			PromptKz:   "Go тіліндегі тапсырыс бойынша жасалған тип интерфейсті қалай іске асырады?",
			Options: []SeedOption{
				{TextEn: "Implicitly, by implementing all the methods declared in the interface.", TextRu: "Неявно — путем реализации всех методов, объявленных в интерфейсе.", TextKz: "Жасырын түрде, интерфейсте жарияланған барлық әдістерді іске асыру арқылы.", IsCorrect: true},
				{TextEn: "Explicitly, using the `implements` keyword.", TextRu: "Явно, с использованием ключевого слова `implements`.", TextKz: "Анық түрде `implements` кілтсөзін пайдалана отырып.", IsCorrect: false},
				{TextEn: "By inheriting from a base interface struct.", TextRu: "Путем наследования от базового интерфейса struct.", TextKz: "Негізгі интерфейс құрылымынан мұрагерлік алу арқылы.", IsCorrect: false},
				{TextEn: "By registering the type with the interface package.", TextRu: "Путем регистрации типа в пакете интерфейсов.", TextKz: "Интерфейс пакетіндегі типті тіркеу арқылы.", IsCorrect: false},
			},
		},

		// HARD MCQ
		{
			Type:       "mcq",
			Difficulty: "hard",
			TitleEn:    "Go Scheduler M:N Model",
			TitleRu:    "Планировщик Go: модель M:N",
			TitleKz:    "Go жоспарлаушысы: M:N моделі",
			PromptEn:   "In Go's concurrency scheduler (G-M-P model), what does the letter 'P' represent?",
			PromptRu:   "Что обозначает буква «P» в планировщике параллелизма языка Go (модель G-M-P)?",
			PromptKz:   "Go-ның параллельділікті жоспарлаушысында (G-M-P моделінде) «P» әрпі нені білдіреді?",
			Options: []SeedOption{
				{TextEn: "Processor: a logical resource representing a context required to execute Go code.", TextRu: "Процессор: логический ресурс, представляющий собой контекст, необходимый для выполнения кода на языке Go.", TextKz: "Процессор: Go кодын орындау үшін қажетті контекстті білдіретін логикалық ресурс.", IsCorrect: true},
				{TextEn: "Platform: the operating system thread abstraction.", TextRu: "Платформа: абстракция потоков в операционной системе.", TextKz: "Платформа: операциялық жүйенің ағын абстракциясы.", IsCorrect: false},
				{TextEn: "Program: the main package memory boundary.", TextRu: "Программа: граница памяти основного пакета.", TextKz: "Бағдарлама: негізгі пакеттің жады шегі.", IsCorrect: false},
				{TextEn: "Process: a standard Unix system process.", TextRu: "Процесс: стандартный системный процесс Unix.", TextKz: "Процесс: стандартты Unix жүйесінің процесі.", IsCorrect: false},
			},
		},
		{
			Type:       "mcq",
			Difficulty: "hard",
			TitleEn:    "Sync Map Use Case",
			TitleRu:    "Применение sync.Map",
			TitleKz:    "sync.Map қолдану жағдайлары",
			PromptEn:   "When is it recommended to use sync.Map instead of a regular map with a sync.Mutex or sync.RWMutex?",
			PromptRu:   "В каких случаях рекомендуется использовать sync.Map вместо обычной карты с sync.Mutex или sync.RWMutex?",
			PromptKz:   "sync.Mutex немесе sync.RWMutex қолданылатын кәдімгі картаның орнына sync.Map-ты қашан пайдалану ұсынылады?",
			Options: []SeedOption{
				{TextEn: "When the entry set is stable and the write rate is very low compared to read rate, or when multiple goroutines read, write, and overwrite entries for disjoint sets of keys.", TextRu: "Когда набор записей является стабильным, а скорость записи значительно ниже скорости чтения, либо когда несколько goroutine выполняют чтение, запись и перезапись записей для непересекающихся наборов ключей.", TextKz: "Кіру жиынтығы тұрақты болғанда және жазу жылдамдығы оқу жылдамдығымен салыстырғанда өте төмен болғанда, немесе бірнеше goroutine әртүрлі кілт жиынтықтары үшін жазбаларды оқып, жазып және үстіне жазғанда.", IsCorrect: true},
				{TextEn: "Whenever multiple goroutines access a map concurrently, regardless of the read/write ratio.", TextRu: "Всякий раз, когда несколько goroutine одновременно обращаются к карте, независимо от соотношения операций чтения и записи.", TextKz: "Көптеген goroutine-дер картаға бір уақытта қол жеткізгенде, оқу/жазу қатынасына қарамастан.", IsCorrect: false},
				{TextEn: "When you need to sort the keys of the map concurrently.", TextRu: "Когда требуется параллельная сортировка ключей карты.", TextKz: "Карта кілттерін бір уақытта сұрыптау қажет болғанда.", IsCorrect: false},
				{TextEn: "When the map stores a small number of keys (less than 100) that change frequently.", TextRu: "Когда в карте хранится небольшое количество ключей (менее 100), которые часто изменяются.", TextKz: "Карта жиі өзгеретін аз ғана (100-ден аз) кілттерді сақтағанда.", IsCorrect: false},
			},
		},
		{
			Type:       "mcq",
			Difficulty: "hard",
			TitleEn:    "Goroutine Stack Allocation",
			TitleRu:    "Распределение стека горутины",
			TitleKz:    "Горутина стегін бөлу",
			PromptEn:   "What is the initial stack size allocated for a goroutine in Go?",
			PromptRu:   "Каков начальный размер стека, выделяемый для горутины в языке Go?",
			PromptKz:   "Go-да goroutine үшін бөлінген бастапқы стек көлемі қанша?",
			Options: []SeedOption{
				{TextEn: "2 KB", TextRu: "2 КБ", TextKz: "2 КБ", IsCorrect: true},
				{TextEn: "8 KB", TextRu: "8 КБ", TextKz: "8 КБ", IsCorrect: false},
				{TextEn: "1 MB", TextRu: "1 МБ", TextKz: "1 МБ", IsCorrect: false},
				{TextEn: "4 KB", TextRu: "4 КБ", TextKz: "4 КБ", IsCorrect: false},
			},
		},

		// EASY TEXT
		{
			Type:            "text",
			Difficulty:      "easy",
			TitleEn:         "Map Lookup Zero Value",
			TitleRu:         "Нулевое значение при поиске в карте",
			TitleKz:         "Картада іздегендегі нөлдік мән",
			PromptEn:        "What is the second, optional boolean value returned when retrieving an element from a Go map (e.g., val, ok := m[key])?",
			PromptRu:        "Какое второе, необязательное булево значение возвращается при извлечении элемента из карты в Go (например, val, ok := m[key])?",
			PromptKz:        "Go картасынан элементті алған кезде қайтарылатын екінші, міндетті емес boolean мәні қандай (мысалы, val, ok := m[key])?",
			AcceptedAnswers: []string{"ok", "exists", "found", "present"},
		},
		{
			Type:            "text",
			Difficulty:      "easy",
			TitleEn:         "Go Garbage Collection",
			TitleRu:         "Сборка мусора в Go",
			TitleKz:         "Go-да қоқыс жинау",
			PromptEn:        "What style of garbage collector does Go use? (Hint: uses three colors)",
			PromptRu:        "Какой тип сборщика мусора используется в Go? (Подсказка: используются три цвета)",
			PromptKz:        "Go қандай стильдегі қоқыс жинағышты пайдаланады? (Көрсеткіш: үш түсті пайдаланады)",
			AcceptedAnswers: []string{"tri-color mark and sweep", "tricolor mark and sweep", "tri-color", "tricolor"},
		},
		{
			Type:            "text",
			Difficulty:      "easy",
			TitleEn:         "Pointer Zero Value",
			TitleRu:         "Нулевое значение указателя",
			TitleKz:         "Көрсеткіштің нөлдік мәні",
			PromptEn:        "What is the zero value of a pointer, slice, channel, or map in Go?",
			PromptRu:        "Что означает нулевое значение указателя, среза, канала или карты в Go?",
			PromptKz:        "Go тілінде көрсеткіштің, тілімшенің, арнаның немесе картаның нөлдік мәні нені білдіреді?",
			AcceptedAnswers: []string{"nil"},
		},

		// MEDIUM TEXT
		{
			Type:            "text",
			Difficulty:      "medium",
			TitleEn:         "Context Cancellation Propagation",
			TitleRu:         "Распространение отмены контекста",
			TitleKz:         "Контекстті болдырмауды тарату",
			PromptEn:        "Which standard Go library function is used to create a context that can be cancelled manually?",
			PromptRu:        "Какая стандартная функция библиотеки Go используется для создания контекста, который можно отменить вручную?",
			PromptKz:        "Қай стандартты Go кітапхана функциясы қолмен тоқтатуға болатын контекстті жасау үшін қолданылады?",
			AcceptedAnswers: []string{"context.WithCancel", "WithCancel"},
		},
		{
			Type:            "text",
			Difficulty:      "medium",
			TitleEn:         "String Length in Bytes",
			TitleRu:         "Длина строки в байтах",
			TitleKz:         "Жолдың ұзындығы байтпен",
			PromptEn:        "What built-in Go function returns the length of a string in bytes?",
			PromptRu:        "Какая встроенная функция языка Go возвращает длину строки в байтах?",
			PromptKz:        "Go тіліндегі қандай кіріктірілген функция жолдың ұзындығын байтпен қайтарады?",
			AcceptedAnswers: []string{"len"},
		},
		{
			Type:            "text",
			Difficulty:      "medium",
			TitleEn:         "Append Capacity",
			TitleRu:         "Ёмкость при добавлении",
			TitleKz:         "Қосу кезіндегі сыйымдылық",
			PromptEn:        "If you append an item to a slice that is at full capacity, the capacity is typically doubled. What built-in Go function creates a new slice under the hood?",
			PromptRu:        "Если добавить элемент в слайс, который уже заполнен полностью, его емкость, как правило, удваивается. Какая встроенная функция языка Go создаёт новый слайс «под капотом»?",
			PromptKz:        "Егер толық сыйымдылығы бар слайсқа элемент қоссаңыз, сыйымдылық әдетте екі есе ұлғаяды. Go тіліндегі қандай кіріктірілген функция жаңа слайсты артында жасап шығарады?",
			AcceptedAnswers: []string{"append"},
		},

		// HARD TEXT
		{
			Type:            "text",
			Difficulty:      "hard",
			TitleEn:         "Go Heap Profiler Tool",
			TitleRu:         "Инструмент профилирования кучи Go",
			TitleKz:         "Go үйме профайлері құралы",
			PromptEn:        "Which tool in the Go toolchain is used to analyze profiling data (e.g., heap allocations or CPU profiles)?",
			PromptRu:        "Какой инструмент из набора инструментов Go используется для анализа данных профилирования (например, выделения памяти в куче или профилей ЦП)?",
			PromptKz:        "Go құралдар тізбегінде профильдік деректерді (мысалы, хип бөліністері немесе CPU профильдері) талдау үшін қандай құрал қолданылады?",
			AcceptedAnswers: []string{"pprof", "go tool pprof"},
		},
		{
			Type:            "text",
			Difficulty:      "hard",
			TitleEn:         "Select Non-blocking Default",
			TitleRu:         "Неблокирующий select по умолчанию",
			TitleKz:         "Блоктамайтын select әдепкі",
			PromptEn:        "Which keyword can be added to a select block in Go to make it non-blocking?",
			PromptRu:        "Какое ключевое слово можно добавить в блок `select` в языке Go, чтобы сделать его неблокирующим?",
			PromptKz:        "Go тілінде select блогына қай кілтсөзді қосуға болады, оны блоктық емес ету үшін?",
			AcceptedAnswers: []string{"default"},
		},
		{
			Type:            "text",
			Difficulty:      "hard",
			TitleEn:         "Unbuffered Channel Behavior",
			TitleRu:         "Поведение канала без буфера",
			TitleKz:         "Буферленбеген арнаның мінез-құлқы",
			PromptEn:        "If you send a value to an unbuffered channel when no goroutine is waiting to receive it, what state does the sending goroutine enter?",
			PromptRu:        "Если отправить значение по небуферизованному каналу, когда ни одна горутина не ожидает его получения, в какое состояние переходит отправляющая горутина?",
			PromptKz:        "Егер ешбір goroutine оны қабылдауға күтіп тұрған жоқ кезде мәнін буферленбеген арнаға жіберсеңіз, жіберуші goroutine қандай күйге өтеді?",
			AcceptedAnswers: []string{"blocked", "blocking", "block"},
		},
	}

	for _, p := range problemsToSeed {
		var problemID string
		err = conn.QueryRow(ctx, `
			INSERT INTO problems (type, difficulty, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, true)
			RETURNING id
		`, p.Type, p.Difficulty, p.TitleEn, p.TitleRu, p.TitleKz, p.PromptEn, p.PromptRu, p.PromptKz).Scan(&problemID)
		if err != nil {
			log.Fatalf("Failed to insert problem %q: %v\n", p.TitleEn, err)
		}

		if p.Type == "mcq" {
			for _, opt := range p.Options {
				_, err = conn.Exec(ctx, `
					INSERT INTO problem_options (problem_id, text_en, text_ru, text_kz, is_correct)
					VALUES ($1, $2, $3, $4, $5)
				`, problemID, opt.TextEn, opt.TextRu, opt.TextKz, opt.IsCorrect)
				if err != nil {
					log.Fatalf("Failed to insert option %q for problem %q: %v\n", opt.TextEn, p.TitleEn, err)
				}
			}
		} else if p.Type == "text" {
			for _, ans := range p.AcceptedAnswers {
				_, err = conn.Exec(ctx, `
					INSERT INTO problem_accepted_answers (problem_id, answer_text)
					VALUES ($1, $2)
				`, problemID, ans)
				if err != nil {
					log.Fatalf("Failed to insert accepted answer %q for problem %q: %v\n", ans, p.TitleEn, err)
				}
			}
		}
	}

	log.Printf("Successfully seeded %d problems.\n", len(problemsToSeed))
	log.Println("Seeding completed successfully!")
}
