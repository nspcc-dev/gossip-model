package model

import (
	"errors"
)

type (
	Network struct {
		Topology   map[int]int           // defines map of available nodes
		LinkMatrix map[int]map[int]bool  // Connectivity Matrix
		History    map[int]map[int][]int // history of all propagation changes
		generated  map[int]map[int]bool  // extra structure for history based algorithms.
	}
)

func (n Network) GetHistoryEpoch(id int, epoch int) []int {
	return n.History[epoch][id]
}

func (n *Network) SetHistoryEpoch(id int, epoch int, history []int) {
	if _, ok := n.History[epoch]; !ok {
		n.History[epoch] = make(map[int][]int, len(n.Topology))
	}
	n.History[epoch][id] = history
}

func (n Network) IsNetworkFilled() bool {
	for _, v := range n.Topology {
		if v == 0 {
			return false
		}
	}
	return true
}

func (n *Network) VisitNode(i int) error {
	if i < 0 || i >= int(len(n.Topology)) {
		return errors.New("index is out of sample rangle")
	}
	n.Topology[i]++
	return nil
}

func (n Network) CountCoverage() (result int) {
	result = 0
	for _, v := range n.Topology {
		if v != 0 {
			result++
		}
	}
	return
}

func (n Network) IsLinkExist(node1 int, node2 int) bool {
	return n.LinkMatrix[node1][node2]
}

func SampleNetwork(size int, probability float64) (Network, error) {
	if size <= 0 {
		return Network{}, errors.New("sample size must be greater than zero")
	}
	netmap := Network{
		Topology:   make(map[int]int, size),
		LinkMatrix: make(map[int]map[int]bool, size),
		History:    make(map[int]map[int][]int, size),
		generated:  make(map[int]map[int]bool, size),
	}

	for i := 0; i < size; i++ {
		netmap.Topology[i] = 0 // to make map with needed len
		netmap.generated[i] = make(map[int]bool)
		netmap.generated[i][i] = true
	}

	for i := 0; i < size; i++ {
		netmap.LinkMatrix[i] = make(map[int]bool)
		for j := 0; j < size; j++ {
			netmap.LinkMatrix[i][j] = r.Float64() <= probability || i == j
		}

//		fmt.Printf("node-%d: %v\n", i, netmap.LinkMatrix[i])
	}
	return netmap, nil
}
