-- +goose Up
-- SQL in this section is executed when the migration is applied.
DROP TABLE IF EXISTS client_app;

CREATE TABLE client_app (
    id integer NOT NULL,
    "clientId" character varying(50) NOT NULL,
    "clientSecret" character varying(100) NOT NULL,
    name character varying(20) NOT NULL,
    status userstatus DEFAULT 'ACTIVE'::userstatus,
    created timestamp without time zone DEFAULT now() NOT NULL,
    "lastModified" timestamp with time zone,
    version integer DEFAULT 1 NOT NULL
);

CREATE SEQUENCE client_app_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE client_app_id_seq OWNED BY client_app.id;

ALTER TABLE ONLY client_app ALTER COLUMN id SET DEFAULT nextval('client_app_id_seq'::regclass);

ALTER TABLE ONLY client_app
    ADD CONSTRAINT client_app_pkey PRIMARY KEY (id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS client_app;
