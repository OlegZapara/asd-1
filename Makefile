mem ?= 512m
in ?= 1

sort:
	@go run cmd/sort/main.go $(in)

docker-build:
	@docker build -t "asd-1" .

docker-run:
	@docker run -it --memory=$(mem) --memory-swap=$(mem) --name="asd-1" "asd-1" $(in)

pull:
	@docker cp `docker ps -lq`:files/out/$(in).txt files/out/$(in).txt

docker-clean:
	@docker rm "asd-1"

limited: docker-build docker-run pull docker-clean

gen:
	@go run cmd/gen/main.go $(in)

clean:
	@rm files/out/*

clean-in:
	@rm files/in/*

clean-all: clean clean-in