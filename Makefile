build:
	go vet
	go build -o bin/go-docs

run:
	-rm -f bin/go-docs
	go build -o bin/go-docs
	./bin/kill-go-docs
	./bin/go-docs &
	
clean:
	-rm -f bin/go-docs
	-rm -f bin/sqlite.db
	go build -o bin/go-docs
	./bin/kill-go-docs
	./bin/go-docs &

tests:
	-rm -f bin/go-docs
	-rm -f bin/test.db
	go build -o bin/go-docs
	./bin/kill-go-docs
	./bin/go-docs --db bin/test.db --dbdata sql/testdata.sql &
	sleep 1
	./test/tests.sh
	@echo "PASS"

