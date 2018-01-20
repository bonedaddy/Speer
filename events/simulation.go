package events

import (
  "fmt"
)

type Simulation struct {
  stopped chan interface {}
  time    int
  EventQueue
}

func NewLazySimulation() (s Simulation) {
  s = Simulation{
    make(chan interface {}),
    0,
    NewLazyEventQueue(),
  }
  return
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
        s.time = event.timestamp
        // event.receiver.Receive(event.payload)
      }
    }
  }
}