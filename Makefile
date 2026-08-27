.PHONY: dev

dev:
	docker-compose -f config/docker/dev/docker-compose.yml up

build-dev:
	docker-compose -f config/docker/dev/docker-compose.yml up --build

run_db:
	docker exec -it payment-system-db psql -U postgres -d payment-system