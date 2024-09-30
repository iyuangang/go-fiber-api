# API Documentation

## Endpoints

### Get All Users

- **URL:** `/api/user/`
- **Method:** `GET`
- **Description:** Retrieves a list of all users.
- **Response:**
  - **Status:** `200 OK`
  - **Body:** 
    ```json
    [
      {
        "id": 1,
        "name": "John Doe",
        "email": "john.doe@example.com",
        "created_at": "2023-10-01T12:00:00Z",
        "updated_at": "2023-10-01T12:00:00Z"
      },
      ...
    ]
    ```

### Get User by ID

- **URL:** `/api/user/{id}`
- **Method:** `GET`
- **Description:** Retrieves a user by their ID.
- **Parameters:**
  - `id` (path) - ID of the user.
- **Response:**
  - **Status:** `200 OK`
  - **Body:**
    ```json
    {
      "id": 1,
      "name": "John Doe",
      "email": "john.doe@example.com",
      "created_at": "2023-10-01T12:00:00Z",
      "updated_at": "2023-10-01T12:00:00Z"
    }
    ```
  - **Errors:**
    - `404 Not Found` - User does not exist.

### Create New User

- **URL:** `/api/user/`
- **Method:** `POST`
- **Description:** Creates a new user.
- **Body:**
  ```json
  {
    "name": "Jane Smith",
    "email": "jane.smith@example.com"
  }
  ```
- **Response:**
  - **Status:** `201 Created`
  - **Body:**
    ```json
    {
      "id": 2,
      "name": "Jane Smith",
      "email": "jane.smith@example.com",
      "created_at": "2023-10-01T12:05:00Z",
      "updated_at": "2023-10-01T12:05:00Z"
    }
    ```
  - **Errors:**
    - `400 Bad Request` - Invalid input.

### Update User

- **URL:** `/api/user/{id}`
- **Method:** `PUT`
- **Description:** Updates an existing user.
- **Parameters:**
  - `id` (path) - ID of the user.
- **Body:**
  ```json
  {
    "name": "Jane Doe",
    "email": "jane.doe@example.com"
  }
  ```
- **Response:**
  - **Status:** `200 OK`
  - **Body:**
    ```json
    {
      "id": 2,
      "name": "Jane Doe",
      "email": "jane.doe@example.com",
      "created_at": "2023-10-01T12:05:00Z",
      "updated_at": "2023-10-01T12:10:00Z"
    }
    ```
  - **Errors:**
    - `400 Bad Request` - Invalid input.
    - `404 Not Found` - User does not exist.

### Delete User

- **URL:** `/api/user/{id}`
- **Method:** `DELETE`
- **Description:** Deletes a user.
- **Parameters:**
  - `id` (path) - ID of the user.
- **Response:**
  - **Status:** `204 No Content`
  - **Errors:**
    - `404 Not Found` - User does not exist.
