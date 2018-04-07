package model

import (
  . "github.com/danalex97/Speer/overlay"
  "github.com/danalex97/Speer/interfaces"
  "math/rand"
)

type Query interfaces.Query

type DHTQueryGenerator interface {
  Next() Query
  // the key for a store query is empty, thus
  // allowing the SDK layer to handle it
}

const MaxQuerySize int = 100

type DHTQuery struct {
  key   string // the key of the node
  size  int    // size of key to be transfered in MB
  node  string // the node which sends/stores the query
  store bool   // store/retrieve
}

func NewDHTQuery(key string, size int, node string, store bool) Query {
  q := new(DHTQuery)
  q.key = key
  q.size = size
  q.node = node
  q.store = store
  return q
}

func (q *DHTQuery) Key() string {
  return q.key
}

func (q *DHTQuery) Size() int {
  return q.size
}

func (q *DHTQuery) Node() string {
  return q.node
}

func (q *DHTQuery) Store() bool {
  return q.store
}

type DHTLedger struct {
  queries   []Query
  bootstrap Bootstrap
}

func NewDHTLedger(bootstrap Bootstrap) *DHTLedger {
  ledger := new(DHTLedger)
  ledger.queries = []Query{}
  ledger.bootstrap = bootstrap
  return ledger
}

func randomKey() string {
  const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

  b := make([]byte, 30)
  for i := range b {
    b[i] = letterBytes[rand.Int63() % int64(len(letterBytes))]
  }
  return string(b)
}


func (l *DHTLedger) Next() Query {
  node := l.bootstrap.Join("")
  size := rand.Intn(MaxQuerySize)
  store := len(l.queries) == 0 || rand.Float32() > 0.5
  key   := randomKey()

  if !store {
    // this is generated uniformly as there are no leaves yet
    // and the history has only 'store' queries
    idx := rand.Intn(len(l.queries))
    key  = l.queries[idx].Key()
  }

  query := NewDHTQuery(key, size, node, store)
  if store {
    l.queries = append(l.queries, query)
  }

  return query
}
