--
-- PostgreSQL database dump
--

-- Dumped from database version 14.15 (Ubuntu 14.15-0ubuntu0.22.04.1)
-- Dumped by pg_dump version 15.0

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: -
--

-- *not* creating schema, since initdb creates it


--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: user_role; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.user_role AS ENUM (
    'admin',
    'editor',
    'viewer'
);


--
-- Name: quote; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.quote AS (
	id uuid,
	text text,
	author character varying(256),
	datetime timestamp with time zone,
	creator_id uuid,
	creator_name character varying(32),
	creator_login character varying(32),
	creator_role public.user_role,
	creator_telegram_id bigint
);


--
-- Name: user; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public."user" AS (
	id uuid,
	name character varying(32),
	login character varying(32),
	role public.user_role,
	telegram_id bigint
);


--
-- Name: quote_add(uuid, text, character varying, timestamp with time zone); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.quote_add(user_id uuid, quote_text text, quote_author character varying DEFAULT NULL::character varying, quote_datetime timestamp with time zone DEFAULT statement_timestamp(), OUT new_quote public.quote) RETURNS public.quote
    LANGUAGE plpgsql
    AS $$
DECLARE
	curr_user_role user_role;
	new_quote_id uuid;
BEGIN
	SELECT "role" FROM users WHERE id=user_id INTO curr_user_role;
	IF curr_user_role IS NULL THEN
		RAISE '401:Чел, ты кто? Я тебя не знаю';
	END IF;
	IF curr_user_role='viewer' THEN
		RAISE '403:Увы и ах, такие права у тебя ещё не скачаны :(';
	END IF;
	INSERT INTO quotes (text, author, datetime, creator_id)
	VALUES (quote_text, quote_author, quote_datetime, user_id)
	RETURNING id INTO new_quote_id;
	SELECT * FROM quote_get(user_id, new_quote_id) INTO new_quote;
END
$$;


--
-- Name: quote_delete(uuid, uuid); Type: PROCEDURE; Schema: public; Owner: -
--

CREATE PROCEDURE public.quote_delete(IN user_id uuid, IN quote_id uuid)
    LANGUAGE plpgsql
    AS $$
DECLARE
	quote_temp "quote";
	curr_user_role user_role;
BEGIN
	SELECT "role" FROM users WHERE id=user_id INTO curr_user_role;
	IF curr_user_role IS NULL THEN
		RAISE '401:Чел, ты кто? Я тебя не знаю';
	END IF;
	SELECT * FROM quote_get(user_id, quote_id) INTO quote_temp;
	IF quote_temp IS NULL THEN
		RAISE '404:Такую цитату ещё не сказанули';
	END IF;
	IF curr_user_role='viewer' OR quote_temp.creator_id!=user_id AND curr_user_role!='admin' THEN
		RAISE '403:Увы и ах, такие права у тебя ещё не скачаны :(';
	END IF;
	UPDATE quotes
	SET is_removed=true
	WHERE id=quote_id;
END
$$;


--
-- Name: quote_get(uuid, uuid); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.quote_get(user_id uuid, quote_id uuid, OUT ret_quote public.quote) RETURNS public.quote
    LANGUAGE plpgsql
    AS $$
BEGIN
	PERFORM 1 FROM users WHERE id=user_id;
	IF NOT FOUND THEN
		RAISE '401:Чел, ты кто? Я тебя не знаю';
	END IF;
	SELECT * FROM quotes_get(user_id) WHERE id=quote_id INTO ret_quote;
	if ret_quote IS NULL THEN
		RAISE '404:Такую цитату ещё не сказанули';
	END IF;
END
$$;


--
-- Name: quote_update(uuid, uuid, text, character varying, timestamp with time zone); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.quote_update(user_id uuid, quote_id uuid, new_text text DEFAULT NULL::text, new_author character varying DEFAULT NULL::character varying, new_datetime timestamp with time zone DEFAULT NULL::timestamp with time zone, OUT updated_quote public.quote) RETURNS public.quote
    LANGUAGE plpgsql
    AS $$
DECLARE
	quote_temp "quote";
	curr_user_role user_role;
BEGIN
	SELECT "role" FROM users WHERE id=user_id INTO curr_user_role;
	IF curr_user_role IS NULL THEN
		RAISE '401:Чел, ты кто? Я тебя не знаю';
	END IF;
	SELECT * FROM quote_get(user_id, quote_id) INTO quote_temp;
	IF quote_temp IS NULL THEN
		RAISE '404:Такую цитату ещё не сказанули';
	END IF;
	IF curr_user_role='viewer' OR quote_temp.creator_id!=user_id AND curr_user_role!='admin' THEN
		RAISE '403:Увы и ах, такие права у тебя ещё не скачаны :(';
	END IF;
	UPDATE quotes
	SET
		text=coalesce(new_text, text),
		author=coalesce(new_author, author),
		datetime=coalesce(new_datetime, datetime)
	WHERE id=quote_id;
	SELECT * FROM quote_get(user_id, quote_id) INTO updated_quote;
