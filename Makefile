run-api:
	docker compose up api

db-web-ui:
	docker compose up db_ui

down:
	docker compose down --remove-orphans