-- Auto-generated from questions.md for topic 'piscine'
BEGIN;

-- Clean up existing piscine questions to ensure idempotency
DELETE FROM problems WHERE topic = 'piscine';

-- Question 1: Printing & Program Structure (easy, mcq)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('61489e7a-f8ff-5812-a1dd-c17b52876cb4', 'mcq', 'easy', 'piscine', 'Printing & Program Structure', 'Вывод и структура программы', 'Шығару және бағдарлама құрылымы', 'Which package declaration is required at the top of a Go file to produce a standalone executable program?', 'Какую декларацию пакета необходимо разместить в начале файла на языке Go, чтобы получить автономную исполняемую программу?', 'Go файлының басында жеке орындалатын бағдарлама алу үшін қандай пакет жариялауы қажет?', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('61489e7a-f8ff-5812-a1dd-c17b52876cb4', 'mcq', '`package run`', '`package run`', '`package run`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('61489e7a-f8ff-5812-a1dd-c17b52876cb4', 'mcq', '`package main`', '`package main`', '`package main`', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('61489e7a-f8ff-5812-a1dd-c17b52876cb4', 'mcq', '`package root`', '`package root`', '`package root`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('61489e7a-f8ff-5812-a1dd-c17b52876cb4', 'mcq', '`package exec`', '`package exec`', '`package exec`', false);

-- Question 2: Printing & Program Structure (medium, mcq)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('af3a681a-c0e6-5027-b64c-ea74bb4af5f1', 'mcq', 'medium', 'piscine', 'Printing & Program Structure', 'Вывод и структура программы', 'Шығару және бағдарлама құрылымы', 'Which of the following statements regarding the function `fmt.Println()` in Go are TRUE? (Select all that apply)', 'Какие из приведённых ниже утверждений относительно функции `fmt.Println()` в языке программирования Go являются верными? (Выберите все правильные варианты)', 'Go тіліндегі `fmt.Println()` функциясына қатысты келесі мәлімдемелердің қайсысы дұрыс? (Барлығын таңдаңыз)', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('af3a681a-c0e6-5027-b64c-ea74bb4af5f1', 'mcq', 'It automatically separates multiple arguments with spaces.', 'Он автоматически разделяет несколько аргументов пробелами.', 'Ол бірнеше аргументтерді автоматты түрде бос орындармен бөледі.', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('af3a681a-c0e6-5027-b64c-ea74bb4af5f1', 'mcq', 'It automatically appends a newline character (`\n`) at the end of the printed output.', 'Она автоматически добавляет символ новой строки (`\n`) в конец выводимого результата.', 'Ол басылып шығарылған нәтиженің соңына автоматты түрде жаңа жол таңбасын (`\n`) қосады.', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('af3a681a-c0e6-5027-b64c-ea74bb4af5f1', 'mcq', 'It can only accept arguments of type `string`.', 'Он может принимать только аргументы типа `string`.', 'Ол тек `string` типіндегі аргументтерді қабылдай алады.', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('af3a681a-c0e6-5027-b64c-ea74bb4af5f1', 'mcq', 'Its name starts with a capital ''P'' because it is exported from the `fmt` package.', 'Его название начинается с заглавной буквы «P», поскольку он экспортируется из пакета `fmt`.', 'Оның атауы `fmt` пакетінен экспортталғандықтан үлкен «P» әрпімен басталады.', true);

-- Question 3: Printing & Program Structure (hard, text)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('847601e0-ef47-5b36-88fb-2ae71bbca525', 'text', 'hard', 'piscine', 'Printing & Program Structure', 'Вывод и структура программы', 'Шығару және бағдарлама құрылымы', 'In Go, what letter casing (uppercase or lowercase) must the first letter of a function name have to be exported and callable from another package (like `Println` in package `fmt`)?', 'В языке Go с какой буквы (заглавной или строчной) должно начинаться имя функции, чтобы она была экспортируемой (доступной) из другого пакета (например, `Println` в пакете `fmt`)?', 'Go тілінде функцияның атауы басқа пакеттен (мысалы, `fmt` пакетіндегі `Println`) экспортталып шақырылуы үшін оның бірінші әрпі қандай (бас әріп немесе кіші әріп) болуы керек?', true);
INSERT INTO problem_accepted_answers (problem_id, problem_type, answer_text_en, answer_text_ru, answer_text_kz)
VALUES ('847601e0-ef47-5b36-88fb-2ae71bbca525', 'text', 'uppercase', 'заглавная буква', 'бас әріп');
INSERT INTO problem_accepted_answers (problem_id, problem_type, answer_text_en, answer_text_ru, answer_text_kz)
VALUES ('847601e0-ef47-5b36-88fb-2ae71bbca525', 'text', 'capital', 'заглавная', 'бас әріппен');
INSERT INTO problem_accepted_answers (problem_id, problem_type, answer_text_en, answer_text_ru, answer_text_kz)
VALUES ('847601e0-ef47-5b36-88fb-2ae71bbca525', 'text', 'capital letter', 'с заглавной буквы', 'үлкен әріп');
INSERT INTO problem_accepted_answers (problem_id, problem_type, answer_text_en, answer_text_ru, answer_text_kz)
VALUES ('847601e0-ef47-5b36-88fb-2ae71bbca525', 'text', 'uppercase letter', 'прописная буква', 'бас әріп');
INSERT INTO problem_accepted_answers (problem_id, problem_type, answer_text_en, answer_text_ru, answer_text_kz)
VALUES ('847601e0-ef47-5b36-88fb-2ae71bbca525', 'text', 'upper', 'верхний регистр', 'жоғарғы регистр');

-- Question 4: Variables, Constants & Types (easy, mcq)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('61efde92-1dca-5e9b-919e-4cb60496bca5', 'mcq', 'easy', 'piscine', 'Variables, Constants & Types', 'Переменные, константы и типы', 'Айнымалылар, тұрақтылар және типтер', 'What is the default zero value of an uninitialized variable declared as `var count int` in Go?', 'Каково значение по умолчанию (нуль) для неинициализированной переменной, объявленной в языке Go как `var count int`?', 'Go тілінде `var count int` деп жарияланған бастапқы мәні нөлге тең емес айнымалының әдепкі нөлдік мәні қандай?', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('61efde92-1dca-5e9b-919e-4cb60496bca5', 'mcq', '`nil`', '`nil`', '`nil`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('61efde92-1dca-5e9b-919e-4cb60496bca5', 'mcq', '`-1`', '`-1`', '`-1`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('61efde92-1dca-5e9b-919e-4cb60496bca5', 'mcq', '`0`', '`0`', '`0`', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('61efde92-1dca-5e9b-919e-4cb60496bca5', 'mcq', '`undefined`', '`undefined`', '`undefined`', false);

-- Question 5: Variables, Constants & Types (medium, mcq)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('1b6ee0f9-f49b-55ec-90d1-ec7bb6d28ac0', 'mcq', 'medium', 'piscine', 'Variables, Constants & Types', 'Переменные, константы и типы', 'Айнымалылар, тұрақтылар және типтер', 'What happens when the following code snippet is compiled in Go?

```go
x := 10
x := 20
fmt.Println(x)
```', 'Что произойдет, если скомпилировать следующий фрагмент кода на языке Go?

```go
x := 10
x := 20
fmt.Println(x)
```', 'Төмендегі код үзіндісі Go тілінде компиляцияланғанда не болады?

```go
x := 10
x := 20
fmt.Println(x)
```', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('1b6ee0f9-f49b-55ec-90d1-ec7bb6d28ac0', 'mcq', 'It prints `20`', 'Выводится `20`', 'Ол 20-ды басып шығарады.', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('1b6ee0f9-f49b-55ec-90d1-ec7bb6d28ac0', 'mcq', 'It prints `10`', 'Выводится `10`', 'Ол 10 деп шығарады.', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('1b6ee0f9-f49b-55ec-90d1-ec7bb6d28ac0', 'mcq', 'Compile-time error: `no new variables on left side of :=`', 'Ошибка компиляции: `нет новых переменных в левой части выражения :=`', 'Құрастыру кезіндегі қате: `=:= операторының сол жақ жағында жаңа айнымалылар жоқ`', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('1b6ee0f9-f49b-55ec-90d1-ec7bb6d28ac0', 'mcq', 'Runtime panic', 'Паника во время выполнения', 'Жүгіру кезіндегі паника', false);

-- Question 6: Variables, Constants & Types (hard, text)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('4216bb48-3977-553b-8508-65700c0787af', 'text', 'hard', 'piscine', 'Variables, Constants & Types', 'Переменные, константы и типы', 'Айнымалылар, тұрақтылар және типтер', 'Which keyword in Go is used to declare an immutable value whose value is fixed at compile time and cannot be reassigned during execution?', 'Какое ключевое слово в языке Go используется для объявления неизменяемого значения, величина которого фиксируется на этапе компиляции и не может быть переопределена во время выполнения?', 'Go тілінде компиляция кезінде мәні бекітіліп, орындау кезінде қайта тағайындалмайтын өзгермейтін мәнді жариялау үшін қандай кілтсөз қолданылады?', true);
INSERT INTO problem_accepted_answers (problem_id, problem_type, answer_text_en, answer_text_ru, answer_text_kz)
VALUES ('4216bb48-3977-553b-8508-65700c0787af', 'text', 'const', 'const', 'тұрақты');
INSERT INTO problem_accepted_answers (problem_id, problem_type, answer_text_en, answer_text_ru, answer_text_kz)
VALUES ('4216bb48-3977-553b-8508-65700c0787af', 'text', 'const keyword', 'ключевое слово const', 'const кілтсөзі');

-- Question 7: Functions & Operators (easy, mcq)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('fd89655b-91fe-5a65-bfe6-683699bebe54', 'mcq', 'easy', 'piscine', 'Functions & Operators', 'Функции и операторы', 'Функциялар мен операторлар', 'What is the evaluated result of the integer division expression `7 / 2` in Go?', 'Каков результат вычисления выражения целочисленного деления `7 / 2` в языке Go?', 'Go тілінде 7 / 2 бүтін сандық бөлу операторының бағаланған нәтижесі қандай?', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('fd89655b-91fe-5a65-bfe6-683699bebe54', 'mcq', '`3.5`', '`3.5`', '`3.5`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('fd89655b-91fe-5a65-bfe6-683699bebe54', 'mcq', '`3`', '`3`', '`3`', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('fd89655b-91fe-5a65-bfe6-683699bebe54', 'mcq', '`4`', '`4`', '`4`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('fd89655b-91fe-5a65-bfe6-683699bebe54', 'mcq', 'Compile-time error', 'Ошибка компиляции', 'Құрастыру кезіндегі қате', false);

-- Question 8: Functions & Operators (medium, mcq)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('231124dd-7e16-56ef-a9c5-7658e7e7f021', 'mcq', 'medium', 'piscine', 'Functions & Operators', 'Функции и операторы', 'Функциялар мен операторлар', 'What does the following function return when called as `calc(14, 4)`?

```go
func calc(a, b int) (int, int) {
    return a / b, a % b
}
```', 'Какое значение возвращает следующая функция при вызове `calc(14, 4)`?

```go
func calc(a, b int) (int, int) {
    return a / b, a % b
}
```', 'Төмендегі функция `calc(14, 4)` деп шақырылғанда қандай мәнін қайтарады?

```go
func calc(a, b int) (int, int) {
    return a / b, a % b
}
```', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('231124dd-7e16-56ef-a9c5-7658e7e7f021', 'mcq', '`3, 2`', '`3, 2`', '`3, 2`', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('231124dd-7e16-56ef-a9c5-7658e7e7f021', 'mcq', '`3.5, 2`', '`3.5, 2`', '`3.5, 2`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('231124dd-7e16-56ef-a9c5-7658e7e7f021', 'mcq', '`2, 3`', '`2, 3`', '`2, 3`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('231124dd-7e16-56ef-a9c5-7658e7e7f021', 'mcq', '`4, 2`', '`4, 2`', '`4, 2`', false);

-- Question 9: Functions & Operators (hard, mcq)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('92976d89-9d79-5edd-b2c3-23db44da8537', 'mcq', 'hard', 'piscine', 'Functions & Operators', 'Функции и операторы', 'Функциялар мен операторлар', 'Which of the following code snippets will produce a COMPILE-TIME error in Go? (Select all that apply)', 'Какой из приведённых ниже фрагментов кода вызовет ошибку на этапе компиляции в языке Go? (Выберите все правильные варианты)', 'Төмендегі код үзінділерінің қайсысы Go тілінде компиляция кезінде қате тудырады? (Барлығын таңдаңыз)', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('92976d89-9d79-5edd-b2c3-23db44da8537', 'mcq', '`var a int = 5; var b float64 = 2.0; var c = a + b`', '`var a int = 5; var b float64 = 2.0; var c = a + b`', '`var a int = 5; var b float64 = 2.0; var c = a + b`', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('92976d89-9d79-5edd-b2c3-23db44da8537', 'mcq', '`var a int = 5; var b float64 = 2.0; var c = float64(a) + b`', '`var a int = 5; var b float64 = 2.0; var c = float64(a) + b`', '`var a int = 5; var b float64 = 2.0; var c = float64(a) + b`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('92976d89-9d79-5edd-b2c3-23db44da8537', 'mcq', '`var x int = 10; _ = x`', '`var x int = 10; _ = x`', '`var x int = 10; _ = x`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('92976d89-9d79-5edd-b2c3-23db44da8537', 'mcq', '`func f() int { if false { return 1 } }`', '`func f() int { if false { return 1 } }`', '`func f() int { if false { return 1 } }`', true);

-- Question 10: Strings & Runes (easy, mcq)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('8c74ee3c-c9a3-5b8b-9622-e1e5eee2ff27', 'mcq', 'easy', 'piscine', 'Strings & Runes', 'Строки и руны', 'Жолдар мен руналар', 'What does the built-in function `len(s)` return when passed a string `s` in Go?', 'Что возвращает встроенная функция `len(s)` при передаче ей строки `s` в языке Go?', 'Go тілінде `len(s)` кіріктірілген функциясы `s` жолды бергенде қандай мәні қайтарады?', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('8c74ee3c-c9a3-5b8b-9622-e1e5eee2ff27', 'mcq', 'The number of bytes in the string', 'Количество байтов в строке', 'Сызықтағы байттар саны', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('8c74ee3c-c9a3-5b8b-9622-e1e5eee2ff27', 'mcq', 'The number of Unicode characters / runes', 'Количество символов Unicode / рун', 'Юникод таңбаларының / руналардың саны', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('8c74ee3c-c9a3-5b8b-9622-e1e5eee2ff27', 'mcq', 'The number of words in the string', 'Количество слов в строке', 'Сызықтағы сөздер саны', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('8c74ee3c-c9a3-5b8b-9622-e1e5eee2ff27', 'mcq', 'The memory capacity of the underlying string buffer', 'Объём памяти базового буфера строк', 'Негізгі жол буферінің жады сыйымдылығы', false);

-- Question 11: Strings & Runes (medium, mcq)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('02a892b5-329d-5b10-97ca-e75950364be4', 'mcq', 'medium', 'piscine', 'Strings & Runes', 'Строки и руны', 'Жолдар мен руналар', 'What happens when the following Go code is compiled and executed?

```go
score := 100
text := "Score: " + score
fmt.Println(text)
```', 'Что произойдет при компиляции и выполнении следующего кода на Go?

```go
score := 100
text := "Score: " + score
fmt.Println(text)
```', 'Төмендегі Go коды компиляцияланып, орындалғанда не болады?

```go
score := 100
text := "Score: " + score
fmt.Println(text)
```', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('02a892b5-329d-5b10-97ca-e75950364be4', 'mcq', 'It prints `Score: 100`', 'Выводится: `Score: 100`', 'Ол «`Score: 100`» деп басады.', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('02a892b5-329d-5b10-97ca-e75950364be4', 'mcq', 'Compile-time error: `cannot use score (variable of type int) as string value`', 'Ошибка компиляции: `cannot use score (variable of type int) as string value`', 'Құрастыру кезіндегі қате: `cannot use score (variable of type int) as string value`', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('02a892b5-329d-5b10-97ca-e75950364be4', 'mcq', 'It prints `Score: d`', 'Выводит: `Score: d`', 'Ол «`Score: d`» деп басады.', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('02a892b5-329d-5b10-97ca-e75950364be4', 'mcq', 'Runtime panic', 'Паника во время выполнения', 'Жүгіру кезіндегі паника', false);

-- Question 12: Strings & Runes (hard, text)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('0a551232-bf61-5159-bdbd-348692392a75', 'text', 'hard', 'piscine', 'Strings & Runes', 'Строки и руны', 'Жолдар мен руналар', 'What integer value is returned by `len("Привет")` in Go, given that each of the 6 Cyrillic characters requires 2 bytes in UTF-8 encoding?', 'Какое целое значение возвращает выражение `len("Привет")` в языке Go, учитывая, что каждый из 6 кириллических символов занимает 2 байта в кодировке UTF-8?', 'Go тілінде 6 кирилликалық таңбаның әрқайсысы UTF-8 кодтауында 2 байт алатынын ескерсек, `len("Привет")` функциясы қандай бүтін санды қайтарады?', true);
INSERT INTO problem_accepted_answers (problem_id, problem_type, answer_text_en, answer_text_ru, answer_text_kz)
VALUES ('0a551232-bf61-5159-bdbd-348692392a75', 'text', '12', '12', 'Он екі');
INSERT INTO problem_accepted_answers (problem_id, problem_type, answer_text_en, answer_text_ru, answer_text_kz)
VALUES ('0a551232-bf61-5159-bdbd-348692392a75', 'text', '12 bytes', '12 байт', '12 байт');

-- Question 14: Conditions (medium, mcq)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('816d1062-1cc9-50e9-a496-4ba13544a3c5', 'mcq', 'medium', 'piscine', 'Conditions', 'Условия', 'Шарттар', 'What is the output of the following code snippet?

```go
x := 15
n := x % 4
if n == 0 {
    fmt.Println("Alpha")
} else if n == 3 {
    fmt.Println("Beta")
} else {
    fmt.Println("Gamma")
}
```', 'Каков результат выполнения следующего фрагмента кода?

```go
x := 15
n := x % 4
if n == 0 {
    fmt.Println("Alpha")
} else if n == 3 {
    fmt.Println("Beta")
} else {
    fmt.Println("Gamma")
}
```', 'Төмендегі код үзіндісінің шығысы қандай?

```go
x := 15
n := x % 4
if n == 0 {
    fmt.Println("Alpha")
} else if n == 3 {
    fmt.Println("Beta")
} else {
    fmt.Println("Gamma")
}
```', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('816d1062-1cc9-50e9-a496-4ba13544a3c5', 'mcq', '`Alpha`', '`Alpha`', '`Alpha`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('816d1062-1cc9-50e9-a496-4ba13544a3c5', 'mcq', '`Beta`', '`Beta`', '`Beta`', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('816d1062-1cc9-50e9-a496-4ba13544a3c5', 'mcq', '`Gamma`', '`Gamma`', '`Gamma`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('816d1062-1cc9-50e9-a496-4ba13544a3c5', 'mcq', 'Compile-time error', 'Ошибка компиляции', 'Құрастыру кезіндегі қате', false);

-- Question 15: Conditions (hard, mcq)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('7388b6c0-8c60-5eb0-9b5b-45033faf1193', 'mcq', 'hard', 'piscine', 'Conditions', 'Условия', 'Шарттар', 'Which of the following statements regarding `if` statements in Go are TRUE? (Select all that apply)', 'Какие из приведенных ниже утверждений относительно операторов «`if`» в языке Go являются верными? (Выберите все правильные варианты)', 'Go тіліндегі `if`-мәлімдемелері туралы келесі мәлімдемелердің қайсысы дұрыс? (Барлығын таңдаңыз)', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('7388b6c0-8c60-5eb0-9b5b-45033faf1193', 'mcq', 'The condition expression in an `if` statement must evaluate strictly to a `bool`.', 'Выражение условия в операторе «`if`» должно давать в результате строго значение типа «`bool`».', '`if` операторындағы шарттық өрнек міндетті түрде дәл `bool` мәніне бағалануы тиіс.', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('7388b6c0-8c60-5eb0-9b5b-45033faf1193', 'mcq', 'Non-zero integer numbers like `1` or non-empty strings are implicitly evaluated as true in `if` conditions.', 'Целые числа, отличные от нуля, такие как `1`, или непустые строки неявно оцениваются как «true» в условиях типа `if`.', '`1` сияқты нөлге тең емес бүтін сандар немесе бос емес жолдар `if` шарттарында жасырын түрде true деп есептеледі.', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('7388b6c0-8c60-5eb0-9b5b-45033faf1193', 'mcq', 'The body block of an `if` statement must always be enclosed in curly braces `{}`.', 'Тело оператора `if` всегда должно заключаться в фигурные скобки `{}`.', '`if` операторының денесі әрдайым дөңгелек жақшалармен қоршалуы тиіс `{}`.', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('7388b6c0-8c60-5eb0-9b5b-45033faf1193', 'mcq', 'Parentheses `()` around the condition expression in an `if` statement are mandatory.', 'Скобки `()`, заключающие выражение условия в операторе `if`, являются обязательными.', '`if` операторындағы шарттық өрнектің айналасындағы `()` жақшалары міндетті.', false);

-- Question 16: Loops (easy, mcq)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('861ecbea-c13f-5f0e-bbb4-fa77b2a668b4', 'mcq', 'easy', 'piscine', 'Loops', 'Циклы', 'Циклдер', 'Which loop keyword exists in Go to construct counting loops, while-style loops, and infinite loops?', 'Какое ключевое слово цикла существует в языке Go для построения циклов с подсчётом, циклов типа «while» и бесконечных циклов?', 'Go тілінде санау циклдарын, while-стильді циклдарын және шексіз циклдарын құру үшін қандай цикл кілттік сөзі бар?', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('861ecbea-c13f-5f0e-bbb4-fa77b2a668b4', 'mcq', '`while`', '`while`', '`while`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('861ecbea-c13f-5f0e-bbb4-fa77b2a668b4', 'mcq', '`for`', '`for`', '`for`', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('861ecbea-c13f-5f0e-bbb4-fa77b2a668b4', 'mcq', '`loop`', '`loop`', '`loop`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('861ecbea-c13f-5f0e-bbb4-fa77b2a668b4', 'mcq', '`repeat`', '`repeat`', '`repeat`', false);

-- Question 17: Loops (hard, mcq)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('d0134e22-3996-5674-a064-db6bfe5b1ac4', 'mcq', 'hard', 'piscine', 'Loops', 'Циклы', 'Циклдер', 'What will be printed when the following code executes?

```go
sum := 0
for i := 1; i <= 5; i++ {
    if i%2 == 0 {
        continue
    }
    sum += i
}
fmt.Println(sum)
```', 'Что будет выведено на экран при выполнении следующего кода?

```go
sum := 0
for i := 1; i <= 5; i++ {
    if i%2 == 0 {
        continue
    }
    sum += i
}
fmt.Println(sum)
```', 'Төмендегі код орындалғанда не басылады?

```go
sum := 0
for i := 1; i <= 5; i++ {
    if i%2 == 0 {
        continue
    }
    sum += i
}
fmt.Println(sum)
```', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('d0134e22-3996-5674-a064-db6bfe5b1ac4', 'mcq', '`6`', '`6`', '`6`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('d0134e22-3996-5674-a064-db6bfe5b1ac4', 'mcq', '`9`', '`9`', '`9`', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('d0134e22-3996-5674-a064-db6bfe5b1ac4', 'mcq', '`15`', '`15`', '`15`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('d0134e22-3996-5674-a064-db6bfe5b1ac4', 'mcq', '`4`', '`4`', '`4`', false);

-- Question 18: Loops (medium, text)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('94f96053-6107-50b4-beba-c81a65af0394', 'text', 'medium', 'piscine', 'Loops', 'Циклы', 'Циклдер', 'When iterating over a string using `for index, val := range "Go"`, what is the Go data type of the variable `val`?', 'При итерации по строке с использованием ` в качестве индекса: val := range "Go"`, каков тип данных переменной `val` в языке Go?', 'Стринг бойынша `for index, val := range "Go"` циклымен өту кезінде `val` айнымалының Go деректер типі қандай?', true);
INSERT INTO problem_accepted_answers (problem_id, problem_type, answer_text_en, answer_text_ru, answer_text_kz)
VALUES ('94f96053-6107-50b4-beba-c81a65af0394', 'text', 'rune', 'Руне', 'Руна');
INSERT INTO problem_accepted_answers (problem_id, problem_type, answer_text_en, answer_text_ru, answer_text_kz)
VALUES ('94f96053-6107-50b4-beba-c81a65af0394', 'text', 'int32', '32-разрядное целое число', '32-биттік бүтін сан');
INSERT INTO problem_accepted_answers (problem_id, problem_type, answer_text_en, answer_text_ru, answer_text_kz)
VALUES ('94f96053-6107-50b4-beba-c81a65af0394', 'text', 'rune (int32)', 'rune (int32)', 'руна (int32)');

-- Question 19: Arrays (easy, mcq)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('2ae45349-5079-5c7b-9209-e4f636a9f0f0', 'mcq', 'easy', 'piscine', 'Arrays', 'Массивы', 'Массивтер', 'Which of the following declarations correctly creates an array of 5 integers initialized to their zero values in Go?', 'Какой из приведённых ниже операторов в языке Go правильно создаёт массив из 5 целых чисел, инициализированных нулями?', 'Төмендегі жарияланымдардың қайсысы Go тілінде 5 бүтін саннан тұратын массивті олардың мәндерін нөлге теңей отырып дұрыс жасайды?', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('2ae45349-5079-5c7b-9209-e4f636a9f0f0', 'mcq', '`var a [5]int`', '`var a [5]int`', '`var a [5]int`', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('2ae45349-5079-5c7b-9209-e4f636a9f0f0', 'mcq', '`var a []int = [5]int{}`', '`var a []int = [5]int{}`', '`var a []int = [5]int{}`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('2ae45349-5079-5c7b-9209-e4f636a9f0f0', 'mcq', '`a := array(5, int)`', '`a := array(5, int)`', '`a := array(5, int)`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('2ae45349-5079-5c7b-9209-e4f636a9f0f0', 'mcq', '`var a int[5]`', '`var a int[5]`', '`var a int[5]`', false);

-- Question 20: Arrays (hard, mcq)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('73ccbf91-d46a-5c72-aedd-ca59e0149218', 'mcq', 'hard', 'piscine', 'Arrays', 'Массивы', 'Массивтер', 'What is the output of the following code?

```go
a := [3]int{1, 2, 3}
b := a
b[0] = 99
fmt.Println(a[0], b[0])
```', 'Каков результат выполнения следующего кода?

```go
a := [3]int{1, 2, 3}
b := a
b[0] = 99
fmt.Println(a[0], b[0])
```', 'Төмендегі кодтың шығысы қандай?

```go
a := [3]int{1, 2, 3}
b := a
b[0] = 99
fmt.Println(a[0], b[0])
```', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('73ccbf91-d46a-5c72-aedd-ca59e0149218', 'mcq', '`99 99`', '`99 99`', '`99 99`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('73ccbf91-d46a-5c72-aedd-ca59e0149218', 'mcq', '`1 99`', '`1 99`', '`1 99`', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('73ccbf91-d46a-5c72-aedd-ca59e0149218', 'mcq', '`1 1`', '`1 1`', '`1 1`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('73ccbf91-d46a-5c72-aedd-ca59e0149218', 'mcq', 'Compile-time error', 'Ошибка компиляции', 'Құрастыру кезіндегі қате', false);

-- Question 21: Arrays (hard, mcq)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('532a412c-49ff-5c92-ba8c-04f4c36e6386', 'mcq', 'hard', 'piscine', 'Arrays', 'Массивы', 'Массивтер', 'Which of the following statements about arrays in Go are TRUE? (Select all that apply)', 'Какие из приведённых ниже утверждений о массивах в языке Go являются верными? (Выберите все подходящие варианты)', 'Go тіліндегі массивтер туралы келесі мәлімдемелердің қайсысы дұрыс? (Барлығын таңдаңыз)', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('532a412c-49ff-5c92-ba8c-04f4c36e6386', 'mcq', 'The length of an array is fixed and is part of its type (e.g. `[3]int` and `[4]int` are distinct types).', 'Длина массива фиксирована и является частью его типа (например, `[3]int` и `[4]int` — это разные типы).', 'Массивтің ұзындығы тұрақты және оның типінің құрамдас бөлігі болып табылады (мысалы, `[3]int` және `[4]int` – әртүрлі типтер).', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('532a412c-49ff-5c92-ba8c-04f4c36e6386', 'mcq', 'Two arrays of identical type and length containing comparable elements can be compared with `==`.', 'Два массива одного и того же типа и длины, содержащие сопоставимые элементы, можно сравнить с помощью операторов `==`.', 'Түрі мен ұзындығы бірдей және салыстырылатын элементтерден тұратын екі массивті `==` операторымен салыстыруға болады.', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('532a412c-49ff-5c92-ba8c-04f4c36e6386', 'mcq', 'Arrays can dynamically expand in length at runtime using `append()`.', 'Длина массивов может динамически увеличиваться во время выполнения с помощью `append()`.', 'Массивтер орындау уақытында `append()` әдісін пайдаланып ұзындығын динамикалық түрде ұлғайта алады.', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('532a412c-49ff-5c92-ba8c-04f4c36e6386', 'mcq', 'Passing an array as an argument to a function creates a complete copy of the array unless passed by pointer.', 'При передаче массива в качестве аргумента функции создается полная копия массива, если только он не передается по указателю.', 'Функцияға массивті аргумент ретінде бергенде, оны көрсеткіш арқылы бермесеңіз, массивтің толық көшірмесі жасалады.', true);

-- Question 22: Slices (easy, mcq)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('f1699759-82cd-52ff-a01a-348c297a9423', 'mcq', 'easy', 'piscine', 'Slices', 'Слайсы', 'Тілімдер', 'Which built-in function is used to append elements to the end of a slice in Go?', 'Какая встроенная функция используется в Go для добавления элементов в конец фрагмента?', 'Go-да массивтің бір бөлігінің соңына элементтерді қосу үшін қай кіріктірілген функция қолданылады?', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('f1699759-82cd-52ff-a01a-348c297a9423', 'mcq', '`push()`', '`push()`', '`push()`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('f1699759-82cd-52ff-a01a-348c297a9423', 'mcq', '`append()`', '`append()`', '`append()`', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('f1699759-82cd-52ff-a01a-348c297a9423', 'mcq', '`insert()`', '`insert()`', '`insert()`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('f1699759-82cd-52ff-a01a-348c297a9423', 'mcq', '`extend()`', '`extend()`', '`extend()`', false);

-- Question 23: Slices (hard, mcq)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('0bcc97c5-1100-51d2-b435-8d446252086e', 'mcq', 'hard', 'piscine', 'Slices', 'Слайсы', 'Тілімдер', 'What is printed by the following code snippet?

```go
s := []int{10, 20, 30, 40}
sub := s[1:3]
sub[0] = 99
fmt.Println(s[1])
```', 'Что выведет на экран следующий фрагмент кода?

```go
s := []int{10, 20, 30, 40}
sub := s[1:3]
sub[0] = 99
fmt.Println(s[1])
```', 'Төмендегі код үзіндісі нені басып шығарады?

```go
s := []int{10, 20, 30, 40}
sub := s[1:3]
sub[0] = 99
fmt.Println(s[1])
```', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('0bcc97c5-1100-51d2-b435-8d446252086e', 'mcq', '`20`', '`20`', '`20`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('0bcc97c5-1100-51d2-b435-8d446252086e', 'mcq', '`99`', '`99`', '`99`', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('0bcc97c5-1100-51d2-b435-8d446252086e', 'mcq', '`10`', '`10`', '`10`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('0bcc97c5-1100-51d2-b435-8d446252086e', 'mcq', '`30`', '`30`', '`30`', false);

-- Question 24: Slices (hard, mcq)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('d23c21ab-b552-50a3-8861-039c72c62361', 'mcq', 'hard', 'piscine', 'Slices', 'Слайсы', 'Тілімдер', 'Which of the following operations are valid when working with slices in Go? (Select all that apply)', 'Какие из перечисленных ниже операций допустимы при работе со срезами в языке Go? (Выберите все подходящие варианты)', 'Go-да кесінділермен жұмыс істегенде төмендегі операциялардың қайсысы жарамды? (Барлығын таңдаңыз)', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('d23c21ab-b552-50a3-8861-039c72c62361', 'mcq', 'Adding new elements to a slice using `append(slice, element)`.', 'Добавление новых элементов в фрагмент с помощью команды «`append(slice, element)`».', '`append(slice, element)` арқылы слайсқа жаңа элементтер қосу.', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('d23c21ab-b552-50a3-8861-039c72c62361', 'mcq', 'Getting the number of elements in a slice using `len(slice)`.', 'Определение количества элементов в фрагменте с помощью функции ``len(slice)``.', '`len(slice)` қолдана отырып, кесіндідегі элементтер санын алу.', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('d23c21ab-b552-50a3-8861-039c72c62361', 'mcq', 'Directly comparing two slices with multiple elements using the `==` operator (e.g. `slice1 == slice2`).', 'Прямое сравнение двух срезов, содержащих несколько элементов, с помощью оператора «`==`» (например, `slice1 == slice2`).', '`==` операторы арқылы бірнеше элементі бар екі тілімді тікелей салыстыру (мысалы: `slice1 == slice2`).', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('d23c21ab-b552-50a3-8861-039c72c62361', 'mcq', 'Creating a sub-slice from an existing slice using sub-slicing syntax `slice[start:end]`.', 'Создание подсреза на основе существующего среза с использованием синтаксиса подсреза `slice[start:end]`.', '`slice[start:end]` синтаксисін пайдаланып, бар слайстан кіші слайс жасау.', true);

-- Question 28: Pointers (easy, mcq)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('68b9456d-2d99-5c84-b0ec-bf1065c3eae4', 'mcq', 'easy', 'piscine', 'Pointers', 'Указатели', 'Көрсеткіштер', 'Which operator is used to obtain the memory address of an existing variable `x` in Go?', 'Какой оператор используется в языке Go для получения адреса в памяти существующей переменной `x`?', 'Go тілінде бар `x` айнымалының жады мекенжайын алу үшін қандай оператор қолданылады?', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('68b9456d-2d99-5c84-b0ec-bf1065c3eae4', 'mcq', '`*x`', '`*x`', '`*x`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('68b9456d-2d99-5c84-b0ec-bf1065c3eae4', 'mcq', '`&x`', '`&x`', '`&x`', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('68b9456d-2d99-5c84-b0ec-bf1065c3eae4', 'mcq', '`@x`', '`@x`', '`@x`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('68b9456d-2d99-5c84-b0ec-bf1065c3eae4', 'mcq', '`addr(x)`', '`addr(x)`', '`addr(x)`', false);

-- Question 29: Pointers (medium, mcq)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('6d3b1f2e-b93b-5a09-8907-aaa3a5a48b66', 'mcq', 'medium', 'piscine', 'Pointers', 'Указатели', 'Көрсеткіштер', 'What is printed by the following code?

```go
x := 42
p := &x
*p = 100
fmt.Println(x)
```', 'Что выведется при выполнении следующего кода?

```go
x := 42
p := &x
*p = 100
fmt.Println(x)
```', 'Төмендегі код нені басып шығарады?

```go
x := 42
p := &x
*p = 100
fmt.Println(x)
```', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('6d3b1f2e-b93b-5a09-8907-aaa3a5a48b66', 'mcq', '`42`', '`42`', '`42`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('6d3b1f2e-b93b-5a09-8907-aaa3a5a48b66', 'mcq', '`100`', '`100`', '`100`', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('6d3b1f2e-b93b-5a09-8907-aaa3a5a48b66', 'mcq', 'The memory address of `x`', 'Адрес в памяти `x`', 'x-тің жад мекенжайы', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('6d3b1f2e-b93b-5a09-8907-aaa3a5a48b66', 'mcq', '`nil`', '`nil`', '`nil`', false);

-- Question 30: Pointers (hard, mcq)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('fc2bb4c7-336a-53f1-811e-4565b2c367c9', 'mcq', 'hard', 'piscine', 'Pointers', 'Указатели', 'Көрсеткіштер', 'What happens when executing the following code in Go? (Select all that apply)

```go
var p *int
fmt.Println(*p)
```', 'Что произойдет при выполнении следующего кода на языке Go? (Выберите все подходящие варианты)

```go
var p *int
fmt.Println(*p)
```', 'Go тілінде келесі кодты орындағанда не болады? (Барлығына сәйкес келетіндерін таңдаңыз)

```go
var p *int
fmt.Println(*p)
```', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('fc2bb4c7-336a-53f1-811e-4565b2c367c9', 'mcq', 'A runtime panic occurs during execution.', 'Во время выполнения программы возникает паника.', 'Орындау кезінде runtime паникасы пайда болады.', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('fc2bb4c7-336a-53f1-811e-4565b2c367c9', 'mcq', 'The panic message contains `invalid memory address or nil pointer dereference`.', 'Сообщение о панике содержит `неверный адрес памяти или разыменование нулевого указателя`.', 'Паникалық хабарлама `жарамсыз жад мекенжайы немесе nil көрсеткішін дереференциациялау` дегенді қамтиды.', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('fc2bb4c7-336a-53f1-811e-4565b2c367c9', 'mcq', 'It safely outputs `0`.', 'Он корректно выводит `0`.', 'Ол қауіпсіз түрде 0 шығарады.', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('fc2bb4c7-336a-53f1-811e-4565b2c367c9', 'mcq', 'The Go compiler rejects the code at compile-time with an uninitialized pointer error.', 'Компилятор Go отклоняет этот код на этапе компиляции, выдавая ошибку «неинициализированный указатель».', 'Go компиляторы кодты компиляция кезінде бастапқыланбаған көрсеткіш қатесімен қабылдамайды.', false);

-- Question 31: Structs & Methods (easy, mcq)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('0cd08bc3-0ff3-53cd-8c50-1ee3c0fab8ff', 'mcq', 'easy', 'piscine', 'Structs & Methods', 'Структуры и методы', 'Құрылымдар мен әдістер', 'Which syntax correctly declares a struct type named `Book` containing a string field `Title` in Go?', 'Какой синтаксис правильно объявляет в языке Go тип структуры с именем `Book`, содержащий строковое поле `Title`?', 'Go тілінде `Book` атты struct типінің `Title` атты жол өрісін дұрыс жариялайтын қай синтаксис?', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('0cd08bc3-0ff3-53cd-8c50-1ee3c0fab8ff', 'mcq', '`type Book struct { Title string }`', '`type Book struct { Title string }`', '`type Book struct { Title string }`', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('0cd08bc3-0ff3-53cd-8c50-1ee3c0fab8ff', 'mcq', '`struct Book { Title string }`', '`struct Book { Title string }`', '`struct Book { Title string }`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('0cd08bc3-0ff3-53cd-8c50-1ee3c0fab8ff', 'mcq', '`class Book { Title string }`', '`class Book { Title string }`', '`class Book { Title string }`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('0cd08bc3-0ff3-53cd-8c50-1ee3c0fab8ff', 'mcq', '`var Book = struct { Title string }`', '`var Book = struct { Title string }`', '`var Book = struct { Title string }`', false);

-- Question 34: Errors & Error Handling (easy, mcq)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('70081cca-957d-5c40-8f57-b10f2a5ddf52', 'mcq', 'easy', 'piscine', 'Errors & Error Handling', 'Ошибки и обработка ошибок', 'Қателер және қателерді өңдеу', 'What value does an idiomatic Go function return for its `error` return value to indicate success (no error)?', 'Какое значение возвращает идиоматическая функция на Go в качестве значения возврата `error`, чтобы указать успешное выполнение (отсутствие ошибки)?', 'Идиоматикалық Go функциясы `error` нәтижесінде сәтті (қатесіз) болғанын көрсету үшін қандай мән қайтарады?', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('70081cca-957d-5c40-8f57-b10f2a5ddf52', 'mcq', '`""` (empty string)', '`""` (пустая строка)', '"" (бос жол)', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('70081cca-957d-5c40-8f57-b10f2a5ddf52', 'mcq', '`0`', '`0`', '`0`', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('70081cca-957d-5c40-8f57-b10f2a5ddf52', 'mcq', '`nil`', '`nil`', '`nil`', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('70081cca-957d-5c40-8f57-b10f2a5ddf52', 'mcq', '`false`', '`false`', '`false`', false);

-- Question 36: Errors & Error Handling (hard, mcq)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('b2982305-c6d0-5a64-91b4-52d393a32cf4', 'mcq', 'hard', 'piscine', 'Errors & Error Handling', 'Ошибки и обработка ошибок', 'Қателер және қателерді өңдеу', 'Which of the following statements about errors and error handling in Go are TRUE? (Select all that apply)', 'Какие из приведённых ниже утверждений об ошибках и обработке ошибок в языке Go являются ПРАВИЛЬНЫМИ? (Выберите все подходящие варианты)', 'Go тіліндегі қателер мен қателерді өңдеу туралы келесі мәлімдемелердің қайсысы дұрыс? (Барлығын таңдаңыз)', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('b2982305-c6d0-5a64-91b4-52d393a32cf4', 'mcq', 'Functions that can produce an error conventionally return the `error` as their last return value.', 'Функции, которые могут вызвать ошибку, по соглашению возвращают значение типа «`error`» в качестве последнего возвращаемого значения.', 'Қате тудыруы мүмкін функциялар дәстүрлі түрде соңғы нәтиже ретінде `error`-ті қайтарады.', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('b2982305-c6d0-5a64-91b4-52d393a32cf4', 'mcq', 'The condition `if err != nil` is the standard idiomatic way to check if an error occurred.', 'Условие «`if err != nil`» — это стандартный идиоматический способ проверки наличия ошибки.', '`if err != nil` – қатенің болғанын тексерудің стандартты идиоматикалық тәсілі.', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('b2982305-c6d0-5a64-91b4-52d393a32cf4', 'mcq', 'A new custom error value can be created using the `errors.New()` function.', 'Новое пользовательское значение ошибки можно создать с помощью функции ``errors.New()``.', '`errors.New()` функциясын пайдаланып жаңа теңшелетін қате мәнін жасауға болады.', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('b2982305-c6d0-5a64-91b4-52d393a32cf4', 'mcq', 'Go uses structured `try`, `catch`, and `finally` blocks for exception handling.', 'В языке Go для обработки исключений используются структурированные блоки `try`, `catch` и `finally`.', 'Go ерекшеліктерді өңдеу үшін құрылымдық `try`, `catch` және `finally` блоктерін пайдаланады.', false);

-- Question 37: Command-Line Arguments & Stdin (easy, mcq)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('4050fcb1-45ea-58f1-be33-49e077373943', 'mcq', 'easy', 'piscine', 'Command-Line Arguments & Stdin', 'Аргументы командной строки и Stdin', 'Командалық жол аргументтері және Stdin', 'What does `os.Args[0]` represent when executing a compiled Go command-line program?', 'Что означает `os.Args[0]` при запуске скомпилированной программы Go, запускаемой из командной строки?', 'Құрастырылған Go командалық жолдағы бағдарламаны орындау кезінде `os.Args[0]` нені білдіреді?', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('4050fcb1-45ea-58f1-be33-49e077373943', 'mcq', 'The first user argument passed after the executable name', 'Первый аргумент пользователя, передаваемый после имени исполняемого файла', 'Орындалатын файл атынан кейін берілетін бірінші пайдаланушы аргументі', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('4050fcb1-45ea-58f1-be33-49e077373943', 'mcq', 'The path or command used to invoke the program itself', 'Путь или команда, используемая для запуска самой программы', 'Бағдарламаны шақыру үшін қолданылатын жол немесе команда', true);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('4050fcb1-45ea-58f1-be33-49e077373943', 'mcq', 'The total number of arguments passed to the program', 'Общее количество аргументов, переданных программе', 'Бағдарламаға берілген аргументтердің жалпы саны', false);
INSERT INTO problem_options (problem_id, problem_type, text_en, text_ru, text_kz, is_correct)
VALUES ('4050fcb1-45ea-58f1-be33-49e077373943', 'mcq', 'An empty string', 'Пустая строка', 'Бос жол', false);

-- Question 39: Command-Line Arguments & Stdin (hard, text)
INSERT INTO problems (id, type, difficulty, topic, title_en, title_ru, title_kz, prompt_en, prompt_ru, prompt_kz, is_published)
VALUES ('b42f4a8a-6792-5500-932e-fb93ddc89af5', 'text', 'hard', 'piscine', 'Command-Line Arguments & Stdin', 'Аргументы командной строки и Stdin', 'Командалық жол аргументтері және Stdin', 'If a compiled Go program is invoked in the terminal as `./cli process -v file.txt`, what is the exact integer value returned by `len(os.Args)`?', 'Если скомпилированная программа на языке Go запускается в терминале следующим образом: `./cli process -v file.txt`, то какое точное целое значение возвращает функция `len(os.Args)`?', 'Егер компиляцияланған Go бағдарламасы терминалда `./cli process -v file.txt` ретінде іске қосылса, `len(os.Args)` қандай дәл бүтін санды қайтарады?', true);
INSERT INTO problem_accepted_answers (problem_id, problem_type, answer_text_en, answer_text_ru, answer_text_kz)
VALUES ('b42f4a8a-6792-5500-932e-fb93ddc89af5', 'text', '4', '4', '4');
INSERT INTO problem_accepted_answers (problem_id, problem_type, answer_text_en, answer_text_ru, answer_text_kz)
VALUES ('b42f4a8a-6792-5500-932e-fb93ddc89af5', 'text', 'four', 'четыре', 'төрт');

COMMIT;
