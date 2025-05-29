.PHONY: build
build:
    go build -buildmode=plugin -o sample-plugin.so main.go

.PHONY: clean
clean:
    rm -f *.so

.PHONY: install
install: build
    # Installation commands if needed