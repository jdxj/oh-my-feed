.PHONY: tag
tag:
	@yq -i '.services.oh-my-feed.image = "$(DOCKER_TAG)"' $(DOCKER)/docker-compose.yml

.PHONY: up
up: tag
	@docker compose -f $(DOCKER)/docker-compose.yml up -d --build
	@git restore $(DOCKER)/docker-compose.yml

.PHONY: stop
stop: tag
	@docker compose -f $(DOCKER)/docker-compose.yml stop
	@git restore $(DOCKER)/docker-compose.yml
