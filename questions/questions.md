# Go Programming Course: Comprehensive Assessment Questions

## Overview
- **Total Questions:** 31
- **Course Quests Covered:** All active core modules from `go/` and `go-content/`
- **Difficulties & Time Limits:**
  - **Easy:** 30 seconds
  - **Medium:** 45 seconds
  - **Hard:** 60 seconds
- **Question Types:** Multiple Choice Questions (Single Choice & Multiple Selection) and Text Questions with Multiple Possible Answers
- **Languages:** English, Russian (Русский), Kazakh (Қазақша) — translated via DeepL

---

## Question 1
- **Topic:** Printing & Program Structure
- **Difficulty:** Easy (30s)
- **Format:** Multiple Choice (Single Answer)

### English Version
**Question:** Which package declaration is required at the top of a Go file to produce a standalone executable program?

**Options:**
- [ ] **A)** `package run`
- [x] **B)** `package main`
- [ ] **C)** `package root`
- [ ] **D)** `package exec`

**Correct Answer:** B
**Explanation:** In Go, every standalone executable program must begin with `package main` and contain a `main()` function as its entry point.

### Русская версия (Russian)
**Вопрос:** Какую декларацию пакета необходимо разместить в начале файла на языке Go, чтобы получить автономную исполняемую программу?

**Варианты ответа:**
- [ ] **A)** `package run`
- [x] **B)** `package main`
- [ ] **C)** `package root`
- [ ] **D)** `package exec`

**Правильный ответ:** B
**Объяснение:** В языке Go каждая автономная исполняемая программа должна начинаться с `package main` и содержать функцию `main()` в качестве точки входа.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Go файлының басында жеке орындалатын бағдарлама алу үшін қандай пакет жариялауы қажет?

**Жауап нұсқалары:**
- [ ] **A)** `package run`
- [x] **B)** `package main`
- [ ] **C)** `package root`
- [ ] **D)** `package exec`

**Дұрыс жауап:** B
**Түсіндірме:** Go-да әрбір өзіндік орындалатын бағдарлама `package main`-нен басталып, кіру нүктесі ретінде `main()` функциясын қамтуы тиіс.

---
## Question 2
- **Topic:** Printing & Program Structure
- **Difficulty:** Medium (45s)
- **Format:** Multiple Choice (Multiple Correct Answers)

### English Version
**Question:** Which of the following statements regarding the function `fmt.Println()` in Go are TRUE? (Select all that apply)

**Options:**
- [x] **A)** It automatically separates multiple arguments with spaces.
- [x] **B)** It automatically appends a newline character (`\n`) at the end of the printed output.
- [ ] **C)** It can only accept arguments of type `string`.
- [x] **D)** Its name starts with a capital 'P' because it is exported from the `fmt` package.

**Correct Answer:** A, B, D
**Explanation:** `fmt.Println` accepts values of any type, formats them separated by spaces, adds a newline at the end, and must start with an uppercase letter to be exported from package `fmt`.

### Русская версия (Russian)
**Вопрос:** Какие из приведённых ниже утверждений относительно функции `fmt.Println()` в языке программирования Go являются верными? (Выберите все правильные варианты)

**Варианты ответа:**
- [x] **A)** Он автоматически разделяет несколько аргументов пробелами.
- [x] **B)** Она автоматически добавляет символ новой строки (`\n`) в конец выводимого результата.
- [ ] **C)** Он может принимать только аргументы типа `string`.
- [x] **D)** Его название начинается с заглавной буквы «P», поскольку он экспортируется из пакета `fmt`.

**Правильный ответ:** A, B, D
**Объяснение:** `fmt.Println` принимает значения любого типа, форматирует их с разделителями в виде пробелов, добавляет символ новой строки в конце и должен начинаться с заглавной буквы, чтобы его можно было экспортировать из пакета `fmt`.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Go тіліндегі `fmt.Println()` функциясына қатысты келесі мәлімдемелердің қайсысы дұрыс? (Барлығын таңдаңыз)

**Жауап нұсқалары:**
- [x] **A)** Ол бірнеше аргументтерді автоматты түрде бос орындармен бөледі.
- [x] **B)** Ол басылып шығарылған нәтиженің соңына автоматты түрде жаңа жол таңбасын (`\n`) қосады.
- [ ] **C)** Ол тек `string` типіндегі аргументтерді қабылдай алады.
- [x] **D)** Оның атауы `fmt` пакетінен экспортталғандықтан үлкен «P» әрпімен басталады.

**Дұрыс жауап:** A, B, D
**Түсіндірме:** fmt.Println кез келген типтегі мәндерді қабылдайды, оларды бос орындармен бөліп форматтайды, соңына жаңа жол белгісін қосады және `fmt` пакетінен экспортталу үшін алғашқы әрпі бас әріп болуы тиіс.

---
## Question 3
- **Topic:** Printing & Program Structure
- **Difficulty:** Hard (60s)
- **Format:** Text Question with Multiple Possible Answers

### English Version
**Question:** In Go, what letter casing (uppercase or lowercase) must the first letter of a function name have to be exported and callable from another package (like `Println` in package `fmt`)?

**Acceptable / Possible Answers:**
- `uppercase`
- `capital`
- `capital letter`
- `uppercase letter`
- `upper`

**Standard Answer:** uppercase / capital letter
**Explanation:** In Go, visibility outside a package is determined by capitalization. Identifiers starting with an uppercase (capital) letter are exported (public); lowercase identifiers remain private to the package.

### Русская версия (Russian)
**Вопрос:** В языке Go с какой буквы (заглавной или строчной) должно начинаться имя функции, чтобы она была экспортируемой (доступной) из другого пакета (например, `Println` в пакете `fmt`)?

**Приемлемые / возможные варианты ответа:**
- `заглавная буква`
- `заглавная`
- `с заглавной буквы`
- `прописная буква`
- `верхний регистр`

**Основной ответ:** заглавная буква (uppercase)
**Объяснение:** В языке Go видимость за пределами пакета определяется регистром букв. Идентификаторы, начинающиеся с заглавной буквы, экспортируются (являются общедоступными); идентификаторы, написанные строчными буквами, остаются частными для данного пакета.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Go тілінде функцияның атауы басқа пакеттен (мысалы, `fmt` пакетіндегі `Println`) экспортталып шақырылуы үшін оның бірінші әрпі қандай (бас әріп немесе кіші әріп) болуы керек?

**Қабылданатын / мүмкін жауап нұсқалары:**
- `бас әріп`
- `бас әріппен`
- `үлкен әріп`
- `жоғарғы регистр`

**Негізгі жауап:** бас әріп (uppercase)
**Түсіндірме:** Go тілінде пакеттен тыс көрінімділік бас әріппен белгіленеді. Бас әріппен басталатын идентификаторлар экспортталады (жария); кіші әріппен басталатын идентификаторлар пакет ішінде жекеше болып қалады.

---

## Question 4
- **Topic:** Variables, Constants & Types
- **Difficulty:** Easy (30s)
- **Format:** Multiple Choice (Single Answer)

### English Version
**Question:** What is the default zero value of an uninitialized variable declared as `var count int` in Go?

**Options:**
- [ ] **A)** `nil`
- [ ] **B)** `-1`
- [x] **C)** `0`
- [ ] **D)** `undefined`

**Correct Answer:** C
**Explanation:** In Go, variables declared without explicit initial values are automatically initialized to their zero value. For integer types, the zero value is `0`.

### Русская версия (Russian)
**Вопрос:** Каково значение по умолчанию (нуль) для неинициализированной переменной, объявленной в языке Go как `var count int`?

**Варианты ответа:**
- [ ] **A)** `nil`
- [ ] **B)** `-1`
- [x] **C)** `0`
- [ ] **D)** `undefined`

**Правильный ответ:** C
**Объяснение:** В языке Go переменные, объявленные без явных начальных значений, автоматически инициализируются нулевым значением. Для целочисленных типов нулевое значение — `0`.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Go тілінде `var count int` деп жарияланған бастапқы мәні нөлге тең емес айнымалының әдепкі нөлдік мәні қандай?

