-- +goose Up
-- SQL in this section is executed when the migration is applied.
DROP TABLE IF EXISTS log_dolphin;

CREATE TABLE log_dolphin (
    id integer NOT NULL,
    user_id character varying(20) NOT NULL,
    event_type character varying(20) NOT NULL,
    log_data jsonb,
    created timestamp without time zone DEFAULT now() NOT NULL
);

CREATE SEQUENCE log_dolphin_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE log_dolphin_id_seq OWNED BY log_dolphin.id;

ALTER TABLE ONLY log_dolphin ALTER COLUMN id SET DEFAULT nextval('log_dolphin_id_seq'::regclass);

ALTER TABLE ONLY log_dolphin
    ADD CONSTRAINT log_dolphin_pkey PRIMARY KEY (id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS log_dolphin;
