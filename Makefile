SRC    := $(shell find . -name '*.go')
SUBDIR :=

.PHONY: all clean test run build upgrade help $(SUBDIR)

all: $(SUBDIR) 		# default action
	@[ -f .git/hooks/pre-commit ] || pre-commit install --install-hooks
	@git config commit.template .git-commit-template

clean: $(SUBDIR)	# clean-up environment
	@find . -name '*.sw[po]' -delete

test:				# run test
	go test -v ./...
	gofmt -w -s $(SRC)

run:				# run in the local environment
	go run cmd/rt/main.go -vv

build:				# build the binary/library
	go mod tidy
	go build -o bin/rt -ldflags "-w -s" cmd/rt/main.go

upgrade:			# upgrade all the necessary packages
	pre-commit autoupdate

help:				# show this message
	@printf "Usage: make [OPTION]\n"
	@printf "\n"
	@perl -nle 'print $$& if m{^[\w-]+:.*?#.*$$}' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?#"} {printf "    %-18s %s\n", $$1, $$2}'

$(SUBDIR):
	$(MAKE) -C $@ $(MAKECMDGOALS)
