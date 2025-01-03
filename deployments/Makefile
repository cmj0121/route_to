.PHONY: all clean test run stop build upgrade help

all: 			# default action

clean:			# clean-up environment

test:			# run test

run:			# run in the local environment
	docker compose up -d

stop:			# stop the running container
	docker compose stop

build:			# build the binary/library
	docker compose build

upgrade:		# upgrade all the necessary packages

help:			# show this message
	@printf "Usage: make [OPTION]\n"
	@printf "\n"
	@perl -nle 'print $$& if m{^[\w-]+:.*?#.*$$}' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?#"} {printf "    %-18s %s\n", $$1, $$2}'
