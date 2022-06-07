BEGIN;

SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = ON;
SET check_function_bodies = FALSE;
SET client_min_messages = WARNING;
SET search_path = public, extensions;
SET default_tablespace = '';
SET default_with_oids = FALSE;


CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE public.courses
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    symbol     TEXT,
    buy        float8,
    rub        float8,
    created_at TIMESTAMPTZ
);

CREATE TABLE public.cbr
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    charcode   TEXT,
    name       TEXT,
    value      float8,
    created_at TIMESTAMPTZ
);


CREATE TABLE public.btctocbr
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT,
    value      float8,
    created_at TIMESTAMPTZ
);

ALTER TABLE public.btctocbr ADD UNIQUE (name);

COMMIT;