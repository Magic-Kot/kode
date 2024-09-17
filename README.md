# A test task from KODE

### Launching the application in Goland:
- `docker run --name my-postgres -e POSTGRES_PASSWORD=12345 -p 5432:5432 -d postgres`  
- `docker run -p 6379:6379 --name my-redis-stack redis/redis-stack:latest`  
- Edit in the config.yml file "server: host: "127.0.0.1"", "postgres: host: "127.0.0.1"", "redis: host: "127.0.0.1"";  
- Select "main" in the file.go" reading configuration from "internal/config/config.yml";  
- `go build cmd/main.go`  
- `./main`  

If the application is launched for the first time, you need to apply migrations to the database:  
`goose -dir . postgres "postgresql://postgres:12345@127.0.0.1:5432/postgres?sslmode=disable" up`

### Launching an application from docker-compose:
- Edit in the config.yml file "server: host: "0.0.0.0"", "postgres: host: "postgresql"", "redis: host: "redis"";  
- Select "main" in the file.go" reading configuration from "config.yml";   
- `docker compose -f docker-compose.yml up`  

## API

The postman collection file is "Kode.postman_collection".

1) User registration. Request Example:  
   `host:port`

| Path       | Method | Request                                                   | Description  |
|------------|--------|-----------------------------------------------------------|--------------|
| `/sign-up` | POST   | Body: `{"login": "username", "password": "userpassword"}` | Registration |

When registering, the login must contain from 4 to 20 characters, the password - from 6 to 20.  

2) Authorization. Request Example:  
   `host:port/auth`

| Path       | Method | Request                                                                                 | Description    |
|------------|--------|-----------------------------------------------------------------------------------------|----------------|
| `/sign-in` | POST   | Query Params: `GUID=guid`<br/>Body: `{"login": "username", "password": "userpassword"}` | Authorization  |
| `/refresh` | POST   | Cookie: `refreshToken=token; Path=/auth/refresh; HttpOnly;`                             | Refresh tokens |
 
When logging in, the Login and password must contain from 1 to 20 characters.  

3) Working with notes. Request Example:  
   `host:port/note`

| Path   | Method | Request                                                            | Description       |
|--------|--------|--------------------------------------------------------------------|-------------------|
| `/add` | POST   | Header: `Authorization: token`<br/>Body: `{"texts": ["username"]}` | Add a note        |
| `/get` | GET    | Header: `Authorization: token`                                     | Get all the notes |


4) Working with the user. Request Example:  
   `host:port/user`

| Path      | Method | Request                                                                                                                                | Description                           |
|-----------|--------|----------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------|
| `/get`    | GET    | Header: `Authorization: token`                                                                                                         | Getting user data                     |
| `/update` | PUT    | Header: `Authorization: token`<br/>Body: `{"login": "username", "name": "name", "surname": "surname", "age": "age", "email": "email"}` | Changing the login or other user data |
| `/delete` | DELETE | Header: `Authorization: token`                                                                                                         | Deleting the user                     |

The login must contain from 4 to 20 characters.