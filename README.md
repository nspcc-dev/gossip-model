## Demo

[![asciicast](https://asciinema.org/a/Lr0L3Jy2FElKZmXEEb3BxIAod.svg)](https://asciinema.org/a/Lr0L3Jy2FElKZmXEEb3BxIAod)

# Push-Gossip model

Simulated model for efficiency evaluation of push-gossip protocol. 

## Description

This model simulates synchronous push-gossip process in distributed 
system with mesh topology. Model has variety of parameters: 
network size, fan-out size, number of experiments. Application executes 
several experiments and calculates time in hops when all network nodes
got the propagated data.  

## Usage
You can build application or use it inside docker container. 
To build use command `go build`, with `git` and `gcc` pre-installed. 

We recommend to run model inside docker container. For interactive mode 
use `make repl` which builds environment and starts `gossipmodel` with 
`-i` parameter

```
$ make repl
Sending build context to Docker daemon  375.3kB
Step 1/8 : FROM golang:alpine as builder
 ---> f56365ec0638
Step 2/8 : RUN apk add --no-cache git gcc musl-dev
. . .
Successfully built 0fabd6489c97
Successfully tagged gossip-model-image:latest
Interactive push-gossip model runner
Print 'run' and fill model parameters
>>> run
Network size: 100
Fan-out size: 10
Number of experiments: 20
-----------
Size: 100 Fan-out: 10
3:17 (85.00%)  4:2 (10.00%)  inf:1 (5.00%)
Reused avg: 600
31.962821ms
>>>  
```

Another way is to run application in silent mode
```
$ make up
Sending build context to Docker daemon  375.8kB           
Step 1/8 : FROM golang:alpine as builder   
 ---> f56365ec0638                       
. . . 
/ #  gossipmodel -s 100 -f 9 -c 1000
Size: 100 Fan-out: 9
3:743 (74.30%)  4:256 (25.60%)  inf:1 (0.10%)
Reused avg: 544
412.675255ms
```

Application outputs model params and simulation results. 
`3:743(74.3%)` in example above means that 743 out of 1000 
experiments (74.3%) were finished in 3 propagation hops. 

## License

This project is licensed under the GPL v3.0 License - see the 
[LICENSE.md](LICENSE.md) file for details
