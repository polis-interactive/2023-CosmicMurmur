package bus

import (
	"errors"
	"fmt"
	"github.com/polis-interactive/2023-CosmicMurmur/internal/domain"
	"log"
	"sync"
	"time"
)

type callbackProxy struct {
	id   int
	fn   func()
	next *callbackProxy
}

type eventBus struct {
	lastEventId int
	pubSubMap   map[domain.Event]*callbackProxy
	idEventMap  map[int]domain.Event
	eventQueue  chan domain.Event
	mu          *sync.RWMutex
	wg          *sync.WaitGroup
	shutdowns   chan struct{}
}

var _ domain.EventHandler = (*eventBus)(nil)

func newEventHandler() *eventBus {
	return &eventBus{
		lastEventId: 0,
		pubSubMap:   make(map[domain.Event]*callbackProxy),
		idEventMap:  make(map[int]domain.Event),
		eventQueue:  nil,
		mu:          &sync.RWMutex{},
		wg:          &sync.WaitGroup{},
		shutdowns:   nil,
	}
}

func (b *eventBus) SubscribeToEvent(e domain.Event, fn func()) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	log.Println(fmt.Sprintf(
		"Bus, EventHandler, SubscribeToEvent: subscribing to %s, giving id %d",
		e.String(),
		b.lastEventId,
	))
	newCallback := &callbackProxy{
		id:   b.lastEventId,
		fn:   fn,
		next: nil,
	}
	b.idEventMap[b.lastEventId] = e
	if callback, ok := b.pubSubMap[e]; ok {
		for callback.next != nil {
			callback = callback.next
		}
		callback.next = newCallback
	} else {
		b.pubSubMap[e] = &callbackProxy{
			id:   -1,
			fn:   nil,
			next: newCallback,
		}
	}
	b.lastEventId += 1
	return b.lastEventId
}

func (b *eventBus) UnsubscribeToEvent(handlerId int) {
	log.Println(fmt.Sprintf("Bus, EventHandler, UnsubscribeToEvent: unsubscribing id %d", handlerId))
	b.mu.Lock()
	defer b.mu.Unlock()
	// invalid handlerId
	if handlerId > b.lastEventId || handlerId < 0 {
		return
	}
	// short circuit if we've already removed the id
	if _, ok := b.idEventMap[handlerId]; !ok {
		return
	}
	eventType := b.idEventMap[handlerId]
	delete(b.idEventMap, handlerId)

	// short circuit if the pubSub map for event type
	// is missing
	if _, ok := b.pubSubMap[eventType]; !ok {
		return
	}
	callback := b.pubSubMap[eventType]
	// search for the callback with id handlerId
	// to remove
	var removeCallback *callbackProxy = nil
	for callback.next != nil {
		if callback.next.id == handlerId {
			removeCallback = callback.next
			break
		}
		callback = callback.next
	}
	// found the callback we were looking for; go ahead
	// and delete it
	if removeCallback != nil {
		callback.next = removeCallback.next
	}
	// cleanup if necessary; can be an orphaned pubSub entry,
	// or we only had the removed callback in the list
	if callback.id == -1 && callback.next == nil {
		delete(b.pubSubMap, eventType)
	}
}

func (b *eventBus) HandleEvent(e domain.Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.eventQueue != nil {
		b.eventQueue <- e
	}
}

func (b *eventBus) startupEventLoop() {
	log.Println("Bus, EventHandler, startupEventLoop: starting")
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.shutdowns == nil {
		b.shutdowns = make(chan struct{})
		b.wg.Add(1)
		go b.runEventLoop()
	}
}

func (b *eventBus) shutdownEventLoop() {
	log.Println("Bus, EventHandler, startupEventLoop: shutting down")
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.shutdowns != nil {
		close(b.shutdowns)
		b.wg.Wait()
		b.shutdowns = nil
	}
}

func (b *eventBus) runEventLoop() {
	for {
		err := b.doRunEventLoop()
		if err != nil {
			log.Println(fmt.Sprintf("Bus EventHandler, runEventLoop: received error; %s", err.Error()))
		}
		select {
		case _, ok := <-b.shutdowns:
			if !ok {
				goto CloseEventLoop
			}
		case <-time.After(5 * time.Second):
			log.Println("Bus EventHandler, runEventLoop: retrying loop")
		}
	}
CloseEventLoop:
	log.Println("Bus EventHandler, runEventLoop: closed")
	b.wg.Done()
}

func (b *eventBus) doRunEventLoop() error {
	// create event queue
	func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		if b.eventQueue != nil {
			b.eventQueue = nil
		}
		// shouldn't have too many events, 20 seems like a sane default
		b.eventQueue = make(chan domain.Event, 20)
	}()
	for {
		select {
		case _, ok := <-b.shutdowns:
			if !ok {
				return nil
			}
		case e, ok := <-b.eventQueue:
			if !ok {
				return errors.New("event queue unexpectedly closed")
			} else {
				err := b.doHandleEvent(e)
				if err != nil {
					return errors.New(fmt.Sprintf(
						"error handling event: %s; recieved error: %s",
						e.String(),
						err.Error(),
					))
				}
			}
		}
	}
}

func (b *eventBus) doHandleEvent(e domain.Event) (err error) {
	b.mu.RLock()
	defer func() {
		// recover from panic if one occurred. Set err to nil otherwise.
		if recover() != nil {
			err = errors.New("array index out of bounds")
		}
		b.mu.RUnlock()
	}()
	if _, ok := b.pubSubMap[e]; !ok {
		log.Println(fmt.Sprintf("Bus EventHandler, doHandleEvent: no handlers for event %s", e.String()))
		return nil
	}
	callback := b.pubSubMap[e]
	for callback != nil {
		if callback.id != -1 {
			callback.fn()
		}
		callback = callback.next
	}
	return nil
}
