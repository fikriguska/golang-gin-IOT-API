.PHONY: restart_db test testv

restart_db:
	docker-compose down
	docker-compose up -d 

testv:
	go test -v ./... -count=1

test:
	go test ./... -count=1