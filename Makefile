run:
	cd ./frontend-vue && npm run build
	go -C backend-go run .

env:
	go run .

test:
	cd ./backend-go && \
	go test -coverprofile=cover.out && \
	go tool cover -html=cover.out
	
test-race:
	cd ./backend-go && go test -race -coverprofile=cover.out && go tool cover -html=cover.out