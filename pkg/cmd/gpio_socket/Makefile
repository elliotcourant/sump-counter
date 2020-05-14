.PHONY: build

build:
	GOOS=linux GOARCH=arm64 go build -o ./bin/gpio_socket github.com/elliotcourant/sump-counter/pkg/cmd/gpio_socket

deploy: build
	scp ./install.sh ./bin/gpio_socket ./gpio_socket.service ubuntu@192.168.0.93:/home/ubuntu/bin
	ssh -t ubuntu@192.168.0.93 'sudo chmod +x ./bin/install.sh && sudo ./bin/install.sh'

