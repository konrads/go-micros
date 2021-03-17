CREATE TABLE IF NOT EXISTS star (
    id                text,
    name              text,
    alias             text[],
    constellation     text,
    coordinates       float[],
    distance          float,
    apparentMagnitude float,
    CONSTRAINT id_pk PRIMARY KEY(id)
);