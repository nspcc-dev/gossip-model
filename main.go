package main

import (
	"flag"
	"fmt"
	"gossipmodel/model"
	"math/rand"
	"os"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/abiosoft/ishell"
)

var (
	R *rand.Rand
)

type TestInfo struct {
	size         int
	fanout       int
	numexp       int
	leader_node  int
	debug        bool
	clusters     model.ClusterList
	default_prob float64
}

func (info *TestInfo) String() string {
	var str string
	str += fmt.Sprintf("Size:                 %d\n", info.size)
	str += fmt.Sprintf("Fanout:               %d\n", info.fanout)
	str += fmt.Sprintf("NumExp:               %d\n", info.numexp)
	str += fmt.Sprintf("LeaderNode:           %d\n", info.leader_node)
	str += fmt.Sprintf("Debug:                %t\n", info.debug)
	str += fmt.Sprintf("Default probability:  %f\n", info.default_prob)
	if len(info.clusters) > 0 {
		for id, v := range info.clusters {
			str += fmt.Sprintf("Cluster-%d:            %f   capacity=%d\n", id+1, v.Prob, v.Capacity)
		}
	} else {
		str += fmt.Sprintf("There are no clusters\n")
	}

	return str
}

func jobWorker(job chan struct{}, c *model.EpochCounter, wg *sync.WaitGroup, testInfo TestInfo) {
	defer func() {
		wg.Done()
	}()
	for range job {
		// Experimental routine starts here
		netmap, err := model.SampleNetwork(testInfo.size, testInfo.clusters, testInfo.default_prob)
		if err != nil {
			panic(err)
		}
		err = netmap.VisitNode(testInfo.leader_node)
		if err != nil {
			panic(err)
		}

		i := -1
		reused := 0

		for !netmap.IsNetworkFilled() {
			i++
			if i > len(netmap.Topology) {
				// debug only
				if testInfo.debug {
					fmt.Println("Found infinite cycle!")
					for epochNum := 0; epochNum < len(netmap.Topology); epochNum++ {
						if epoch, ok := netmap.History[epochNum]; ok {
							fmt.Println("Epoch:", epochNum+1)
							for nodeid, data := range epoch {
								fmt.Printf("  Node:#%d %v\n", nodeid, data)
							}
						}
					}
					os.Exit(1)
				}
				c.InfCounter++
				break
			}
			// Here we calling gossip algorithm
			stat := netmap.RunEpochNaiveOnce(testInfo.fanout, i)
			reused += stat.Reused
		}
		if netmap.IsNetworkFilled() {
			c.Inc(i)
			c.AddRe(reused)
		}
	}
}

func runExperiment(testInfo TestInfo) {

	fmt.Println(testInfo.String())

	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start))
	}()

	workerCount := runtime.NumCPU() + runtime.NumCPU()/2

	wg := new(sync.WaitGroup)
	wg.Add(workerCount)

	c := model.EpochCounter{
		Mu:         new(sync.Mutex),
		Counter:    make(map[int]int),
		ReCounter:  0,
		InfCounter: 0,
	}

	jobs := make(chan struct{}, testInfo.numexp)

	for j := 0; j < testInfo.numexp; j++ {
		jobs <- struct{}{}
	}

	for j := 0; j < workerCount; j++ {
		go jobWorker(jobs, &c, wg, testInfo)
	}
	close(jobs)
	wg.Wait()

	if testInfo.debug {
		dataString := ""
		for i := 0; i < 20; i++ {
			if try, ok := c.Counter[i]; ok {
				dataString += fmt.Sprintf("%d;", try)
			} else {
				dataString += "0;"
			}
		}
		fmt.Printf("%d;%d;%s\n", testInfo.size, testInfo.fanout, dataString)
	} else {
		fmt.Printf("Size: %d Fan-out: %d\n", testInfo.size, testInfo.fanout)
		hopNumbers := make([]int, 0, len(c.Counter))
		for hop := range c.Counter {
			hopNumbers = append(hopNumbers, hop)
		}
		sort.Ints(hopNumbers) //sort by key
		for _, ind := range hopNumbers {
			fmt.Printf("%d:%d (%.2f%%)  ", ind+1, c.Counter[ind],
				float32(c.Counter[ind])/float32(testInfo.numexp)*100)
		}
		fmt.Printf("inf:%d (%.2f%%)\n", c.InfCounter, float32(c.InfCounter)/float32(testInfo.numexp)*100)
		fmt.Printf("Reused avg: %d\n", c.ReCounter/testInfo.numexp)
	}
}

