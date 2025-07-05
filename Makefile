default:
	go run ./cmd/main.go

dev:
	docker compose up -d
	air

mongosh:
	docker exec -it vsc_go-mimic_mongodb mongosh -u root -p example
