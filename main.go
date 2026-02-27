package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/telebot.v3"
)

type Movie struct {
	Rating  int
	Id      int
	Title   string
	Watched bool
}

var movies map[int64][]Movie

func main() {
	godotenv.Load()

	movies = make(map[int64][]Movie)
	LoadMovies()
	pref := telebot.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	b.Handle("/start", func(c telebot.Context) error {
		menu := &telebot.ReplyMarkup{ResizeKeyboard: true}

		btnAdd := menu.Text("➕ Добавить фильм")
		btnList := menu.Text("📋 Список фильмов")
		btnRandom := menu.Text("🎲 Случайный фильм")

		menu.Reply(
			menu.Row(btnAdd),
			menu.Row(btnList, btnRandom),
		)

		return c.Send("Привет! Выбери действие:", menu)
	})
	b.Handle("➕ Добавить фильм", func(c telebot.Context) error {
		return c.Send("Напиши /add Название фильма")
	})

	b.Handle("📋 Список фильмов", func(c telebot.Context) error {
		return c.Send(GetMovieList(c.Sender().ID))
	})

	b.Handle("🎲 Случайный фильм", func(c telebot.Context) error {
		userMovies := movies[c.Sender().ID]
		var unwatched []Movie
		for _, m := range userMovies {
			if !m.Watched {
				unwatched = append(unwatched, m)
			}
		}
		if len(unwatched) == 0 {
			return c.Send("Нет непросмотренных фильмов!")
		}
		random := unwatched[rand.Intn(len(unwatched))]
		return c.Send(fmt.Sprintf("Смотри сегодня: %s", random.Title))
	})

	b.Handle("/add", func(c telebot.Context) error {
		title := strings.TrimPrefix(c.Message().Text, "/add ")
		userID := c.Sender().ID
		userMovies := movies[userID]

		// находим максимальный ID
		maxID := 0
		for _, m := range userMovies {
			if m.Id > maxID {
				maxID = m.Id
			}
		}

		movie := Movie{
			Id:      maxID + 1,
			Title:   title,
			Watched: false,
		}
		movies[userID] = append(movies[userID], movie)
		SaveMovies()
		return c.Send("Фильм добавлен!")
	})
	b.Handle("/list", func(ctx telebot.Context) error {
		return ctx.Send(GetMovieList(ctx.Sender().ID))
	})
	b.Handle("/search", func(ctx telebot.Context) error {
		title := strings.TrimPrefix(ctx.Message().Text, "/search ")
		userMovies := movies[ctx.Sender().ID]
		var results []Movie
		for _, movie := range userMovies {
			if strings.Contains(strings.ToLower(movie.Title), strings.ToLower(title)) {
				results = append(results, movie)
			}
		}
		if len(results) == 0 {
			return ctx.Send("Фильмы не найдены!")
		}
		var list string
		for _, movie := range results {
			status := "не просмотрен"
			if movie.Watched {
				status = "просмотрен"
			}
			list += fmt.Sprintf("%d. %s - %s - рейтинг: %d\n", movie.Id, movie.Title, status, movie.Rating)
		}
		return ctx.Send(list)
	})

	b.Handle("/watched", func(c telebot.Context) error {
		idStr := strings.TrimPrefix(c.Message().Text, "/watched ")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Send("Неверный ID!")
		}
		WatchMovie(c.Sender().ID, id)
		return c.Send("Фильм отмечен как просмотренный!")
	})
	b.Handle("/delete", func(ctx telebot.Context) error {
		idStr := strings.TrimPrefix(ctx.Message().Text, "/delete ")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return ctx.Send("Неверный ID!")
		}
		DeleteMovie(ctx.Sender().ID, id)
		return ctx.Send("Фильм удален!")
	})
	b.Handle("/rate", func(ctx telebot.Context) error {
		args := strings.Split(strings.TrimPrefix(ctx.Message().Text, "/rate "), " ")
		if len(args) != 2 {
			return ctx.Send("Использование: /rate <ID> <рейтинг>")
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return ctx.Send("Неверный ID!")
		}
		rating, err := strconv.Atoi(args[1])
		if err != nil || rating < 1 || rating > 5 {
			return ctx.Send("Рейтинг должен быть от 1 до 5!")
		}
		RateMovie(ctx.Sender().ID, id, rating)
		return ctx.Send("Фильм оценен!")
	})
	b.Start()
}

func SaveMovies() {
	data, err := json.Marshal(movies)
	if err != nil {
		log.Println("Error marshaling movies:", err)
		return
	}
	err = os.WriteFile("movies.json", data, 0644)
	if err != nil {
		log.Println("Error writing movies file:", err)
	}
}
func LoadMovies() {
	data, err := os.ReadFile("movies.json")
	// превратить data обратно в список через json.Unmarshal
	if err != nil {
		log.Println("Error reading movies file:", err)
		return
	}
	err = json.Unmarshal(data, &movies)
	if err != nil {
		log.Println("Error unmarshaling movies:", err)
	}
}

func WatchMovie(userID int64, id int) {
	userMovies := movies[userID]
	for i, movie := range userMovies {
		if movie.Id == id {
			movies[userID][i].Watched = true
			SaveMovies()
			return
		}
	}
}

func DeleteMovie(userID int64, id int) {
	userMovies := movies[userID]
	for i, movie := range userMovies {
		if movie.Id == id {
			movies[userID] = append(movies[userID][:i], movies[userID][i+1:]...)
			SaveMovies()
			return
		}
	}
}

func RateMovie(userID int64, id int, rating int) {
	userMovies := movies[userID]
	for i, movie := range userMovies {
		if movie.Id == id {
			movies[userID][i].Rating = rating
			SaveMovies()
			return
		}
	}
}

func GetMovieList(userID int64) string {
	userMovies := movies[userID]
	if len(userMovies) == 0 {

		return "Список пуст!"
	}
	var list string
	for _, movie := range userMovies {
		status := "не просмотрен"
		if movie.Watched {
			status = "просмотрен"
		}
		list += fmt.Sprintf("%d. %s - %s - рейтинг: %d\n", movie.Id, movie.Title, status, movie.Rating)
	}
	return list
}
