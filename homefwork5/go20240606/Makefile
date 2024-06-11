PROGRAM = "./cmd/url.go"

all:
	make swag
	make fmt
	make run

swag:
	cp ./cmd/url.go ./main.go
	swag fmt
	swag init
	mv ./main.go ./cmd/url.go

init:
	go mod init short-url
	go install github.com/swaggo/swag/cmd/swag@latest
	go get -u github.com/swaggo/echo-swagger
	docker run -d -p 5672:5672 --hostname my-rabbit --name some-rabbit rabbitmq:3

test:
	go clean -testcache
	go test -v ./... -cover

fmt:
	go mod tidy
	go fmt ./...

build:
	make fmt
	# go build -o url $(PROGRAM)
	# CGO_ENABLED=0 GOOS=linux env allow to build a linux binary, that can run in alpine
	CGO_ENABLED=0 GOOS=linux go build -v -o url $(PROGRAM)

run:
	go run $(PROGRAM)

docker:
	docker build -t url .
	docker run -d -p 8080:8080 url
	# tag the image
	docker tag url leon1234858/url:latest
	# push the image
	docker push leon1234858/url:latest

k8s:
	kubectl apply -f k8s.yaml
	echo "http://127.0.0.1:8080/swagger/index.html"
	# kubectl delete -f k8s.yaml -R
	# remove persistent volume
	# kubectl delete persistentvolumeclaim data-postgres-0
	# kubectl delete persistentvolumeclaim data-mongodb-0
	# kubectl get pvc