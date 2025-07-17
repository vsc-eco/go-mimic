default:
	go run ./cmd/main.go

dev:
	air
	rm /tmp/go-mimic

check:
	go build -o /tmp/go-mimic ./cmd/main.go && echo "server build ok."
	go test ./... -v

mongosh:
	docker exec -it vsc_go-mimic_mongodb mongosh -u root -p example
