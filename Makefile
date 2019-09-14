docker_name := two_services_sample
docker_tag := 1.0
container_name := table_time_docker



docker:
	docker build -t ${docker_name}:${docker_tag} -f Dockerfile ./

run:
	docker run \
		-p 80:80 \
		--rm -it \
		--name ${container_name} \
		${docker_name}:${docker_tag}
daemon:
	docker run \
		-p 80:80 \
		--rm -d -it \
		--name ${container_name} \
		${docker_name}:${docker_tag}

Nolan:
	docker exec -it ${container_name} bash


stop:
	docker stop ${container_name}

logs:
	docker logs ${container_name}


delete-container:
	docker images
	docker rmi ${docker_name}:${docker_tag}
	docker images
