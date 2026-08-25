package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/jackc/pgx/v5"

	"server/internal/locale"
)

type BoardCellData struct {
	Name   locale.Text
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

	_, err = conn.Exec(ctx, "TRUNCATE games, board_cells, problems, problem_options, problem_accepted_answers, submissions, game_events CASCADE")
	if err != nil {
		log.Fatalf("Failed to truncate tables: %v\n", err)
	}
	log.Println("Database truncated successfully.")

	defaultGameID := "00000000-0000-0000-0000-000000000000"
	_, err = conn.Exec(ctx, "INSERT INTO games (id, status) VALUES ($1, $2)", defaultGameID, "lobby")
	if err != nil {
		log.Fatalf("Failed to insert default game: %v\n", err)
	}
	log.Println("Seeded default game.")

	boardCells := make(map[int]BoardCellData)
	boardCells[0] = BoardCellData{Name: locale.New("Deploy", "Развертывание", "Орнату"), Type: "deploy", Params: map[string]interface{}{"bonus": 200}}
	boardCells[8] = BoardCellData{Name: locale.New("Code Freeze", "Заморозка кода", "Кодты тоқтату"), Type: "code_freeze", Params: map[string]interface{}{}}
	boardCells[16] = BoardCellData{Name: locale.New("Coffee Break", "Перерыв на кофе", "Кофе үзілісі"), Type: "coffee_break", Params: map[string]interface{}{}}
	boardCells[24] = BoardCellData{Name: locale.New("Deadline", "Крайний срок", "Соңғы мерзім"), Type: "deadline", Params: map[string]interface{}{}}

	effects := []BoardCellData{
		{Name: locale.New("XP Gain (S)", "Получение XP (S)", "XP ұпайларын алу (S)"), Type: "xp_gain", Params: map[string]interface{}{"size": "S", "amount": 50}},
		{Name: locale.New("Mystery Event", "Таинственное событие", "Құпия оқиға"), Type: "mystery", Params: map[string]interface{}{}},
		{Name: locale.New("XP Loss (S)", "Потеря XP (S)", "XP ұпайларын жоғалту (S)"), Type: "xp_loss", Params: map[string]interface{}{"size": "S", "amount": 30}},
		{Name: locale.New("Double XP", "Двойной XP", "XP екі есе"), Type: "double_xp", Params: map[string]interface{}{}},
		{Name: locale.New("XP Gain (M)", "Получение XP (M)", "XP ұпайларын алу (M)"), Type: "xp_gain", Params: map[string]interface{}{"size": "M", "amount": 100}},
		{Name: locale.New("Skip Next Turn", "Пропустить следующий ход", "Келесі кезеңді өткізу"), Type: "skip_turn", Params: map[string]interface{}{}},
		{Name: locale.New("XP Gain (S)", "Получение XP (S)", "XP ұпайларын алу (S)"), Type: "xp_gain", Params: map[string]interface{}{"size": "S", "amount": 50}},
		{Name: locale.New("XP Gain (M)", "Получение XP (M)", "XP ұпайларын алу (M)"), Type: "xp_gain", Params: map[string]interface{}{"size": "M", "amount": 100}},
		{Name: locale.New("Teleport", "Телепортация", "Телепорт"), Type: "teleport", Params: map[string]interface{}{"target": 18}},
		{Name: locale.New("XP Loss (M)", "Потеря XP (M)", "XP ұпайларын жоғалту (M)"), Type: "xp_loss", Params: map[string]interface{}{"size": "M", "amount": 60}},
		{Name: locale.New("Free Pass", "Бесплатный проход", "Тегін өту"), Type: "free_pass", Params: map[string]interface{}{}},
		{Name: locale.New("XP Gain (L)", "Получение XP (L)", "XP ұпайларын алу (L)"), Type: "xp_gain", Params: map[string]interface{}{"size": "L", "amount": 200}},
		{Name: locale.New("Mystery Event", "Таинственное событие", "Құпия оқиға"), Type: "mystery", Params: map[string]interface{}{}},
		{Name: locale.New("Special Bonus Challenge", "Специальное бонусное испытание", "Арнайы бонус сынағы"), Type: "bonus_challenge", Params: map[string]interface{}{}},
		{Name: locale.New("XP Gain (S)", "Получение XP (S)", "XP ұпайларын алу (S)"), Type: "xp_gain", Params: map[string]interface{}{"size": "S", "amount": 50}},
		{Name: locale.New("Mystery Event", "Таинственное событие", "Құпия оқиға"), Type: "mystery", Params: map[string]interface{}{}},
		{Name: locale.New("XP Loss (S)", "Потеря XP (S)", "XP ұпайларын жоғалту (S)"), Type: "xp_loss", Params: map[string]interface{}{"size": "S", "amount": 30}},
		{Name: locale.New("Double XP", "Двойной XP", "XP екі есе"), Type: "double_xp", Params: map[string]interface{}{}},
		{Name: locale.New("XP Gain (M)", "Получение XP (M)", "XP ұпайларын алу (M)"), Type: "xp_gain", Params: map[string]interface{}{"size": "M", "amount": 100}},
		{Name: locale.New("Skip Next Turn", "Пропустить следующий ход", "Келесі кезеңді өткізу"), Type: "skip_turn", Params: map[string]interface{}{}},
		{Name: locale.New("XP Gain (S)", "Получение XP (S)", "XP ұпайларын алу (S)"), Type: "xp_gain", Params: map[string]interface{}{"size": "S", "amount": 50}},
		{Name: locale.New("XP Gain (M)", "Получение XP (M)", "XP ұпайларын алу (M)"), Type: "xp_gain", Params: map[string]interface{}{"size": "M", "amount": 100}},
		{Name: locale.New("Teleport", "Телепортация", "Телепорт"), Type: "teleport", Params: map[string]interface{}{"target": 2}},
		{Name: locale.New("XP Loss (M)", "Потеря XP (M)", "XP ұпайларын жоғалту (M)"), Type: "xp_loss", Params: map[string]interface{}{"size": "M", "amount": 60}},
		{Name: locale.New("Free Pass", "Бесплатный проход", "Тегін өту"), Type: "free_pass", Params: map[string]interface{}{}},
		{Name: locale.New("XP Gain (L)", "Получение XP (L)", "XP ұпайларын алу (L)"), Type: "xp_gain", Params: map[string]interface{}{"size": "L", "amount": 200}},
		{Name: locale.New("Mystery Event", "Таинственное событие", "Құпия оқиға"), Type: "mystery", Params: map[string]interface{}{}},
		{Name: locale.New("Special Bonus Challenge", "Специальное бонусное испытание", "Арнайы бонус сынағы"), Type: "bonus_challenge", Params: map[string]interface{}{}},
	}

	effectIdx := 0
	for i := 0; i < 32; i++ {
		if i == 0 || i == 8 || i == 16 || i == 24 {
			continue
		}
		boardCells[i] = effects[effectIdx]
		effectIdx++
	}

	for i := 0; i < 32; i++ {
		cell := boardCells[i]
		if cell.Params == nil {
			cell.Params = map[string]interface{}{}
		}
		paramsJSON, err := json.Marshal(cell.Params)
		if err != nil {
			log.Fatalf("Failed to marshal cell params: %v\n", err)
		}
		_, err = conn.Exec(ctx, "INSERT INTO board_cells (game_id, cell_index, name_en, name_ru, name_kz, type, params) VALUES ($1, $2, $3, $4, $5, $6, $7)", defaultGameID, i, cell.Name.En, cell.Name.Ru, cell.Name.Kz, cell.Type, paramsJSON)
		if err != nil {
			log.Fatalf("Failed to insert board cell %d: %v\n", i, err)
		}
	}
	log.Printf("Seeded 32 board cells for game %s.\n", defaultGameID)

	type SeedOption struct {
		Text      locale.Text
		IsCorrect bool
	}
	type SeedProblem struct {
		Type            string
		Difficulty      string
		Title           locale.Text
		Prompt          locale.Text
		Options         []SeedOption
		AcceptedAnswers []locale.Text
	}

	problemsToSeed := []SeedProblem{
		{
			Type:       "mcq",
			Difficulty: "easy",
			Title:      locale.New("Go Variable Declaration", "Объявление переменных в Go", "Go айнымалысын жариялау"),
			Prompt:     locale.New("Which of the following is the correct short variable declaration in Go?", "Какое из следующих заявлений является правильным кратким объявлением переменной в Go?", "Go тілінде келесілердің қайсысы дұрыс қысқа айнымалы жариялау?"),
			Options: []SeedOption{
				{Text: locale.New("x := 10", "x := 10", "x := 10"), IsCorrect: true},
				{Text: locale.New("x = 10", "x = 10", "x = 10"), IsCorrect: false},
				{Text: locale.New("var x = 10", "var x = 10", "var x = 10"), IsCorrect: false},
				{Text: locale.New("let x = 10", "let x = 10", "let x = 10"), IsCorrect: false},
			},
		},
		{
			Type:       "mcq",
			Difficulty: "easy",
			Title:      locale.New("Slices vs Arrays in Go", "Слайсы и массивы в Go", "Go-дегі тілімдер мен массивтер"),
			Prompt:     locale.New("What is the key difference in declaration between a slice and an array in Go?", "В чём заключается ключевое различие в объявлении слайса и массива в Go?", "Go тілінде slice пен array жариялануындағы негізгі айырмашылық неде?"),
			Options: []SeedOption{
				{Text: locale.New("An array has a fixed size specified in its type, while a slice has a dynamic size.", "Массив имеет фиксированный размер, заданный его типом, тогда как слэйс имеет динамический размер.", "Массивтің түрінде көрсетілген тұрақты өлшемі бар, ал слайстың өлшемі динамикалық."), IsCorrect: true},
				{Text: locale.New("A slice is declared with parentheses `()` and an array with brackets `[]`.", "Слэйс объявляется с помощью круглых скобок `()`, а массив — с помощью квадратных скобок `[]`.", "Слайс `()` скобкаларымен, ал массив `[]` скобкаларымен жарияланады."), IsCorrect: false},
				{Text: locale.New("Slices are value types, while arrays are reference types.", "Слэйсы являются типами значений, тогда как массивы — типами ссылок.", "Слайстар мәндер типі болып табылады, ал массивтер сілтеме типі."), IsCorrect: false},
			},
		},
		{
			Type:       "mcq",
			Difficulty: "easy",
			Title:      locale.New("Go Package Main", "Главный пакет в Go", "Go пакетінің main"),
			Prompt:     locale.New("Which package must be declared for a Go file to contain the entry point of an executable program?", "Какой пакет необходимо объявить, чтобы файл Go содержал точку входа исполняемой программы?", "Go файлында орындалатын бағдарламаның кіру нүктесін қамту үшін қай пакет жариялануы керек?"),
			Options: []SeedOption{
				{Text: locale.New("main", "main", "main"), IsCorrect: true},
				{Text: locale.New("root", "root", "root"), IsCorrect: false},
				{Text: locale.New("start", "start", "start"), IsCorrect: false},
				{Text: locale.New("global", "global", "global"), IsCorrect: false},
			},
		},
		{
			Type:       "mcq",
			Difficulty: "medium",
			Title:      locale.New("Go Channel Directional Constraint", "Ограничение направления каналов в Go", "Go арнасының бағыттық шектеуі"),
			Prompt:     locale.New("How do you define a write-only channel of type int in a Go function signature?", "Как определить канал типа int, доступный только для записи, в сигнатуре функции на Go?", "Go функциясының қолтаңбасында int типінің тек жазуға арналған арнасын қалай анықтайсыз?"),
			Options: []SeedOption{
				{Text: locale.New("chan<- int", "chan<- int", "chan<- int"), IsCorrect: true},
				{Text: locale.New("<-chan int", "<-chan int", "<-chan int"), IsCorrect: false},
				{Text: locale.New("chan int", "chan int", "chan int"), IsCorrect: false},
				{Text: locale.New("write chan int", "write chan int", "write chan int"), IsCorrect: false},
			},
		},
		{
			Type:       "mcq",
			Difficulty: "medium",
			Title:      locale.New("Defer Statement Execution Order", "Порядок выполнения оператора Defer", "Defer операторының орындалу реті"),
			Prompt:     locale.New("In what order are deferred function calls executed in Go?", "В каком порядке выполняются отложенные вызовы функций в Go?", "Go-да кейінге қалдырылған функция шақырулары қандай ретпен орындалады?"),
			Options: []SeedOption{
				{Text: locale.New("LIFO (Last-In, First-Out)", "LIFO (последний вошёл — первый вышел)", "LIFO (соңғы кірген, бірінші шыққан)"), IsCorrect: true},
				{Text: locale.New("FIFO (First-In, First-Out)", "FIFO (первый вошёл — первый вышел)", "FIFO (алғаш кірген, алғаш шыққан)"), IsCorrect: false},
				{Text: locale.New("Random order", "Случайный порядок", "Кездейсоқ рет"), IsCorrect: false},
				{Text: locale.New("Concurrent order", "Параллельный порядок", "Параллель рет"), IsCorrect: false},
			},
		},
		{
			Type:       "mcq",
			Difficulty: "medium",
			Title:      locale.New("Interface Implementation in Go", "Реализация интерфейсов в Go", "Go-дегі интерфейсті іске асыру"),
			Prompt:     locale.New("How does a custom type in Go implement an interface?", "Как пользовательский тип в Go реализует интерфейс?", "Go-да арнайы тип интерфейсті қалай іске асырады?"),
			Options: []SeedOption{
				{Text: locale.New("Implicitly, by implementing all the methods declared in the interface.", "Неявно, путём реализации всех методов, объявленных в интерфейсе.", "Жасырын түрде, интерфейсте жарияланған барлық әдістерді іске асыру арқылы."), IsCorrect: true},
				{Text: locale.New("Explicitly, using the `implements` keyword.", "Явно, с помощью ключевого слова `implements`.", "Анық түрде, `implements` кілтсөзін қолдану арқылы."), IsCorrect: false},
				{Text: locale.New("By inheriting from a base interface struct.", "Путем наследования от базовой структуры интерфейса.", "Негізгі интерфейс құрылымынан мұрагерлік алу арқылы."), IsCorrect: false},
				{Text: locale.New("By registering the type with the interface package.", "Путем регистрации типа в пакете интерфейса.", "Типті интерфейс пакетіне тіркеу арқылы."), IsCorrect: false},
			},
		},
		{
			Type:       "mcq",
			Difficulty: "hard",
			Title:      locale.New("Go Scheduler M:N Model", "Модель M:N планировщика в Go", "Go жоспарлаушысының M:N моделі"),
			Prompt:     locale.New("In Go's concurrency scheduler (G-M-P model), what does the letter 'P' represent?", "Что обозначает буква «P» в модели планировщика параллелизма языка Go (модель G-M-P)?", "Go-ның параллельділік жоспарлаушысында (G-M-P моделінде) 'P' әрпі нені білдіреді?"),
			Options: []SeedOption{
				{Text: locale.New("Processor: a logical resource representing a context required to execute Go code.", "Процессор: логический ресурс, представляющий контекст, необходимый для выполнения кода на языке Go.", "Процессор: Go кодын орындау үшін қажетті контекстті білдіретін логикалық ресурс."), IsCorrect: true},
				{Text: locale.New("Platform: the operating system thread abstraction.", "Платформа: абстракция потоков в операционной системе.", "Платформа: операциялық жүйенің ағын абстракциясы."), IsCorrect: false},
				{Text: locale.New("Program: the main package memory boundary.", "Программа: граница памяти основного пакета.", "Бағдарлама: негізгі пакеттің жады шекарасы."), IsCorrect: false},
				{Text: locale.New("Process: a standard Unix system process.", "Процесс: стандартный системный процесс Unix.", "Процесс: стандартты Unix жүйесінің процесі."), IsCorrect: false},
			},
		},
		{
			Type:       "mcq",
			Difficulty: "hard",
			Title:      locale.New("Sync Map Use Case", "Пример использования синхронизированной карты", "Sync Map қолдану жағдайы"),
			Prompt:     locale.New("When is it recommended to use sync.Map instead of a regular map with a sync.Mutex or sync.RWMutex?", "В каких случаях рекомендуется использовать sync.Map вместо обычной карты с sync.Mutex или sync.RWMutex?", "sync.Mutex немесе sync.RWMutex-пен қорғалған қарапайым map орнына sync.Map-ты қашан пайдалану ұсынылады?"),
			Options: []SeedOption{
				{Text: locale.New("When the entry set is stable and the write rate is very low compared to read rate, or when multiple goroutines read, write, and overwrite entries for disjoint sets of keys.", "Когда набор записей стабилен, а скорость записи очень низкая по сравнению со скоростью чтения, или когда несколько goroutine читают, записывают и перезаписывают записи для непересекающихся наборов ключей.", "Кіру жиынтығы тұрақты болғанда және жазу жылдамдығы оқу жылдамдығына қарағанда өте төмен болғанда, немесе бірнеше goroutine әртүрлі кілт жиынтықтары үшін жазбаларды оқып, жазып және үстіне жазғанда."), IsCorrect: true},
				{Text: locale.New("Whenever multiple goroutines access a map concurrently, regardless of the read/write ratio.", "Всякий раз, когда несколько горутин одновременно обращаются к карте, независимо от соотношения чтения и записи.", "Бірнеше goroutine картаға бір уақытта қол жеткізгенде, оқу/жазу қатынасына қарамастан."), IsCorrect: false},
				{Text: locale.New("When you need to sort the keys of the map concurrently.", "Когда необходимо одновременно сортировать ключи карты.", "Карта кілттерін бір уақытта сұрыптау қажет болғанда."), IsCorrect: false},
				{Text: locale.New("When the map stores a small number of keys (less than 100) that change frequently.", "Когда карта хранит небольшое количество ключей (менее 100), которые часто изменяются.", "Картада жиі өзгеретін аз мөлшерде (100-ден аз) кілт сақталғанда."), IsCorrect: false},
			},
		},
		{
			Type:       "mcq",
			Difficulty: "hard",
			Title:      locale.New("Goroutine Stack Allocation", "Распределение памяти в стеке горутин", "Goroutine стек бөлінісі"),
			Prompt:     locale.New("What is the initial stack size allocated for a goroutine in Go?", "Каков начальный размер стека, выделяемый для goroutine в Go?", "Go-да goroutine үшін бөлінген бастапқы стек көлемі қанша?"),
			Options: []SeedOption{
				{Text: locale.New("2 KB", "2 KB", "2 KB"), IsCorrect: true},
				{Text: locale.New("8 KB", "8 KB", "8 KB"), IsCorrect: false},
				{Text: locale.New("1 MB", "1 MB", "1 MB"), IsCorrect: false},
				{Text: locale.New("4 KB", "4 KB", "4 KB"), IsCorrect: false},
			},
		},
		{
			Type:       "text",
			Difficulty: "easy",
			Title:      locale.New("Map Lookup Zero Value", "Нулевое значение при поиске в карте", "Map іздеуіндегі нөлдік мән"),
			Prompt:     locale.New("What is the second, optional boolean value returned when retrieving an element from a Go map (e.g., val, ok := m[key])?", "Какое второе, необязательное булево значение возвращается при извлечении элемента из карты Go (например, val, ok := m[key])?", "Go картасынан элементті алған кезде қайтарылатын екінші міндетті емес boolean мәні қандай (мысалы, val, ok := m[key])?"),
			AcceptedAnswers: []locale.Text{
				locale.New("ok", "ok", "ok"),
				locale.New("exists", "существует", "бар"),
				locale.New("found", "найден", "табылды"),
				locale.New("present", "присутствует", "бар"),
			},
		},
		{
			Type:       "text",
			Difficulty: "easy",
			Title:      locale.New("Go Garbage Collection", "Сборка мусора в Go", "Go қоқыс жинау"),
			Prompt:     locale.New("What style of garbage collector does Go use? (Hint: uses three colors)", "Какой тип сборщика мусора используется в Go? (Подсказка: используются три цвета)", "Go қандай типтегі қоқыс жинаушыны (garbage collector) қолданады? (Нұсқау: үш түсті пайдаланады)"),
			AcceptedAnswers: []locale.Text{
				locale.New("tri-color mark and sweep", "трёхцветный метод «mark and sweep»", "үш түсті белгілеу және тазалау"),
				locale.New("tricolor mark and sweep", "трёхцветный метод «mark and sweep»", "үш түсті белгілеу және тазалау"),
				locale.New("tri-color", "трёхцветный", "үш түсті"),
				locale.New("tricolor", "трёхцветный", "үш түсті"),
			},
		},
		{
			Type:       "text",
			Difficulty: "easy",
			Title:      locale.New("Pointer Zero Value", "Нулевое значение указателя", "Көрсеткіштің нөлдік мәні"),
			Prompt:     locale.New("What is the zero value of a pointer, slice, channel, or map in Go?", "Каково нулевое значение указателя, слайса, канала или карты в Go?", "Go-да көрсеткіш, жолақ, арна немесе картаның нөлдік мәні не?"),
			AcceptedAnswers: []locale.Text{
				locale.New("nil", "nil", "nil"),
			},
		},
		{
			Type:       "text",
			Difficulty: "medium",
			Title:      locale.New("Context Cancellation Propagation", "Распространение отмены контекста", "Контексттің жойылуын тарату"),
			Prompt:     locale.New("Which standard Go library function is used to create a context that can be cancelled manually?", "Какая стандартная функция библиотеки Go используется для создания контекста, который можно отменить вручную?", "Қай стандартты Go кітапхана функциясы қолмен тоқтатуға болатын контекстті жасау үшін қолданылады?"),
			AcceptedAnswers: []locale.Text{
				locale.New("context.WithCancel", "context.WithCancel", "context.WithCancel"),
				locale.New("WithCancel", "WithCancel", "WithCancel"),
			},
		},
		{
			Type:       "text",
			Difficulty: "medium",
			Title:      locale.New("String Length in Bytes", "Длина строки в байтах", "Сымбаның байтпен ұзындығы"),
			Prompt:     locale.New("What built-in Go function returns the length of a string in bytes?", "Какая встроенная функция Go возвращает длину строки в байтах?", "Қай кіріктірілген Go функциясы жолды байтпен өлшеу үшін қолданылады?"),
			AcceptedAnswers: []locale.Text{
				locale.New("len", "len", "len"),
			},
		},
		{
			Type:       "text",
			Difficulty: "medium",
			Title:      locale.New("Append Capacity", "Емкость при добавлении", "Append сыйымдылығы"),
			Prompt:     locale.New("If you append an item to a slice that is at full capacity, the capacity is typically doubled. What built-in Go function creates a new slice under the hood?", "Если добавить элемент в слайс, который заполнен до предела, его вместимость, как правило, удваивается. Какая встроенная функция Go создаёт новый слайс «под капотом»?", "Егер толық сыйымдылығы бар массивке элемент қоссаңыз, сыйымдылық әдетте екі еселенеді. Қандай кіріктірілген Go функциясы ішкі деңгейде жаңа массив жасайды?"),
			AcceptedAnswers: []locale.Text{
				locale.New("append", "append", "append"),
			},
		},
		{
			Type:       "text",
			Difficulty: "hard",
			Title:      locale.New("Go Heap Profiler Tool", "Инструмент профилирования кучи в Go", "Go Heap Profiler құралы"),
			Prompt:     locale.New("Which tool in the Go toolchain is used to analyze profiling data (e.g., heap allocations or CPU profiles)?", "Какой инструмент из набора средств разработки Go используется для анализа данных профилирования (например, выделения памяти в куче или профилей ЦП)?", "Go құралдар тізбегіндегі қандай құрал профильдік деректерді (мысалы, хип бөліністері немесе CPU профильдері) талдау үшін қолданылады?"),
			AcceptedAnswers: []locale.Text{
				locale.New("pprof", "pprof", "pprof"),
				locale.New("go tool pprof", "go tool pprof", "go tool pprof"),
			},
		},
		{
			Type:       "text",
			Difficulty: "hard",
			Title:      locale.New("Select Non-blocking Default", "Неблокирующее поведение по умолчанию для Select", "Блоктамайтын әдепкі таңдау"),
			Prompt:     locale.New("Which keyword can be added to a select block in Go to make it non-blocking?", "Какое ключевое слово можно добавить в блок `select` в Go, чтобы сделать его неблокирующим?", "Go тілінде select блогына оны блокталмайтын (non-blocking) ету үшін қандай кілтсөзді қосуға болады?"),
			AcceptedAnswers: []locale.Text{
				locale.New("default", "default", "default"),
			},
		},
		{
			Type:       "text",
			Difficulty: "hard",
			Title:      locale.New("Unbuffered Channel Behavior", "Поведение небуферизованных каналов", "Буферленбеген арнаның мінез-құлқы"),
			Prompt:     locale.New("If you send a value to an unbuffered channel when no goroutine is waiting to receive it, what state does the sending goroutine enter?", "Если отправить значение в небуферизованный канал, когда ни одна горутина не ожидает его получения, в какое состояние переходит отправляющая горутина?", "Егер ешбір goroutine оны қабылдауға күтіп тұрған жоқ кезде мәнін буферленбеген арнаға жіберсеңіз, жіберуші goroutine қандай күйге өтеді?"),
			AcceptedAnswers: []locale.Text{
				locale.New("blocked", "заблокирован", "бұғатталды"),
				locale.New("blocking", "блокировка", "бұғаттау"),
				locale.New("block", "блок", "блок"),
			},
		},
	}

	for _, p := range problemsToSeed {
		var problemID string
		err = conn.QueryRow(ctx, `
			INSERT INTO problems (type, difficulty, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, true)
			RETURNING id
		`, p.Type, p.Difficulty, p.Title.En, p.Title.Ru, p.Title.Kz, p.Prompt.En, p.Prompt.Ru, p.Prompt.Kz).Scan(&problemID)
		if err != nil {
			log.Fatalf("Failed to insert problem %q: %v\n", p.Title.En, err)
		}

		if p.Type == "mcq" {
			for _, opt := range p.Options {
				_, err = conn.Exec(ctx, `
					INSERT INTO problem_options (problem_id, text_en, text_ru, text_kz, is_correct)
					VALUES ($1, $2, $3, $4, $5)
				`, problemID, opt.Text.En, opt.Text.Ru, opt.Text.Kz, opt.IsCorrect)
				if err != nil {
					log.Fatalf("Failed to insert option %q for problem %q: %v\n", opt.Text.En, p.Title.En, err)
				}
			}
		} else if p.Type == "text" {
			for _, ans := range p.AcceptedAnswers {
				_, err = conn.Exec(ctx, `
					INSERT INTO problem_accepted_answers (problem_id, answer_text_en, answer_text_ru, answer_text_kz)
					VALUES ($1, $2, $3, $4)
				`, problemID, ans.En, ans.Ru, ans.Kz)
				if err != nil {
					log.Fatalf("Failed to insert accepted answer %q for problem %q: %v\n", ans.En, p.Title.En, err)
				}
			}
		}
	}

	log.Printf("Successfully seeded %d problems.\n", len(problemsToSeed))
	log.Println("Seeding completed successfully!")
}
