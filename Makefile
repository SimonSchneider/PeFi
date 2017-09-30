DCPROD=docker-compose.prod.yml
DCLOCAL=docker-compose.local.yml

local-build:
	docker-compose -f $(DCLOCAL) pull

local-up:
	docker-compose -f $(DCLOCAL) -p pefi-local up

local-down:
	docker-compose -f $(DCLOCAL) -p pefi-local down

prod-build:
	docker-compose -f $(DCPROD) build
	docker-compose -f $(DCPROD) pull redis postgres

prod-up:
	docker-compose -f $(DCPROD) up

prod-down:
	docker-compose -f $(DCPROD) down