**Жауап нұсқалары:**
- [ ] **A)** `nil`
- [ ] **B)** `-1`
- [x] **C)** `0`
- [ ] **D)** `undefined`

**Дұрыс жауап:** C
**Түсіндірме:** Go тілінде нақты бастапқы мәні көрсетілмеген айнымалылар автоматты түрде нөлге теңестіріледі. Бүтін сан типтері үшін нөл мәні 0-ге тең.

---
## Question 5
- **Topic:** Variables, Constants & Types
- **Difficulty:** Medium (45s)
- **Format:** Multiple Choice (Single Answer)

### English Version
**Question:** What happens when the following code snippet is compiled in Go?

```go
x := 10
x := 20
fmt.Println(x)
```

**Options:**
- [ ] **A)** It prints `20`
- [ ] **B)** It prints `10`
- [x] **C)** Compile-time error: `no new variables on left side of :=`
- [ ] **D)** Runtime panic

**Correct Answer:** C
**Explanation:** The short variable declaration operator `:=` requires at least one new variable on its left-hand side within the current scope. Reassigning an existing variable must use the assignment operator `=`.

### Русская версия (Russian)
**Вопрос:** Что произойдет, если скомпилировать следующий фрагмент кода на языке Go?

```go
x := 10
x := 20
fmt.Println(x)
```

**Варианты ответа:**
- [ ] **A)** Выводится `20`
- [ ] **B)** Выводится `10`
- [x] **C)** Ошибка компиляции: `нет новых переменных в левой части выражения :=`
- [ ] **D)** Паника во время выполнения

**Правильный ответ:** C
**Объяснение:** Оператор краткого объявления переменных `:=` требует наличия по крайней мере одной новой переменной в левой части выражения в пределах текущей области видимости. Для переопределения существующей переменной необходимо использовать оператор присваивания `=`.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Төмендегі код үзіндісі Go тілінде компиляцияланғанда не болады?

```go
x := 10
x := 20
fmt.Println(x)
```

**Жауап нұсқалары:**
- [ ] **A)** Ол 20-ды басып шығарады.
- [ ] **B)** Ол 10 деп шығарады.
- [x] **C)** Құрастыру кезіндегі қате: `=:= операторының сол жақ жағында жаңа айнымалылар жоқ`
- [ ] **D)** Жүгіру кезіндегі паника

**Дұрыс жауап:** C
**Түсіндірме:** Қысқа айнымалы жариялау операторы `:=` ағымдағы ауқымда сол жақ жағында кем дегенде бір жаңа айнымалы болуын талап етеді. Бар айнымалыны қайта тағайындау үшін тағайындау операторы `=` қолданылуы тиіс.

---
## Question 6
- **Topic:** Variables, Constants & Types
- **Difficulty:** Hard (60s)
- **Format:** Text Question with Multiple Possible Answers

### English Version
**Question:** Which keyword in Go is used to declare an immutable value whose value is fixed at compile time and cannot be reassigned during execution?

**Acceptable / Possible Answers:**
- `const`
- `const keyword`

**Standard Answer:** const
**Explanation:** The `const` keyword declares compile-time constants in Go. Once declared, constant values cannot be reassigned or modified at runtime.

### Русская версия (Russian)
**Вопрос:** Какое ключевое слово в языке Go используется для объявления неизменяемого значения, величина которого фиксируется на этапе компиляции и не может быть переопределена во время выполнения?

**Приемлемые / возможные варианты ответа:**
- `const`
- `ключевое слово `const``

**Основной ответ:** const
**Объяснение:** Ключевое слово `const` используется в языке Go для объявления констант, определяемых на этапе компиляции. После объявления значения констант нельзя пересчитывать или изменять во время выполнения.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Go тілінде компиляция кезінде мәні бекітіліп, орындау кезінде қайта тағайындалмайтын өзгермейтін мәнді жариялау үшін қандай кілтсөз қолданылады?

**Қабылданатын / мүмкін жауап нұсқалары:**
- `тұрақты`
- `const кілтсөзі`

**Негізгі жауап:** const
**Түсіндірме:** Go тіліндегі `const` кілтсөзі компиляция кезіндегі тұрақтыларды жариялайды. Бір рет жарияланғаннан кейін тұрақты мәндерді орындау уақытында қайта тағайындауға немесе өзгертуге болмайды.

---
## Question 7
- **Topic:** Functions & Operators
- **Difficulty:** Easy (30s)
- **Format:** Multiple Choice (Single Answer)

### English Version
**Question:** What is the evaluated result of the integer division expression `7 / 2` in Go?

**Options:**
- [ ] **A)** `3.5`
- [x] **B)** `3`
- [ ] **C)** `4`
- [ ] **D)** Compile-time error

**Correct Answer:** B
**Explanation:** In Go, dividing two integer values performs integer division, which truncates any decimal portion toward zero, yielding `3`.

### Русская версия (Russian)
**Вопрос:** Каков результат вычисления выражения целочисленного деления `7 / 2` в языке Go?

**Варианты ответа:**
- [ ] **A)** `3.5`
- [x] **B)** `3`
- [ ] **C)** `4`
- [ ] **D)** Ошибка компиляции

**Правильный ответ:** B
**Объяснение:** В языке Go при делении двух целочисленных значений выполняется целочисленное деление, при котором десятичная часть округляется в сторону нуля, в результате чего получается `3`.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Go тілінде 7 / 2 бүтін сандық бөлу операторының бағаланған нәтижесі қандай?

**Жауап нұсқалары:**
- [ ] **A)** `3.5`
- [x] **B)** `3`
- [ ] **C)** `4`
- [ ] **D)** Құрастыру кезіндегі қате

**Дұрыс жауап:** B
**Түсіндірме:** Go-да екі бүтін санды бөлу бүтін сандық бөлуді орындайды, ол ондық бөлшекті нөлге қарай кесіп тастап, нәтижесінде 3 береді.

---
## Question 8
- **Topic:** Functions & Operators
- **Difficulty:** Medium (45s)
- **Format:** Multiple Choice (Single Answer)

### English Version
**Question:** What does the following function return when called as `calc(14, 4)`?

```go
func calc(a, b int) (int, int) {
    return a / b, a % b
}
```

**Options:**
- [x] **A)** `3, 2`
- [ ] **B)** `3.5, 2`
- [ ] **C)** `2, 3`
- [ ] **D)** `4, 2`

**Correct Answer:** A
**Explanation:** Integer division `14 / 4` produces `3`, and the modulo operator `14 % 4` produces the remainder `2`.

### Русская версия (Russian)
**Вопрос:** Какое значение возвращает следующая функция при вызове `calc(14, 4)`?

```go
func calc(a, b int) (int, int) {
    return a / b, a % b
}
```

**Варианты ответа:**
- [x] **A)** `3, 2`
- [ ] **B)** `3.5, 2`
- [ ] **C)** `2, 3`
- [ ] **D)** `4, 2`

**Правильный ответ:** A
**Объяснение:** Целочисленное деление `14 / 4` даёт в результате `3`, а оператор модуля `14 % 4` даёт остаток `2`.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Төмендегі функция `calc(14, 4)` деп шақырылғанда қандай мәнін қайтарады?

```go
func calc(a, b int) (int, int) {
    return a / b, a % b
}
```

**Жауап нұсқалары:**
- [x] **A)** `3, 2`
- [ ] **B)** `3.5, 2`
- [ ] **C)** `2, 3`
- [ ] **D)** `4, 2`

**Дұрыс жауап:** A
**Түсіндірме:** Бүтін сандық бөлу `14 / 4` нәтижесінде `3` шығады, ал модуль операторы `14 % 4` қалдықты `2` береді.

---
## Question 9
- **Topic:** Functions & Operators
- **Difficulty:** Hard (60s)
- **Format:** Multiple Choice (Multiple Correct Answers)

### English Version
**Question:** Which of the following code snippets will produce a COMPILE-TIME error in Go? (Select all that apply)

