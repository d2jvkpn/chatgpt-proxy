gitBranch = $(shell git rev-parse --abbrev-ref HEAD)
gitTime = $(shell date +'%FT%T%:z')

build:
	echo ">>> git branch: $(gitBranch), git time: $(gitTime)"
	mkdir -p target
	go build -o target/main main.go

run:
	mkdir -p target
	go build -o target/main main.go
	./target/main