func main() {
	repl := flag.Bool("i", false, "interactive mode")
	testInfo := TestInfo{
		size:         *flag.Int("s", 100, "size of network map"),
		fanout:       *flag.Int("f", 10, "size of fanout value"),
		leader_node:  *flag.Int("n", 0, "index of leader node"),
		numexp:       *flag.Int("c", 10, "number of experiments"),
		default_prob: *flag.Float64("p", 0.5, "default probability of node connections"),
		debug:        *flag.Bool("debug", false, "debug mode"),
	}

	flag.Var(&testInfo.clusters, "k", "clusters data. e.g.: -k 0.5/100 -k 0.8/20 ...")
	flag.Parse()

	if *repl {
		fmt.Println("Interactive push-gossip model runner")
		fmt.Println("Print 'run' and fill model parameters")
		shell := ishell.New()
		shell.AddCmd(&ishell.Cmd{
			Name: "run",
			Help: "run gossip experiment",
			Func: func(c *ishell.Context) {
				var testInfo TestInfo
				var err error = nil

				// disable the '>>>' for cleaner same line input.
				c.ShowPrompt(false)
				defer c.ShowPrompt(true) // yes, revert after login.

				c.Print("Network size: ")
				testInfo.size, err = strconv.Atoi(c.ReadLine())
				if err != nil || testInfo.size <= 0 {
					c.Println("Incorrect network size")
					return
				}

				c.Print("Fan-out size: ")
				testInfo.fanout, err = strconv.Atoi(c.ReadLine())
				if err != nil || testInfo.fanout <= 0 || testInfo.fanout > testInfo.size-1 {
					c.Println("Incorrect fan-out size")
					return
				}

				c.Print("Number of experiments: ")
				testInfo.numexp, err = strconv.Atoi(c.ReadLine())
				if err != nil || testInfo.size <= 0 {
					c.Println("Incorrect number of experiments")
					return
				}

				c.Print("Default probability of node connections: ")
				testInfo.default_prob, err = strconv.ParseFloat(c.ReadLine(), 64)
				if err != nil || testInfo.default_prob < 0 || testInfo.default_prob > 1 {
					c.Println("Incorrect probability")
					return
				}

				c.Print("Number of clusters: ")
				clusters_num, err := strconv.ParseUint(c.ReadLine(), 10, 64)
				if err != nil {
					c.Println(err)
					return
				}

				var i uint64
				for i = 0; i < clusters_num; i++ {
					c.Printf("Settings for Cluster-%d:\n", i+1)
					c.Printf("\tProbability: ")
					var prob float64
					var capacity uint64

					prob, err = strconv.ParseFloat(c.ReadLine(), 64)
					if err != nil || prob < 0 || prob > 1 {
						c.Println("Incorrect probability")
						return
					}

					c.Printf("\tCapacity: ")
					capacity, err = strconv.ParseUint(c.ReadLine(), 10, 64)
					if err != nil {
						c.Println("Incorrect probability")
						return
					}

					testInfo.clusters = append(testInfo.clusters, &model.Cluster{Prob: prob, Capacity: capacity})
				}

				testInfo.leader_node = 0
				testInfo.debug = false

				c.Println("-----------")
				runExperiment(testInfo)
			},
		})
		shell.Run()
	} else {
		runExperiment(testInfo)
	}
}