**Options:**
- [x] **A)** `var a int = 5; var b float64 = 2.0; var c = a + b`
- [ ] **B)** `var a int = 5; var b float64 = 2.0; var c = float64(a) + b`
- [ ] **C)** `var x int = 10; _ = x`
- [x] **D)** `func f() int { if false { return 1 } }`

**Correct Answer:** A, D
**Explanation:** In (A), Go does not allow arithmetic operations between different types without explicit conversion (mismatched types). In (D), the function has a return type of `int`, but the control flow can reach the end without a `return` statement, causing a 'missing return' compile error. (B) uses explicit casting, and (C) silences the unused variable check via `_`.

### Русская версия (Russian)
**Вопрос:** Какой из приведённых ниже фрагментов кода вызовет ошибку на этапе компиляции в языке Go? (Выберите все правильные варианты)

**Варианты ответа:**
- [x] **A)** `var a int = 5; var b float64 = 2.0; var c = a + b`
- [ ] **B)** `var a int = 5; var b float64 = 2.0; var c = float64(a) + b`
- [ ] **C)** `var x int = 10; _ = x`
- [x] **D)** `func f() int { if false { return 1 } }`

**Правильный ответ:** A, D
**Объяснение:** В варианте (A) язык Go не допускает арифметических операций между типами разных классов без явного преобразования (несовместимые типы). В варианте (D) функция имеет тип возвращаемого значения `int`, но поток управления может достичь конца без оператора `return`, что приводит к ошибке компиляции «missing return». В варианте (B) используется явное приведение типов, а в варианте (C) проверка на наличие неиспользуемых переменных подавляется с помощью `_`.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Төмендегі код үзінділерінің қайсысы Go тілінде компиляция кезінде қате тудырады? (Барлығын таңдаңыз)

**Жауап нұсқалары:**
- [x] **A)** `var a int = 5; var b float64 = 2.0; var c = a + b`
- [ ] **B)** `var a int = 5; var b float64 = 2.0; var c = float64(a) + b`
- [ ] **C)** `var x int = 10; _ = x`
- [x] **D)** `func f() int { if false { return 1 } }`

**Дұрыс жауап:** A, D
**Түсіндірме:** (A)-да Go әртүрлі типтер арасындағы арифметикалық операцияларды ашық түрлендірусіз (сәйкессіз типтер) рұқсат етпейді. (D)-де функцияның қайтарылатын типі `int`, бірақ басқару ағыны `return` операторы жоқ-ақ соңына дейін жетеді, бұл 'missing return' компиляциялық қатеге әкеледі. (B) анық конвертацияны қолданады, ал (C) `_` арқылы пайдаланылмаған айнымалыны тексеруді өшіреді.

---
## Question 10
- **Topic:** Strings & Runes
- **Difficulty:** Easy (30s)
- **Format:** Multiple Choice (Single Answer)

### English Version
**Question:** What does the built-in function `len(s)` return when passed a string `s` in Go?

**Options:**
- [x] **A)** The number of bytes in the string
- [ ] **B)** The number of Unicode characters / runes
- [ ] **C)** The number of words in the string
- [ ] **D)** The memory capacity of the underlying string buffer

**Correct Answer:** A
**Explanation:** In Go, strings are read-only byte slices. Therefore, `len(s)` returns the number of bytes, which may differ from the number of characters if multi-byte UTF-8 runes are present.

### Русская версия (Russian)
**Вопрос:** Что возвращает встроенная функция `len(s)` при передаче ей строки `s` в языке Go?

**Варианты ответа:**
- [x] **A)** Количество байтов в строке
- [ ] **B)** Количество символов Unicode / рун
- [ ] **C)** Количество слов в строке
- [ ] **D)** Объём памяти базового буфера строк

**Правильный ответ:** A
**Объяснение:** В языке Go строки представляют собой байтовые срезы, доступные только для чтения. Поэтому функция `len(s)` возвращает количество байтов, которое может отличаться от количества символов при наличии многобайтовых рун кодировки UTF-8.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Go тілінде `len(s)` кіріктірілген функциясы `s` жолды бергенде қандай мәні қайтарады?

**Жауап нұсқалары:**
- [x] **A)** Сызықтағы байттар саны
- [ ] **B)** Юникод таңбаларының / руналардың саны
- [ ] **C)** Сызықтағы сөздер саны
- [ ] **D)** Негізгі жол буферінің жады сыйымдылығы

**Дұрыс жауап:** A
**Түсіндірме:** Go-да жолдар – тек оқуға арналған байттық тілімдер. Сондықтан `len(s)` функциясы байттар санын қайтарады, ол көпбайтты UTF-8 руналары болған жағдайда таңбалар санынан өзгеше болуы мүмкін.

---
## Question 11
- **Topic:** Strings & Runes
- **Difficulty:** Medium (45s)
- **Format:** Multiple Choice (Single Answer)

### English Version
**Question:** What happens when the following Go code is compiled and executed?

```go
score := 100
text := "Score: " + score
fmt.Println(text)
```

**Options:**
- [ ] **A)** It prints `Score: 100`
- [x] **B)** Compile-time error: `cannot use score (variable of type int) as string value`
- [ ] **C)** It prints `Score: d`
- [ ] **D)** Runtime panic

**Correct Answer:** B
**Explanation:** In Go, you cannot directly concatenate a string and an integer with the `+` operator without explicit conversion (e.g., using `strconv.Itoa(score)` or `fmt.Sprintf`).

### Русская версия (Russian)
**Вопрос:** Что произойдет при компиляции и выполнении следующего кода на Go?

```go
score := 100
text := "Score: " + score
fmt.Println(text)
```

**Варианты ответа:**
- [ ] **A)** Выводится: `Score: 100`
- [x] **B)** Ошибка компиляции: `cannot use score (variable of type int) as string value`
- [ ] **C)** Выводит: `Score: d`
- [ ] **D)** Паника во время выполнения

**Правильный ответ:** B
**Объяснение:** В языке Go нельзя напрямую объединять строку и целое число с помощью оператора «`+`» без явного преобразования (например, с помощью «`strconv.Itoa(score)`» или «`fmt.Sprintf`»).

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Төмендегі Go коды компиляцияланып, орындалғанда не болады?

```go
score := 100
text := "Score: " + score
fmt.Println(text)
```

**Жауап нұсқалары:**
- [ ] **A)** Ол «`Score: 100`» деп басады.
- [x] **B)** Құрастыру кезіндегі қате: `cannot use score (variable of type int) as string value`
- [ ] **C)** Ол «`Score: d`» деп басады.
- [ ] **D)** Жүгіру кезіндегі паника

**Дұрыс жауап:** B
**Түсіндірме:** Go тілінде `+` операторы арқылы жол мен бүтін санды ашық түрлендірусіз тікелей біріктіруге болмайды (мысалы, `strconv.Itoa(score)` немесе `fmt.Sprintf` қолдану арқылы).

---

## Question 12
- **Topic:** Strings & Runes
- **Difficulty:** Hard (60s)
- **Format:** Text Question with Multiple Possible Answers

### English Version
**Question:** What integer value is returned by `len("Привет")` in Go, given that each of the 6 Cyrillic characters requires 2 bytes in UTF-8 encoding?

**Acceptable / Possible Answers:**
- `12`
- `12 bytes`

**Standard Answer:** 12
**Explanation:** In UTF-8, each Russian Cyrillic letter takes 2 bytes. Since 'Привет' contains 6 Cyrillic letters and `len()` counts bytes, the result is `6 * 2 = 12`.

### Русская версия (Russian)
**Вопрос:** Какое целое значение возвращает выражение `len("Привет")` в языке Go, учитывая, что каждый из 6 кириллических символов занимает 2 байта в кодировке UTF-8?

**Приемлемые / возможные варианты ответа:**
- `12`
- `12 байт`

