package examples

import (
	"fmt"
	. "github.com/danalex97/Speer/interfaces"
	"github.com/danalex97/Speer/structs"
	"runtime"
	"sync"
)

type SimpleTorrent struct {
	sync.Mutex

	id         string
	ids        []string
	links      map[string]Link
	transport  Transport
}

type idBroadcast struct {
	ids []string
}

/* Interface functions. */
func (s *SimpleTorrent) OnJoin() {
	if s.transport == nil {
		return
	}

	go func() {
		for {
			// check links
			for _, l := range s.links {
				select {
				case m, _ := <-l.Download():
					fmt.Println(s.id, "data", m)
				default:
					continue
				}
			}

			// check for control messages
			select {

			case m, ok := <-s.transport.ControlRecv():
				if !ok {
					continue
				}

				switch msg := m.(type) {
				case idBroadcast:
					s.updateIds(msg.ids)
					fmt.Println(s.id, "received", msg.ids)
				}

			default:
				runtime.Gosched()
			}
		}
	}()

	go func() {
		// broadcast neighbours
		for _, id := range s.ids {
			if id != s.id {
				if !s.transport.ControlPing(id) {
					continue
				}

				s.transport.ControlSend(id, idBroadcast{s.ids})
			}
		}

		fmt.Println("Done", s.id)
	}()
}

func (s *SimpleTorrent) OnLeave() {
}

func (s *SimpleTorrent) New(util NodeUtil) Node {
	// Constructor that assumes the UnreliableNode component is filled in
	node := new(SimpleTorrent)

	node.id = util.Id()
	node.ids = []string{node.id, util.Join()}
	node.transport = util.Transport()
	node.links = map[string]Link{}

	return node
}

/* Local functions */
func (s *SimpleTorrent) updateIds(ids []string) {
	allIds := make(map[string]bool)
	for _, id := range ids {
		allIds[id] = true
	}
	for _, id := range s.ids {
		allIds[id] = true
	}

	s.ids = []string{}
	for id, _ := range allIds {
		s.ids = append(s.ids, id)

		if id == s.id {
			continue
		}

		// register link if not registered
		if _, ok := s.links[id]; !ok {
			s.links[id] = s.transport.Connect(id)

			// if the link is new, we broadcast our list again
			if !s.transport.ControlPing(id) {
				continue
			}
			s.transport.ControlSend(id, idBroadcast{s.ids})

			// send a big packet
			s.links[id].Upload(Data{structs.RandomKey(), 1000})
		}
	}
}
