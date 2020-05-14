.PHONY: dry

dry:
	kubectl apply -f k8s.yaml --dry-run=true

build:
	go generate ./...
	GOOS=linux GOARCH=arm GOARM=7 go build -o ./bin/sump-boi github.com/elliotcourant/sump-counter/pkg/cmd/sump-boi

docker: build
	docker build --platform linux/arm64 -t 192.168.0.93:32000/sump-boi:latest .
	docker push 192.168.0.93:32000/sump-boi:latest

undeploy:
	kubectl delete -f k8s.yaml --ignore-not-found=true

deploy: docker config
	kubectl scale --replicas=0 deploy/sump-boi || true # shutdown the previous versions
	kubectl apply -f k8s.yaml

clean:
	rm -rfd bin

config:
	kubectl create configmap sump-boi-config --from-file=config.yaml --dry-run -o yaml | kubectl apply -f -

secrets:
	kubectl create secret generic google-application-credentials --from-file=$$GOOGLE_APPLICATION_CREDENTIALS --dry-run -o yaml | kubectl apply -f -
	kubectl create secret generic postgresql-credentials --from-literal=postgres-password=$$POSTGRES_PASSWORD --dry-run -o yaml | kubectl apply -f -
