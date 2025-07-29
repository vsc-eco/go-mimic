build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o go-mimic ./cmd/main.go

dev:
	[ -z $(docker ps | grep vsc_go-mimic) ] && docker compose up -d
	export $(cat '${PWD}'/.env | xargs)
	go run ./cmd/main.go

check:
	go build -o /tmp/go-mimic ./cmd/main.go && echo "server build ok."
	go test ./... -v

mongosh:
	docker exec -it vsc_go-mimic_mongodb mongosh -u root -p example
