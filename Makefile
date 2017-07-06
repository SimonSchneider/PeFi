build: up build/pefi/pefi build/pefi/static

api: build/api/pefi-api

up:
	docker-compose up -d
	docker-compose exec run go get -d -v .

down:
	docker-compose down

build/pefi/pefi: *.go
	docker-compose exec run go get -d -v .
	docker-compose exec run bash -c "CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $@ ."
	docker-compose exec run chmod +x build/pefi/pefi

build/pefi/static: static/
	docker-compose exec run rm -rf $@
	docker-compose exec run cp -r static/ $@

create_image: build/
	cd build && docker build -t simonschneider/pefi .

run: 
	docker-compose exec run go run *.go

clean: up
	docker-compose exec run rm -rf build/pefi

build/api/pefi-api: api/main/pefi-api.go
	docker-compose exec run go get -d -v pefi/api/main
	docker-compose exec run bash -c "CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $@ pefi/api/main/"
	docker-compose exec run chmod +x build/api/pefi-api

create_api_img: build/api/pefi-api
	cd build && docker build -f Dockerfile.api -t simonschneider/pefi:api .
