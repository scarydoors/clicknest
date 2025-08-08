package workerutil

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type Shutdowner interface {
	Shutdown(context.Context) error
}

type Service struct {
	Name       string
	Shutdowner Shutdowner
}

type ShutdownError struct {
	Name string
	err  error
}

func (e *ShutdownError) Error() string {
	return fmt.Sprintf("unable to shutdown gracefully: %s", e.err)
}

func (e *ShutdownError) Unwrap() error {
	return e.err
}

func ShutdownServices(ctx context.Context, services ...Service) error {
	errChan := make(chan *ShutdownError, len(services))
	var wg sync.WaitGroup

	for _, service := range services {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := service.Shutdowner.Shutdown(ctx); err != nil {
				err := &ShutdownError{
					Name: service.Name,
					err:  err,
				}
				errChan <- err
			}
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
