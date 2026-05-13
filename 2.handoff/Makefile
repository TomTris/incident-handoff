run:
	go run .

test:
	go test -coverprofile=cover.out && go tool cover -html=cover.out
	
test-race:
	go test -race -coverprofile=cover.out && go tool cover -html=cover.out