0. build requires make

- sudo apt install make

1. sqlite-3 requires gcc, CGO_ENABLED=1

- sudo apt install gcc
- export CGO_ENABLED=1
- go get github.com/mattn/go-sqlite3
- go install github.com/mattn/go-sqlite3

2. env config

- source .setup

3. build

- go mod tidy
- make build

4. make and run

- make run 

5. or erase db and run

- make clean

6. run tests with make

- sudo apt install jq wget
- make tests
