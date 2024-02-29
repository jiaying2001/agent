package launcher

import (
	"github.com/jiaying2001/agent/config"
	"github.com/jiaying2001/agent/harvester"
	"github.com/jiaying2001/agent/log"
	"github.com/jiaying2001/agent/parser"
)

var L = Launcher{
	make(map[string]*harvester.Harvester),
	make(map[string]parser.Parser),
}

type Launcher struct {
	Workers map[string]*harvester.Harvester
	Parsers map[string]parser.Parser
}

func (l *Launcher) Launch() {
	// Register parsers
	l.Parsers["RFC5424"] = &parser.RFC5424Parser{}
	// Get configs
	configs := config.GetConfigs()
	// Create workers
	l.Load(configs)
	// Run
	l.StartWorkers()
}

func (l *Launcher) Load(configs *[]harvester.Harvester) {
	l.interruptAll() // Remove all the configs
	// Plug configs into the launcher
	for _, hc := range *configs {
		hc_, exist := l.Workers[hc.Path]
		if exist {
			log.Logger.Info("Harvester already exists for the path: " + hc_.Path)
		} else {
			log.Logger.Info("Created a harvester for the path " + hc.Path)
			l.Workers[hc.Path] = harvester.CreateHarvester(hc.Path, hc.FileFormat)
		}
	}
}

func (l *Launcher) interruptAll() {
	log.Logger.Info("Started interrupting all harvesters")
	for _, hc := range l.Workers {
		hc.Interrupt = true
		hc.Shutdown.Wait()
	}
	l.Workers = make(map[string]*harvester.Harvester)
	log.Logger.Info("Interrupted all harvesters")
}

func (l *Launcher) StartWorkers() {
	for _, hc := range l.Workers {
		switch hc.State {
		case harvester.Created:
			hc.State = harvester.Running
			log.Logger.Info("Started a harvester for the path " + hc.Path)
			go func(hc *harvester.Harvester) {
				defer hc.Shutdown.Done() // Synchronize with main thread
				hc.Shutdown.Add(1)
				for !hc.Interrupt { // Interrupt when sets true. Note that only main thread write to it
					if _, ok := l.Parsers[hc.FileFormat]; ok {
						hc.Run(l.Parsers[hc.FileFormat]) // Run the main logic
					} else {
						log.Logger.Error("Error reading logs from " + hc.Path + ": File format " + hc.FileFormat + " is not supported")
						break
					}
				}
				hc.HandleInterrupt() // Handle interrupt
			}(hc)
		}
	}
}
