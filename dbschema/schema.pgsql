-- vim: ft=sql:ts=4:sw=4:et
-- TRY5 DB SCHEMA
--
-- (c) ACEB,SAU 2015

-------------------------------------------------------
--                                                   --
-- Instala el esquema por defecto para la base de    --
-- datos Try5                                        --
--                                                   --
-- IMPORTANTE! DEBE EJECUTARSE COMO USUARIO postgres --
--     su - postgres -c "psql < schema.pgsql"        --
-------------------------------------------------------

SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = off;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET escape_string_warning = off;

CREATE DATABASE try5db WITH ENCODING = 'UTF8';
ALTER DATABASE try5db OWNER TO try5adm;

\connect try5db

BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;
COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';
SET search_path = public, pg_catalog;
SET default_tablespace = '';
SET default_with_oids = false;

-- ----------------------------
--  Table structure for "users"
-- ----------------------------
CREATE TABLE IF NOT EXISTS users (
    id        SERIAL,
    uid       VARCHAR(36),
    email     VARCHAR(100),
    name      VARCHAR(200),
    password  VARCHAR(60),
    active    BOOLEAN,
    gravatar  VARCHAR(60),
    created   TIMESTAMP NOT NULL DEFAULT NOW(),
    updated   TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted   TIMESTAMP,

    CONSTRAINT users_pkey PRIMARY KEY (id)
)
WITH (OIDS=FALSE);
ALTER TABLE users OWNER TO try5adm;
CREATE INDEX user_idx ON users USING btree (id);
CREATE INDEX user_email_idx ON users USING btree (email);

CREATE TABLE rbac_role (
    id SERIAL NOT NULL PRIMARY KEY,
    slug VARCHAR(256) UNIQUE NOT NULL,
    name VARCHAR(256),
    description TEXT DEFAULT '',
    parameters JSONB DEFAULT '[]',
    created TIMESTAMP  NOT NULL DEFAULT NOW(),
    updated TIMESTAMP  NOT NULL
)
WITH (OIDS=FALSE);
ALTER TABLE public.rbac_role OWNER TO try5adm;
CREATE INDEX rbac_role_idx ON rbac_role USING btree (id, name);

CREATE TABLE rbac_grant (
    id SERIAL NOT NULL PRIMARY KEY,
    from_role INT,
    to_role INT,
    assigment JSONB NOT NULL DEFAULT '{}',

    CONSTRAINT memberships_granted_fkey
        FOREIGN KEY (from_role)
        REFERENCES rbac_role (id)
        ON DELETE CASCADE NOT DEFERRABLE,
    CONSTRAINT members_fkey
        FOREIGN KEY (to_role)
        REFERENCES rbac_role (id)
        ON DELETE CASCADE NOT DEFERRABLE
)
WITH (OIDS=FALSE);
ALTER TABLE public.rbac_grant OWNER TO try5adm;






--
-- Data for Name: accounts; Type: TABLE DATA; Schema: public; Owner: try5adm
--

COPY users (id, uid, name, email, password, active, gravatar, created, updated) FROM stdin (DELIMITER ',');
1,ce30ed61-6b5d-4136-95a3-ab11e3e97d87,Test User,user@test.com,$2a$10$T/tj9OCnQ4XUf7qcVsQsIuV9AxQgHaoaNxSOEnvdGdm.BEPpEG56e,true,\N,2013-08-18 17:46:23.748705,2013-08-18 17:46:23.748705
\.

ALTER SEQUENCE IF EXISTS users_id_seq RESTART WITH 2;

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;

COMMIT;
