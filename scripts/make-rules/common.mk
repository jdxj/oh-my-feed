DOCKER := docker
OUTPUT := output

DOCKER_TAG := jdxj/oh-my-feed:test-$(shell git rev-parse --short HEAD)

.PHONY: clean
clean:
	@rm -rf $(OUTPUT)
	@rm -rf $(DOCKER)/*.out
