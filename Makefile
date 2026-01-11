all: build

build:
	#
	# Building project for current OS and processor architecture...
	#
	go build -ldflags '-s -w' -trimpath -buildvcs
	#
	# Done.
	#

build_all:
	#
	# Cleaning old build files...
	#
	rm -rf dist/
	mkdir dist/
	#
	# Building for Apple ARM x64...
	#
	GOOS=darwin GOARCH=arm64 go build -ldflags '-s -w' -trimpath -buildvcs
	tar -czf dist/cloudflare-speedtest.apple-arm64.tar.gz cloudflare-speedtest
	rm cloudflare-speedtest
	#
	# Building for Linux AMD x64...
	#
	GOOS=linux GOARCH=amd64 go build -ldflags '-s -w' -trimpath -buildvcs
	tar -czf dist/cloudflare-speedtest.linux-amd64.tar.gz cloudflare-speedtest
	rm cloudflare-speedtest
	#
	# Building for Linux ARM x64 (Raspberry Pi)...
	#
	GOOS=linux GOARCH=arm64 go build -ldflags '-s -w' -trimpath -buildvcs
	tar -czf dist/cloudflare-speedtest.linux-arm64.tar.gz cloudflare-speedtest
	rm cloudflare-speedtest
	#
	# Building for Windows x64...
	#
	GOOS=windows GOARCH=amd64 go build -ldflags '-s -w' -trimpath -buildvcs
	zip -q dist/cloudflare-speedtest.windows-x64.zip cloudflare-speedtest.exe
	rm cloudflare-speedtest.exe
	#
	# Done.
	#

update:
	#
	# Updating GO modules...
	#
	go get -u ./...
	#
	# Tidying GO modules...
	#
	go mod tidy

format:
	#
	# Formatting GO files...
	#
	find . -type f -name '*.go' -exec go fmt '{}' \;

clean:
	#
	# Cleaning built files...
	#
	rm -rf dist/
