.PHONY: build db api test up down

build:
	docker-compose build

db:
	docker-compose up -d db

api:
	docker-compose up -d product-api

test:
	docker-compose run --rm test

up: db api

down:
	docker-compose down -v
