goose -dir internal/migrate/postgres/migrations postgres "postgresql://postgres:postgres@127.0.0.1:5436/identic?sslmode=disable" down
goose -dir internal/migrate/postgres/migrations create new_table sql
scp -r ./dist administrator@route:/home/administrator/apps/identic
npx vite-bundle-visualizer

export DOCKER_API_VERSION=1.44
