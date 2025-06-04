# Build for Windows
GOOS=windows GOARCH=amd64 go build -o app.exe

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o app

# Build for current platform (automatic)
go build