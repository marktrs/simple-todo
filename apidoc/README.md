# Project: simple-todo

# 📁 Collection: User

## End-point: Create user

### Method: POST

> ```
> http://localhost:3000/api/users
> ```

### Body (**raw**)

```json
{
  "username": "test05",
  "password": "5555"
}
```

### 🔑 Authentication bearer

| Param | value                                                                                                                                                | Type   |
| ----- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | ------ |
| token | eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNTkyNDE1NTIwLCJpZGVudGl0eSI6ImVuZGVyIn0.B9iKWiVn2TMF7uXLLq6t2axOxuvXFhG_WMMYsDREeDY | string |

### Response: 200

```json
{
  "message": "Created user",
  "status": "success",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODI4MDExMDgsInVzZXJfaWQiOiI3NWIxZjVjYi03MTJkLTQxMjYtYmFjZC1kMjJlNWVhZjAzMGYiLCJ1c2VybmFtZSI6InRlc3QwNSJ9.BDF62foGs07GZ2f-FNfMnnW6HqaTxkmI0xqKdkTVuT8",
  "user": {
    "id": "75b1f5cb-712d-4126-bacd-d22e5eaf030f",
    "created_at": "2023-04-27T03:45:08.660779+07:00",
    "updated_at": "2023-04-27T03:45:08.660779+07:00",
    "DeletedAt": null,
    "username": "test05",
    "password": "$2a$14$VQWEEWcTaE0ZUHoCSkH3reiDW6HvR1WNJmglKWOsx/FStPJiuxIOi",
    "tasks": null
  }
}
```

⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

# 📁 Collection: Auth

## End-point: Login

### Method: POST

> ```
> http://localhost:3000/api/auth/login
> ```

### Body (**raw**)

```json
{
  "username": "test01",
  "password": "1111"
}
```

### Response: 200

```json
{
  "status": "success",
  "message": "login success",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2ODI4MDEzMTQsInVzZXJfaWQiOiI1YmRlMjIxMy0xMjExLTQ1MDEtYjY2Yi0wOWE3NWI0MGMxMWYiLCJ1c2VybmFtZSI6InRlc3QwMSJ9.-1IYuB53Y3jyCwZ0lGNfEL4xdt3tpQTJrb-8wZH3zlY",
  "user": {
    "id": "5bde2213-1211-4501-b66b-09a75b40c11f",
    "created_at": "2023-04-27T03:44:42.843778+07:00",
    "updated_at": "2023-04-27T03:44:42.843778+07:00",
    "DeletedAt": null,
    "username": "test01",
    "password": "$2a$14$xhrlRSxbCgGi5Hfj7b/i1e6/UfQoXiY4/2mllEkzybY18cNo2qdbG",
    "tasks": null
  }
}
```

⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

# 📁 Collection: Tasks

## End-point: Create task

### Method: POST

> ```
> http://localhost:3000/api/tasks
> ```

### Body (**raw**)

```json
{
  "message": ""
}
```

### 🔑 Authentication bearer

| Param | value | Type   |
| ----- | ----- | ------ |
| token |       | string |

### Response: 200

```json
{
  "status": "success",
  "message": "task created",
  "task": {
    "id": "c245040a-5381-4c84-ac54-23b16c45c81c",
    "created_at": "2023-04-27T03:50:01.646019+07:00",
    "updated_at": "2023-04-27T03:50:01.646019+07:00",
    "deleted_at": null,
    "message": "",
    "completed": false,
    "completed_at": "0001-01-01T00:00:00Z",
    "user_id": "5bde2213-1211-4501-b66b-09a75b40c11f"
  }
}
```

⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: Get all tasks

### Method: GET

> ```
> http://localhost:3000/api/tasks
> ```

### 🔑 Authentication bearer

| Param | value | Type   |
| ----- | ----- | ------ |
| token |       | string |

### Response: 200

```json
{
  "status": "success",
  "tasks": []
}
```

⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: Update a task

### Method: PUT

> ```
> http://localhost:3000/api/tasks/
> ```

### Body (**raw**)

```json
{
  "message": "",
  "completed": true
}
```

### 🔑 Authentication bearer

| Param | value | Type   |
| ----- | ----- | ------ |
| token |       | string |

### Response: 200

```json
{
  "status": "success",
  "message": "task updated",
  "task": {
    "id": "f7b63290-f884-4a1d-97e0-9175fc134e83",
    "created_at": "2023-04-27T03:51:23.926139+07:00",
    "updated_at": "2023-04-27T03:51:47.363191+07:00",
    "deleted_at": null,
    "message": "Use the virtual IB array, then you can quantify the solid state card!",
    "completed": true,
    "completed_at": "0001-01-01T00:00:00Z",
    "user_id": "5bde2213-1211-4501-b66b-09a75b40c11f"
  }
}
```

⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: Delete a task

### Method: DELETE

> ```
> http://localhost:3000/api/tasks/
> ```

### 🔑 Authentication bearer

| Param | value | Type   |
| ----- | ----- | ------ |
| token |       | string |

### Response: 200

```json
{
  "status": "success",
  "message": "task deleted",
  "task": {
    "id": "f7b63290-f884-4a1d-97e0-9175fc134e83",
    "created_at": "2023-04-27T03:51:23.926139+07:00",
    "updated_at": "2023-04-27T03:51:47.363191+07:00",
    "deleted_at": "2023-04-27T03:52:26.689946+07:00",
    "message": "Use the virtual IB array, then you can quantify the solid state card!",
    "completed": true,
    "completed_at": "0001-01-01T06:42:04+06:42",
    "user_id": "5bde2213-1211-4501-b66b-09a75b40c11f"
  }
}
```

⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: Health check

### Method: GET

> ```
> http://localhost:3000/api/health
> ```

### Response: 200

```json
{
  "message": "API is running",
  "status": "success"
}
```

⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃
