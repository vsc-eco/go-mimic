default:
	go run ./cmd/main.go

dev:
	air

check:
	go build -o /tmp/go-mimic ./cmd/main.go && echo "build ok."

mongosh:
	docker exec -it vsc_go-mimic_mongodb mongosh -u root -p example
