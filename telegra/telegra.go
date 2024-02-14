package telegra

import (
	"net/http"
	"strings"
	"time"

	"github.com/alehano/translit"
)

type Telegra struct {
	out chan *Page // канал для передачи результата

}

type Page struct {
	Title   string
	Address string
	Images  []Image
}

type Image struct {
	Url string
}

var (
	base = "https://telegra.ph"
)

func New() *Telegra {
	return &Telegra{}
}

func (t *Telegra) Search(query string) {
	// удаляем из текста все символы кроме букв и цифр и пробелы заменяем на -
	text := removeSpecialChars(query)
	// преобразовать запрос запрос в транслит
	trans := translit.ToLatin(text, translit.RussianASCII)

	// запускаем цикл по количеству дней в году
	for i := 1; i <= 365; i++ {
		// формируем адрес страницы
		url := base + "/" + trans + "-" + getDate(i)
		// получаем страницу
		t.GetPage(url)

	}
}

func (t *Telegra) GetPage(url string) {
	// загружаем страницу
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	// если статус ответа не 200, то выходим
	if resp.StatusCode != 200 {
		return
	}

	// загружаем страницу добавляем в канал

}

// getDate - преобразует номер дня в дату типа месяц-день
func getDate(day int) string {
	// Предполагаем, что день начинается с 1 января текущего года
	t := time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
	// Добавляем количество дней к этой дате
	date := t.AddDate(0, 0, day-1)
	// Возвращаем дату в формате месяц-день
	return date.Format("01-02")
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
