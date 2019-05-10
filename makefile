ifdef NO_CACHE
ifeq ($(NO_CACHE), $(filter $(NO_CACHE), true))
	CACHE_OPTION := "--no-cache"
else
	CACHE_OPTION :=
endif
endif

build:
	@docker-compose build $(CACHE_OPTION)

env-up:
	@docker-compose up -d

env-down:
	@docker-compose down

.PHONY: help
help:
	@echo 'Environments:'
	@echo ' build 				- Build the rabbitmq image'
	@echo '	env-up				- Launch the environment with RabbitMQ'
	@echo '	env-down			- Stop the environment'
