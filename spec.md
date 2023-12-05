# User API Spec

## Register User API

Endpoint :  POST http://34.101.246.52:8080/register

Request Body :

```json
{
  "username" : "imam",
  "email": "imam@example.com",  
  "password" : "rahasia",
  "name" : "Imam Ahmad Fahrezi"
}
```

Response Body Success :

```json
{
  "status": "201",
  "message": "success registered",
  "data" : {
    "email": "imam@example.com",
    "name" : "Imam Ahmad Fahrezi"
  }
}
```

Response Body Error :

```json
{
  "status": "",
  "message": "email or username already exist"
}
```

## Login User API

Endpoint : POST http://34.101.246.52:8080/register

Request Body :

```json
{
"email" : "imam@example.com",
"password" : "rahasia"
}
```

Response Body Success :
data tersebut maksudnya token
```json
{
  "status": "success",
  "message": "ok",
  "data": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDE3NTA4MjQsImlhdCI6MTcwMTc0NzIyNCwic3ViIjoiMyJ9.Mkv8F2gFXUVsbolhc_H2jRJGwXuIb7RmNOW8q3TTUcE"
}
```

Response Body Error :

```json
{
  "status": "unauthorized",
  "message": "Invalid email or password"
}
```

## Update User API

Endpoint : PATCH /api/users/current

Headers :
- Authorization : token

Request Body :

```json
{
"username" : "imamnew"
}
```

Response Body Success :

```json
{
"data" : {
"username" : "imam",
"name" : "Imam Ahmad Fahrezi new"
}
}
```

Response Body Error :

```json
{
"errors" : "Name length max 100"
}
```

## Get User API

Endpoint : GET /api/users/current

Headers :
- Authorization : token

Response Body Success:

```json
{
"data" : {
"username" : "imam",
"name" : "Imam Ahmad Fahrezi"
}
}
```

Response Body Error :

```json
{
"errors" : "Unauthorized"
}
```

## Logout User API

Endpoint : DELETE /api/users/logout

Headers :
- Authorization : token

Response Body Success :

```json
{
"data" : "OK"
}
```

Response Body Error :

```json
{
"errors" : "Unauthorized"
}
```