# Auth API Documentation

## User Endpoints

### Sign-in User
- **URL**: `/sign-in`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "username": "name",
    "password": "Password123"
  }
  ```
- **Response**:
  - **Status**: `201 Created`
  - **Body**:
    ```json
    {
      "state": {
          "status": "Success"
      },
      "data": {
          "userId": 3
      }
    }
    ```

### Log-in User
- **URL**: `/log-in`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "username": "name",
    "password": "Password123"
  }
  ```
- **Response**:
  - **Status**: `200 OK`
  - **Body**:
    ```json
    {
        "state": {
            "status": "Success"
        },
        "data": {
            "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE3NDE4ODYxNDgsInVzZXJJZCI6OH0.qQ8BK4qA0UtXYZ4L-f5QDJBYaU-Tvm2Iq-xbnxXnVdQ"
        }
    }
    ```

## Error Responses
- **Status**: `400 Bad Request`
  - **Body**:
    ```json
    {
      "status": "error",
      "message": "failed to bind request"
    }
    ```
- **Status**: `404 Not Found`
  - **Body**:
    ```json
    {
      "status": "error",
      "message": "resource not found"
    }
    ```
- **Status**: `500 Internal Server Error`
  - **Body**:
    ```json
    {
      "status": "error",
      "message": "unexpected server error"
    }
    ```