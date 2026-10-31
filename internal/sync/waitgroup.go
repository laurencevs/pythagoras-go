package sync

import "sync"

type WaitGroup struct {
	funcs []func()
	wg    sync.WaitGroup
}

// Queue fns to run when Go() is called.
func (wg *WaitGroup) Add(fns ...func()) {
	wg.funcs = append(wg.funcs, fns...)
	wg.wg.Add(len(fns))
}

// Run each of wg's functions in its own goroutine, closing done when all have
// finished.
func (wg *WaitGroup) Go() (done <-chan struct{}) {
	wgDone := make(chan struct{})
	go func() {
		wg.wg.Wait()
		close(wgDone)
	}()
	for _, f := range wg.funcs {
		go func() {
			f()
			wg.wg.Done()
		}()
	}
	return wgDone
}
