.PHONY : all clean format test cover

PACKAGES=$(shell go list ./... | grep -v "vendor")

format:
	find . -name "*.go" -not -path "./vendor/*" -not -path ".git/*" | xargs gofmt -s -d -w

clean:
	[ -f user-service-osx ] && rm user-service-osx || true
	[ -f user-service-linux ] && rm user-service-linux || true
	[ -f user-service32.exe ] && rm user-service32.exe || true
	[ -f user-service64.exe ] && rm user-service64.exe || true
	[ -f coverage.txt ] && rm coverage.txt || true
	rm ./coverages/*.txt
	
user-service-osx: main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags '-s -w' -o $@

user-service-linux: main.go
	GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -o $@

user-service64.exe: main.go
	GOOS=windows GOARCH=amd64 go build -ldflags '-s -w' -o $@

user-service32.exe: main.go
	GOOS=windows GOARCH=386 go build -ldflags '-s -w' -o $@

user-service-windows: user-service64.exe b2c32.exe

user-service: user-service-osx user-service-linux user-service-windows

test:
	$(foreach pkg, $(PACKAGES), \
	go test $(pkg);)

cover:
	@echo "mode: cover" > coverage.txt
	@echo "make coverprofile"

	$(foreach pkg, $(PACKAGES), \
	go test -coverprofile=coverage.out -covermode=atomic $(pkg); \
	tail -n +2 coverage.out >> coverage.txt;)
	rm coverage.out