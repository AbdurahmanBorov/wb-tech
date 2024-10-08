CREATE TABLE IF NOT EXISTS orders (
    order_uid           TEXT PRIMARY KEY,
    track_number        TEXT,
    entry               TEXT,
    locale              TEXT,
    internal_signature  TEXT,
    customer_id         TEXT,
    delivery_service    TEXT,
    shardkey            TEXT,
    sm_id               INT,
    date_created        TIMESTAMP,
    oof_shard           TEXT
);

CREATE TABLE IF NOT EXISTS delivery (
    order_uid   TEXT REFERENCES orders(order_uid),
    name        TEXT,
    phone       TEXT,
    zip         TEXT,
    city        TEXT,
    address     TEXT,
    region      TEXT,
    email       TEXT
);

CREATE TABLE IF NOT EXISTS payment (
    order_uid       TEXT REFERENCES orders(order_uid),
    transaction     TEXT,
    request_id      TEXT,
    currency        TEXT,
    provider        TEXT,
    amount          INT,
    payment_dt      BIGINT,
    bank            TEXT,
    delivery_cost   INT,
    goods_total     INT,
    custom_fee      INT
);

CREATE TABLE IF NOT EXISTS order_items (
    order_uid       TEXT REFERENCES orders(order_uid),
    chrt_id         BIGSERIAL PRIMARY KEY,
    track_number    TEXT,
    price           INT,
    rid             TEXT,
    name            TEXT,
    sale            INT,
    size            TEXT,
    total_price     INT,
    nm_id           INT,
    brand           TEXT,
    status          INT
);