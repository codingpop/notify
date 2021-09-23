package notification

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Notification struct {
	url        string
	interval   time.Duration
	mu         sync.Mutex
	queue      *list.List // messages to be sent
	sem        chan int
	loopSignal chan struct{}
	errors     chan error
}

var ErrRequest = errors.New("request failed")

func New(ctx context.Context, url string, interval time.Duration, poolSize int, errors chan error) *Notification {
	n := &Notification{
		url:        url,
		interval:   interval,
		sem:        make(chan int, poolSize),
		queue:      list.New(),
		loopSignal: make(chan struct{}, 1),
		errors:     errors,
	}

	// Start the loop
	go n.loop(ctx)

	return n
}

func (n *Notification) Enqueue(msg string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.queue.PushBack(msg)

	n.tickleLoop()
}

func (n *Notification) loop(ctx context.Context) {
	for {
		select {
		case <-n.loopSignal:
			n.dequeue(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (n *Notification) dequeue(ctx context.Context) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.queue.Len() == 0 {
		return
	}
	select {
	case n.sem <- 1:
		element := n.queue.Front()
		n.queue.Remove(element)
		msg := element.Value.(string)

		go func() {
			err := n.postMessage(ctx, msg)
			if err != nil {
				n.errors <- err
			}
		}()

	default:
		fmt.Println("Goroutine pool maxed out!")
		fmt.Println("Processing will resume once the pool is freed")
	}
}

func (n *Notification) postMessage(ctx context.Context, msg string) (retErr error) {
	// make room for another request
	defer func() {
		<-n.sem
		n.tickleLoop()
	}()

	body := strings.NewReader(msg)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.url, body)
	if err != nil {
		return err
	}

	n.respectInterval()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			retErr = err
		}
	}()

	if resp.StatusCode >= http.StatusBadRequest {
		return ErrRequest
	}

	return nil
}

func (n *Notification) tickleLoop() {
	select {
	case n.loopSignal <- struct{}{}:
	default:
		// To avoid blocking
	}
}

func (n *Notification) respectInterval() {
	<-time.NewTimer(n.interval).C
}
