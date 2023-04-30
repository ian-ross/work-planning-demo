
-- +migrate Up

CREATE TABLE IF NOT EXISTS worker (
  id        SERIAL   PRIMARY KEY,
  email     TEXT     NOT NULL,
  name      TEXT     NOT NULL,
  is_admin  BOOLEAN  NOT NULL,
  password  TEXT     NOT NULL
);

CREATE INDEX worker_email_idx ON worker(email);


CREATE TABLE IF NOT EXISTS shift (
  id          SERIAL       PRIMARY KEY,
  start_time  TIMESTAMPTZ  NOT NULL,
  end_time    TIMESTAMPTZ  NOT NULL,
  capacity    INTEGER      NOT NULL
);

CREATE INDEX shift_start_time_idx ON shift(start_time);


CREATE TABLE IF NOT EXISTS shift_assignment (
  worker_id  INTEGER  NOT NULL REFERENCES worker(id) ON DELETE CASCADE,
  shift_id   INTEGER  NOT NULL REFERENCES shift(id) ON DELETE CASCADE,

  CONSTRAINT shift_assignment_unique UNIQUE (worker_id, shift_id)
);

CREATE INDEX shift_assignment_worker_idx ON shift_assignment(worker_id);
CREATE INDEX shift_assignment_shift_idx ON shift_assignment(shift_id);


-- +migrate Down

DROP TABLE IF EXISTS shift_assignment;
DROP TABLE IF EXISTS shift;
DROP TABLE IF EXISTS worker;
