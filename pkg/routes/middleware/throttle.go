package middleware

import (
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

// NewThrottleGroup will throttle a request per-ip.  Of course, this won't work
// for anything behind a load-balancer (or at least will scale differently). But
// can add some security for things that have direct-use of the application
// to slow down certain types of brute-force attacks
func NewThrottleGroup(perIP int, d time.Duration) echo.MiddlewareFunc {
	sema := make(map[string]chan int)
	var mux sync.Mutex

	waitFor := func(ip string) {
		mux.Lock()
		s := sema[ip]
		if s == nil {
			s = make(chan int, perIP)
			sema[ip] = s
		}
		mux.Unlock()

		s <- 1
	}

	release := func(ip string) {
		mux.Lock()
		s := sema[ip]
		if s == nil {
			panic("There isn't a releaseable semaphore")
		}

		// This isn't purely correct, because the
		// semaphore could add-one above after checking len(0), but it's good enough and won't break
		<-s

		if len(s) == 0 {
			delete(sema, ip)
		}
		mux.Unlock()
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			waitFor(ip)
			time.Sleep(d)
			err := next(c)
			release(ip)
			return err
		}
	}
}
