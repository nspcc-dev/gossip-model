package model

import (
	"errors"
	"strconv"
	"strings"
)

type (
	Cluster struct {
		Prob     float64
		Capacity uint64
	}

	ClusterList []*Cluster

	Node struct {
		data       int
		cluster_id int // -1 means that cluster undefined
	}

	Network struct {
		Topology   map[int]*Node         // defines map of available nodes
		LinkMatrix map[int]map[int]bool  // Connectivity Matrix
		History    map[int]map[int][]int // history of all propagation changes
		generated  map[int]map[int]bool  // extra structure for history based algorithms.
	}
)

func (cl *ClusterList) Set(value string) error {
	arr := strings.Split(value, "/")
	if len(arr) != 2 {
		return errors.New("Bad format of cluster data")
	}

	if prob, err := strconv.ParseFloat(arr[0], 64); err != nil {
		return err
	} else if capacity, err := strconv.ParseUint(arr[1], 10, 64); err != nil {
		return err
	} else if prob < 0 || prob > 1 {
		return errors.New("Probability should be in the range: [0..1]")
	} else {
		*cl = append(*cl, &Cluster{prob, capacity})
		return nil
	}
}

func (cl *ClusterList) String() string {
	return ""
}

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
		if v.data == 0 {
			return false
		}
	}
	return true
}

func (n *Network) VisitNode(i int) error {
	if i < 0 || i >= int(len(n.Topology)) {
		return errors.New("index is out of sample rangle")
	}
	n.Topology[i].data++
	return nil
}

func (n Network) CountCoverage() (result int) {
	result = 0
	for _, v := range n.Topology {
		if v.data != 0 {
			result++
		}
	}
	return
}

func (n Network) IsLinkExist(node1 int, node2 int) bool {
	return n.LinkMatrix[node1][node2]
}

func (n Network) GetClusterConnectivity(node1 int, node2 int, clusters ClusterList, default_prob float64) bool {
	if n.Topology[node1].cluster_id == n.Topology[node2].cluster_id && n.Topology[node1].cluster_id > 0 {
		prob := clusters[n.Topology[node1].cluster_id].Prob
		return r.Float64() <= prob || node1 == node2
	} else {
		return r.Float64() <= default_prob || node1 == node2
	}
}

func AllocateClusterId(clusters ClusterList) int {
	for id, cluster := range clusters {
		if cluster.Capacity > 0 {
			cluster.Capacity--
			return id
		}
	}
	return -1
}

func SampleNetwork(size int, clusters ClusterList, default_prob float64) (Network, error) {
	if size <= 0 {
		return Network{}, errors.New("sample size must be greater than zero")
	}
	netmap := Network{
		Topology:   make(map[int]*Node, size),
		LinkMatrix: make(map[int]map[int]bool, size),
		History:    make(map[int]map[int][]int, size),
		generated:  make(map[int]map[int]bool, size),
	}

	for i := 0; i < size; i++ {
		netmap.Topology[i] = &Node{data: 0, cluster_id: AllocateClusterId(clusters)}
		netmap.generated[i] = make(map[int]bool)
		netmap.generated[i][i] = true
	}

	for i := 0; i < size; i++ {
		netmap.LinkMatrix[i] = make(map[int]bool)
		for j := 0; j < size; j++ {
			netmap.LinkMatrix[i][j] = netmap.GetClusterConnectivity(i, j, clusters, default_prob)
		}
	}
	return netmap, nil
}
