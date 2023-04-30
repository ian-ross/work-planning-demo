1. Create user:

```
sudo adduser --system --no-create-home --disabled-login planning-demo
```

2. Create database:

```
createdb -h localhost -U postgres planning_demo
```

3. Create database user:

```
psql -h localhost -U postgres
...
postgres=# CREATE USER planning_demo WITH PASSWORD '<INSERT-PASSWORD>';
CREATE ROLE
postgres=# GRANT ALL ON DATABASE planning_demo TO planning_demo;
GRANT
```

4. Copy deployment files to server:

```
/etc/caddy/conf.d/planning-demo.skybluetrades.net
/etc/planning-demo.env
/etc/systemd/system/planning-demo.service
work-planning-demo executable => /opt/planning-demo/planning-demo
```

 - Fix domain in `/etc/caddy/conf.d/planning-demo.skybluetrades.net`
   (and set up Caddy generally if it's not being used).
 - Update Postgres user and password and authentication key in
   `/etc/planning-demo.env` (use the Postgres password set up above
   and make a random string for the authentication key).
 - Make sure the executable is executable by the `planning-demo` user.
 