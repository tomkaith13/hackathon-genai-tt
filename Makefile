run:
	go run ./main.go
docker-build:
	docker build -t github.com/tomkaith13/hackathon-genai-tt .
docker-run:
	docker run -p 8080:8080 github.com/tomkaith13/hackathon-genai-tt

clean:
	docker compose down
up:
	make docker-build && docker compose up
restart:
	make clean && make up
