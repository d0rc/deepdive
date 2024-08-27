package graphs

import (
	"fmt"
	"sort"
)

// Graph represents a directed graph
type Graph struct {
	nodes map[string]bool
	edges map[string]map[string]string
}

// NewGraph returns a new empty graph
func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[string]bool),
		edges: make(map[string]map[string]string),
	}
}

// AddEdge adds an edge between two nodes with a description
func (g *Graph) AddEdge(node1, node2, edgeDescription string) {
	g.addNode(node1)
	g.addNode(node2)
	if _, ok := g.edges[node1]; !ok {
		g.edges[node1] = make(map[string]string)
	}
	g.edges[node1][node2] = edgeDescription
}

func (g *Graph) addNode(node string) {
	g.nodes[node] = true
}

// RenderMermaid renders the graph in Mermaid format
func (g *Graph) RenderMermaid() string {
	var mermaid string
	mermaid += "graph LR\n"
	nodeIds := make(map[string]string)
	i := 0
	for _, node := range g.getSortedNodes() {
		nodeId := fmt.Sprintf("node%d", i)
		nodeIds[node] = nodeId
		mermaid += fmt.Sprintf("    %s[%s]\n", nodeId, node)
		i++
	}
	for node1, edges := range g.edges {
		for node2, description := range edges {
			mermaid += fmt.Sprintf("    %s -->|%s| %s\n", nodeIds[node1], description, nodeIds[node2])
		}
	}
	return mermaid
}

func (g *Graph) getSortedNodes() []string {
	nodes := make([]string, 0, len(g.nodes))
	for node := range g.nodes {
		nodes = append(nodes, node)
	}
	sort.Strings(nodes)
	return nodes
}
