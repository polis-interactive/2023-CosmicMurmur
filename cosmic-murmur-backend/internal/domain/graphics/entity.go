package graphics

import (
	"errors"
	"fmt"
	"github.com/polis-interactive/2023-CosmicMurmur/internal/types"
	"github.com/polis-interactive/go-lighting-utils/pkg/graphicsShader"
	"io/ioutil"
	"log"
	"path"
	"strings"
	"sync"
	"time"
)

type Graphics struct {
	s  *service
	mu *sync.RWMutex

	shaderPath            string
	pixelSize             int
	defaultReloadOnUpdate bool
	defaultShader         string
	defaultFrequency      time.Duration

	shaderList            graphicsShader.ShaderIdentifiers
	pb                    *types.PixelBuffer
	ud                    graphicsShader.UniformDict
	gs                    *graphicsShader.GraphicsShader
	runningReloadOnUpdate bool
	runningShader         string
	runningFrequency      time.Duration
	lastTimeStep          time.Time
}

func newGraphics(s *service, cfg Config) (*Graphics, error) {
	programName := cfg.GetProgramName()
	shaderPath, err := graphicsShader.GetShaderPathIfAvailable(programName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(
			"Graphics Service, newGraphics: couldn't find shader path in program %s", programName,
		))
	}

	return &Graphics{
		s:  s,
		mu: &sync.RWMutex{},

		shaderPath:            shaderPath,
		pixelSize:             cfg.GetGraphicsPixelSize(),
		defaultReloadOnUpdate: cfg.GetGraphicsReloadOnUpdate(),
		defaultShader:         cfg.GetGraphicsDefaultShader(),
		defaultFrequency:      cfg.GetGraphicsFrequency(),

		shaderList:            nil,
		pb:                    nil,
		ud:                    nil,
		gs:                    nil,
		runningReloadOnUpdate: false,
		runningShader:         "",
		runningFrequency:      time.Minute,
	}, nil
}

func (g *Graphics) runMainLoop() {
	for {
		err := g.runGraphicsLoop()
		if err != nil {
			log.Println(fmt.Sprintf("Graphics, runMainLoop: received error; %s", err.Error()))
			g.s.bus.EmitGraphicsCrashed()
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
	defer g.cleanupGraphicsLoop()
	err := g.setupGraphicsLoop()
	if err != nil {
		return err
	}
	for {
		dur := g.getRunningFrequency()
		select {
		case _, ok := <-g.s.shutdowns:
			if !ok {
				return nil
			}
		case <-time.After(dur):
			g.stepTime()
			err = g.tryReloadShader()
			if err != nil {
				return err
			}
			err = g.doRunShader()
			if err != nil {
				return err
			}
			err = g.gs.ReadToPixels(g.pb.GetUnsafePointer())
			if err != nil {
				return err
			}
			g.s.bus.EmitGraphicsReady()
		}
	}
}

func (g *Graphics) setupGraphicsLoop() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	shaders, err := g.getShaders()
	if err != nil {
		return err
	} else if len(shaders) == 0 {
		return errors.New(fmt.Sprintf("No shaders found in path %s", g.shaderPath))
	}
	g.shaderList = shaders

	grid := g.s.bus.GetGridDimensions()
	gridWidth := grid.MaxX - grid.MinX + 1
	gridHeight := grid.MaxY - grid.MinY + 1

	gridWidth = gridWidth * g.pixelSize
	gridHeight = gridHeight * g.pixelSize
	g.pb = types.NewPixelBuffer(gridWidth, gridHeight, grid.MinX, grid.MinY, g.pixelSize)

	g.ud = make(graphicsShader.UniformDict)
	g.lastTimeStep = time.Now()
	g.ud["time"] = 0.0
	g.ud["pixel"] = float32(g.pixelSize)

	gs, err := graphicsShader.NewGraphicsShader(
		g.shaderPath, int32(gridWidth), int32(gridHeight), g.ud, g.mu,
	)
	if err != nil {
		return err
	}
	g.gs = gs
	err = g.gs.AttachShaders(g.shaderList)
	if err != nil {
		return err
	}
	err = g.initializeVariables()
	if err != nil {
		return err
	}
	return nil
}

func (g *Graphics) getShaders() (graphicsShader.ShaderIdentifiers, error) {
	files, err := ioutil.ReadDir(g.shaderPath)
	if err != nil {
		return nil, err
	}
	shaderNameList := make(graphicsShader.ShaderIdentifiers)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		ext := path.Ext(f.Name())
		if ext != ".frag" {
			continue
		}
		fName := strings.TrimSuffix(f.Name(), ext)
		shaderNameList[graphicsShader.ShaderKey(fName)] = fName
	}
	return shaderNameList, nil
}

func (g *Graphics) initializeVariables() error {
	var ok bool
	var reloadOnUpdate bool
	reloadOnUpdate, ok = g.s.repo.GetGraphicsReloadOnUpdate()
	if !ok {
		log.Println("Graphics, initializeVariables: no update preference found, using default")
		reloadOnUpdate = g.defaultReloadOnUpdate
	}
	g.runningReloadOnUpdate = reloadOnUpdate
	var shaderName string
	shaderName, ok = g.s.repo.GetGraphicsShader()
	if !ok {
		log.Println("Graphics, initializeVariables: no shader saved, using default")
		shaderName = g.defaultShader
	}
	if _, ok = g.shaderList[graphicsShader.ShaderKey(shaderName)]; !ok {
		var firstShaderFound string
		for _, v := range g.shaderList {
			firstShaderFound = v
			break
		}
		log.Println(fmt.Sprintf("Couldn't find default shader %s; using %s", shaderName, firstShaderFound))
		shaderName = firstShaderFound
	}
	g.runningShader = shaderName
	err := g.gs.SetShader(graphicsShader.ShaderKey(g.runningShader))
	if err != nil {
		return err
	}
	var frequency time.Duration
	frequency, ok = g.s.repo.GetGraphicsFrequency()
	if !ok {
		log.Println("Graphics, initializeVariables: no frequency, using default")
		frequency = g.defaultFrequency
	}
	g.runningFrequency = frequency
	return nil
}

func (g *Graphics) getRunningFrequency() time.Duration {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.runningFrequency
}

func (g *Graphics) tryReloadShader() error {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if g.runningReloadOnUpdate {
		return g.gs.ReloadShader()
	}
	return nil
}

func (g *Graphics) doRunShader() error {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.gs.RunShader()
}

func (g *Graphics) stepTime() {
	g.mu.Lock()
	defer g.mu.Unlock()
	nt := time.Now()
	timeMultiplier := 1.0
	elapsed := nt.Sub(g.lastTimeStep).Seconds() * timeMultiplier
	g.ud["time"] += float32(elapsed)
	g.lastTimeStep = nt
}

func (g *Graphics) setShader(shader graphicsShader.ShaderKey) error {
	g.runningShader = string(shader)
	if g.gs != nil {
		return g.gs.SetShader(shader)
	}
	return errors.New("setShaderError")
}

func (g *Graphics) cleanupGraphicsLoop() {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.gs != nil {
		g.gs.Cleanup()
		g.gs = nil
	}
	if g.pb != nil {
		g.pb = nil
	}
	if g.shaderList != nil {
		g.shaderList = nil
	}
	if g.ud != nil {
		g.ud = nil
	}
}
