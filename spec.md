# User API Spec

## Register User API

Endpoint :  POST http://34.101.246.52:8080/register

Request Body :

```json
{
  "username" : "imam",
  "email": "imam@example.com",  
  "password" : "rahasia" //min 8 char
}
```

Response Body Success :

```json
{
  "status": "success",
  "message": "Registration successful",
  "data": {
    "id": "12345",
    "username": "imam",
    "email": "imam@example.com"
  }
}
```

Response Body Error :

```json
{
  "status": "error",
  "message": "Registration failed",
  "error_details": "Email or Username address is already in use"
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
  "message": "Login successful",
  "data": {
    "id": "12345",
    "username": "imam",
    "email" : "imam@example.com",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiAiMTIzNDUiLCAiaWF0IjogMTYyMzEyMzUwMH0.H6MQUMR1Jvh7zxP3kW6VXWd7OlvGp7sFcpj2ZDqkNKk"
  }
}

```

Response Body Error :

```json
{
  "status": "error",
  "message": "Login failed",
  "error_details": "Invalid username or Invalid password"
}

```

## Update User API

Endpoint : PATCH /api/users/current

Headers :
- Authorization : token
- Content-Type : form-data

Request Body :

```json

{
  "profile-picture":"file.img",
  "username" : "updated_username"
}
```

Response Body Success :

```json
{
  "status": "success",
  "message": "User information updated successfully",
  "data": {
    "id": "12345",
    "username": "updated_username",
    "profile-picture":"file.img"
  }
}

```

Response Body Error :

```json
{
  "status": "error",
  "message": "Failed to update user information",
  "error_details": ""
}

```

## Get User API

Endpoint : GET /api/users/current

Headers :
- Authorization : token

Response Body Success:

```json
{
  "status": "success",
  "message": "User information retrieved successfully",
  "data": {
    "user_id": "12345",
    "username": "john_doe",
    "email": "john.doe@example.com"
    // additional user information
  }
}

```

Response Body Error :

```json
{
  "status": "error",
  "message": "Failed to retrieve user information",
  "error_details": [
    {
      "code": "UNAUTHORIZED",
      "message": "Invalid or expired token"
    }
    // additional error details if needed
  ]
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