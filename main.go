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

func jobWorker(job chan struct{}, c *model.EpochCounter, wg *sync.WaitGroup, s int, f int, iid int, d bool) {
	defer func() {
		wg.Done()
	}()
	for range job {
		// Experimental routine starts here
		netmap, err := model.SampleNetwork(s)
		if err != nil {
			panic(err)
		}
		err = netmap.VisitNode(iid)
		if err != nil {
			panic(err)
		}

		i := -1
		reused := 0

		for !netmap.IsNetworkFilled() {
			i++
			if i > len(netmap.Topology) {
				// debug only
				if d {
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
			stat := netmap.RunEpochNaiveOnce(f, i)
			reused += stat.Reused
		}
		if netmap.IsNetworkFilled() {
			c.Inc(i)
			c.AddRe(reused)
		}
	}
}

func runExperiment(size int, fanout int, numexp int, initid int, debug bool) {
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

	jobs := make(chan struct{}, numexp)

	for j := 0; j < numexp; j++ {
		jobs <- struct{}{}
	}

	for j := 0; j < workerCount; j++ {
		go jobWorker(jobs, &c, wg, size, fanout, initid, debug)
	}
	close(jobs)
	wg.Wait()

	if debug {
		dataString := ""
		for i := 0; i < 20; i++ {
			if try, ok := c.Counter[i]; ok {
				dataString += fmt.Sprintf("%d;", try)
			} else {
				dataString += "0;"
			}
		}
		fmt.Printf("%d;%d;%s\n", size, fanout, dataString)
	} else {
		fmt.Printf("Size: %d Fan-out: %d\n", size, fanout)
		hopNumbers := make([]int, 0, len(c.Counter))
		for hop := range c.Counter {
			hopNumbers = append(hopNumbers, hop)
		}
		sort.Ints(hopNumbers) //sort by key
		for _, ind := range hopNumbers {
			fmt.Printf("%d:%d (%.2f%%)  ", ind+1, c.Counter[ind],
				float32(c.Counter[ind])/float32(numexp)*100)
		}
		fmt.Printf("inf:%d (%.2f%%)\n", c.InfCounter, float32(c.InfCounter)/float32(numexp)*100)
		fmt.Printf("Reused avg: %d\n", c.ReCounter/numexp)
	}
}

func main() {
	sampleSize := flag.Int("s", 100, "size of network map")
	fanoutSize := flag.Int("f", 10, "size of fanout value")
	initialNode := flag.Int("n", 0, "index of leader node")
	numExperiments := flag.Int("c", 10, "number of experiments")
	debug := flag.Bool("debug", false, "debug mode")
	repl := flag.Bool("i", false, "interactive mode")
	flag.Parse()

	if *repl {
		fmt.Println("Interactive push-gossip model runner")
		fmt.Println("Print 'run' and fill model parameters")
		shell := ishell.New()
		shell.AddCmd(&ishell.Cmd{
			Name: "run",
			Help: "run gossip experiment",
			Func: func(c *ishell.Context) {
				// disable the '>>>' for cleaner same line input.
				c.ShowPrompt(false)
				defer c.ShowPrompt(true) // yes, revert after login.

				c.Print("Network size: ")
				netsize, err := strconv.Atoi(c.ReadLine())
				if err != nil || netsize <= 0 {
					c.Println("Incorrect network size")
					return
				}

				c.Print("Fan-out size: ")
				fanout, err := strconv.Atoi(c.ReadLine())
				if err != nil || fanout <= 0 || fanout > netsize-1 {
					c.Println("Incorrect fan-out size")
					return
				}

				c.Print("Number of experiments: ")
				expnum, err := strconv.Atoi(c.ReadLine())
				if err != nil || netsize <= 0 {
					c.Println("Incorrect number of experiments")
					return
				}

				c.Println("-----------")
				runExperiment(netsize, fanout, expnum, 0, false)
			},
		})
		shell.Run()
	} else {
		runExperiment(*sampleSize, *fanoutSize, *numExperiments, *initialNode, *debug)
	}
}
