.PHONY:docker_build_processor docker_build_server docker_compose_run

docker_compose_run:
	docker compose -f build/docker-compose.yaml up

docker_compose_build:
	docker compose -f build/docker-compose.yaml up --build
