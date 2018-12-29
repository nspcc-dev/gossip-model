package model

import (
	"crypto/rand"
	"encoding/binary"
	mrand "math/rand"
)

type (
	CryptoSource struct {
		buf [8]byte
	}
)

var (
	r = mrand.New(&CryptoSource{})
)

func (s *CryptoSource) Int63() int64 {
	rand.Read(s.buf[:])
	return int64(binary.BigEndian.Uint64(s.buf[:]) & (1<<63 - 1))
}

func (s CryptoSource) Seed(seed int64) {
	panic("seed")

}

func (n *Network) ChooseNodesCheck(fanout int, exclude map[int]bool) []int {
	if fanout > len(n.Topology) {
		return []int{}
	}
	var nodes []int
	if len(n.Topology)-len(exclude) <= len(n.Topology)/2 {
		// if re-random is way too long
		candidates := r.Perm(len(n.Topology))
		for i := 0; len(nodes) < fanout && i < len(n.Topology); i++ {
			if !exclude[candidates[i]] {
				nodes = append(nodes, candidates[i])
			}
		}
	} else {
		// if re-random is fast
		alreadySelected := make(map[int]bool, fanout)
		for len(nodes) < fanout {
			candidate := r.Intn(len(n.Topology))
			if !exclude[candidate] && !alreadySelected[candidate] {
				nodes = append(nodes, candidate)
				alreadySelected[candidate] = true
			}
		}
	}
	return nodes
}
