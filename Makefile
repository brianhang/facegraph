.PHONY: webserver

clean:
	rm -rf bin
	
webserver:
	mkdir -p bin && cd cmd/webserver && go build -o ../../bin/webserver main.go

dev:
	make webserver && WEBSERVER_PORT=8080 ./bin/webserver
