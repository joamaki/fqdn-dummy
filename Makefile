.PHONY: all clean

all:
	go mod tidy
	go build .

clean:
	rm -f fqdn-dummy

