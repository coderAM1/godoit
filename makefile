.PHONY: test startPostgres stopPostgres

test:
	go test ./... --short

startPostgres:
	docker compose -f integration/pgint/docker-compose.yml --project-directory . up --detach

stopPostgres:
	docker compose -f integration/pgint/docker-compose.yml --project-directory . down