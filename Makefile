MB_KEY=put_here_the_key

.PHONY: all compile run tagbox_up tagbox_down

all: compile run

run:
	./visualsearch

compile: 
	go build

tagbox_up:
	docker run \
		-e "MB_KEY=$(MB_KEY)" \
		-e "MB_WORKERS=$(MB_WORKERS)" \
		-d -p 8080:8080 --name tagbox -t machinebox/tagbox

tagbox_down:
	docker rm -f tagbox 2>/dev/null || true


