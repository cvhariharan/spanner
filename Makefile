install: build
	go install

build: main.go
	pkger -o codegen -include /codegen/templates

