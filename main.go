package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Цвета
const (
	Red    = "\033[31m" // красный
	Green  = "\033[32m" // зеленый
	Yellow = "\033[33m" // желтый
	Blue   = "\033[34m" // синий
	Reset  = "\033[0m"  // возврат к начальному цвету
)

// Структура для сохранения результатов
type Result struct {
	Date     string `json:"date"`     // Дата
	Outcome  string `json:"outcome"`  // "Победа" или "Проигрыш"
	Attempts int    `json:"attempts"` // Сколько попыток потрачено
}

func main() {
	rand.Seed(time.Now().UnixNano())
	bestResults()

	for {
		playGame()
		if !Replay() {
			fmt.Println(Green + "Спасибо за игру! 🔥" + Reset)
			break
		}
	}
}

// Вывод предыдущих результатов
func bestResults() {
	results := loadResults()
	if len(results) == 0 {
		fmt.Println(Yellow + "Нет предыдущих результатов." + Reset)
		return
	}

	fmt.Println(Yellow + "Последние результаты:" + Reset)
	for i, r := range results {
		fmt.Printf("%d) %s — %s (%d попыток)\n", i+1, r.Date, r.Outcome, r.Attempts)
	}
	fmt.Println()
}

// Загрузка результатов из файла
func loadResults() []Result {
	results := []Result{}
	file, err := os.Open("results.json")
	if err != nil {
		return results
	}
	defer file.Close()

	json.NewDecoder(file).Decode(&results)
	return results
}

// Сохранение результатов в фай
func saveResult(outcome string, attempts int) {
	results := loadResults()
	results = append(results, Result{
		Date:     time.Now().Format("2006-01-02 15:04:05"),
		Outcome:  outcome,
		Attempts: attempts,
	})

	file, err := os.Create("results.json")
	if err != nil {
		fmt.Println("Ошибка сохранения результатов:", err)
		return
	}
	defer file.Close()

	json.NewEncoder(file).Encode(results)
}

// Основная функция игры
func playGame() {
	maxNumber, attempts := difficultySelection()
	secret := generateNumber(maxNumber)
	guesses := []int{}

	fmt.Printf(Yellow+"Игра началась! Я загадал число от 1 до %d.\n"+Reset, maxNumber)

	for i := 1; i <= attempts; i++ {
		guess := userInput(i, attempts)
		guesses = append(guesses, guess)

		if guess == secret {
			fmt.Println(Green + "Вы угадали!🙌" + Reset)
			saveResult("Победа", i)
			return
		}

		showClue(guess, secret)
		fmt.Print("Прошлые попытки: ")
		fmt.Println(guesses)
		fmt.Println()
	}

	fmt.Println(Red+"Вы проиграли! 😢"+Reset, "Секретное число было:", secret)
	saveResult("Проигрыш", attempts)
}

// Выбор уровня сложности
func difficultySelection() (maxNumber int, attempts int) {
	for {
		fmt.Println(Blue + "Выбери уровень сложности:" + Reset)
		fmt.Println("1 — Easy (1–50, 15 попыток)")
		fmt.Println("2 — Medium (1–100, 10 попыток)")
		fmt.Println("3 — Hard (1–200, 5 попыток)")
		fmt.Print("Ввод: ")

		var diffInput int
		_, err := fmt.Scan(&diffInput)
		if err != nil || (diffInput < 1 || diffInput > 3) {
			fmt.Println(Red + "Некорректный выбор уровня сложности! Попробуй снова." + Reset)
			clearStdin()
			continue
		}

		switch diffInput {
		case 1:
			return 50, 15
		case 2:
			return 100, 10
		case 3:
			return 200, 5
		}
	}
}

// Генерация секретного числа
func generateNumber(max int) int {
	return rand.Intn(max) + 1
}

// Пользовательский ввод
func userInput(attempt, total int) int {
	for {
		fmt.Printf(Yellow+"Попытка %d/%d: "+Reset, attempt, total)
		fmt.Print("Введи число: ")

		var guess int
		_, err := fmt.Scan(&guess)
		if err != nil {
			fmt.Println(Red + "Ошибка: введи именно ЧИСЛО!" + Reset)
			clearStdin()
			continue
		}
		return guess
	}
}

// Подсказки и сравнение с секретным числом
func showClue(guess, secret int) {
	diff := abs(guess - secret)

	if diff <= 5 {
		fmt.Println("🔥 Горячо!")
	} else if diff <= 15 {
		fmt.Println("🙂 Тепло!")
	} else {
		fmt.Println("❄️ Холодно.")
	}

	if guess < secret {
		fmt.Println("Секретное число больше 👆")
	} else {
		fmt.Println("Секретное число меньше 👇")
	}
}

// Запуститть новую игру
func Replay() bool {
	fmt.Print(Blue + "Хочешь сыграть ещё раз? (y/n): " + Reset)
	var again string
	fmt.Scan(&again)
	if again == "y" || again == "Y" {
		return true
	}
	return false

}

// Вспомогательная функция - МОДУЛЬ ЧИСЛА (разница для подсказок)
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Очистка stdin при неправильном пользовательском вводе
func clearStdin() {
	var tmp string
	fmt.Scanln(&tmp)
}
