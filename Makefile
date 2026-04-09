include .env
export

env-build:
	@docker compose build todoapp-postgres

env-up:
	@docker compose up -d todoapp-postgres

env-down:
	@docker compose down todoapp-postgres

env-cleanup:
	@read -p "Очистить все volume файлы окружения? Опасность утери данных [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
	  docker compose down todoapp-postgres port-forwarder&& \

	  sudo rm -rf out/pgdata &&\
	  echo "Файлы окружения очищены"; \
  	else \
		echo "Очистка отменена"; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
  		echo "Отсутствует seq. make migrate-create seq=init"; \
  	  	exit 1; \
  	fi; \
	docker compose run --rm todoapp-postgres-migrate \
	create \
	-ext sql \
	-dir /migrations \
	-seq "$(seq)"

migrate-up:
	make migrate-action action=up

migrate-down:
	make migrate-action action=down

migrate-force:
	@read -p "Введите версию для force: " v; \
	make migrate-action action="force $$v"

migrate-action:
	@if [ -z "$(action)" ]; then \
	   echo "Отсутствует action. Пример: make migrate-action action=up"; \
	   exit 1; \
	fi; \
	docker compose run --rm \
	   -e POSTGRES_USER=$(POSTGRES_USER) \
	   -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	   -e POSTGRES_DB=$(POSTGRES_DB) \
	   todoapp-postgres-migrate \
	    -path /migrations \
	    -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@todoapp-postgres:5432/$(POSTGRES_DB)?sslmode=disable" \
	    $(action)

log-cleanup:
	@read -p "Очистить все log файлы? Опасность утери логов [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
	  sudo rm -rf out/logs &&\
	  echo "Файлы логов очищены"; \
	else \
		echo "Очистка отменена"; \
	fi

todoapp-run:
	go mod tidy && \
	go run cmd/todoapp/main.go


#Only for local development
up:
	make env-up && \
	make env-port-forward && \
	sudo chmod -R 777 /out && \
	make todoapp-run

todoapp-deploy:
	@docker compose up -d --build todoapp frontend

todoapp-deploy-prod:
	@docker compose -f docker-compose.yaml -f docker-compose.prod.yaml up -d --build todoapp frontend

todoapp-undeploy:
	@docker compose down todoapp

todoapp-undeploy-prod:
	@docker compose -f docker-compose.yaml -f docker-compose.prod.yaml down todoapp frontend

swagger-gen:
	@docker compose run --rm swagger \
		init \
		-g cmd/todoapp/main.go \
		-o docs \
		--parseInternal \
		--parseDependency