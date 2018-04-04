package capacity

import (
  "testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

type node struct {
  up   int
  down int
}

type link struct {
  from int
  to   int
}

func (n *node) Up() int {
  return n.up
}

func (n *node) Down() int {
  return n.down
}

func buildGraph(nodes []node, idxs []link) (*scheduler, []node, []Link) {
  s := NewScheduler(0).(*scheduler)

  links := []Link{}
  for _, l := range idxs {
    link := NewPerfectLink(&nodes[l.from], &nodes[l.to])
    links = append(links, link)
    s.RegisterLink(link)
  }
  for _, status := range s.linkStatus {
    status.active = true
  }

  return s, nodes, links
}

func checkCapacity(t *testing.T, s *scheduler, link Link, cap float64) {
  if s.linkStatus[link] == nil {
    t.Fatalf("Link not found!")
  }
  assertEqual(t, s.linkStatus[link].capacity, cap)
}

func TestUpdCapacityTwoNodes(t *testing.T) {
  s, _, links := buildGraph([]node{
    node{10, 10},
    node{10, 10},
  }, []link{
    link{0, 1},
  })

  s.updCapacity()

  checkCapacity(t, s, links[0], 10)
}
