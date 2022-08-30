.PHONY: restart_db test

restart_db:
	docker-compose down
	docker-compose up -d 

test:
	go test -v ./... -count=1