CMDS = commit make_api_key panos_init

.PHONY: subdirs $(CMDS)

subdirs: build_dir $(CMDS)

build_dir:
	mkdir -p build/

$(CMDS):
	go build -o build/$@ ./cmd/$@/main.go

clean:
	rm -rf build/