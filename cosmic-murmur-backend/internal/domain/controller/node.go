package controller

import (
	"errors"
	"fmt"
	"github.com/jsimonetti/go-artnet"
	"github.com/jsimonetti/go-artnet/packet"
	"github.com/polis-interactive/2023-CosmicMurmur/internal/types"
	"log"
	"net"
	"sync"
	"time"
)

type node struct {
	c               *controller
	address         string
	universeNumbers []int
	universePackets map[int]*packet.ArtDMXPacket

	shutdowns chan struct{}
	wg        *sync.WaitGroup
	mu        *sync.RWMutex

	pollChan chan struct{}
	sendChan chan int
	conn     net.Conn
}

func newNode(c *controller, definition types.NodeDefinition) *node {
	n := &node{
		c:               c,
		address:         definition.Address,
		universeNumbers: definition.Universes,
		universePackets: make(map[int]*packet.ArtDMXPacket),

		shutdowns: nil,
		wg:        &sync.WaitGroup{},
		mu:        &sync.RWMutex{},

		pollChan: nil,
		sendChan: nil,
		conn:     nil,
	}
	for _, u := range definition.Universes {
		address := artnet.Address{
			Net:    uint8(u >> 8 & 0xFF),
			SubUni: uint8(u & 0xFF),
		}
		artNetPacket := &packet.ArtDMXPacket{
			Sequence: 0,
			SubUni:   address.SubUni,
			Net:      address.Net,
			Length:   0,
			Data:     [512]byte{},
		}
		n.universePackets[u] = artNetPacket
		// add node to controller maps
		c.universeBufferMap[u] = &artNetPacket.Data
		c.universeNodeMap[u] = n
	}
	return n
}

func (n *node) startup() {
	if n.shutdowns == nil {
		n.shutdowns = make(chan struct{})
		n.wg.Add(1)
		go n.runMainLoop()
	}
}

func (n *node) shutdown() {
	if n.shutdowns != nil {
		close(n.shutdowns)
		// if n.conn is trying to write, force it to close
		n.cleanupNodeLoop()
		n.wg.Wait()
		n.shutdowns = nil
	}
}

func (n *node) reset() {
	n.shutdown()
	n.startup()
}

func (n *node) getDefinition() types.NodeDefinition {
	return types.NodeDefinition{
		Address:   n.address,
		Universes: n.universeNumbers,
	}
}

func (n *node) runMainLoop() {
	defer func() {
		log.Println("controller node, runMainLoop, Main Loop: closed")
		n.wg.Done()
	}()
	for {
		err := n.runSendLoop()
		if err != nil {
			log.Println(fmt.Sprintf("controller, node, runMainLoop: received error; %s", err.Error()))
		}
		select {
		case _, ok := <-n.shutdowns:
			if !ok {
				return
			}
		case <-time.After(1 * time.Second):
			log.Println("controller, node, runMainLoop: retrying sender")
		}
	}
}

func (n *node) runSendLoop() error {
	defer n.cleanupNodeLoop()
	err := n.setupNodeLoop()
	if err != nil {
		return err
	}
	for {
		select {
		case _, ok := <-n.shutdowns:
			if !ok {
				return nil
			}
		case _, ok := <-n.pollChan:
			if !ok {
				return errors.New("poll chan unexpectedly closed")
			}
			// handle poll
		case u, ok := <-n.sendChan:
			if !ok {
				return errors.New("send chan unexpectedly closed")
			}
			err = n.sendUniverseUpdate(u)
			if err != nil {
				return errors.New(fmt.Sprintf("couldn't send full universe update %d", u))
			}
		}
	}
}

func (n *node) sendUniverseUpdate(u int) error {
	// universe guaranteed to be on node because it's coordinated by service
	b, err := n.universePackets[u].MarshalBinary()
	if err != nil {
		return err
	}
	return n.doSendPacket(b)
}

func (n *node) doSendPacket(packet []byte) error {
	for len(packet) > 0 {
		b, err := n.conn.Write(packet)
		if err != nil {
			return err
		} else if b == 0 {
			return errors.New("couldn't write push payload")
		}
		packet = packet[b:]
	}
	return nil
}

func (n *node) setupNodeLoop() error {
	n.mu.Lock()
	defer n.mu.Unlock()
	addr := fmt.Sprintf("%s:%d", n.address, packet.ArtNetPort)
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return err
	}
	n.conn = conn
	n.pollChan = make(chan struct{}, 5)
	n.sendChan = make(chan int, 10)
	return nil
}

func (n *node) cleanupNodeLoop() {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.conn != nil {
		_ = n.conn.Close()
	}
	n.pollChan = nil
	n.sendChan = nil
}
