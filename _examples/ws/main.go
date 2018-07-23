package main

import (
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/m0t0k1ch1/nebula/graph"
	"github.com/m0t0k1ch1/nebula/utils"
)

const (
	n        = 100
	kAvg     = 4 // must be even
	p        = 0.1
	filePath = "./ws.dot"
)

func newID(id int) graph.ID {
	return graph.StringID(strconv.Itoa(id))
}

func newRandomID(excludes map[graph.ID]bool) (id graph.ID) {
	ok := false
	for !ok {
		id = newID(rand.Intn(n))
		if _, ok := excludes[id]; ok {
			continue
		}
		ok = true
	}
	return id
}

func pickRandomID(ends map[graph.ID]graph.Node) (id graph.ID) {
	i, target := 0, rand.Intn(len(ends))
	for id, _ = range ends {
		if i == target {
			break
		}
		i++
	}
	return
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	g, err := createGraph()
	if err != nil {
		panic(err)
	}

	if err := writeGraph(g); err != nil {
		panic(err)
	}
}

func createGraph() (graph.Graph, error) {
	g := graph.NewUndirected()

	// add nodes
	for i := 0; i < n; i++ {
		g.AddNode(graph.NewNode(strconv.Itoa(i)))
	}

	// add edges
	for i := 0; i < n; i++ {
		for j := i + 1; j <= i+kAvg/2; j++ {
			idTail, idHead := i, j
			if idHead >= n {
				idHead -= n
			}

			if err := g.AddEdge(newID(idTail), newID(idHead), 1.0); err != nil {
				return nil, err
			}
		}
	}

	targetsNum := int(math.Floor(float64(n) * float64(kAvg/2) * p))
	targets := map[graph.ID]map[graph.ID]bool{}
	cnt := 0

	// pick target edges
	for cnt < targetsNum {
		idTail := newRandomID(nil)

		heads, err := g.GetHeads(idTail)
		if err != nil {
			return nil, err
		}
		idHead := pickRandomID(heads)

		if _, ok := targets[idTail]; ok {
			if _, ok := targets[idTail][idHead]; ok {
				continue
			}
		}
		if _, ok := targets[idHead]; ok {
			if _, ok := targets[idHead][idTail]; ok {
				continue
			}
		}

		if _, ok := targets[idTail]; ok {
			targets[idTail][idHead] = true
		} else {
			targets[idTail] = map[graph.ID]bool{idHead: true}
		}

		cnt++
	}

	// switch edges
	for idTail, targetHeads := range targets {
		for idHeadOld, _ := range targetHeads {
			heads, err := g.GetHeads(idTail)
			if err != nil {
				return nil, err
			}
			excludes := map[graph.ID]bool{idTail: true}
			for _, node := range heads {
				excludes[node.ID()] = true
			}

			if err := g.RemoveEdge(idTail, idHeadOld); err != nil {
				return nil, err
			}

			idHeadNew := newRandomID(excludes)
			if err := g.AddEdge(idTail, idHeadNew, 1.0); err != nil {
				return nil, err
			}
		}
	}

	return g, nil
}

func writeGraph(g graph.Graph) error {
	dg, err := utils.NewDOTGraph(g)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write([]byte(dg.String())); err != nil {
		return err
	}

	return nil
}
