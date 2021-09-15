package pagination

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrConditionNeverMet = errors.New("looped all items and the condition was never met")
)

type Page interface {
	Num() int
	Size() int
	Total() int
	Content() interface{}
}

// Pageable should be implemented by evey struct that supports pagination
type Pageable interface {
	GetPage(ctx context.Context, page, limit int) (Page, error)
	NextPage(ctx context.Context, page Page) (Page, error)
}

// HasNext is a helper function which checks whether there are more pages to fetch
func HasNext(page Page) bool {
	return page.Num() < page.Total()
}

type UntilTrueFunc func(interface{}) (bool, error)

// LoopUntil takes a pageable and loops over it until untilFunc returns true or an error.
func LoopUntil(ctx context.Context, pageable Pageable, untilFunc UntilTrueFunc) error {
	page, err := pageable.GetPage(ctx, 1, 10)
	if err != nil {
		return err
	}

	content := reflect.ValueOf(page.Content())
	switch content.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		stopped := false
		for i := 0; i < reflect.ValueOf(content).Len(); i++ {
			stop, err := untilFunc(content.Index(i))
			if err != nil {
				return err
			}
			if stop {
				stopped = true
				break
			}
		}
		if !stopped {
			return ErrConditionNeverMet
		}
	default:
		panic(fmt.Sprintf("The page content is supposed to be of type slice or array but was %T", page.Content()))
	}
	return nil
}

type CancelFunc func()

// AsChan takes a Pageable and returns its Pageable.Content via a channel until there are no more pages or
// CancelFunc gets called by the consumer
func AsChan(ctx context.Context, pageable Pageable) (chan interface{}, CancelFunc) {
	consumer := make(chan interface{})
	done := make(chan interface{})
	cancel := func() {
		close(done)
	}

	go func() {
		defer close(consumer)
		defer close(done)
		_ = LoopUntil(ctx, pageable, func(i interface{}) (bool, error) {
			select {
			case consumer <- i:
			case _, ok := <-done:
				if !ok {
					return true, nil
				}
			}
			return false, nil
		})
	}()
	return consumer, cancel
}