**Основной ответ:** 12
**Объяснение:** В кодировке UTF-8 каждая буква русской кириллицы занимает 2 байта. Поскольку слово «Привет» содержит 6 кириллических букв, а функция `len()` подсчитывает количество байтов, результат равен `6 * 2 = 12`.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Go тілінде 6 кирилликалық таңбаның әрқайсысы UTF-8 кодтауында 2 байт алатынын ескерсек, `len("Привет")` функциясы қандай бүтін санды қайтарады?

**Қабылданатын / мүмкін жауап нұсқалары:**
- `Он екі`
- `12 байт`

**Негізгі жауап:** 12
**Түсіндірме:** UTF-8-де әрбір орыс кириллица әріпі 2 байт орын алады. 'Привет' сөзінде 6 орыс кириллица әріптері бар және `len()` функциясы байттарды есептейтіндіктен, нәтиже 6 * 2 = 12 болады.

---
## Question 14
- **Topic:** Conditions
- **Difficulty:** Medium (45s)
- **Format:** Multiple Choice (Single Answer)

### English Version
**Question:** What is the output of the following code snippet?

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
```

**Options:**
- [ ] **A)** `Alpha`
- [x] **B)** `Beta`
- [ ] **C)** `Gamma`
- [ ] **D)** Compile-time error

**Correct Answer:** B
**Explanation:** Evaluating `15 % 4` assigns `3` to `n`. The condition `n == 0` is false, but `else if n == 3` evaluates to true, printing `Beta`.

### Русская версия (Russian)
**Вопрос:** Каков результат выполнения следующего фрагмента кода?

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
```

**Варианты ответа:**
- [ ] **A)** `Alpha`
- [x] **B)** `Beta`
- [ ] **C)** `Gamma`
- [ ] **D)** Ошибка компиляции

