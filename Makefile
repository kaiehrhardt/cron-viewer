IMAGE_NAME?=cron-viewer
IMAGE_TAG?=latest

docker-build:
	@docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

run: docker-build
	@docker run -it --rm -p 8080:8080 $(IMAGE_NAME):$(IMAGE_TAG)

run-test-config: docker-build
	@docker run -it --rm -p 8080:8080 -v $(shell pwd)/config/config-test.yml:/etc/cron-viewer/config.yml $(IMAGE_NAME):$(IMAGE_TAG)