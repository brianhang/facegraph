.PHONY: webserver

webserver:
	mkdir -p bin && cd cmd/webserver && go build -o ../../bin/webserver main.go