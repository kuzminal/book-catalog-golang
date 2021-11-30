package recilience

import (
	"book-catalog/core"
	"errors"
	"github.com/sony/gobreaker"
	"log"
	"sync"
	"time"
)

var CB *gobreaker.CircuitBreaker

func init() {
	var st gobreaker.Settings
	st.Name = "HTTP GET"
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.6
	}
	st.OnStateChange = func(_ string, from gobreaker.State, to gobreaker.State) {
		// Handler for every state change. We'll use for debugging purpose
		log.Println("state changed from", from.String(), "to", to.String())
	}
	// When to flush counters int the Closed state
	st.Interval = 7 * time.Second
	// Time to switch from Open to Half-open
	st.Timeout = 10 * time.Second

	CB = gobreaker.NewCircuitBreaker(st)
	log.Println("CB initialized")
}

type Circuit func(book *core.Book) error

func Breaker(circuit Circuit, failureThreshold uint) Circuit {
	var consecutiveFailures int = 0
	var lastAttempt = time.Now()
	var m sync.RWMutex
	return func(book *core.Book) error {
		m.RLock() // Установить "блокировку чтения"
		d := consecutiveFailures - int(failureThreshold)
		if d >= 0 {
			shouldRetryAt := lastAttempt.Add(time.Second * 2 << d)
			if !time.Now().After(shouldRetryAt) {
				m.RUnlock()
				log.Println("if unreach")
				return errors.New("service unreachable")
			}
		}
		m.RUnlock()          // Освободить блокировку чтения
		err := circuit(book) // Послать запрос, как обычно
		m.Lock()             // Заблокировать общие ресурсы
		defer m.Unlock()
		lastAttempt = time.Now() // Зафиксировать время попытки
		if err != nil {          // Если Circuit вернула ошибку,
			log.Println("if error")
			log.Println(err)
			consecutiveFailures++ // увеличить счетчик ошибок
			return err            // и вернуть ошибку
		}
		consecutiveFailures = 0 // Сбросить счетчик ошибок
		return nil
	}
}
