build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o gomimic ./cmd/gomimic

dev:
	[ -z $(docker ps | grep vsc_go-mimic) ] && docker compose up -d
	export $(cat '${PWD}'/.env | xargs)
	go run ./cmd/gomimic

check:
	go build -o /tmp/gomimic ./cmd/gomimic && echo "gomimic build ok."
	go build -o /tmp/example ./cmd/example && echo "example build ok."
	go test -tags=e2e ./... -v -count=1

mongosh:
	docker exec -it vsc_go-mimic_mongodb mongosh -u root -p example
