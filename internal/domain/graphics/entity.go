package graphics

import (
	"fmt"
	"log"
	"time"
)

type Graphics struct {
	s *service
}

func newGraphics(s *service) (*Graphics, error) {
	return &Graphics{
		s: s,
	}, nil
}

func (g *Graphics) runMainLoop() {
	for {
		err := g.runGraphicsLoop()
		if err != nil {
			log.Println(fmt.Sprintf("Graphics, runMainLoop: received error; %s", err.Error()))
		}
		select {
		case _, ok := <-g.s.shutdowns:
			if !ok {
				goto CloseMainLoop
			}
		case <-time.After(5 * time.Second):
			log.Println("Graphics, Main Loop: retrying window")
		}
	}
CloseMainLoop:
	log.Println("Graphics runMainLoop, Main Loop: closed")
	g.s.wg.Done()
}

func (g *Graphics) runGraphicsLoop() error {
	return nil
}
