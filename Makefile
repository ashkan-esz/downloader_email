update_swagger:
	swag init -g cmd/main.go --output docs

build_rabbitmq:
	docker image build --network=host -t rabbitmq ./docker/rabbitmq

run_rabbitmq:
	#docker run -d --hostname rabbitmq --name rabbitmq -p 15672:15672 -p 5672:5672 --network rabbitnet -e RABBITMQ_DEFAULT_USER=user -e RABBITMQ_DEFAULT_PASS=password rabbitmq
	docker run --restart unless-stopped --network=host --hostname rabbitmq --name rabbitmq -p 15672:15672 -p 5672:5672 -e RABBITMQ_DEFAULT_USER=user -e RABBITMQ_DEFAULT_PASS=password rabbitmq

run_dev:
	clear && swag fmt && make update_swagger && go run cmd/main.go

up:
	docker-compose up --build

build:
	docker image build --network=host -t downloader_email .

run:
	docker run --rm --network=host --name downloader_email -p 587:587 --env-file ./.env downloader_email

push-image:
	docker tag downloader_email ashkanaz2828/downloader_email
	docker push ashkanaz2828/downloader_email

.PHONY: update_swagger build_rabbitmq run_rabbitmq run_dev up build run push-image