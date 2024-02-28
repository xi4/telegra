package job

// Job описывает задание, которое должно быть выполнено
type Job struct {
	Payload Payload
}

type Payload struct {
	Id  uint
	Url string
}
