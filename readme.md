# Telego
[*] move validators and bindings to separate middleware
[*] move parse id to separate middleware
[*] assets
[*] access authorization
[*] add proper migrations
[*] login
[] transform response
[] error message
[] add redis
[] todos

## Create a new migration
```
$ migrate create -ext sql -dir app/db/migrations -seq ${create_users_table}
```