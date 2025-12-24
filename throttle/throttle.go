package throttle

import (
	"sync"
	"time"
)

// global bucket state
var (
	mu     sync.Mutex
	rate   int       // байт в секунду (0 — блокировать всё)
	tokens float64   // токены в байтах
	last   time.Time // время последнего refill
)

func init() {
	last = time.Now()
}

// SetUploadLimit устанавливает скорость в Mbit/s
func SetUploadLimit(mbit int) {
	mu.Lock()
	rate = mbit * 1024 * 1024 / 8
	// при изменении лимита моментально накапливаем полную «ёмкость»
	tokens = float64(rate)
	last = time.Now()
	mu.Unlock()
}

// StopUpload — эквивалент скорости=0
func StopUpload() {
	mu.Lock()
	rate = 0
	tokens = 0
	last = time.Now()
	mu.Unlock()
}

// GetUploadLimit возвращает текущий лимит в байт/сек
func GetUploadLimit() int {
	mu.Lock()
	defer mu.Unlock()
	return rate
}

// Take блокирует (с помощью time.Sleep) до тех пор, пока
// не окажется в бакете >= n токенов (байт).
func Take(n int) {
	for {
		mu.Lock()
		now := time.Now()
		elapsed := now.Sub(last).Seconds()
		last = now
		// накапливаем токены
		tokens += elapsed * float64(rate)
		if tokens > float64(rate) {
			tokens = float64(rate) // capacity = rate
		}
		// если хватило — забираем и выходим
		if rate > 0 && tokens >= float64(n) {
			tokens -= float64(n)
			mu.Unlock()
			return
		}
		// рассчитаем паузу
		r := rate
		mu.Unlock()

		var wait time.Duration
		if r > 0 {
			// нужно накапливать недостающее
			mu.Lock()
			deficit := float64(n) - tokens
			mu.Unlock()
			if deficit < 0 {
				deficit = 0
			}
			wait = time.Duration(deficit / float64(r) * float64(time.Second))
		} else {
			// speed==0 — просто ждём, пока не установят >0
			wait = 100 * time.Millisecond
		}
		time.Sleep(wait)
	}
}
