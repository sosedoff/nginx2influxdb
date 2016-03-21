PROGRAM = nginx2influxdb
TARGETS = darwin/amd64 linux/amd64

build:
	go build -o ./bin/$(PROGRAM)

all:
	gox \
		-osarch="$(TARGETS)" \
		-output="./bin/$(PROGRAM)_{{.OS}}_{{.Arch}}"

clean:
	rm -rf ./bin/*