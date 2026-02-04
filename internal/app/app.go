package app

import "sync"

type App struct {
	MsgChannel chan []byte
	RWMutext   *sync.RWMutex
	WG         *sync.WaitGroup
}

func NewApp(msgChann chan []byte, rwmut *sync.RWMutex, wg *sync.WaitGroup) *App {
	return &App{
		MsgChannel: msgChann,
		RWMutext:   rwmut,
		WG:         wg,
	}
}
