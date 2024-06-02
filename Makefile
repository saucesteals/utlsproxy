.PHONY: run utlsproxy


all: utlsproxy

# note: marking "utlsproxy" as phony because we want to recompile
# each time in case we made changes locally to a dependency.

utlsproxy:
	go build -o utlsproxy certstore.go main.go

run: utlsproxy
	./utlsproxy