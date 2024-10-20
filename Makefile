all: format mod test dirty

dirty:
	git diff --exit-code

mod:
	go mod tidy

format: license
	gofmt -w .

test:
	go test -v ./...

license:
	licenser apply -r "Notedown Authors"
