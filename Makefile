go_mod_name=$(shell go list -m)
git_version=$(shell git describe --always --tags --dirty)
go_build_ldflags=-X '$(go_mod_name)/internal/buildinfo.version=$(git_version)'

.PHONY: test
test:
	go test ./...

.PHONY: install/7hlc
install/7hlc:
	go install \
	-ldflags="$(go_build_ldflags)" \
	./cmd/7hlc/

.PHONY: out/7hlc
out/7hlc: | out
	go build -o $@ \
	-ldflags="$(go_build_ldflags)" \
	./cmd/7hlc/

out:
	mkdir -p $@
