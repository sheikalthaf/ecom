# BUILD APP

`env GOOS=linux GOARCH=amd64 go build ecom.com`
`GOOS=wasip1 GOARCH=wasm go build -o ecom.com`
use https://github.com/xxjwxc/gormt this package for migrating existing DB

Run ./gormt

# SERVER SETUP

```bash
sudo apt-get update
sudo apt-get install ffmpeg
```

# RUN FIBER APP

to list the screens
`screen -ls`

to quit the app
`screen -X -S fiber quit`

to run the app
`screen -S fiber -dm ./fiber`

`screen -X -S fiber quit && screen -S fiber -dm ./fiber`
