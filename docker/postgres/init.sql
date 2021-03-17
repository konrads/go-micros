CREATE TABLE IF NOT EXISTS ports (
    id          text,
    name        text,
    coordinates float[],
    city        text,
    province    text,
    country     text,
    alias       text[],
    regions     text[],
    timezone    text,
    unlocs      text[],
    code        text,
    CONSTRAINT id_pk PRIMARY KEY(id)
);