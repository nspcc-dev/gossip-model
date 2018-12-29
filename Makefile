.PHONY: build up repl

build:
	@docker build . -t gossip-model-image

up: build
	@docker run -it --rm gossip-model-image:latest

repl: build
	@docker run -it --rm gossip-model-image:latest gossipmodel -i

