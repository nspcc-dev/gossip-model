package model

/*
	Here defined different algorithms for push gossip processing
	In model we use RunEpochNaiveOnce algorithm
*/

//	If node has data, choose F other nodes and propagate info. Do it once in lifetime.
//	Topology notation:
//  `- -1 : Node has data, already propagated
//  `-  0 : Node has not data
//	`-  1 : Node has data, ready to propagate
func (n *Network) RunEpochNaiveOnce(fanout int, epoch int) Stat {
	var s Stat

	newVotes := make(map[int]int, len(n.Topology))

	for ind, v := range n.Topology {
		if v == 1 {
			voted := n.ChooseNodesCheck(fanout, n.generated[ind], ind)
			n.SetHistoryEpoch(ind, epoch, voted)
			s.Sent += len(voted)
			for _, vote := range voted {
				newVotes[vote]++
			}
			n.Topology[ind] = -1
		}
	}

	for node, repeated := range newVotes {
		if n.Topology[node] != 0 {
			s.Reused++
		} else {
			n.Topology[node] = 1
		}
		if repeated > 1 {
			s.Reused += repeated - 1
		}
	}
	s.Coverage = n.CountCoverage()
	return s
}

/*
	There are also different processing approaches that can be used in a model
*/

//	If node has data, choose F other nodes and propagate info.
//	Do it forever until some service will stop it.
//	Topology notation:
//  `-  0 : Node has not data
//	`-  1 : Node has data, ready to propagate
func (n *Network) RunEpochNaiveForever(fanout int, epoch int) Stat {
	var s Stat

	newVotes := make(map[int]int, len(n.Topology))

	for ind, v := range n.Topology {
		if v == 1 {
			voted := n.ChooseNodesCheck(fanout, n.generated[ind], ind)
			n.SetHistoryEpoch(ind, epoch, voted)
			s.Sent += len(voted)
			for _, vote := range voted {
				newVotes[vote]++
			}
		}
	}
	for node, repeated := range newVotes {
		if n.Topology[node] == 1 {
			s.Reused++
		}
		n.Topology[node] = 1
		if repeated > 1 {
			s.Reused += repeated - 1
		}
	}
	s.Coverage = n.CountCoverage()
	return s
}

// 	Improvement of simple algorithm where node do not send message
// 	to another node twice.
//	Topology notation:
//  `-  0 : Node has not data
//	`-  1 : Node has data, ready to propagate
func (n *Network) RunEpochNaiveForeverMemorise(fanout int, epoch int) Stat {
	var s Stat

	newVotes := make(map[int]int, len(n.Topology))

	for ind, v := range n.Topology {
		if v == 1 {
			voted := n.ChooseNodesCheck(fanout, n.generated[ind], ind)
			n.SetHistoryEpoch(ind, epoch, voted)
			s.Sent += len(voted)
			for _, vote := range voted {
				n.generated[ind][vote] = true
				newVotes[vote]++
			}
		}
	}
	for node, repeated := range newVotes {
		if n.Topology[node] == 1 {
			s.Reused++
		}
		n.Topology[node] = 1
		if repeated > 1 {
			s.Reused += repeated - 1
		}
	}
	s.Coverage = n.CountCoverage()
	return s
}

//	Only one node propagate data without memorising .
//	Topology notation:
//  `-  0 : Node has not data
//	`-  1 : Node has data, ready to propagate
func (n *Network) RunEpochCentralised(fanout int, epoch int) Stat {
	var s Stat

	newVotes := make(map[int]int, len(n.Topology))

	voted := n.ChooseNodesCheck(fanout, n.generated[0], 0)
	n.SetHistoryEpoch(0, epoch, voted)
	s.Sent += len(voted)
	for _, vote := range voted {
		newVotes[vote]++
	}

	for node, repeated := range newVotes {
		if n.Topology[node] == 1 {
			s.Reused++
		}
		n.Topology[node] = 1
		if repeated > 1 {
			s.Reused += repeated - 1
		}
	}
	s.Coverage = n.CountCoverage()
	return s
}

//	Only one node propagate data with memorising .
//	Topology notation:
//  `-  0 : Node has not data
//	`-  1 : Node has data, ready to propagate
func (n *Network) RunEpochCentralisedMemorise(fanout int, epoch int) Stat {
	var s Stat

	newVotes := make(map[int]int, len(n.Topology))

	voted := n.ChooseNodesCheck(fanout, n.generated[0], 0)
	n.SetHistoryEpoch(0, epoch, voted)
	s.Sent += len(voted)
	for _, vote := range voted {
		newVotes[vote]++
	}

	for node, repeated := range newVotes {
		if n.Topology[node] == 1 {
			s.Reused++
		}
		n.Topology[node] = 1
		n.generated[0][node] = true
		if repeated > 1 {
			s.Reused += repeated - 1
		}
	}
	s.Coverage = n.CountCoverage()
	return s
}

//	If node has data, choose F other nodes and propagate info. Do it once in lifetime.
//  Also send vector of parent nodes to exclude them in choosing process.
//	Topology notation:
//  `- -1 : Node has data, already propagated
//  `-  0 : Node has not data
//	`-  1 : Node has data, ready to propagate
func (n *Network) RunEpochVectorOnce(fanout int, epoch int) Stat {
	var s Stat

	newVotes := make(map[int]int, len(n.Topology))
	voters := make(map[int][]int, len(n.Topology))

	for ind, v := range n.Topology {
		if v == 1 {
			voted := n.ChooseNodesCheck(fanout, n.generated[ind], ind)
			n.SetHistoryEpoch(ind, epoch, voted)
			s.Sent += len(voted)
			for _, vote := range voted {
				newVotes[vote]++

				if _, ok := voters[vote]; !ok {
					voters[vote] = make([]int, 0, len(n.Topology))
				}
				voters[vote] = append(voters[vote], ind)
			}
			n.Topology[ind] = -1
		}
	}

	for node, repeated := range newVotes {
		if n.Topology[node] != 0 {
			s.Reused++
		} else {
			voter := voters[node][0]
			for k := range n.generated[voter] {
				n.generated[node][k] = true
			}
			n.Topology[node] = 1
		}
		if repeated > 1 {
			s.Reused += repeated - 1
		}
	}
	s.Coverage = n.CountCoverage()
	return s
}
