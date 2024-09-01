test-build-and-run:
	docker-compose -f docker-compose.dev.yml up -d --build

test-run:
	docker-compose -f docker-compose.dev.yml up -d

test-stop:
	docker-compose -f docker-compose.dev.yml down

test-restart:
	make test-stop
	make test-build-and-run
