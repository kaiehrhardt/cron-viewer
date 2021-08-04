IMAGE_NAME?=cron-viewer
IMAGE_TAG?=latest

docker-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

run:
	docker run -it --rm $(IMAGE_NAME):$(IMAGE_TAG)

run-test-config:
	docker run -it --rm -v $(shell pwd)/config-test.yml:/etc/cron-viewer/config.yml $(IMAGE_NAME):$(IMAGE_TAG)