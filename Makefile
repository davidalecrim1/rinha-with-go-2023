# For testing locally
test-build-and-run:
	docker-compose -f docker-compose.dev.yml up -d --build

test-run:
	docker-compose -f docker-compose.dev.yml up -d

test-stop:
	docker-compose -f docker-compose.dev.yml down

test-restart:
	make test-stop
	make test-build-and-run

# For testing with full environment (limits, nginx, etx)
build-and-run:
	docker-compose -f docker-compose.qa.yml up -d --build

run:
	docker-compose -f docker-compose.qa.yml up -d

stop:
	docker-compose -f docker-compose.qa.yml down

restart:
	make stop
	make build-and-run
