default: build

build:
	go mod vendor
	docker-compose up -d --build
	docker exec mygram_api /bin/mygram_migrate up
