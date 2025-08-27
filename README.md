# Go Auth
> Simple authentication and permission checker service

### Requirements:
#### Docker:
 ![docker](https://badgen.net/static/docker/@latest/purple)<br/>
 You can install Docker <a href="https://docs.docker.com/engine/install/">there</a>

### Installing:
1. Clone repository 
2. Go to ./auth-service/config and setup config.yaml as specified in config.example.yaml <br/>
   Note: it is crucial to set repository:host value equal to "auth-database"<br/>

3. Start the application<br/>
  - for Windows users:<br/>
    ```bash
    docker-compose build
    ```
    ```bash
    docker-compose up
    ```
  
  - for Linux users:<br/>
    ```bash
    docker compose build
    ```
    ```bash
    docker compose up
    ```

### GRPC methods
- **Register**: registers user in auth service and returns pair of jwt tokens(refresh and access)
- **Login**: returns pair of jwt tokens(refresh and access) if user is registered
- **Refresh**: returns pair of jwt tokens(refresh and access) if refresh token is valid
- **Validate**: validates access token

### Requests examples
- **Register**:
```
{
  "email": "Roberto@gmail.com",
  "first_name": "Andrey",
  "last_name": "Svishev",
  "password": "lorem_ipsum",
  "username": "Youtubelover"
}
```
- **Login**:
```
{
  "email": "Roberto@gmail.com",
  "password": "lorem_ipsum"
}
```
- **Refresh**:
```
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTYxMzU2MzIsImlkIjowLCJ1c2VybmFtZSI6IiIsImZpcnN0X25hbWUiOiIiLCJsYXN0X25hbWUiOiIiLCJpc19hZG1pbiI6ZmFsc2V9.DfqoXRbnabTbL8tL6HWej7lW_AciEV8mNriMlLxiM-E"
}
```
- **Validate**:
```
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTYwNTY0MzIsImlkIjowLCJ1c2VybmFtZSI6ImV1IiwiZmlyc3RfbmFtZSI6Im1hZ25hIGV1IGN1bHBhIHZlbGl0IiwibGFzdF9uYW1lIjoiZXUgZnVnaWF0IG1pbmltIGFsaXF1YSIsImlzX2FkbWluIjpmYWxzZX0.lZZWH37PaWM2c0_MsamWjiW4qt6PzEkjTYQzQ5BgNWk"
}
```

### Response example
```
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTYwNTY0MzIsImlkIjowLCJ1c2VybmFtZSI6ImV1IiwiZmlyc3RfbmFtZSI6Im1hZ25hIGV1IGN1bHBhIHZlbGl0IiwibGFzdF9uYW1lIjoiZXUgZnVnaWF0IG1pbmltIGFsaXF1YSIsImlzX2FkbWluIjpmYWxzZX0.lZZWH37PaWM2c0_MsamWjiW4qt6PzEkjTYQzQ5BgNWk",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTYxMzU2MzIsImlkIjowLCJ1c2VybmFtZSI6IiIsImZpcnN0X25hbWUiOiIiLCJsYXN0X25hbWUiOiIiLCJpc19hZG1pbiI6ZmFsc2V9.DfqoXRbnabTbL8tL6HWej7lW_AciEV8mNriMlLxiM-E"
}
```
