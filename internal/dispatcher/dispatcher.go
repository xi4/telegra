package dispatcher

import (
	"fmt"
	"strings"
	"telegra/internal/db"
	"telegra/internal/job"
	w "telegra/internal/worker"
	"time"
)

var (
	base = "https://telegra.ph"
)

// Dispatcher описывает диспетчера, который управляет работниками
type Dispatcher struct {
	WorkerQueue chan chan job.Job
	MaxWorkers  int
	JobQueue    chan job.Job
	db          *db.DB
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	db := db.New()
	workerQueue := make(chan chan job.Job, maxWorkers)
	return &Dispatcher{WorkerQueue: workerQueue, MaxWorkers: maxWorkers, db: db}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.MaxWorkers; i++ {
		worker := w.NewWorker(d.WorkerQueue, i, d.db)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for job := range d.JobQueue {
		_job := job // Create a new variable and assign the value of job to it
		go func() {
			workerJobQueue := <-d.WorkerQueue
			workerJobQueue <- _job
		}()
	}
}

func (d *Dispatcher) Search(query string) {
	// удаляем из текста все символы кроме букв и цифр и пробелы заменяем на -
	query = removeSpecialChars(query)
	// преобразовать запрос запрос в транслит
	text := translateToLatin(query)
	// сохраняем слово в базе
	word_id := d.db.AddWord(query, text)
	// создаем базовый url
	baseUrl := base + "/" + text
	for s := time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.UTC); s.Year() == time.Now().Year(); s = s.AddDate(0, 0, 1) {
		url := fmt.Sprintf("%s-%02d-%02d", baseUrl, s.Month(), s.Day())
		d.JobQueue <- job.Job{Payload: job.Payload{Url: url, Id: word_id}}
	}
}

// removeSpecialChars - удаляет из текста все символы кроме букв и цифр
func removeSpecialChars(text string) string {
	// массив символов, которые нужно удалить
	specialChars := []string{":", ",", ".", "!", "?", "(", ")", "«", "»", "—"}
	// удаляем все символы из массива
	for _, char := range specialChars {
		text = strings.ReplaceAll(text, char, "")
	}
	// возвращаем текст
	return text
}

func translateToLatin(text string) string {
	// text делаем в нижнем регистре
	text = strings.ToLower(text)
	// массив символов, которые нужно заменить
	chars := map[string]string{
		"а": "a", "б": "b", "в": "v", "г": "g", "д": "d", "е": "e", "ё": "e", "ж": "zh", "з": "z", "и": "i", "й": "i", "к": "k", "л": "l", "м": "m", "н": "n", "о": "o", "п": "p", "р": "r", "с": "s", "т": "t", "у": "u", "ф": "f", "х": "kh", "ц": "ts", "ч": "ch", "ш": "sh", "щ": "shch", "ы": "y", "ъ": "ie", "э": "e", "ю": "iu", "я": "ia",
	}
	// заменяем все символы из массива
	for k, v := range chars {
		text = strings.ReplaceAll(text, k, v)
	}

	// пробелы заменяем на -
	text = strings.ReplaceAll(text, " ", "-")

	// возвращаем текст
	return text

}
