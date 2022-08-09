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

SET search_path = public, pg_catalog;

--
-- Data for Name: basicauth; Type: TABLE DATA; Schema: public; Owner: bhinneka
--

INSERT INTO basicauth (id, username, password, created, "lastModified", version) VALUES (1, 'bhinneka-microservices-b13714-5312115', '626869-6e6e65-6b6120-6d656e-746172-692064-696d656e-736900', '2018-01-26 17:54:04.630235+07', NULL, 1);


--
-- Name: basicauth_id_seq; Type: SEQUENCE SET; Schema: public; Owner: bhinneka
--

SELECT pg_catalog.setval('basicauth_id_seq', 1, true);


--
-- PostgreSQL database dump complete
--

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
