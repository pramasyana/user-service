-- +goose Up
-- SQL in this section is executed when the migration is applied.

--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.3
-- Dumped by pg_dump version 9.6.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

--
-- Name: gender; Type: TYPE; Schema: public; Owner: bhinneka
--

CREATE TYPE gender AS ENUM (
    'MALE',
    'FEMALE'
);


ALTER TYPE gender OWNER TO bhinneka;

--
-- Name: userstatus; Type: TYPE; Schema: public; Owner: bhinneka
--

CREATE TYPE userstatus AS ENUM (
    'INACTIVE',
    'ACTIVE',
    'BLOCKED'
);


ALTER TYPE userstatus OWNER TO bhinneka;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: basicauth; Type: TABLE; Schema: public; Owner: bhinneka
--

CREATE TABLE basicauth (
    id integer NOT NULL,
    username character varying(50) NOT NULL,
    password character varying(255) NOT NULL,
    created timestamp with time zone DEFAULT now() NOT NULL,
    "lastModified" timestamp with time zone,
    version integer DEFAULT 1 NOT NULL
);


ALTER TABLE basicauth OWNER TO bhinneka;

--
-- Name: basicauth_id_seq; Type: SEQUENCE; Schema: public; Owner: bhinneka
--

CREATE SEQUENCE basicauth_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE basicauth_id_seq OWNER TO bhinneka;

--
-- Name: basicauth_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bhinneka
--

ALTER SEQUENCE basicauth_id_seq OWNED BY basicauth.id;


--
-- Name: city; Type: TABLE; Schema: public; Owner: bhinneka
--

CREATE TABLE city (
    id character varying NOT NULL,
    "provinceId" character varying NOT NULL,
    name character varying NOT NULL,
    created timestamp with time zone,
    "lastModified" timestamp with time zone,
    version integer
);


ALTER TABLE city OWNER TO bhinneka;

--
-- Name: district; Type: TABLE; Schema: public; Owner: bhinneka
--

CREATE TABLE district (
    id character varying NOT NULL,
    "cityId" character varying NOT NULL,
    name character varying NOT NULL,
    created timestamp with time zone,
    "lastModified" timestamp with time zone,
    version integer
);


ALTER TABLE district OWNER TO bhinneka;

--
-- Name: island; Type: TABLE; Schema: public; Owner: bhinneka
--

CREATE TABLE island (
    id character varying NOT NULL,
    name character varying NOT NULL,
    created timestamp with time zone,
    "lastModified" timestamp with time zone,
    version integer
);


ALTER TABLE island OWNER TO bhinneka;

--
-- Name: member; Type: TABLE; Schema: public; Owner: bhinneka
--

CREATE TABLE member (
    id character varying NOT NULL,
    "firstName" character varying(100) NOT NULL,
    "lastName" character varying(100),
    email character varying(100) NOT NULL,
    gender gender,
    mobile character varying,
    phone character varying(15),
    ext character varying(8),
    "birthDate" date,
    password character varying(255),
    salt character varying(255),
    province character varying(50),
    "provinceId" character varying(15),
    city character varying(50),
    "cityId" character varying(15),
    district character varying(50),
    "districtId" character varying(15),
    "subDistrict" character varying(50),
    "subDistrictId" character varying(15),
    "zipCode" character varying(6),
    address text,
    "jobTitle" character varying(100),
    department character varying(50),
    "facebookId" character varying(100),
    "googleId" character varying(100),
    "azureId" character varying(100),
    "isAdmin" boolean DEFAULT false,
    "isStaff" boolean DEFAULT false,
    "signUpFrom" character varying(15),
    status userstatus DEFAULT 'ACTIVE'::userstatus,
    token character varying(255),
    "lastLogin" timestamp with time zone,
    "lastBlocked" timestamp with time zone,
    created timestamp with time zone DEFAULT now(),
    "lastModified" timestamp with time zone,
    version integer DEFAULT 1 NOT NULL
);


ALTER TABLE member OWNER TO bhinneka;

--
-- Name: province; Type: TABLE; Schema: public; Owner: bhinneka
--

