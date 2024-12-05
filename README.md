# mailer_sender

This simple microservice sends mails to users of [lenic](https://github.com/Anacardo89/lenic):
- upon Register, for account activation
- per user request, for password recovery

## Setup:
- install [go](https://go.dev/doc/install)
- setup the yaml config files `internal/config`
- run `go mod tidy` to fetch dependencies
- make sure [lenic](https://github.com/Anacardo89/lenic) is running
- inside `/cmd` run `gp build` to compile, or `go run .` to run with out compiling
- if you built it, run the executable
- you can now send mails
