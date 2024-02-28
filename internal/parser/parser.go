package telegrapkg

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type TelegraPkg struct {
	Answer []Page
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

func New() *TelegraPkg {
	return &TelegraPkg{}
}

func (t *TelegraPkg) Search(query string) {
	// удаляем из текста все символы кроме букв и цифр и пробелы заменяем на -
	text := removeSpecialChars(query)
	// преобразовать запрос запрос в транслит
	trans := translateToLatin(text)
	baseUrl := base + "/" + trans
	t.IterateDays(baseUrl)
}

func (t *TelegraPkg) GetPage(url string, ch chan<- Page) {
	resp, err := http.Get(url)
	if err != nil {
		ch <- Page{}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			ch <- Page{}
			return
		}
		// ищем заголовок страницы
		title := doc.Find("title").Text()
		var page Page
		page.Title = title
		page.Address = url
		// ищем все изображения на странице
		doc.Find("img").Each(func(i int, s *goquery.Selection) {
			src, _ := s.Attr("src")
			var img Image

			img.Url = url + src
			page.Images = append(page.Images, img)
		})
		ch <- page
	} else {
		ch <- Page{}
	}
}

func (t *TelegraPkg) IterateDays(baseURL string) {
	ch := make(chan Page)
	for d := time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.UTC); d.Year() == time.Now().Year(); d = d.AddDate(0, 0, 1) {
		url := fmt.Sprintf("%s-%02d-%02d", baseURL, d.Month(), d.Day())
		go t.GetPage(url, ch)
	}
	for range time.Tick(time.Second) {
		select {
		case page := <-ch:
			// если страница пустая, то пропускаем
			if page.Address == "" {
				continue
			}
			t.Answer = append(t.Answer, page)
		default:
			break
		}
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
	// а → a, б → b, в → v, г → g, д → d, е → e, ё → e, ж → zh, з → z, и → i, й → i, к → k, л → l, м → m, н → n, о → o, п → p, р → r, с → s, т → t, у → u, ф → f, х → кh, ц → ts, ч → сh, ш → sh, щ → shch, ы → y, ъ → ie, э → e, ю → iu, я → ia.
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
