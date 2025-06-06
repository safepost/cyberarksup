# Build for Windows
GOOS=windows GOARCH=amd64 go build -o CyberarkSupervision.exe

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o CyberarkSupervision

# Build for current platform (automatic)
go build