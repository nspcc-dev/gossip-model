package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func prepareNetwork(size int) (Network, error) {
	var err error
	net, err := SampleNetwork(size)
	if err != nil {
		return net, err
	}
	err = net.VisitNode(0)
	if err != nil {
		return net, err
	}
	return net, nil
}

func TestNetwork_RunEpochNaiveOnce(t *testing.T) {
	net, err := prepareNetwork(10)
	require.NoError(t, err)
	net.RunEpochNaiveOnce(2, 0)
	require.Equal(t, 3, net.CountCoverage())

	net, err = prepareNetwork(10)
	require.NoError(t, err)
	net.RunEpochNaiveOnce(9, 0)
	require.Equal(t, 10, net.CountCoverage())
}

func TestNetwork_ChooseNodesCheck(t *testing.T) {
	net, err := prepareNetwork(10)
	require.NoError(t, err)

	require.Len(t, net.ChooseNodesCheck(2, nil), 2)
	require.Len(t, net.ChooseNodesCheck(10, nil), 10)
	require.Len(t, net.ChooseNodesCheck(12, nil), 0)

	r1 := net.ChooseNodesCheck(2, map[int]bool{0: true, 1: true, 2: true, 3: true, 4: true, 6: true, 7: true, 9: true})
	require.Contains(t, r1, 5)
	require.Contains(t, r1, 8)
}