**Правильный ответ:** B
**Объяснение:** При вычислении выражения ``15 % 4`` переменной ``3`` присваивается значение ``n``. Условие ``n == 0`` является ложным, но выражение ``else if n == 3`` вычисляется как истинное, в результате чего выводится ``Beta``.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Төмендегі код үзіндісінің шығысы қандай?

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
```

**Жауап нұсқалары:**
- [ ] **A)** `Alpha`
- [x] **B)** `Beta`
- [ ] **C)** `Gamma`
- [ ] **D)** Құрастыру кезіндегі қате

**Дұрыс жауап:** B
**Түсіндірме:** `15 % 4` бағалау кезінде `3` мәнін `n`-ге тағайындайды. `n == 0` шарты жалған, бірақ `else if n == 3` шындыққа теңеліп, `Beta`-ді шығарады.

---

## Question 15
- **Topic:** Conditions
- **Difficulty:** Hard (60s)
- **Format:** Multiple Choice (Multiple Correct Answers)

### English Version
**Question:** Which of the following statements regarding `if` statements in Go are TRUE? (Select all that apply)

**Options:**
- [x] **A)** The condition expression in an `if` statement must evaluate strictly to a `bool`.
- [ ] **B)** Non-zero integer numbers like `1` or non-empty strings are implicitly evaluated as true in `if` conditions.
- [x] **C)** The body block of an `if` statement must always be enclosed in curly braces `{}`.
- [ ] **D)** Parentheses `()` around the condition expression in an `if` statement are mandatory.

**Correct Answer:** A, C
**Explanation:** Go has no implicit truthy/falsy coercion (conditions must be of type `bool`), curly braces `{}` are always required even for single-line bodies, and parentheses around the condition are not required in Go.

### Русская версия (Russian)
**Вопрос:** Какие из приведенных ниже утверждений относительно операторов «`if`» в языке Go являются верными? (Выберите все правильные варианты)

**Варианты ответа:**
- [x] **A)** Выражение условия в операторе «`if`» должно давать в результате строго значение типа «`bool`».
- [ ] **B)** Целые числа, отличные от нуля, такие как `1`, или непустые строки неявно оцениваются как «true» в условиях типа `if`.
- [x] **C)** Тело оператора `if` всегда должно заключаться в фигурные скобки `{}`.
- [ ] **D)** Скобки `()`, заключающие выражение условия в операторе `if`, являются обязательными.

**Правильный ответ:** A, C
**Объяснение:** В языке Go отсутствует неявная приведение значений к «истинным» или «ложным» (условия должны быть типа `bool`), фигурные скобки `{}` всегда обязательны даже для однострочных тел, а скобки вокруг условия в Go не требуются.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Go тіліндегі `if`-мәлімдемелері туралы келесі мәлімдемелердің қайсысы дұрыс? (Барлығын таңдаңыз)

**Жауап нұсқалары:**
- [x] **A)** `if` операторындағы шарттық өрнек міндетті түрде дәл `bool` мәніне бағалануы тиіс.
- [ ] **B)** `1` сияқты нөлге тең емес бүтін сандар немесе бос емес жолдар `if` шарттарында жасырын түрде true деп есептеледі.
- [x] **C)** `if` операторының денесі әрдайым дөңгелек жақшалармен қоршалуы тиіс `{}`.
- [ ] **D)** `if` операторындағы шарттық өрнектің айналасындағы `()` жақшалары міндетті.

**Дұрыс жауап:** A, C
**Түсіндірме:** Go-да жасырын truthy/falsy түрлендіру жоқ (шарттар `bool` типінен болуы тиіс), бір жолдық блоктарда да `{}` дөңгелек жақшалары міндетті және шарттың айналасындағы дөңгелек жақшалар қажет емес.

---

## Question 16
- **Topic:** Loops
- **Difficulty:** Easy (30s)
- **Format:** Multiple Choice (Single Answer)

### English Version
**Question:** Which loop keyword exists in Go to construct counting loops, while-style loops, and infinite loops?

**Options:**
- [ ] **A)** `while`
- [x] **B)** `for`
- [ ] **C)** `loop`
- [ ] **D)** `repeat`

**Correct Answer:** B
**Explanation:** Go has only a single loop keyword: `for`. It unifies all loop constructs in the language.

### Русская версия (Russian)
**Вопрос:** Какое ключевое слово цикла существует в языке Go для построения циклов с подсчётом, циклов типа «while» и бесконечных циклов?

**Варианты ответа:**
- [ ] **A)** `while`
- [x] **B)** `for`
- [ ] **C)** `loop`
- [ ] **D)** `repeat`

**Правильный ответ:** B
**Объяснение:** В языке Go есть только одно ключевое слово для циклов: `for`. Оно объединяет все конструкции циклов в языке.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Go тілінде санау циклдарын, while-стильді циклдарын және шексіз циклдарын құру үшін қандай цикл кілттік сөзі бар?

**Жауап нұсқалары:**
- [ ] **A)** `while`
- [x] **B)** `for`
- [ ] **C)** `loop`
- [ ] **D)** `repeat`

**Дұрыс жауап:** B
**Түсіндірме:** Go тілінде циклдерді ұйымдастыру үшін тек бір ғана кілтсөз бар: `for`. Ол тілдегі барлық циклдік конструкцияларды біріктіреді.

---
## Question 17
- **Topic:** Loops
- **Difficulty:** Hard (60s)
- **Format:** Multiple Choice (Single Answer)

### English Version
**Question:** What will be printed when the following code executes?

```go
sum := 0
for i := 1; i <= 5; i++ {
    if i%2 == 0 {
        continue
    }
    sum += i
}
fmt.Println(sum)
```

**Options:**
- [ ] **A)** `6`
- [x] **B)** `9`
- [ ] **C)** `15`
- [ ] **D)** `4`

**Correct Answer:** B
**Explanation:** The `continue` statement skips even values (2 and 4). Only odd values (1, 3, 5) are added: 1 + 3 + 5 = 9.

### Русская версия (Russian)
**Вопрос:** Что будет выведено на экран при выполнении следующего кода?

```go
sum := 0
for i := 1; i <= 5; i++ {
    if i%2 == 0 {
        continue
    }
    sum += i
}
fmt.Println(sum)
```

**Варианты ответа:**
- [ ] **A)** `6`
- [x] **B)** `9`
- [ ] **C)** `15`
- [ ] **D)** `4`

**Правильный ответ:** B
**Объяснение:** В выражении `continue` пропускаются четные числа (2 и 4). Складываются только нечетные числа (1, 3, 5): 1 + 3 + 5 = 9.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Төмендегі код орындалғанда не басылады?

```go
sum := 0
for i := 1; i <= 5; i++ {
    if i%2 == 0 {
        continue
    }
    sum += i
}
fmt.Println(sum)
```

**Жауап нұсқалары:**
- [ ] **A)** `6`
- [x] **B)** `9`
- [ ] **C)** `15`
- [ ] **D)** `4`

**Дұрыс жауап:** B
**Түсіндірме:** continue операторы тіпті сандарды (2 және 4) өткізіп жібереді. Тек тақ сандар (1, 3, 5) қосылады: 1 + 3 + 5 = 9.

---
## Question 18
- **Topic:** Loops
- **Difficulty:** Medium (45s)
- **Format:** Text Question with Multiple Possible Answers

### English Version
**Question:** When iterating over a string using `for index, val := range "Go"`, what is the Go data type of the variable `val`?

**Acceptable / Possible Answers:**
- `rune`
- `int32`
- `rune (int32)`

**Standard Answer:** rune (or int32)
**Explanation:** When ranging over a string, Go iterates over Unicode code points, decoding UTF-8 characters as `rune` (which is an alias for `int32`).

### Русская версия (Russian)
**Вопрос:** При итерации по строке с использованием ` в качестве индекса: val := range "Go"`, каков тип данных переменной `val` в языке Go?

**Приемлемые / возможные варианты ответа:**
- `Руне`
- `32-разрядное целое число`
- `rune (int32)`

**Основной ответ:** rune (or int32)
**Объяснение:** При обработке строки Go проходит по кодовым точкам Unicode, декодируя символы UTF-8 в виде `rune` (что является псевдонимом для `int32`).

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Стринг бойынша `for index, val := range "Go"` циклымен өту кезінде `val` айнымалының Go деректер типі қандай?

**Қабылданатын / мүмкін жауап нұсқалары:**
- `Руна`
- `32-биттік бүтін сан`
- `руна (int32)`

**Негізгі жауап:** rune (or int32)
**Түсіндірме:** Желідегі аралықты анықтау кезінде Go Unicode кодының таңбалары бойынша өтіп, UTF-8 таңбаларын `rune` (ол `int32`-дің синонимі) ретінде декодтайды.

---
## Question 19
- **Topic:** Arrays
- **Difficulty:** Easy (30s)
- **Format:** Multiple Choice (Single Answer)

### English Version
**Question:** Which of the following declarations correctly creates an array of 5 integers initialized to their zero values in Go?

**Options:**
- [x] **A)** `var a [5]int`
- [ ] **B)** `var a []int = [5]int{}`
- [ ] **C)** `a := array(5, int)`
- [ ] **D)** `var a int[5]`

**Correct Answer:** A
**Explanation:** In Go, array types specify the size before the element type, so `var a [5]int` creates an array of 5 integers with all elements initialized to 0.

### Русская версия (Russian)
**Вопрос:** Какой из приведённых ниже операторов в языке Go правильно создаёт массив из 5 целых чисел, инициализированных нулями?

**Варианты ответа:**
- [x] **A)** `var a [5]int`
- [ ] **B)** `var a []int = [5]int{}`
- [ ] **C)** `a := array(5, int)`
- [ ] **D)** `var a int[5]`

**Правильный ответ:** A
**Объяснение:** В Go типы массивов указывают размер перед типом элемента, поэтому `var a [5]int` создает массив из 5 целых чисел, все элементы которого инициализированы значением 0.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Төмендегі жарияланымдардың қайсысы Go тілінде 5 бүтін саннан тұратын массивті олардың мәндерін нөлге теңей отырып дұрыс жасайды?

**Жауап нұсқалары:**
- [x] **A)** `var a [5]int`
- [ ] **B)** `var a []int = [5]int{}`
- [ ] **C)** `a := array(5, int)`
- [ ] **D)** `var a int[5]`

**Дұрыс жауап:** A
**Түсіндірме:** Go тілінде массив типтері элемент типінен бұрын өлшемін көрсетеді, сондықтан `var a [5]int` барлық элементтері 0-ге теңестірілген 5 бүтін сандар массивін жасайды.

---
## Question 20
- **Topic:** Arrays
- **Difficulty:** Hard (60s)
- **Format:** Multiple Choice (Single Answer)

### English Version
**Question:** What is the output of the following code?

```go
a := [3]int{1, 2, 3}
b := a
b[0] = 99
fmt.Println(a[0], b[0])
```

**Options:**
- [ ] **A)** `99 99`
- [x] **B)** `1 99`
- [ ] **C)** `1 1`
- [ ] **D)** Compile-time error

**Correct Answer:** B
**Explanation:** Arrays in Go are value types. The assignment `b := a` copies all elements of `a` into `b`. Changing `b[0]` does not affect `a[0]`.

### Русская версия (Russian)
**Вопрос:** Каков результат выполнения следующего кода?

```go
a := [3]int{1, 2, 3}
b := a
b[0] = 99
fmt.Println(a[0], b[0])
```

**Варианты ответа:**
- [ ] **A)** `99 99`
- [x] **B)** `1 99`
- [ ] **C)** `1 1`
- [ ] **D)** Ошибка компиляции

**Правильный ответ:** B
**Объяснение:** Массивы в Go являются типами значений. Присваивание `b := a` копирует все элементы `a` в `b`. Изменение `b[0]` не влияет на `a[0]`.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Төмендегі кодтың шығысы қандай?

```go
a := [3]int{1, 2, 3}
b := a
b[0] = 99
fmt.Println(a[0], b[0])
```

**Жауап нұсқалары:**
- [ ] **A)** `99 99`
- [x] **B)** `1 99`
- [ ] **C)** `1 1`
- [ ] **D)** Құрастыру кезіндегі қате

**Дұрыс жауап:** B
**Түсіндірме:** Go-дағы массивтер – мәндер типтері. `b := a` тағайындау операторы `a`-ның барлық элементтерін `b`-ге көшіреді. `b[0]`-ды өзгерту `a[0]`-ға әсер етпейді.

---
## Question 21
- **Topic:** Arrays
- **Difficulty:** Hard (60s)
- **Format:** Multiple Choice (Multiple Correct Answers)

### English Version
**Question:** Which of the following statements about arrays in Go are TRUE? (Select all that apply)

**Options:**
- [x] **A)** The length of an array is fixed and is part of its type (e.g. `[3]int` and `[4]int` are distinct types).
- [x] **B)** Two arrays of identical type and length containing comparable elements can be compared with `==`.
- [ ] **C)** Arrays can dynamically expand in length at runtime using `append()`.
- [x] **D)** Passing an array as an argument to a function creates a complete copy of the array unless passed by pointer.

**Correct Answer:** A, B, D
**Explanation:** Arrays in Go have a fixed size that is part of their type signature, are comparable with `==` if elements are comparable, and are passed by value (copying all elements). `append()` only works on slices, not fixed arrays.

### Русская версия (Russian)
**Вопрос:** Какие из приведённых ниже утверждений о массивах в языке Go являются верными? (Выберите все подходящие варианты)

**Варианты ответа:**
- [x] **A)** Длина массива фиксирована и является частью его типа (например, `[3]int` и `[4]int` — это разные типы).
- [x] **B)** Два массива одного и того же типа и длины, содержащие сопоставимые элементы, можно сравнить с помощью операторов `==`.
- [ ] **C)** Длина массивов может динамически увеличиваться во время выполнения с помощью `append()`.
- [x] **D)** При передаче массива в качестве аргумента функции создается полная копия массива, если только он не передается по указателю.

**Правильный ответ:** A, B, D
**Объяснение:** Массивы в Go имеют фиксированный размер, который является частью их сигнатуры типа; они сопоставимы с помощью оператора `==`, если их элементы сопоставимы, и передаются по значению (с копированием всех элементов). Функция `append()` работает только со слайсами, а не с фиксированными массивами.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Go тіліндегі массивтер туралы келесі мәлімдемелердің қайсысы дұрыс? (Барлығын таңдаңыз)

**Жауап нұсқалары:**
- [x] **A)** Массивтің ұзындығы тұрақты және оның типінің құрамдас бөлігі болып табылады (мысалы, `[3]int` және `[4]int` – әртүрлі типтер).
- [x] **B)** Түрі мен ұзындығы бірдей және салыстырылатын элементтерден тұратын екі массивті `==` операторымен салыстыруға болады.
- [ ] **C)** Массивтер орындау уақытында `append()` әдісін пайдаланып ұзындығын динамикалық түрде ұлғайта алады.
- [x] **D)** Функцияға массивті аргумент ретінде бергенде, оны көрсеткіш арқылы бермесеңіз, массивтің толық көшірмесі жасалады.

**Дұрыс жауап:** A, B, D
**Түсіндірме:** Go-дағы массивтердің өлшемі олардың типтік белгісінің бір бөлігі болып табылады, элементтері салыстырылатын болса, оларды `==` операторымен салыстыруға болады және олар мәндері бойынша беріледі (барлық элементтері көшіріледі). `append()` тек кесінділерде ғана жұмыс істейді, тұрақты массивтерде емес.

---
## Question 22
- **Topic:** Slices
- **Difficulty:** Easy (30s)
- **Format:** Multiple Choice (Single Answer)

### English Version
**Question:** Which built-in function is used to append elements to the end of a slice in Go?

**Options:**
- [ ] **A)** `push()`
- [x] **B)** `append()`
- [ ] **C)** `insert()`
- [ ] **D)** `extend()`

**Correct Answer:** B
**Explanation:** Go provides the built-in `append()` function to add elements to the end of a slice, reallocating underlying capacity if needed.

### Русская версия (Russian)
**Вопрос:** Какая встроенная функция используется в Go для добавления элементов в конец фрагмента?

**Варианты ответа:**
- [ ] **A)** `push()`
- [x] **B)** `append()`
- [ ] **C)** `insert()`
- [ ] **D)** `extend()`

**Правильный ответ:** B
**Объяснение:** В Go предусмотрена встроенная функция `append()`, позволяющая добавлять элементы в конец среза с перераспределением базовой памяти при необходимости.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Go-да массивтің бір бөлігінің соңына элементтерді қосу үшін қай кіріктірілген функция қолданылады?

**Жауап нұсқалары:**
- [ ] **A)** `push()`
- [x] **B)** `append()`
- [ ] **C)** `insert()`
- [ ] **D)** `extend()`

**Дұрыс жауап:** B
**Түсіндірме:** Go тілі кесіндінің соңына элементтерді қосуға арналған кіріктірілген `append()` функциясын ұсынады, қажет болған жағдайда негізгі сыйымдылықты қайта бөледі.

---
## Question 23
- **Topic:** Slices
- **Difficulty:** Hard (60s)
- **Format:** Multiple Choice (Single Answer)

### English Version
**Question:** What is printed by the following code snippet?

```go
s := []int{10, 20, 30, 40}
sub := s[1:3]
sub[0] = 99
fmt.Println(s[1])
```

**Options:**
- [ ] **A)** `20`
- [x] **B)** `99`
- [ ] **C)** `10`
- [ ] **D)** `30`

**Correct Answer:** B
**Explanation:** Slices are views over an underlying array. `sub[0]` points to the exact same element as `s[1]`. Modifying `sub[0]` directly changes `s[1]`.

### Русская версия (Russian)
**Вопрос:** Что выведет на экран следующий фрагмент кода?

```go
s := []int{10, 20, 30, 40}
sub := s[1:3]
sub[0] = 99
fmt.Println(s[1])
```

**Варианты ответа:**
- [ ] **A)** `20`
- [x] **B)** `99`
- [ ] **C)** `10`
- [ ] **D)** `30`

**Правильный ответ:** B
**Объяснение:** Срезы представляют собой представления базового массива. `sub[0]` указывает на тот же самый элемент, что и `s[1]`. Изменение `sub[0]` напрямую приводит к изменению `s[1]`.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Төмендегі код үзіндісі нені басып шығарады?

```go
s := []int{10, 20, 30, 40}
sub := s[1:3]
sub[0] = 99
fmt.Println(s[1])
```

**Жауап нұсқалары:**
- [ ] **A)** `20`
- [x] **B)** `99`
- [ ] **C)** `10`
- [ ] **D)** `30`

**Дұрыс жауап:** B
**Түсіндірме:** Слайстер – бұл негізгі массив бойынша көріністер. `sub[0]` `s[1]`-пен дәл бірдей элементке нұсқайды. `sub[0]`-ді тікелей өзгерту `s[1]`-ді өзгертеді.

---
## Question 24
- **Topic:** Slices
- **Difficulty:** Hard (60s)
- **Format:** Multiple Choice (Multiple Correct Answers)

### English Version
**Question:** Which of the following operations are valid when working with slices in Go? (Select all that apply)

**Options:**
- [x] **A)** Adding new elements to a slice using `append(slice, element)`.
- [x] **B)** Getting the number of elements in a slice using `len(slice)`.
- [ ] **C)** Directly comparing two slices with multiple elements using the `==` operator (e.g. `slice1 == slice2`).
- [x] **D)** Creating a sub-slice from an existing slice using sub-slicing syntax `slice[start:end]`.

**Correct Answer:** A, B, D
**Explanation:** In Go, you can append elements to slices with `append()`, obtain length with `len()`, and create sub-slices with `slice[start:end]`. Slices cannot be compared directly with `==` (except when checking against `nil`).

### Русская версия (Russian)
**Вопрос:** Какие из перечисленных ниже операций допустимы при работе со срезами в языке Go? (Выберите все подходящие варианты)

**Варианты ответа:**
- [x] **A)** Добавление новых элементов в фрагмент с помощью команды «`append(slice, element)`».
- [x] **B)** Определение количества элементов в фрагменте с помощью функции ``len(slice)``.
- [ ] **C)** Прямое сравнение двух срезов, содержащих несколько элементов, с помощью оператора «`==`» (например, `slice1 == slice2`).
- [x] **D)** Создание подсреза на основе существующего среза с использованием синтаксиса подсреза `slice[start:end]`.

**Правильный ответ:** A, B, D
**Объяснение:** В Go можно добавлять элементы в слайсы с помощью `append()`, получать их длину с помощью `len()` и создавать подслайсы с помощью `slice[start:end]`. Слайсы нельзя сравнивать напрямую с помощью `==` (за исключением проверки на равность с помощью `nil`).

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Go-да кесінділермен жұмыс істегенде төмендегі операциялардың қайсысы жарамды? (Барлығын таңдаңыз)

**Жауап нұсқалары:**
- [x] **A)** `append(slice, element)` арқылы слайсқа жаңа элементтер қосу.
- [x] **B)** `len(slice)` қолдана отырып, кесіндідегі элементтер санын алу.
- [ ] **C)** `==` операторы арқылы бірнеше элементі бар екі тілімді тікелей салыстыру (мысалы: `slice1 == slice2`).
- [x] **D)** `slice[start:end]` синтаксисін пайдаланып, бар слайстан кіші слайс жасау.

**Дұрыс жауап:** A, B, D
**Түсіндірме:** Go тілінде сіз `append()` арқылы тілімдерге элементтерді қоса аласыз, `len()` арқылы ұзындығын аласыз және `slice[start:end]` арқылы кіші тілімдер жасай аласыз. Тілімдерді `==` арқылы тікелей салыстыруға болмайды (`nil`-ге теңдікті тексеруді қоспағанда).

---

## Question 28
- **Topic:** Pointers
- **Difficulty:** Easy (30s)
- **Format:** Multiple Choice (Single Answer)

### English Version
**Question:** Which operator is used to obtain the memory address of an existing variable `x` in Go?

**Options:**
- [ ] **A)** `*x`
- [x] **B)** `&x`
- [ ] **C)** `@x`
- [ ] **D)** `addr(x)`

**Correct Answer:** B
**Explanation:** The address-of operator `&` generates a pointer pointing to the memory location of variable `x`.

### Русская версия (Russian)
**Вопрос:** Какой оператор используется в языке Go для получения адреса в памяти существующей переменной `x`?

**Варианты ответа:**
- [ ] **A)** `*x`
- [x] **B)** `&x`
- [ ] **C)** `@x`
- [ ] **D)** `addr(x)`

**Правильный ответ:** B
**Объяснение:** Оператор «адрес-of» `&` создаёт указатель, указывающий на ячейку памяти переменной `x`.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Go тілінде бар `x` айнымалының жады мекенжайын алу үшін қандай оператор қолданылады?

**Жауап нұсқалары:**
- [ ] **A)** `*x`
- [x] **B)** `&x`
- [ ] **C)** `@x`
- [ ] **D)** `addr(x)`

**Дұрыс жауап:** B
**Түсіндірме:** Адреске алу операторы `&` `x` айнымалысының жады мекенжайына нүктелейтін көрсеткішті жасайды.

---
## Question 29
- **Topic:** Pointers
- **Difficulty:** Medium (45s)
- **Format:** Multiple Choice (Single Answer)

### English Version
**Question:** What is printed by the following code?

```go
x := 42
p := &x
*p = 100
fmt.Println(x)
```

**Options:**
- [ ] **A)** `42`
- [x] **B)** `100`
- [ ] **C)** The memory address of `x`
- [ ] **D)** `nil`

**Correct Answer:** B
**Explanation:** The pointer `p` holds the memory address of `x`. Dereferencing and assigning `*p = 100` updates the memory value of `x` directly to `100`.

### Русская версия (Russian)
**Вопрос:** Что выведется при выполнении следующего кода?

```go
x := 42
p := &x
*p = 100
fmt.Println(x)
```

**Варианты ответа:**
- [ ] **A)** `42`
- [x] **B)** `100`
- [ ] **C)** Адрес в памяти `x`
- [ ] **D)** `nil`

**Правильный ответ:** B
**Объяснение:** Указатель `p` содержит адрес в памяти переменной `x`. Развертывание указателя и присваивание `*p = 100` приводит к непосредственному обновлению значения переменной `x` в памяти до `100`.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Төмендегі код нені басып шығарады?

```go
x := 42
p := &x
*p = 100
fmt.Println(x)
```

**Жауап нұсқалары:**
- [ ] **A)** `42`
- [x] **B)** `100`
- [ ] **C)** x-тің жад мекенжайы
- [ ] **D)** `nil`

**Дұрыс жауап:** B
**Түсіндірме:** Пойнтер `p` `x`-тің жады мекенжайын сақтайды. `*p = 100` операторы арқылы мекенжайға сілтеме жасап, оған 100 тағайындау `x`-тің жадыдағы мәнін тікелей 100-ге өзгертеді.

---
## Question 30
- **Topic:** Pointers
- **Difficulty:** Hard (60s)
- **Format:** Multiple Choice (Multiple Correct Answers)

### English Version
**Question:** What happens when executing the following code in Go? (Select all that apply)

```go
var p *int
fmt.Println(*p)
```

**Options:**
- [x] **A)** A runtime panic occurs during execution.
- [x] **B)** The panic message contains `invalid memory address or nil pointer dereference`.
- [ ] **C)** It safely outputs `0`.
- [ ] **D)** The Go compiler rejects the code at compile-time with an uninitialized pointer error.

**Correct Answer:** A, B
**Explanation:** An uninitialized pointer has the zero value `nil`. Dereferencing a `nil` pointer causes a runtime panic (`runtime error: invalid memory address or nil pointer dereference`), not a compile error.

### Русская версия (Russian)
**Вопрос:** Что произойдет при выполнении следующего кода на языке Go? (Выберите все подходящие варианты)

```go
var p *int
fmt.Println(*p)
```

**Варианты ответа:**
- [x] **A)** Во время выполнения программы возникает паника.
- [x] **B)** Сообщение о панике содержит `неверный адрес памяти или разыменование нулевого указателя`.
- [ ] **C)** Он корректно выводит `0`.
- [ ] **D)** Компилятор Go отклоняет этот код на этапе компиляции, выдавая ошибку «неинициализированный указатель».

**Правильный ответ:** A, B
**Объяснение:** Неинициализированный указатель имеет значение нуль `nil`. Разреференция указателя `nil` приводит к панике во время выполнения (`ошибка выполнения: недопустимый адрес памяти или разреференция указателя nil`), а не к ошибке компиляции.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Go тілінде келесі кодты орындағанда не болады? (Барлығына сәйкес келетіндерін таңдаңыз)

```go
var p *int
fmt.Println(*p)
```

**Жауап нұсқалары:**
- [x] **A)** Орындау кезінде runtime паникасы пайда болады.
- [x] **B)** Паникалық хабарлама `жарамсыз жад мекенжайы немесе nil көрсеткішін дереференциациялау` дегенді қамтиды.
- [ ] **C)** Ол қауіпсіз түрде 0 шығарады.
- [ ] **D)** Go компиляторы кодты компиляция кезінде бастапқыланбаған көрсеткіш қатесімен қабылдамайды.

**Дұрыс жауап:** A, B
**Түсіндірме:** Бастапқыланбаған көрсеткіш 0 мәніне, яғни `nil`-ге тең. `nil` көрсеткішінің сілтемесін шақыру компиляция кезінде емес, орындау кезінде паника тудырады (`runtime error: invalid memory address or nil pointer dereference`).

---
## Question 31
- **Topic:** Structs & Methods
- **Difficulty:** Easy (30s)
- **Format:** Multiple Choice (Single Answer)

### English Version
**Question:** Which syntax correctly declares a struct type named `Book` containing a string field `Title` in Go?

**Options:**
- [x] **A)** `type Book struct { Title string }`
- [ ] **B)** `struct Book { Title string }`
- [ ] **C)** `class Book { Title string }`
- [ ] **D)** `var Book = struct { Title string }`

**Correct Answer:** A
**Explanation:** In Go, custom struct types are declared using the `type <Name> struct { ... }` syntax.

### Русская версия (Russian)
**Вопрос:** Какой синтаксис правильно объявляет в языке Go тип структуры с именем `Book`, содержащий строковое поле `Title`?

**Варианты ответа:**
- [x] **A)** `type Book struct { Title string }`
- [ ] **B)** `struct Book { Title string }`
- [ ] **C)** `class Book { Title string }`
- [ ] **D)** `var Book = struct { Title string }`

**Правильный ответ:** A
**Объяснение:** В языке Go пользовательские типы структуры объявляются с помощью синтаксиса `type <Name> struct { ... }`.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Go тілінде `Book` атты struct типінің `Title` атты жол өрісін дұрыс жариялайтын қай синтаксис?

**Жауап нұсқалары:**
- [x] **A)** `type Book struct { Title string }`
- [ ] **B)** `struct Book { Title string }`
- [ ] **C)** `class Book { Title string }`
- [ ] **D)** `var Book = struct { Title string }`

**Дұрыс жауап:** A
**Түсіндірме:** Go тілінде арнайы struct типтері `type <Name> struct { ... }` синтаксисі арқылы жарияланады.

---
## Question 34
- **Topic:** Errors & Error Handling
- **Difficulty:** Easy (30s)
- **Format:** Multiple Choice (Single Answer)

### English Version
**Question:** What value does an idiomatic Go function return for its `error` return value to indicate success (no error)?

**Options:**
- [ ] **A)** `""` (empty string)
- [ ] **B)** `0`
- [x] **C)** `nil`
- [ ] **D)** `false`

**Correct Answer:** C
**Explanation:** The `error` type in Go is an interface. A `nil` error indicates that the operation completed successfully without errors.

### Русская версия (Russian)
**Вопрос:** Какое значение возвращает идиоматическая функция на Go в качестве значения возврата `error`, чтобы указать успешное выполнение (отсутствие ошибки)?

**Варианты ответа:**
- [ ] **A)** `""` (пустая строка)
- [ ] **B)** `0`
- [x] **C)** `nil`
- [ ] **D)** `false`

**Правильный ответ:** C
**Объяснение:** Тип ошибки `error` в языке Go представляет собой интерфейс. Ошибка `nil` означает, что операция завершилась успешно, без ошибок.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Идиоматикалық Go функциясы `error` нәтижесінде сәтті (қатесіз) болғанын көрсету үшін қандай мән қайтарады?

**Жауап нұсқалары:**
- [ ] **A)** "" (бос жол)
- [ ] **B)** `0`
- [x] **C)** `nil`
- [ ] **D)** `false`

**Дұрыс жауап:** C
**Түсіндірме:** Go тіліндегі `error` типі интерфейс болып табылады. `nil` қате операцияның қатесіз сәтті аяқталғанын білдіреді.

---
## Question 36
- **Topic:** Errors & Error Handling
- **Difficulty:** Hard (60s)
- **Format:** Multiple Choice (Multiple Correct Answers)

### English Version
**Question:** Which of the following statements about errors and error handling in Go are TRUE? (Select all that apply)

**Options:**
- [x] **A)** Functions that can produce an error conventionally return the `error` as their last return value.
- [x] **B)** The condition `if err != nil` is the standard idiomatic way to check if an error occurred.
- [x] **C)** A new custom error value can be created using the `errors.New()` function.
- [ ] **D)** Go uses structured `try`, `catch`, and `finally` blocks for exception handling.

**Correct Answer:** A, B, C
**Explanation:** In Go, errors are explicit values returned as the final return value. Success is denoted by `nil`, and errors are checked via `if err != nil`. Custom errors are created with `errors.New()`. Go does not use try/catch exception handling.

### Русская версия (Russian)
**Вопрос:** Какие из приведённых ниже утверждений об ошибках и обработке ошибок в языке Go являются ПРАВИЛЬНЫМИ? (Выберите все подходящие варианты)

**Варианты ответа:**
- [x] **A)** Функции, которые могут вызвать ошибку, по соглашению возвращают значение типа «`error`» в качестве последнего возвращаемого значения.
- [x] **B)** Условие «`if err != nil`» — это стандартный идиоматический способ проверки наличия ошибки.
- [x] **C)** Новое пользовательское значение ошибки можно создать с помощью функции ``errors.New()``.
- [ ] **D)** В языке Go для обработки исключений используются структурированные блоки `try`, `catch` и `finally`.

**Правильный ответ:** A, B, C
**Объяснение:** В языке Go ошибки представляют собой явные значения, возвращаемые в качестве конечного значения возврата. Успех обозначается с помощью `nil`, а ошибки проверяются с помощью `if err != nil`. Пользовательские ошибки создаются с помощью `errors.New()`. В Go не используется обработка исключений по схеме try/catch.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Go тіліндегі қателер мен қателерді өңдеу туралы келесі мәлімдемелердің қайсысы дұрыс? (Барлығын таңдаңыз)

**Жауап нұсқалары:**
- [x] **A)** Қате тудыруы мүмкін функциялар дәстүрлі түрде соңғы нәтиже ретінде `error`-ті қайтарады.
- [x] **B)** `if err != nil` – қатенің болғанын тексерудің стандартты идиоматикалық тәсілі.
- [x] **C)** `errors.New()` функциясын пайдаланып жаңа теңшелетін қате мәнін жасауға болады.
- [ ] **D)** Go ерекшеліктерді өңдеу үшін құрылымдық `try`, `catch` және `finally` блоктерін пайдаланады.

**Дұрыс жауап:** A, B, C
**Түсіндірме:** Go-да қателер соңғы қайтару мәні ретінде нақты мәндер ретінде қайтарылады. Сәттілік `nil` арқылы көрсетіледі, ал қателер `if err != nil` арқылы тексеріледі. Жеке қателер `errors.New()` арқылы жасалады. Go try/catch ерекшеліктерін өңдеуді пайдаланбайды.

---

## Question 37
- **Topic:** Command-Line Arguments & Stdin
- **Difficulty:** Easy (30s)
- **Format:** Multiple Choice (Single Answer)

### English Version
**Question:** What does `os.Args[0]` represent when executing a compiled Go command-line program?

**Options:**
- [ ] **A)** The first user argument passed after the executable name
- [x] **B)** The path or command used to invoke the program itself
- [ ] **C)** The total number of arguments passed to the program
- [ ] **D)** An empty string

**Correct Answer:** B
**Explanation:** In Go and POSIX environments, `os.Args[0]` holds the path or command used to launch the binary. Actual user arguments start at index `1`.

### Русская версия (Russian)
**Вопрос:** Что означает `os.Args[0]` при запуске скомпилированной программы Go, запускаемой из командной строки?

**Варианты ответа:**
- [ ] **A)** Первый аргумент пользователя, передаваемый после имени исполняемого файла
- [x] **B)** Путь или команда, используемая для запуска самой программы
- [ ] **C)** Общее количество аргументов, переданных программе
- [ ] **D)** Пустая строка

**Правильный ответ:** B
**Объяснение:** В средах Go и POSIX переменная `os.Args[0]` содержит путь или команду, использованную для запуска бинарного файла. Фактические аргументы пользователя начинаются с индекса `1`.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Құрастырылған Go командалық жолдағы бағдарламаны орындау кезінде `os.Args[0]` нені білдіреді?

**Жауап нұсқалары:**
- [ ] **A)** Орындалатын файл атынан кейін берілетін бірінші пайдаланушы аргументі
- [x] **B)** Бағдарламаны шақыру үшін қолданылатын жол немесе команда
- [ ] **C)** Бағдарламаға берілген аргументтердің жалпы саны
- [ ] **D)** Бос жол

**Дұрыс жауап:** B
**Түсіндірме:** Go және POSIX орталарында `os.Args[0]` бинарлық файлды іске қосу үшін пайдаланылған жол немесе команданы сақтайды. Нақты пайдаланушы аргументтері 1-индекстен басталады.

---
## Question 39
- **Topic:** Command-Line Arguments & Stdin
- **Difficulty:** Hard (60s)
- **Format:** Text Question with Multiple Possible Answers

### English Version
**Question:** If a compiled Go program is invoked in the terminal as `./cli process -v file.txt`, what is the exact integer value returned by `len(os.Args)`?

**Acceptable / Possible Answers:**
- `4`
- `four`

**Standard Answer:** 4
**Explanation:** `os.Args` contains 4 elements: `os.Args[0] = "./cli"`, `os.Args[1] = "process"`, `os.Args[2] = "-v"`, and `os.Args[3] = "file.txt"`. Thus, `len(os.Args)` is 4.

### Русская версия (Russian)
**Вопрос:** Если скомпилированная программа на языке Go запускается в терминале следующим образом: `./cli process -v file.txt`, то какое точное целое значение возвращает функция `len(os.Args)`?

**Приемлемые / возможные варианты ответа:**
- `4`
- `четыре`

**Основной ответ:** 4
**Объяснение:** `os.Args` содержит 4 элемента: `os.Args[0] = "./cli"`, `os.Args[1] = "process"`, `os.Args[2] = "-v"` и `os.Args[3] = "file.txt"`. Таким образом, `len(os.Args)` равно 4.

### Қазақша нұсқасы (Kazakh)
**Сұрақ:** Егер компиляцияланған Go бағдарламасы терминалда `./cli process -v file.txt` ретінде іске қосылса, `len(os.Args)` қандай дәл бүтін санды қайтарады?

**Қабылданатын / мүмкін жауап нұсқалары:**
- `4`
- `төрт`

**Негізгі жауап:** 4
**Түсіндірме:** os.Args құрамында 4 элемент бар: os.Args[0] = "./cli", os.Args[1] = "process", os.Args[2] = "-v", және os.Args[3] = "file.txt". Сондықтан `len(os.Args)` 4-ке тең.

---
