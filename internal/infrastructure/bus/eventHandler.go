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

type eventHandler struct {
	lastEventId int
	pubSubMap   map[domain.Event]*callbackProxy
	idEventMap  map[int]domain.Event
	eventQueue  chan domain.Event
	mu          *sync.RWMutex
	wg          *sync.WaitGroup
	shutdowns   chan struct{}
}

var _ domain.EventHandler = (*eventHandler)(nil)

func newEventHandler() *eventHandler {
	return &eventHandler{
		lastEventId: 0,
		pubSubMap:   make(map[domain.Event]*callbackProxy),
		idEventMap:  make(map[int]domain.Event),
		eventQueue:  nil,
		mu:          &sync.RWMutex{},
		wg:          &sync.WaitGroup{},
		shutdowns:   nil,
	}
}

func (h *eventHandler) SubscribeToEvent(e domain.Event, fn func()) int {
	h.mu.Lock()
	defer h.mu.Unlock()
	log.Println(fmt.Sprintf(
		"Bus, EventHandler, SubscribeToEvent: subscribing to %s, giving id %d",
		e.String(),
		h.lastEventId,
	))
	newCallback := &callbackProxy{
		id:   h.lastEventId,
		fn:   fn,
		next: nil,
	}
	h.idEventMap[h.lastEventId] = e
	if callback, ok := h.pubSubMap[e]; ok {
		for callback.next != nil {
			callback = callback.next
		}
		callback.next = newCallback
	} else {
		h.pubSubMap[e] = &callbackProxy{
			id:   -1,
			fn:   nil,
			next: newCallback,
		}
	}
	h.lastEventId += 1
	return h.lastEventId
}

func (h *eventHandler) UnsubscribeToEvent(handlerId int) {
	log.Println(fmt.Sprintf("Bus, EventHandler, UnsubscribeToEvent: unsubscribing id %d", handlerId))
	h.mu.Lock()
	defer h.mu.Unlock()
	// invalid handlerId
	if handlerId > h.lastEventId || handlerId < 0 {
		return
	}
	// short circuit if we've already removed the id
	if _, ok := h.idEventMap[handlerId]; !ok {
		return
	}
	eventType := h.idEventMap[handlerId]
	delete(h.idEventMap, handlerId)

	// short circuit if the pubSub map for event type
	// is missing
	if _, ok := h.pubSubMap[eventType]; !ok {
		return
	}
	callback := h.pubSubMap[eventType]
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
		delete(h.pubSubMap, eventType)
	}
}

func (h *eventHandler) HandleEvent(e domain.Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.eventQueue != nil {
		h.eventQueue <- e
	}
}

func (h *eventHandler) startupEventLoop() {
	log.Println("Bus, EventHandler, startupEventLoop: starting")
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.shutdowns == nil {
		h.shutdowns = make(chan struct{})
		h.wg.Add(1)
		go h.runEventLoop()
	}
}

func (h *eventHandler) shutdownEventLoop() {
	log.Println("Bus, EventHandler, startupEventLoop: shutting down")
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.shutdowns != nil {
		close(h.shutdowns)
		h.wg.Wait()
		h.shutdowns = nil
	}
}

func (h *eventHandler) runEventLoop() {
	for {
		err := h.doRunEventLoop()
		if err != nil {
			log.Println(fmt.Sprintf("Bus EventHandler, runEventLoop: received error; %s", err.Error()))
		}
		select {
		case _, ok := <-h.shutdowns:
			if !ok {
				goto CloseEventLoop
			}
		case <-time.After(5 * time.Second):
			log.Println("Bus EventHandler, runEventLoop: retrying loop")
		}
	}
CloseEventLoop:
	log.Println("Bus EventHandler, runEventLoop: closed")
	h.wg.Done()
}

func (h *eventHandler) doRunEventLoop() error {
	// create event queue
	func() {
		h.mu.Lock()
		defer h.mu.Unlock()
		if h.eventQueue != nil {
			h.eventQueue = nil
		}
		// shouldn't have too many events, 20 seems like a sane default
		h.eventQueue = make(chan domain.Event, 20)
	}()
	for {
		select {
		case _, ok := <-h.shutdowns:
			if !ok {
				return nil
			}
		case e, ok := <-h.eventQueue:
			if !ok {
				return errors.New("event queue unexpectedly closed")
			} else {
				err := h.doHandleEvent(e)
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

func (h *eventHandler) doHandleEvent(e domain.Event) (err error) {
	h.mu.RLock()
	defer func() {
		// recover from panic if one occurred. Set err to nil otherwise.
		if recover() != nil {
			err = errors.New("array index out of bounds")
		}
		h.mu.RUnlock()
	}()
	if _, ok := h.pubSubMap[e]; !ok {
		log.Println(fmt.Sprintf("Bus EventHandler, doHandleEvent: no handlers for event %s", e.String()))
		return nil
	}
	callback := h.pubSubMap[e]
	for callback != nil {
		if callback.id != -1 {
			callback.fn()
		}
		callback = callback.next
	}
	return nil
}
