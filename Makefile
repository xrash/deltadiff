.PHONY : build
build :
	go build -o ./bin/deltadiff ./cmd/deltadiff/*.go

.PHONY : run
run : build
	./bin/deltadiff

.PHONY : test
test :
	go test ./...

.PHONY : install
install : 
	cp ./bin/deltadiff /usr/local/bin
