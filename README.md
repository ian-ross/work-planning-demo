# Simple REST API service for shift planning

For now, this is just a simple REST API implementation (in Go) for a
basic work planning service. I might turn it into a platform for
experimenting with some of the heuristic optimization methods
described in Michalewicz & Fogel's [*How To Solve It: Modern
Heuristics*](https://www.amazon.de/-/en/Zbigniew-Michalewicz/dp/3642061346/),
which I've been reading recently as background for another project.

Implemented things:

 - [Echo](https://echo.labstack.com/)-based server derived from
   OpenAPI API specification using
   [oapi-codegen](https://github.com/deepmap/oapi-codegen).
 - JWT authentication with refresh tokens.
 - Pluggable data store interface.
 - In-memory data store for development (use `STORE_URL=memory`).
 - PostgreSQL data store including embedded migrations (use
   `STORE_URL=postgres://whatever`).
 - OpenAPI documentation using Redoc.
 - Some tests (just for the login flow and authentication middleware
   so far).

Missing things:

 - An actual usable API for the problem! It should be possible for
   workers to specify preferences for shifts and for an admin to
   generate a feasible schedule that covers all the required shifts.
   There should also be endpoints to allow admin users to reassign
   shifts, remove shift assignments, reschedule for cases of illness,
   set up other ruled for scheduling, etc.
 - More tests. I'm sure there are things that don't quite work, just
   because I knocked this together pretty quickly.
 - Scheduling algorithms. Some of these kinds of problems have
   polynomial-time algorithms for finding admissible schedules, but
   it's become clear to me from reading Michalewicz & Fogel that as
   soon as you move from simple problems to more realistic and
   interesting problems, where you usually have extra constraints and
   maybe some kind of measure of satisfaction for the final schedule,
   you're going to end up going to have to use some sort of heuristic
   search. Whether that's simulated annealing or some kind of
   evolutionary search depends on the problem details, but it might be
   interesting to experiment with some options.
 - Maybe try a different Go SQL library? I usually use `sqlx`, but a
   lot of people seem to like `pgx` (probably combined with `pgxscan`
   to get something like the nice scanning behavior of `sqlx`). Or try
   an ORM? I've played with `gorm` a bit, but I'm not super keen on
   it.
 
**TODO**:
 - Swagger docs

## Installation requirements

 - Go version: 1.20: it might work with earlier versions, but no
   promises.
 - Mockery v2.26.1 (from
   [here](https://github.com/vektra/mockery/releases/tag/v2.26.1)):
   it's recommended that you install a binary version from that
   release page. Using `go install` might cause you some problems (see
   the Mockery docs for details).

### Database setup

#### In-memory database

Set `STORE_URL=memory` in `.env`

#### PostgreSQL database

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

----

## The basic application

### Requirements

Here are the barebones requirements from the original problem
definition:
 
 - A **worker** has **shifts**.
 - A **shift** is 8 hours long.
 - A **worker** never has two **shifts** on the same day.
 - It is a 24 hour timetable 0-8, 8-16, 16-24.

Starting from there, I slightly went to town on this, because I've
been thinking about scheduling problems, and having a little platform
to experiment with them seemed like it might be useful.

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

OK, that's probably enough to get going.

### Models

Something like this...

Workers:

```go
type WorkerID int64

type Worker struct {
	ID       WorkerID `db:"id"`
	Email    string   `db:"email"`
	Name     string   `db:"name"`
	IsAdmin  bool     `db:"is_admin"`
	Password string   `db:"password"`
}
```

Shifts:

```go
type ShiftID int64

type Shift struct {
	ID        ShiftID   `db:"id"`
	StartTime time.Time `db:"start_time"`
	EndTime   time.Time `db:"end_time"`
	Capacity  int       `db:"capacity"`
}
```

Shift assignments:

```go
type ShiftAssignment struct {
	Worker WorkerID `db:"worker_id"`
	Shift  ShiftID  `db:"shift_id"`
}
```

### Endpoints

Rough notes. At some point I stopped with this and switched over to
working on [the OpenAPI
spec](https://github.com/ian-ross/work-planning-demo/tree/main/spec/openapi.yaml).
  
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
