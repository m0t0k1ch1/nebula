package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/m0t0k1ch1/nebula/graph"
	"github.com/m0t0k1ch1/nebula/utils"
)

const (
	n        = 100
	m0       = 3
	m        = 2 // must be m0 or less
	filePath = "./ba.dot"
)

func newID(id int) graph.ID {
	return graph.ID(strconv.Itoa(id))
}

func pickID(g *graph.Graph) (graph.ID, error) {
	edges := g.GetEdges()

	kSum := 0
	ids := make([]graph.ID, 0, len(edges))
	for id, nodeEdges := range edges {
		kSum += len(nodeEdges)
		ids = append(ids, id)
	}

	weight := rand.Intn(kSum) + 1

	weightSum := 0
	for _, id := range ids {
		weightSum += len(edges[id])
		if weightSum >= weight {
			return id, nil
		}
	}

	return "", fmt.Errorf("failed to pick")
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	g, err := createGraph()
	if err != nil {
		panic(err)
	}

	writeGraphFeatures(g)

	if err := writeGraph(g); err != nil {
		panic(err)
	}
}

func createGraph() (*graph.Graph, error) {
	g := graph.NewUndirected()

	nodesNum := 0

	// add default nodes
	for i := 0; i < m0; i++ {
		g.AddNode(graph.NewNode(strconv.Itoa(nodesNum)))
		nodesNum++
	}

	// add default edges
	for i := 0; i < m0; i++ {
		idTail, idHead := i, i+1
		if idHead >= m0 {
			idHead -= m0
		}

		if err := g.AddEdge(newID(idTail), newID(idHead), 1.0); err != nil {
			return nil, err
		}
	}

	for nodesNum < n {
		nodesNum++
		node := graph.NewNode(strconv.Itoa(nodesNum))

		// add node
		g.AddNode(node)

		// pick target nodes
		picked := map[graph.ID]bool{}
		for len(picked) < m {
			id, err := pickID(g)
			if err != nil {
				return nil, err
			}
			picked[id] = true
		}

		// add edges
		for id, _ := range picked {
			if err := g.AddEdge(node.ID(), id, 1.0); err != nil {
				return nil, err
			}
		}
	}

	return g, nil
}

func writeGraphFeatures(g *graph.Graph) {
	kDist := g.GetIndegreeDistribution()
	sort.Sort(kDist)

	fmt.Println("kAvg:", kDist.CalcAverageDegree())
	fmt.Println("kDist:")
	ks := kDist.GetDegrees()
	for _, k := range ks {
		fmt.Println("-", k, ":", kDist.GetNum(k))
	}
}

func writeGraph(g *graph.Graph) error {
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
