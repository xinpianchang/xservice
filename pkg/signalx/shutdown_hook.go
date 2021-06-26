package signalx

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	shutdownHooks      = make([]ShutdownHook, 0, 64)
	shutdownHooksMutex sync.Mutex
)

type ShutdownHook func(os.Signal)

func AddShutdownHook(hook ShutdownHook) {
	shutdownHooksMutex.Lock()
	defer shutdownHooksMutex.Unlock()
	shutdownHooks = append(shutdownHooks, hook)
}

func ShutdownListen(signals ...os.Signal) {
	ch := make(chan os.Signal, 1)
	if len(signals) == 0 {
		signals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT}
	}
	signal.Notify(ch, signals...)
	sig := <-ch
	var wg sync.WaitGroup
	for i := len(shutdownHooks) - 1; i >= 0; i-- {
		wg.Add(1)
		go func(sig os.Signal, fn func(os.Signal)) {
			defer wg.Done()
			fn(sig)
		}(sig, shutdownHooks[i])
	}
	wg.Wait()
}

// Shutdown direct shutdown without signal
func Shutdown() {
	for i := len(shutdownHooks) - 1; i >= 0; i-- {
		shutdownHooks[i](syscall.SIGTERM)
	}
}
