download-dependencies:
	go mod download
.PHONY: download-dependencies

remove-artifacts:
	rm -f *.tgz
.PHONY: remove-artifacts

package-windows-amd64:
	mkdir listagent_windows_x86_64
	env CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -ldflags="-s -w" -trimpath -o ./listagent_windows_x86_64/listagent.exe main.go
	tar --exclude="*.DS_Store" -czf listagent_windows_x86_64.tgz listagent_windows_x86_64/
	rm -r listagent_windows_x86_64
.PHONY: package-windows-amd64

package-windows-arm64:
	mkdir listagent_windows_arm64
	env CGO_ENABLED=0 GOARCH=arm64 GOOS=windows go build -ldflags="-s -w" -trimpath -o ./listagent_windows_arm64/listagent.exe main.go
	tar --exclude="*.DS_Store" -czf listagent_windows_arm64.tgz listagent_windows_arm64/
	rm -r listagent_windows_arm64
.PHONY: package-windows-amd64

package-darwin-amd64:
	mkdir listagent_darwin_amd64
	env CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -ldflags="-s -w" -trimpath -o ./listagent_darwin_amd64/listagent main.go
	tar --exclude="*.DS_Store" -czf listagent_darwin_amd64.tgz listagent_darwin_amd64/
	rm -r listagent_darwin_amd64
.PHONY: package-windows-amd64

package-darwin-arm64:
	mkdir listagent_darwin_arm64
	env CGO_ENABLED=0 GOARCH=arm64 GOOS=darwin go build -ldflags="-s -w" -trimpath -o ./listagent_darwin_arm64/listagent main.go
	tar --exclude="*.DS_Store" -czf listagent_darwin_arm64.tgz listagent_darwin_arm64/
	rm -r listagent_darwin_arm64
.PHONY: package-windows-arm64

package: remove-artifacts release
.PHONY: package

release: package-windows-amd64 package-windows-arm64 package-darwin-amd64 package-darwin-arm64
.PHONY: release

build-local:
	go build -o $(shell which listagent) main.go
.PHONY: build-local

generate-docs:
	mkdir -p ./docs
	rm -f ./docs/*.md
	go run -tags docs ./... --dir ./docs
.PHONY: generate-docs