CREATE TABLE province (
    id character varying NOT NULL,
    "islandId" character varying NOT NULL,
    name character varying NOT NULL,
    created timestamp with time zone,
    "lastModified" timestamp with time zone,
    version integer
);


ALTER TABLE province OWNER TO bhinneka;

--
-- Name: subdistrict; Type: TABLE; Schema: public; Owner: bhinneka
--

CREATE TABLE subdistrict (
    id character varying NOT NULL,
    "districtId" character varying NOT NULL,
    name character varying NOT NULL,
    created timestamp with time zone,
    "lastModified" timestamp with time zone,
    version integer
);


ALTER TABLE subdistrict OWNER TO bhinneka;

--
-- Name: zipcode; Type: TABLE; Schema: public; Owner: bhinneka
--

CREATE TABLE zipcode (
    id bigint NOT NULL,
    code character varying NOT NULL,
    "cityId" character varying NOT NULL,
    "districtName" character varying NOT NULL,
    "subDistrictName" character varying NOT NULL,
    created timestamp with time zone,
    "lastModified" timestamp with time zone,
    version integer
);


ALTER TABLE zipcode OWNER TO bhinneka;

--
-- Name: zipcode_id_seq; Type: SEQUENCE; Schema: public; Owner: bhinneka
--

CREATE SEQUENCE zipcode_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE zipcode_id_seq OWNER TO bhinneka;

--
-- Name: zipcode_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: bhinneka
--

ALTER SEQUENCE zipcode_id_seq OWNED BY zipcode.id;


--
-- Name: basicauth id; Type: DEFAULT; Schema: public; Owner: bhinneka
--

ALTER TABLE ONLY basicauth ALTER COLUMN id SET DEFAULT nextval('basicauth_id_seq'::regclass);


--
-- Name: zipcode id; Type: DEFAULT; Schema: public; Owner: bhinneka
--

ALTER TABLE ONLY zipcode ALTER COLUMN id SET DEFAULT nextval('zipcode_id_seq'::regclass);


--
-- Name: basicauth basicauth_pkey; Type: CONSTRAINT; Schema: public; Owner: bhinneka
--

ALTER TABLE ONLY basicauth
    ADD CONSTRAINT basicauth_pkey PRIMARY KEY (id);


--
-- Name: city city_pkey; Type: CONSTRAINT; Schema: public; Owner: bhinneka
--

ALTER TABLE ONLY city
    ADD CONSTRAINT city_pkey PRIMARY KEY (id);


--
-- Name: district district_pkey; Type: CONSTRAINT; Schema: public; Owner: bhinneka
--

ALTER TABLE ONLY district
    ADD CONSTRAINT district_pkey PRIMARY KEY (id);


--
-- Name: island island_pkey; Type: CONSTRAINT; Schema: public; Owner: bhinneka
--

ALTER TABLE ONLY island
    ADD CONSTRAINT island_pkey PRIMARY KEY (id);


--
-- Name: member member_pkey; Type: CONSTRAINT; Schema: public; Owner: bhinneka
--

ALTER TABLE ONLY member
    ADD CONSTRAINT member_pkey PRIMARY KEY (id);


--
-- Name: province province_pkey; Type: CONSTRAINT; Schema: public; Owner: bhinneka
--

ALTER TABLE ONLY province
    ADD CONSTRAINT province_pkey PRIMARY KEY (id);


--
-- Name: subdistrict subdistrict_pkey; Type: CONSTRAINT; Schema: public; Owner: bhinneka
--

ALTER TABLE ONLY subdistrict
    ADD CONSTRAINT subdistrict_pkey PRIMARY KEY (id);


--
-- Name: zipcode zipcode_pkey; Type: CONSTRAINT; Schema: public; Owner: bhinneka
--

ALTER TABLE ONLY zipcode
    ADD CONSTRAINT zipcode_pkey PRIMARY KEY (id);


--
-- Name: member_unique; Type: INDEX; Schema: public; Owner: bhinneka
--

CREATE UNIQUE INDEX member_unique ON member USING btree (email);


--
-- PostgreSQL database dump complete
--

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.