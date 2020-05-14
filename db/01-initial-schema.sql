CREATE TABLE pumps (
    pump_id BIGSERIAL PRIMARY KEY,
    name    TEXT     NOT NULL,
    pump     SMALLINT NOT NULL
);

CREATE TABLE cycles (
    cycle_id   BIGSERIAL PRIMARY KEY,
    pump_id    BIGINT    NOT NULL REFERENCES pumps (pump_id),
    start_time TIMESTAMP NOT NULL,
    end_time   TIMESTAMP NULL,
    notified   BOOLEAN   NOT NULL
);