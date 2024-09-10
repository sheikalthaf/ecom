GOOS=linux GOARCH=amd64 go build -o bin/application

chmod +x bin/application

zip -r build.zip bin/application Procfile