.PHONY: run clean

run:
	docker build -t brew-chatbot .
	docker run --rm -p 8080:8080 --env-file .env brew-chatbot

clean:
	docker image prune -f