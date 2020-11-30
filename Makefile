build: main.go
	pkger -o codegen -include /codegen/templates

install: build
	go install