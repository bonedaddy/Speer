package events

import (
  "fmt"
)

type Simulation struct {
  observers []EventObserver
  stopped chan interface {}
  time    int
  EventQueue
}

func NewLazySimulation() (s Simulation) {
  s = Simulation{
    make([]EventObserver, 0),
    make(chan interface {}),
    0,
    NewLazyEventQueue(),
  }
  return
}

func (s *Simulation) RegisterObserver(eventObserver EventObserver) {
  s.observers = append(s.observers, eventObserver)
}

func (s *Simulation) Time() int {
  return s.time
}

func (s *Simulation) Stop() {
  s.stopped <- nil
}

func (s *Simulation) Run() {
  s.time = 0
  for {
    select {
    case <-s.stopped:
      break
    default:
      if event:= s.Pop(); event != nil {
        fmt.Println("Event received at time:", event.timestamp)

        // The event gets dispached to observers
        for _, observer := range(s.observers) {
          observer.EnqueEvent(event)
        }

        s.time = event.timestamp
        receiver := event.receiver

        if receiver == nil {
          continue
        }

        newEvent := receiver.Receive(event)

        if newEvent != nil {
          s.Push(newEvent)
        }
      }
    }
  }
}
