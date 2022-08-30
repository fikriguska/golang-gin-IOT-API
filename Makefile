.PHONY: restart_db

restart_db:
	docker-compose down
	docker-compose up -d 