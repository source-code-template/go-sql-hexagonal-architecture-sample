# go-sql-layer-architecture-sample

#### To run the application
```shell
go run main.go
```

## How to make source code cleaner
- This repository is forked from https://github.com/go-tutorials/go-sql-hexagonal-architecture-sample
- The original repository is written at GO SDK level. It mean we develop the micro service, using GO SDK, without any utility

### How to make sql database adapter cleaner
#### Query data
<table><thead><tr><td>

[GO SDK Only](https://github.com/go-tutorials/go-sql-hexagonal-architecture-sample/blob/main/internal/user/adapter/repository/user_adapter.go)
</td><td>

[GO SDK with utilities](https://github.com/source-code-template/go-sql-hexagonal-architecture-sample/blob/main/internal/user/adapter/repository/user_adapter.go)
</td></tr></thead><tbody><tr><td>

```go
func (r *UserAdapter) Load(ctx context.Context, id string) (*domain.User, error) {
	query := `
		select
			id, 
			username,
			email,
			phone,
			date_of_birth
		from users where id = ?`
	rows, err := r.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user domain.User
		err = rows.Scan(
			&user.Id,
			&user.Username,
			&user.Phone,
			&user.Email,
			&user.DateOfBirth)
		return &user, nil
	}
	return nil, nil
}
```
</td>
<td>

```go
import 	q "github.com/core-go/sql"

func (r *UserAdapter) Load(ctx context.Context, id string) (*User, error) {
	var users []User
	query := fmt.Sprintf(`
		select
			id,
			username,
			email,
			phone,
			date_of_birth
		from users where id = %s limit 1`, q.BuildParam(1))
	err := q.Select(ctx, r.DB, &users, query, id)
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return &users[0], nil
	}
	return nil, nil
}
```
</td></tr></tbody></table>

#### Execute query
<table><thead><tr><td>

[GO SDK Only](https://github.com/go-tutorials/go-sql-hexagonal-architecture-sample/blob/main/internal/user/adapter/repository/user_adapter.go)
</td><td>

[GO SDK with utilities](https://github.com/source-code-template/go-sql-hexagonal-architecture-sample/blob/main/internal/user/adapter/repository/user_adapter.go)
</td></tr></thead><tbody><tr><td>

```go
func (r *UserAdapter) Create(ctx context.Context, user *domain.User) (int64, error) {
	query := `
		insert into users (
			id,
			username,
			email,
			phone,
			date_of_birth)
		values (
			?,
			?,
			?, 
			?,
			?)`
	tx := GetTx(ctx)
	stmt, err := tx.Prepare(query)
	if err != nil {
		return -1, err
	}
	res, err := stmt.ExecContext(ctx,
		user.Id,
		user.Username,
		user.Email,
		user.Phone,
		user.DateOfBirth)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}
```
</td>
<td>

```go
import 	q "github.com/core-go/sql"

func (r *UserAdapter) Create(ctx context.Context, user *User) (int64, error) {
	query, args := q.BuildToInsert("users", user, q.BuildParam)
	tx := q.GetTx(ctx)
	res, err := tx.ExecContext(ctx, query, args...)
	return q.RowsAffected(res, err)
}
```
</td></tr></tbody></table>

### How to make service cleaner
<table><thead><tr><td>

[GO SDK Only](https://github.com/go-tutorials/go-sql-hexagonal-architecture-sample/blob/main/internal/user/service/user_service.go)
</td><td>

[GO SDK with utilities](https://github.com/source-code-template/go-sql-hexagonal-architecture-sample/blob/main/internal/user/service/user_service.go)
</td></tr></thead><tbody><tr><td>

```go
func (s *userService) Create(ctx context.Context, user *User) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return -1, nil
	}
	ctx = context.WithValue(ctx, "tx", tx)
	res, err := s.repository.Create(ctx, user)
	if err != nil {
		er := tx.Rollback()
		if er != nil {
			return -1, er
		}
		return -1, err
	}
	err = tx.Commit()
	return res, err
}
```
</td>
<td>

```go
func (s *userService) Create(ctx context.Context, user *User) (int64, error) {
	ctx, tx, err := q.Begin(ctx, s.db)
	if err != nil {
		return  -1, err
	}
	res, err := s.repository.Create(ctx, user)
	return q.End(tx, res, err)
}
```
</td></tr></tbody></table>

### How to make http handler cleaner
#### Get data
<table><thead><tr><td>

[GO SDK and Mux](https://github.com/go-tutorials/go-sql-hexagonal-architecture-sample/blob/main/internal/user/adapter/handler/user_handler.go)
</td><td>

[GO SDK with utilities and data validation](https://github.com/source-code-template/go-sql-hexagonal-architecture-sample/blob/main/internal/user/adapter/handler/user_handler.go)
</td></tr></thead><tbody><tr><td>

```go
func (h *HttpUserHandler) Load(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}

	user, err := h.service.Load(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, http.StatusOK, user)
}
```
</td>
<td>

```go
func (h *HttpUserHandler) Load(w http.ResponseWriter, r *http.Request) {
	id := sv.GetRequiredParam(w, r)
	if len(id) > 0 {
		res, err := h.service.Load(r.Context(), id)
		sv.RespondModel(w, r, res, err, h.Error, nil)
	}
}
```
</td></tr></tbody></table>

#### Create data 
<table><thead><tr><td>

[GO SDK and Mux without data validation](https://github.com/go-tutorials/go-sql-hexagonal-architecture-sample/blob/main/internal/user/adapter/handler/user_handler.go)
</td><td>

[GO SDK with utilities and data validation](https://github.com/source-code-template/go-sql-hexagonal-architecture-sample/blob/main/internal/user/adapter/handler/user_handler.go)
</td></tr></thead><tbody><tr><td>

```go
func (h *HttpUserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user User
	er1 := json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
		return
	}

	res, er2 := h.service.Create(r.Context(), &user)
	if er2 != nil {
		http.Error(w, er1.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, http.StatusCreated, res)
}
```
</td>
<td>

```go
func (h *HttpUserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user User
	er1 := sv.Decode(w, r, &user)
	if er1 == nil {
		errors, er2 := h.Validate(r.Context(), &user)
		if !sv.HasError(w, r, errors, er2, *h.Status.ValidationError, h.Error, h.Log, h.Resource, h.Action.Create) {
			res, er3 := h.service.Create(r.Context(), &user)
			sv.AfterCreated(w, r, &user, res, er3, h.Status, h.Error, h.Log, h.Resource, h.Action.Create)
		}
	}
}
```
</td></tr></tbody></table>

#### [core-go/search](https://github.com/core-go/search)
- Build the search model at http handler
- Build dynamic SQL for search
  - Build SQL for paging by page index (page) and page size (limit)
  - Build SQL to count total of records
### Search users: Support both GET and POST 
#### POST /users/search
##### *Request:* POST /users/search
In the below sample, search users with these criteria:
- get users of page "1", with page size "20"
- email="tony": get users with email starting with "tony"
- dateOfBirth between "min" and "max" (between 1953-11-16 and 1976-11-16)
- sort by phone ascending, id descending
```json
{
    "page": 1,
    "limit": 20,
    "sort": "phone,-id",
    "email": "tony",
    "dateOfBirth": {
        "min": "1953-11-16T00:00:00+07:00",
        "max": "1976-11-16T00:00:00+07:00"
    }
}
```
##### GET /users/search?page=1&limit=2&email=tony&dateOfBirth.min=1953-11-16T00:00:00+07:00&dateOfBirth.max=1976-11-16T00:00:00+07:00&sort=phone,-id
In this sample, search users with these criteria:
- get users of page "1", with page size "20"
- email="tony": get users with email starting with "tony"
- dateOfBirth between "min" and "max" (between 1953-11-16 and 1976-11-16)
- sort by phone ascending, id descending

#### *Response:*
- total: total of users, which is used to calculate numbers of pages at client 
- list: list of users
```json
{
    "list": [
        {
            "id": "ironman",
            "username": "tony.stark",
            "email": "tony.stark@gmail.com",
            "phone": "0987654321",
            "dateOfBirth": "1963-03-24T17:00:00Z"
        }
    ],
    "total": 1
}
```

## API Design
### Common HTTP methods
- GET: retrieve a representation of the resource
- POST: create a new resource
- PUT: update the resource
- PATCH: perform a partial update of a resource, refer to [service](https://github.com/core-go/service) and [sql](https://github.com/core-go/sql)  
- DELETE: delete a resource

## API design for health check
To check if the service is available.
#### *Request:* GET /health
#### *Response:*
```json
{
    "status": "UP",
    "details": {
        "sql": {
            "status": "UP"
        }
    }
}
```

## API design for users
#### *Resource:* users

### Get all users
#### *Request:* GET /users
#### *Response:*
```json
[
    {
        "id": "spiderman",
        "username": "peter.parker",
        "email": "peter.parker@gmail.com",
        "phone": "0987654321",
        "dateOfBirth": "1962-08-25T16:59:59.999Z"
    },
    {
        "id": "wolverine",
        "username": "james.howlett",
        "email": "james.howlett@gmail.com",
        "phone": "0987654321",
        "dateOfBirth": "1974-11-16T16:59:59.999Z"
    }
]
```

### Get one user by id
#### *Request:* GET /users/:id
```shell
GET /users/wolverine
```
#### *Response:*
```json
{
    "id": "wolverine",
    "username": "james.howlett",
    "email": "james.howlett@gmail.com",
    "phone": "0987654321",
    "dateOfBirth": "1974-11-16T16:59:59.999Z"
}
```

### Create a new user
#### *Request:* POST /users 
```json
{
    "id": "wolverine",
    "username": "james.howlett",
    "email": "james.howlett@gmail.com",
    "phone": "0987654321",
    "dateOfBirth": "1974-11-16T16:59:59.999Z"
}
```
#### *Response:*
- status: configurable; 1: success, 0: duplicate key, 4: error
```json
{
    "status": 1,
    "value": {
        "id": "wolverine",
        "username": "james.howlett",
        "email": "james.howlett@gmail.com",
        "phone": "0987654321",
        "dateOfBirth": "1974-11-16T00:00:00+07:00"
    }
}
```
#### *Fail case sample:* 
- Request:
```json
{
    "id": "wolverine",
    "username": "james.howlett",
    "email": "james.howlett",
    "phone": "0987654321a",
    "dateOfBirth": "1974-11-16T16:59:59.999Z"
}
```
- Response: in this below sample, email and phone are not valid
```json
{
    "status": 4,
    "errors": [
        {
            "field": "email",
            "code": "email"
        },
        {
            "field": "phone",
            "code": "phone"
        }
    ]
}
```

### Update one user by id
#### *Request:* PUT /users/:id
```shell
PUT /users/wolverine
```
```json
{
    "username": "james.howlett",
    "email": "james.howlett@gmail.com",
    "phone": "0987654321",
    "dateOfBirth": "1974-11-16T16:59:59.999Z"
}
```
#### *Response:*
- status: configurable; 1: success, 0: duplicate key, 2: version error, 4: error
```json
{
    "status": 1,
    "value": {
        "id": "wolverine",
        "username": "james.howlett",
        "email": "james.howlett@gmail.com",
        "phone": "0987654321",
        "dateOfBirth": "1974-11-16T00:00:00+07:00"
    }
}
```

### Patch one user by id
Perform a partial update of user. For example, if you want to update 2 fields: email and phone, you can send the request body of below.
#### *Request:* PATCH /users/:id
```shell
PATCH /users/wolverine
```
```json
{
    "email": "james.howlett@gmail.com",
    "phone": "0987654321"
}
```
#### *Response:*
- status: configurable; 1: success, 0: duplicate key, 2: version error, 4: error
```json
{
    "status": 1,
    "value": {
        "email": "james.howlett@gmail.com",
        "phone": "0987654321"
    }
}
```

#### Problems for patch
If we pass a struct as a parameter, we cannot control what fields we need to update. So, we must pass a map as a parameter.
```go
type UserService interface {
    Update(ctx context.Context, user *User) (int64, error)
    Patch(ctx context.Context, user map[string]interface{}) (int64, error)
}
```
We must solve 2 problems:
1. At http handler layer, we must convert the user struct to map, with json format, and make sure the nested data types are passed correctly.
2. At repository layer, from json format, we must convert the json format to database column name

#### Solutions for patch  
At http handler layer, we use [core-go/service](https://github.com/core-go/service), to convert the user struct to map, to make sure we just update the fields we need to update
```go
import server "github.com/core-go/service"

func (h *UserHandler) Patch(w http.ResponseWriter, r *http.Request) {
    var user User
    userType := reflect.TypeOf(user)
    _, jsonMap := sv.BuildMapField(userType)
    body, _ := sv.BuildMapAndStruct(r, &user)
    json, er1 := sv.BodyToJson(r, user, body, ids, jsonMap, nil)

    result, er2 := h.service.Patch(r.Context(), json)
    if er2 != nil {
        http.Error(w, er2.Error(), http.StatusInternalServerError)
        return
    }
    respond(w, result)
}
```

### Delete a new user by id
#### *Request:* DELETE /users/:id
```shell
DELETE /users/wolverine
```
#### *Response:* 1: success, 0: not found, -1: error
```json
1
```

## Common libraries
- [core-go/health](https://github.com/core-go/health): include HealthHandler, HealthChecker, SqlHealthChecker
- [core-go/config](https://github.com/core-go/config): to load the config file, and merge with other environments (SIT, UAT, ENV)
- [core-go/log](https://github.com/core-go/log): log and log middleware

### core-go/health
To check if the service is available, refer to [core-go/health](https://github.com/core-go/health)
#### *Request:* GET /health
#### *Response:*
```json
{
    "status": "UP",
    "details": {
        "sql": {
            "status": "UP"
        }
    }
}
```
To create health checker, and health handler
```go
    db, err := sql.Open(conf.Driver, conf.DataSourceName)
    if err != nil {
        return nil, err
    }

    sqlChecker := s.NewSqlHealthChecker(db)
    healthHandler := health.NewHealthHandler(sqlChecker)
```

To handler routing
```go
    r := mux.NewRouter()
    r.HandleFunc("/health", healthHandler.Check).Methods("GET")
```

### core-go/config
To load the config from "config.yml", in "configs" folder
```go
package main

import "github.com/core-go/config"

type Root struct {
    DB DatabaseConfig `mapstructure:"db"`
}

type DatabaseConfig struct {
    Driver         string `mapstructure:"driver"`
    DataSourceName string `mapstructure:"data_source_name"`
}

func main() {
    var conf Root
    err := config.Load(&conf, "configs/config")
    if err != nil {
        panic(err)
    }
}
```

### core-go/log *&* core-go/middleware
```go
import (
    "github.com/core-go/config"
    "github.com/core-go/log"
    m "github.com/core-go/middleware"
    "github.com/gorilla/mux"
)

func main() {
    var conf app.Root
    config.Load(&conf, "configs/config")

    r := mux.NewRouter()

    log.Initialize(conf.Log)
    r.Use(m.BuildContext)
    logger := m.NewLogger()
    r.Use(m.Logger(conf.MiddleWare, log.InfoFields, logger))
    r.Use(m.Recover(log.ErrorMsg))
}
```
To configure to ignore the health check, use "skips":
```yaml
middleware:
  skips: /health
```
