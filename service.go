package main

import (
	"book-catalog/frontend"
	"book-catalog/storage"
	"log"
	"time"
)

func asyncHandleChannels() {
	for {
		select {
		case book, ok := <-frontend.CreateBookChannel:
			if ok {
				//некий retry на случай ошибки, положу его пока сюда просто чтобы был
				/*err := storage.SaveBook(&book)
				base, capacity := time.Second, time.Minute
				for backoff := base; !strings.Contains(err.Error(), "duplicate key"); backoff <<= 1 {
					if backoff > capacity {
						backoff = capacity
					}
					jitter := rand.Int63n(int64(backoff * 3))
					sleep := base + time.Duration(jitter)
					time.Sleep(sleep)
					err = storage.SaveBook(&book)
				}*/
				//конец повнотрных попыток
				_, err := storage.UpdateBook(&book)
				if err != nil {
					log.Println(err.Error())
				}
			}
		case book, ok := <-frontend.UpdateBookChannel:
			if ok {
				_, err := storage.UpdateBook(&book)
				if err != nil {
					log.Fatal(err.Error())
				}
			}
		case isbn, ok := <-frontend.DeleteBookChannel:
			if ok {
				err := storage.DeleteBook(isbn)
				if err != nil {
					log.Fatal(err.Error())
				}
			}
		case <-time.After(4 * time.Second):
		}
	}
}

func main() {

	go asyncHandleChannels()
	fe := frontend.RestFrontEnd{}
	fe.Start()

}
