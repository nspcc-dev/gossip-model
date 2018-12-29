package model

import (
	"testing"

	. "github.com/onsi/gomega"
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
	g := NewGomegaWithT(t)

	net, err := prepareNetwork(10)
	g.Expect(err).NotTo(HaveOccurred())
	net.RunEpochNaiveOnce(2, 0)
	g.Expect(net.CountCoverage()).To(Equal(3))

	net, err = prepareNetwork(10)
	g.Expect(err).NotTo(HaveOccurred())
	net.RunEpochNaiveOnce(9, 0)
	g.Expect(net.CountCoverage()).To(Equal(10))
}

func TestNetwork_ChooseNodesCheck(t *testing.T) {
	g := NewGomegaWithT(t)

	net, err := prepareNetwork(10)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(net.ChooseNodesCheck(2, nil)).To(HaveLen(2))
	g.Expect(net.ChooseNodesCheck(10, nil)).To(HaveLen(10))
	g.Expect(net.ChooseNodesCheck(12, nil)).To(HaveLen(0))

	r1 := net.ChooseNodesCheck(2, map[int]bool{0: true, 1: true, 2: true, 3: true, 4: true, 6: true, 7: true, 9: true})
	g.Expect(r1).To(ContainElement(5))
	g.Expect(r1).To(ContainElement(8))
}
