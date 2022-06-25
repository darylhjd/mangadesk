x32 = 386
x64 = amd64
ver = v0-7-8

all: windows macos linux

windows:
	@echo "Creating Windows builds";
	env GOOS=windows GOARCH=$(x64) go build -o mangadesk_$(ver)_win64.exe;
	env GOOS=windows GOARCH=$(x32) go build -o mangadesk_$(ver)_win32.exe;

macos:
	@echo "Creating macOS builds";
	env GOOS=darwin GOARCH=$(x64) go build -o mangadesk_$(ver)_mac;

linux:
	@echo "Creating Linux builds";
	env GOOS=linux GOARCH=$(x64) go build -o mangadesk_$(ver)_x64;
	env GOOS=linux GOARCH=$(x32) go build -o mangadesk_$(ver)_x32;