deploy-socket:
	cd pkg/cmd/gpio_socket && make deploy

build:
	cd pkg/cmd/sump-boi && make build

undeploy:
	cd pkg/cmd/sump-boi && make undeploy

deploy:
	cd pkg/cmd/sump-boi && make deploy
