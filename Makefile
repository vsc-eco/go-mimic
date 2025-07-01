default:
	go run ./cmd/main.go

dev:
	air --build.cmd "go build -o /tmp/vsc-go-mimic cmd/main.go" --build.bin "/tmp/vsc-go-mimic"

mongosh:
	docker exec -it vsc_go-mimic_mongodb mongosh -u root -p example
