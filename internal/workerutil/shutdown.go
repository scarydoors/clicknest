package workerutil

import (
	"context"
	"errors"
	"sync"
)

type Shutdowner interface {
	Shutdown(context.Context) error
}

func ShutdownServices(ctx context.Context, shutdowners... Shutdowner) error {
	errChan := make(chan error, len(shutdowners))
	var wg sync.WaitGroup

	for _, shutdowner := range shutdowners {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errChan <- shutdowner.Shutdown(ctx)
		}()
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	var err error
	for errc := range errChan {
		err = errors.Join(err, errc)
	}

	return err
}
