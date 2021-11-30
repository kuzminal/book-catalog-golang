package recilience

import (
	"context"
	"time"
)

// Effector -- это функция, доступ к которой требуется регулировать.
type Effector func(context.Context) (string, error)

// Throttled служит оберткой для Effector. Она принимает те же параметры плюс
// строку "UID", идентифицирующую вызывающего, и возвращает то же значение
// плюс логическое значение true, если вызов не подвергался регулированию.
type Throttled func(context.Context, string) (bool, string, error)

// Структура bucket запоминает информацию о последнем запросе, связанном с UID.
type bucket struct {
	tokens uint
	time   time.Time
}

// Throttle принимает функцию Effector и возвращает функцию Throttled
// с корзиной, соответствующей заданному UID и заполненной максимальным
// количеством жетонов. Корзина постоянно пополняется жетонами со скоростью
// refill каждый интервал времени d.
func Throttle(e Effector, max uint, refill uint, d time.Duration) Throttled {
	// buckets отображает строки UID в конкретные корзины
	buckets := map[string]*bucket{}
	return func(ctx context.Context, uid string) (bool, string, error) {
		b := buckets[uid]
		// Это новая запись! Предполагается, что емкость >= 1.
		if b == nil {
			buckets[uid] = &bucket{tokens: max - 1, time: time.Now()}
			str, err := e(ctx)
			return true, str, err
		}
		// Подсчитать, сколько жетонов можно добавить в корзину, учитывая
		// время, прошедшее с момента предыдущего запроса.
		refillInterval := uint(time.Since(b.time) / d)
		tokensAdded := refill * refillInterval
		currentTokens := b.tokens + tokensAdded
		// Если жетонов недостаточно, вернуть false.
		if currentTokens < 1 {
			return false, "", nil
		}
		// Если корзина пополнилась, запомнить текущее время.
		// Иначе выяснить, когда в последний раз добавлялись жетоны.
		if currentTokens > max {
			b.time = time.Now()
			b.tokens = max - 1
		} else {
			deltaTokens := currentTokens - b.tokens
			deltaRefills := deltaTokens / refill
			deltaTime := time.Duration(deltaRefills) * d
			b.time = b.time.Add(deltaTime)
			b.tokens = currentTokens - 1
		}
		str, err := e(ctx)
		return true, str, err
	}
}
