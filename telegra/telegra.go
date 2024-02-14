package telegra

import (
	"github.com/alehano/translit"
)

type Telegra struct {
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

	// преобразовать запрос запрос в транслит
	trans := translit.ToLatin("Привет, Человек", translit.RussianASCII)

	// сформировать запрос
	url := base + "/" + trans

}
