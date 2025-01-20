## gRPC SSO service

That is simple service you can integrate in your
own project if you want jwt authentication(access, refresh)
mechanism in it.


## Usage

requirements:
- Docker
- Postman(for requests)

At first, you need to add your application into database using migrations:

migration example(up):
```sql
INSERT INTO apps (id, name)
VALUES (1, 'test')
ON CONFLICT DO NOTHING;
```

migrations example(down):
```sql
DELETE FROM apps WHERE id = 1;
```

**Important**! Migration name should look like
(int)*.up.sql and (int)*.down.sql respectively for one migration
where (int) is migration number and * is arbitrary comment.

Then move your migrations into [migrations folder](/internal/storage/migrations/)
and run

```shell
make migrations-up
```

to apply changes.

**To start service**:

```shell
make docker-up
```

it builds up and starts docker container with sso service with SECRET_KEY=default_secret_key.
Or you can run
```shell
docker build -t your_img_name . &&
docker run --rm -d 
	-e SECRET_KEY="your_own_secret_key" 
	-e STORAGE_PATH="sso.db" 
	-e CONFIG_PATH="config.yaml" 
	-p 443:443 your_img_name
```
to set secret key and docker image name

### [Client example](/pkg/client/sso/grpc-client.go)
### Client usage example:

```go
client, err := grpclient.New(...)

isAdmin, err := client.IsAdmin(context.Background, userID)
if err != nil{
    ...
}
...

```
