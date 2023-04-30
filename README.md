# Simple REST API service for shift planning

For now, this is just a simple REST API implementation (in Go) for a
basic work planning service. I might turn it into a platform for
experimenting with some of the heuristic optimization methods
described in Michalewicz & Fogel's *How To Solve It: Modern
Heuristics*, which I've been reading recently as background for
another project.

Implemented things:

 - JWT authentication with refresh tokens.
 - Pluggable data store interface.
 - In-memory data store for development (use `STORE_URL=memory`).
 - PostgreSQL data store including embedded migrations (use
   `STORE_URL=postgres://whatever`).
 - Echo-based server derived from OpenAPI API specification.
 - OpenAPI documentation using Redoc.

**TODO**:
 - Postgres migrations
 - Postgres store methods
 - A few tests
 
## The basic application

### Requirements

Here are the barebones requirements:
 
 - A **worker** has **shifts**.
 - A **shift** is 8 hours long.
 - A **worker** never has two **shifts** on the same day.
 - It is a 24 hour timetable 0-8, 8-16, 16-24.

### Data modelling description

We have **workers** and **shifts**, so we'll probably have a `Worker`
model and a `Shift` model, both of which will have be represented in
the database. **Workers** can be **assigned** to **shifts**, so we'll
probably also have a join table with a model called `ShiftAssignment`
recording the fact that a **worker** is working a particular
**shift**. (We're obviously going to need a join table here because
the **worker**/**shift** relationship is many-to-many: **workers** can
be assigned to multiple **shifts** and any **shift** may have more
than one **worker**.)

*Assumption*: A **worker** maps one-to-one to users of the
application, so the `Worker` model will have login details associated
with it.

*Design decision*: We'll use a simple username/password setup for
login, and will use JWTs for authorization for most API routes (using
the usual access token + refresh token approach).

*Design decision*: We'll make a distinction between admin and
non-admin users, just with a flag in the `Worker` model. The API
endpoints for creating, modifying and deleting existing **workers**
will be accessible only to admin users.

*Assumption*: Any **shift** has a maximmum number of **workers** that
can be assigned to it, which we'll call the **shift**'s "capacity".

*Assumption*: **Shifts** are created "manually" by admin users.




### Models

```go
type Worker struct {
    ID         int
    Email      string
    Password   string
    Name       string
    IsAdmin    bool
    DeletedAt  *time.Time
}
```

```go
type Shift struct {
    ID int
    Day *date.Date
    StartTime *time.Time
    EndTime *time.Time
    Capacity uint
}
```



### Endpoints

```
  POST /auth/login
  {"email": "x@y.com", "password": "blah"} => 200 {"access_token": "...", "refresh_token": "..."}
  => 403 {"message": "Login failed"}
  
  POST /auth/logout
  => 204

  POST /auth/refresh_token
  {"refresh_token": "..."} => 200 {"access_token": "...", "refresh_token": "..."}
  => 403 {"message": "Invalid refresh token"}

All the other routes return 404s for unauthorized users.

Use "Authorization: Bearer <access_token>" header for authentication.

  GET /me
  => 200 {"id": 123, "email": "x@y.com", "name": "Blah", "is_admin": false}

  GET /schedule?date=YYYY-MM-DD&span={week|day}

GET /workers      (admin only)
  => 200 [{id, email, name, is_admin}, ...]
  
GET /workers/:id  (admin only)
  => 200 {id, email, name, is_admin}
  => 404 {"message": "Unknown user ID"}
  
POST /workers      (admin only)
  {"email": "x@y.com", "name": "Blah", "password": "blah", "is_admin": false}
    => 200 {"id": 234, "email": "x@y.com", "name": "Blah", "password": "blah", "is_admin": false}

PUT /workers  (admin only)
  {"id": 234, "email": "x@y.com", "name": "Blah", "password": "new", "is_admin": true}
    => 200 {"id": 234, "email": "x@y.com", "name": "Blah", "password": "new", "is_admin": true}

DELETE /workers/:id  (admin only)
  => 204


POST /shifts (admin only)
  {"day": "YYYY-MM-DD", start_time: "HH:MM", "end_time": "HH:MM", "capacity": 1}
    => 200 {"id": 456, "day": "YYYY-MM-DD", start_time: "HH:MM", "end_time": "HH:MM", "capacity": 1}

GET /shifts?date=YYYY-MM-DD&span={week|day}
  (date defaults to today, span defaults to "week")
  
GET /shifts/:id
  => 200 {"id": 456, "day": "YYYY-MM-DD", start_time: "HH:MM", "end_time": "HH:MM", "capacity": 1}
  => 404 {"message": "Unknown shift ID"}

PUT /shifts (admin only)
  => {}

DELETE /shifts/:id (admin only)
  => 204
  => 404 {"message": "Unknown shift ID"}


POST /shifts/:id/assignment
  => 200 {"date": "YYYY-MM-DD", start_time: "HH:MM", end_time: "HH:MM"}
  => 404 {"message": "Unknown shift ID"}
  => 400 {"message": "Shift has no capacity"}

DELETE /shifts/:id/assignment
  => 204
  => 404 {"message": "Unknown shift ID"}
  => 404 {"message": "User has no assignment for this shift"}
```



# Other things to look at

 - `sqlx` ⇒ `pgx` + `pgxscan`?
 - Swagger UI instead of Redoc?
 
 
# Database setup

## In-memory database

Set `STORE_URL=memory` in `.env`

## PostgreSQL database

Assuming you have a local PostgreSQL instance running... Run `psql` as
user `postgres` (`psql -U postgres`) and do the following:

Create a database and connect to it:

```
postgres=# CREATE DATABASE planning_dev;
CREATE DATABASE
postgres=# \c planning_dev
You are now connected to database "planning_dev" as user "postgres".
```

Create a user and grant them permissions on the database:

```
planning_dev=# CREATE USER planning_dev WITH PASSWORD '<some-password>';
CREATE ROLE
planning_dev=# GRANT ALL ON DATABASE planning_dev TO planning_dev;
GRANT
```

Now put the following `STORE_URL` setting in your `.env`:

```
STORE_URL=postgres://planning_dev:<some-password>@localhost:5432/planning_dev?sslmode=disable
```

where `<some-password>` is the password you used for the
`planning_dev` user.
