package sdk

import (
	"github.com/danalex97/Speer/interfaces"

	"github.com/danalex97/Speer/overlay"
	"github.com/danalex97/Speer/underlay"
	"github.com/danalex97/Speer/events"

	"github.com/danalex97/Speer/logs"

	"fmt"
)

type ISimulation interface {
	interfaces.ISimulation
}

type Simulation struct {
	underlaySimulation *underlay.NetworkSimulation

	latencyMap overlay.LatencyMap
	capacityMap overlay.CapacityMap
	nodes int
	cnode int

	userNodes map[string]SpeerNode

	template interfaces.Node

	progressProperties []events.Receiver
}

type SimulationBuilder struct {
	*Simulation
}

func NewSimulationBuilder(template interfaces.Node) *SimulationBuilder {
	b := &SimulationBuilder{
		Simulation : new(Simulation),
	}

	if template == nil {
		panic("No valid template provided")
	}

	b.template = template
	b.progressProperties = []events.Receiver{}
	b.userNodes = map[string]SpeerNode{}
	b.nodes = -1

	return b
}

func (b *SimulationBuilder) WithParallelSimulation() *SimulationBuilder {
	if b.underlaySimulation == nil {
		panic("Underlay simulation component has to be appended first.")
	}

	b.underlaySimulation.SetParallel(true)
	return b
}

func (b *SimulationBuilder) WithInternetworkUnderlay(
	transitDomains int,
	transitDomainSize int,
	stubDomains int,
	stubDomainSize int,
) *SimulationBuilder {
	network := underlay.NewInternetwork(
		transitDomains,
		transitDomainSize,
		stubDomains,
		stubDomainSize,
	)
	simulation := underlay.NewNetworkSimulation(
		events.NewLazySimulation(),
		network,
	)

	fmt.Printf("Internetwork built with %d nodes.\n", len(network.Routers))
	b.underlaySimulation = simulation
	b.latencyMap = overlay.NewNetworkMap(b.underlaySimulation.Network())

	return b
}

func (b *SimulationBuilder) WithRandomUniformUnderlay(
	nodes int,
	edges int,
	minLatency int,
	maxLatency int,
) *SimulationBuilder {
	network := underlay.NewRandomUniformNetwork(
		nodes,
		edges,
		minLatency,
		maxLatency,
	)
	s := underlay.NewNetworkSimulation(
		events.NewLazySimulation(),
		network,
	)

	b.underlaySimulation = s
	b.latencyMap = overlay.NewNetworkMap(b.underlaySimulation.Network())

	return b
}

func (b *SimulationBuilder) WithProgress(
	progress interfaces.Progress,
	interval int,
) *SimulationBuilder {
	property := events.NewProgressProperty(progress, interval)
	b.progressProperties = append(b.progressProperties, property)

	return b
}

func (b *SimulationBuilder) WithFixedNodes(
	nodes int,
) *SimulationBuilder {
	b.nodes = nodes
	b.cnode = 0
	return b
}

func (b *SimulationBuilder) WithCapacityScheduler(
	interval int,
) *SimulationBuilder {
	b.capacityMap = overlay.NewScheduledCapacityMap(interval)
	return b
}

func (b *SimulationBuilder) WithCapacityNodes(
	nodes int,
	upload int,
	download int,
) *SimulationBuilder {
	if b.nodes == -1 {
		panic("Node number not specified.")
	}
	if b.capacityMap == nil {
		panic("No capacity scheduler provided.")
	}
	limit := b.cnode + nodes
	if b.nodes < limit {
		limit = b.nodes
	}
	for i := b.cnode; i < limit; i++ {
		// assign ID to router
		id := b.latencyMap.NewId()

		// create latency connector
		latencyConnector := overlay.NewUnderlayChan(
			id,
			b.underlaySimulation,
			b.latencyMap,
		)

		// register capacity
		capacityConnector := overlay.NewCapacityConnector(
			upload,
			download,
			b.capacityMap,
		)
		b.capacityMap.AddConnector(id, capacityConnector)

		// register autowired nodes
		newNode := NewAutowiredNode(b.template, overlay.NewSimulatedNode(
			latencyConnector,
			capacityConnector,
			b.latencyMap,
			id,
			b.Time,
		))
		b.userNodes[id] = newNode
	}
	b.cnode = limit
	return b
}

func (b *SimulationBuilder) WithLogs(logsFile string) *SimulationBuilder {
	globalObserver := events.NewGlobalEventObserver()
	b.underlaySimulation.RegisterObserver(globalObserver)

	logger := logs.NewEventMonitor(globalObserver, b.latencyMap, logsFile)

	go logger.GatherEvents()

	return b
}

func (b *SimulationBuilder) Build() ISimulation {
	if (b.underlaySimulation == nil) {
		panic("No underlay simulation provided.")
	}

	return b.Simulation
}

func (s *Simulation) Run() {
	for _, progress := range s.progressProperties {
		event := events.NewEvent(0, nil, progress)
		s.underlaySimulation.Push(event)
	}

	s.capacityMap.Start(s.underlaySimulation)

	go s.underlaySimulation.Run()
	for _, node := range s.userNodes {
		go node.OnJoin()
	}
}

func (s *Simulation) Stop() {
	s.underlaySimulation.Stop()
}

func (s *Simulation) Time() int {
	return s.underlaySimulation.Time();
}
