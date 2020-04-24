.PHONY: build up repl

build-docker:
	@docker build . -t gossip-model-image

up: build-docker
	@docker run -it --rm gossip-model-image:latest

repl: build-docker
	@docker run -it --rm gossip-model-image:latest gossipmodel -i

