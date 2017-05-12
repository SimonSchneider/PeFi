all: build run

up:
	docker-compose up -d --build build

build:
	docker-compose exec build make

run:
	docker-compose up --build run
