package main

import (
	"sync"
	"time"
)

var timers = map[time.Time]bool{}
var timersMutex = sync.Mutex{}

func resetTimers() {
	timersMutex.Lock()
	timers = map[time.Time]bool{}
	timersMutex.Unlock()
}

func addTimer(t time.Time) {
	timersMutex.Lock()
	t = t.Truncate(time.Second)
	timers[t] = true
	timersMutex.Unlock()
}

func runTimers() {
	for {
		timersMutex.Lock()
		now := time.Now()
		now = now.Truncate(time.Second)

		if _, ok := timers[now]; ok {
			go sendUpdate()
		}
		timersMutex.Unlock()

		time.Sleep(time.Second)
	}
}
