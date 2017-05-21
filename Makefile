all: deps test build

test:
	go test -v --cover --race  `glide novendor`

build:
	CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o main .

deps:
	glide install

run:
	./main