END
$$;


--
-- Name: quotes_get(uuid); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.quotes_get(user_id uuid) RETURNS SETOF public.quote
    LANGUAGE plpgsql
    AS $$
BEGIN
	PERFORM 1 FROM users WHERE id=user_id;
	IF NOT FOUND THEN
		RAISE '401:Чел, ты кто? Я тебя не знаю';
	END IF;
	RETURN QUERY SELECT q.id, q.text, q.author, q.datetime, u.id, u.name, u.login, u.role, u.telegram_id
	FROM quotes q
	JOIN users u ON u.id=q.creator_id
	WHERE NOT q.is_removed;
END
$$;


--
-- Name: user_auth(character varying, character varying); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.user_auth(user_login character varying, user_password character varying, OUT ret_user public."user") RETURNS public."user"
    LANGUAGE plpgsql STABLE PARALLEL SAFE
    AS $$
BEGIN
	SELECT u.id, u.name, u.login, u.role, u.telegram_id
	FROM users u
	WHERE u.login=user_login AND u.password=crypt(user_password, u.password)
	INTO ret_user;
	IF ret_user IS NULL THEN
		RAISE '401:Чел, ты кто? Я тебя не знаю';
	END IF;
END
$$;


--
-- Name: user_get(uuid); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.user_get(user_id uuid, OUT ret_user public."user") RETURNS public."user"
    LANGUAGE plpgsql
    AS $$
BEGIN
	SELECT id, name, login, role, telegram_id FROM users WHERE id=user_id INTO ret_user;
	IF ret_user IS NULL THEN
		RAISE '401:Чел, ты кто? Я тебя не знаю';
	END IF;
END
$$;


--
-- Name: user_update(uuid, character varying, character varying, bigint, character varying); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.user_update(user_id uuid, new_name character varying DEFAULT NULL::character varying, new_login character varying DEFAULT NULL::character varying, new_telegram_id bigint DEFAULT NULL::bigint, new_password character varying DEFAULT NULL::character varying, OUT updated_user public."user") RETURNS public."user"
    LANGUAGE plpgsql
    AS $$
BEGIN
	PERFORM 1 FROM users WHERE id=user_id;
	IF NOT FOUND THEN
		RAISE '401:Чел, ты кто? Я тебя не знаю';
	END IF;
	UPDATE users
	SET
		name=coalesce(new_name, name),
		login=coalesce(new_login, login),
		telegram_id=coalesce(new_telegram_id, telegram_id),
		password=coalesce(crypt(new_password, gen_salt('bf')), password)
	WHERE id=user_id;
	SELECT * FROM user_get(user_id) WHERE id=user_id INTO updated_user;
END
$$;


SET default_table_access_method = heap;

--
-- Name: quotes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.quotes (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    text text NOT NULL,
    datetime timestamp with time zone DEFAULT now() NOT NULL,
    author character varying(256) NOT NULL,
    creator_id uuid NOT NULL,
    is_removed boolean DEFAULT false NOT NULL
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(32) NOT NULL,
    telegram_id bigint NOT NULL,
    is_editor boolean DEFAULT false NOT NULL,
    login character varying(32) NOT NULL,
    password text NOT NULL,
    role public.user_role DEFAULT 'viewer'::public.user_role NOT NULL
);


--
-- Name: quotes prm__quotes; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.quotes
    ADD CONSTRAINT prm__quotes PRIMARY KEY (id);


--
-- Name: users prm__users; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT prm__users PRIMARY KEY (id);


--
-- Name: users uni__users__login; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT uni__users__login UNIQUE (login);


--
-- Name: users uni__users__name; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT uni__users__name UNIQUE (name);


--
-- Name: users uni__users__telegram_uid; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT uni__users__telegram_uid UNIQUE (telegram_id);


--
-- Name: idx__quotes__creator_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx__quotes__creator_id ON public.quotes USING hash (creator_id);


--
-- Name: idx__quotes__date; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx__quotes__date ON public.quotes USING btree (datetime DESC NULLS LAST);


--
-- Name: quotes frn__quotes__creator_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.quotes
    ADD CONSTRAINT frn__quotes__creator_id FOREIGN KEY (creator_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE RESTRICT;


--
-- PostgreSQL database dump complete
--

