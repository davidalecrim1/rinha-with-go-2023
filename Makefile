# For testing locally
dev-build-and-run:
	docker compose -f ./docker-compose.dev.yml up -d --build

dev-run:
	docker compose -f ./docker-compose.dev.yml up -d

dev-stop:
	docker compose -f ./docker-compose.dev.yml down

dev-restart:
	make dev-stop
	make dev-build-and-run

# For testing with full environment in Windows or MacOS (limits, nginx, etc)
build-and-run:
	docker compose -f ./docker-compose.yml up -d --build

run:
	docker compose -f ./docker-compose.yml up -d

stop:
	docker compose -f ./docker-compose.yml down

restart:
	make stop
	make build-and-run


# For testing in linux environment with network_mode=host, given bridge mode causes performance issues
linux-build-and-run:
	docker compose -f ./docker-compose.linux.yml up -d --build

linux-run:
	docker compose -f ./docker-compose.linux.yml up -d

linux-stop:
	docker compose -f ./docker-compose.linux.yml down

linux-restart:
	make linux-stop
	make linux-build-and-run
