package worker

import (
	"fmt"
	"net/http"
	"telegra/internal/db"
	"telegra/internal/job"

	"github.com/PuerkitoBio/goquery"
)

// Worker описывает работника, который выполняет задания
type Worker struct {
	id          int
	JobQueue    chan job.Job
	WorkerQueue chan chan job.Job
	quitChan    chan bool
	db          *db.DB
}

// NewWorker создает нового работника
func NewWorker(workerQueue chan chan job.Job, id int, db *db.DB) Worker {
	worker := Worker{
		id:          id,
		JobQueue:    make(chan job.Job),
		WorkerQueue: workerQueue,
		quitChan:    make(chan bool),
		db:          db,
	}

	return worker
}

// Start метод запускает работника на выполнение заданий
func (w Worker) Start() {
	go func() {
		for {
			// Регистрируем свободный рабочий
			w.WorkerQueue <- w.JobQueue
			select {
			case job := <-w.JobQueue:
				// Мы получили задание
				page := w.isCurrentPage(job.Payload.Url)
				if page.Title != "" {
					w.db.AddPage(w.pageToDbPage(page), job.Payload.Id)
					start := 2
					for {
						url := fmt.Sprintf("%s-%02d", job.Payload.Url, start)
						newPage := w.isCurrentPage(url)
						if page.Title == "" {
							break
						}
						w.db.AddPage(w.pageToDbPage(newPage), job.Payload.Id)
						start++
					}
				}
				println("Worker", w.id, "processed job", job.Payload.Url)
			case <-w.quitChan:
				// Мы получили сигнал остановиться
				return
			}
		}
	}()
}

func (w Worker) pageToDbPage(page Page) db.Page {
	var dbPage db.Page
	dbPage.Title = page.Title
	dbPage.Address = page.Address
	for _, img := range page.Images {
		dbPage.Images = append(dbPage.Images, db.Image{Url: img.Url})
	}
	for _, video := range page.Videos {
		dbPage.Videos = append(dbPage.Videos, db.Video{Url: video.Url})
	}
	return dbPage
}

func (w Worker) isCurrentPage(url string) Page {
	resp, err := http.Get(url)
	if err != nil {
		return Page{}
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return Page{}
		}

		title := doc.Find("title").Text()
		var page Page
		page.Title = title
		page.Address = url

		doc.Find("img").Each(func(i int, s *goquery.Selection) {
			src, _ := s.Attr("src")
			var img Image

			img.Url = url + src
			page.Images = append(page.Images, img)
		})

		doc.Find("video").Each(func(i int, s *goquery.Selection) {
			src, _ := s.Attr("src")
			var video Video

			video.Url = url + src
			page.Videos = append(page.Videos, video)
		})

		return page
	}

	return Page{}

}

// Stop сигнализирует работнику остановиться
func (w Worker) Stop() {
	go func() {
		w.quitChan <- true
	}()
}

type Page struct {
	Title   string
	Address string
	Images  []Image
	Videos  []Video
}

type Image struct {
	Url string
}

type Video struct {
	Url string
}
