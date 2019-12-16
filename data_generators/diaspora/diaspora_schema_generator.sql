--
-- PostgreSQL database dump
--

-- Dumped from database version 10.8 (Ubuntu 10.8-1.pgdg16.04+1)
-- Dumped by pg_dump version 10.8 (Ubuntu 10.8-1.pgdg16.04+1)

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
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


--
-- Name: dblink; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS dblink WITH SCHEMA public;


--
-- Name: EXTENSION dblink; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION dblink IS 'connect to other PostgreSQL databases from within a database';


--
-- Name: truncate_tables(character varying); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.truncate_tables(username character varying) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
    statements CURSOR FOR
        SELECT tablename FROM pg_tables
        WHERE tableowner = username AND schemaname = 'public';
BEGIN
    FOR stmt IN statements LOOP
        EXECUTE 'TRUNCATE TABLE ' || quote_ident(stmt.tablename) || ' CASCADE;';
    END LOOP;
END;
$$;


ALTER FUNCTION public.truncate_tables(username character varying) OWNER TO postgres;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: account_deletions; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.account_deletions (
    id integer NOT NULL,
    person_id integer,
    completed_at timestamp without time zone,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.account_deletions OWNER TO diaspora;

--
-- Name: account_deletions_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.account_deletions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.account_deletions_id_seq OWNER TO diaspora;

--
-- Name: account_deletions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.account_deletions_id_seq OWNED BY public.account_deletions.id;


--
-- Name: account_migrations; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.account_migrations (
    id bigint NOT NULL,
    old_person_id integer NOT NULL,
    new_person_id integer NOT NULL,
    completed_at timestamp without time zone,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.account_migrations OWNER TO diaspora;

--
-- Name: account_migrations_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.account_migrations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.account_migrations_id_seq OWNER TO diaspora;

--
-- Name: account_migrations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.account_migrations_id_seq OWNED BY public.account_migrations.id;


--
-- Name: ar_internal_metadata; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.ar_internal_metadata (
    key character varying NOT NULL,
    value character varying,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.ar_internal_metadata OWNER TO diaspora;

--
-- Name: aspect_memberships; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.aspect_memberships (
    id integer NOT NULL,
    aspect_id integer NOT NULL,
    contact_id integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.aspect_memberships OWNER TO diaspora;

--
-- Name: aspect_memberships_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.aspect_memberships_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.aspect_memberships_id_seq OWNER TO diaspora;

--
-- Name: aspect_memberships_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.aspect_memberships_id_seq OWNED BY public.aspect_memberships.id;


--
-- Name: aspect_visibilities; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.aspect_visibilities (
    id integer NOT NULL,
    shareable_id integer NOT NULL,
    aspect_id integer NOT NULL,
    shareable_type character varying DEFAULT 'Post'::character varying NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.aspect_visibilities OWNER TO diaspora;

--
-- Name: aspect_visibilities_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.aspect_visibilities_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.aspect_visibilities_id_seq OWNER TO diaspora;

--
-- Name: aspect_visibilities_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.aspect_visibilities_id_seq OWNED BY public.aspect_visibilities.id;


--
-- Name: aspects; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.aspects (
    id integer NOT NULL,
    name character varying NOT NULL,
    user_id integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    order_id integer,
    chat_enabled boolean DEFAULT false,
    post_default boolean DEFAULT true,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.aspects OWNER TO diaspora;

--
-- Name: aspects_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.aspects_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.aspects_id_seq OWNER TO diaspora;

--
-- Name: aspects_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.aspects_id_seq OWNED BY public.aspects.id;


--
-- Name: authorizations; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.authorizations (
    id integer NOT NULL,
    user_id integer,
    o_auth_application_id integer,
    refresh_token character varying,
    code character varying,
    redirect_uri character varying,
    nonce character varying,
    scopes character varying,
    code_used boolean DEFAULT false,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.authorizations OWNER TO diaspora;

--
-- Name: authorizations_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.authorizations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.authorizations_id_seq OWNER TO diaspora;

--
-- Name: authorizations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.authorizations_id_seq OWNED BY public.authorizations.id;


--
-- Name: blocks; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.blocks (
    id integer NOT NULL,
    user_id integer,
    person_id integer,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.blocks OWNER TO diaspora;

--
-- Name: blocks_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.blocks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.blocks_id_seq OWNER TO diaspora;

--
-- Name: blocks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.blocks_id_seq OWNED BY public.blocks.id;


--
-- Name: chat_contacts; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.chat_contacts (
    id integer NOT NULL,
    user_id integer NOT NULL,
    jid character varying NOT NULL,
    name character varying(255),
    ask character varying(128),
    subscription character varying(128) NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.chat_contacts OWNER TO diaspora;

--
-- Name: chat_contacts_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.chat_contacts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.chat_contacts_id_seq OWNER TO diaspora;

--
-- Name: chat_contacts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.chat_contacts_id_seq OWNED BY public.chat_contacts.id;


--
-- Name: chat_fragments; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.chat_fragments (
    id integer NOT NULL,
    user_id integer NOT NULL,
    root character varying(256) NOT NULL,
    namespace character varying(256) NOT NULL,
    xml text NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.chat_fragments OWNER TO diaspora;

--
-- Name: chat_fragments_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.chat_fragments_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.chat_fragments_id_seq OWNER TO diaspora;

--
-- Name: chat_fragments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.chat_fragments_id_seq OWNED BY public.chat_fragments.id;


--
-- Name: chat_offline_messages; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.chat_offline_messages (
    id integer NOT NULL,
    "from" character varying NOT NULL,
    "to" character varying NOT NULL,
    message text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.chat_offline_messages OWNER TO diaspora;

--
-- Name: chat_offline_messages_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.chat_offline_messages_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.chat_offline_messages_id_seq OWNER TO diaspora;

--
-- Name: chat_offline_messages_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.chat_offline_messages_id_seq OWNED BY public.chat_offline_messages.id;


--
-- Name: comment_signatures; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.comment_signatures (
    comment_id integer NOT NULL,
    author_signature text NOT NULL,
    signature_order_id integer NOT NULL,
    additional_data text,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.comment_signatures OWNER TO diaspora;

--
-- Name: comments; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.comments (
    id integer NOT NULL,
    text text NOT NULL,
    commentable_id integer NOT NULL,
    author_id integer NOT NULL,
    guid character varying NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    likes_count integer DEFAULT 0 NOT NULL,
    commentable_type character varying(60) DEFAULT 'Post'::character varying NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.comments OWNER TO diaspora;

--
-- Name: comments_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.comments_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.comments_id_seq OWNER TO diaspora;

--
-- Name: comments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.comments_id_seq OWNED BY public.comments.id;


--
-- Name: contacts; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.contacts (
    id integer NOT NULL,
    user_id integer NOT NULL,
    person_id integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    sharing boolean DEFAULT false NOT NULL,
    receiving boolean DEFAULT false NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.contacts OWNER TO diaspora;

--
-- Name: contacts_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.contacts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.contacts_id_seq OWNER TO diaspora;

--
-- Name: contacts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.contacts_id_seq OWNED BY public.contacts.id;


--
-- Name: conversation_visibilities; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.conversation_visibilities (
    id integer NOT NULL,
    conversation_id integer NOT NULL,
    person_id integer NOT NULL,
    unread integer DEFAULT 0 NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.conversation_visibilities OWNER TO diaspora;

--
-- Name: conversation_visibilities_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.conversation_visibilities_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.conversation_visibilities_id_seq OWNER TO diaspora;

--
-- Name: conversation_visibilities_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.conversation_visibilities_id_seq OWNED BY public.conversation_visibilities.id;


--
-- Name: conversations; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.conversations (
    id integer NOT NULL,
    subject character varying,
    guid character varying NOT NULL,
    author_id integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.conversations OWNER TO diaspora;

--
-- Name: conversations_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.conversations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.conversations_id_seq OWNER TO diaspora;

--
-- Name: conversations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.conversations_id_seq OWNED BY public.conversations.id;


--
-- Name: invitation_codes; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.invitation_codes (
    id integer NOT NULL,
    token character varying,
    user_id integer,
    count integer,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.invitation_codes OWNER TO diaspora;

--
-- Name: invitation_codes_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.invitation_codes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.invitation_codes_id_seq OWNER TO diaspora;

--
-- Name: invitation_codes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.invitation_codes_id_seq OWNED BY public.invitation_codes.id;


--
-- Name: like_signatures; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.like_signatures (
    like_id integer NOT NULL,
    author_signature text NOT NULL,
    signature_order_id integer NOT NULL,
    additional_data text,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.like_signatures OWNER TO diaspora;

--
-- Name: likes; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.likes (
    id integer NOT NULL,
    positive boolean DEFAULT true,
    target_id integer,
    author_id integer,
    guid character varying,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    target_type character varying(60) NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.likes OWNER TO diaspora;

--
-- Name: likes_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.likes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.likes_id_seq OWNER TO diaspora;

--
-- Name: likes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.likes_id_seq OWNED BY public.likes.id;


--
-- Name: locations; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.locations (
    id integer NOT NULL,
    address character varying,
    lat character varying,
    lng character varying,
    status_message_id integer,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.locations OWNER TO diaspora;

--
-- Name: locations_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.locations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.locations_id_seq OWNER TO diaspora;

--
-- Name: locations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.locations_id_seq OWNED BY public.locations.id;


--
-- Name: mentions; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.mentions (
    id integer NOT NULL,
    mentions_container_id integer NOT NULL,
    person_id integer NOT NULL,
    mentions_container_type character varying NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.mentions OWNER TO diaspora;

--
-- Name: mentions_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.mentions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.mentions_id_seq OWNER TO diaspora;

--
-- Name: mentions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.mentions_id_seq OWNED BY public.mentions.id;


--
-- Name: messages; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.messages (
    id integer NOT NULL,
    conversation_id integer NOT NULL,
    author_id integer NOT NULL,
    guid character varying NOT NULL,
    text text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.messages OWNER TO diaspora;

--
-- Name: messages_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.messages_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.messages_id_seq OWNER TO diaspora;

--
-- Name: messages_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.messages_id_seq OWNED BY public.messages.id;


--
-- Name: notification_actors; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.notification_actors (
    id integer NOT NULL,
    notification_id integer,
    person_id integer,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.notification_actors OWNER TO diaspora;

--
-- Name: notification_actors_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.notification_actors_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.notification_actors_id_seq OWNER TO diaspora;

--
-- Name: notification_actors_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.notification_actors_id_seq OWNED BY public.notification_actors.id;


--
-- Name: notifications; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.notifications (
    id integer NOT NULL,
    target_type character varying,
    target_id integer,
    recipient_id integer NOT NULL,
    unread boolean DEFAULT true NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    type character varying,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.notifications OWNER TO diaspora;

--
-- Name: notifications_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.notifications_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.notifications_id_seq OWNER TO diaspora;

--
-- Name: notifications_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.notifications_id_seq OWNED BY public.notifications.id;


--
-- Name: o_auth_access_tokens; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.o_auth_access_tokens (
    id integer NOT NULL,
    authorization_id integer,
    token character varying,
    expires_at timestamp without time zone,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.o_auth_access_tokens OWNER TO diaspora;

--
-- Name: o_auth_access_tokens_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.o_auth_access_tokens_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.o_auth_access_tokens_id_seq OWNER TO diaspora;

--
-- Name: o_auth_access_tokens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.o_auth_access_tokens_id_seq OWNED BY public.o_auth_access_tokens.id;


--
-- Name: o_auth_applications; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.o_auth_applications (
    id integer NOT NULL,
    user_id integer,
    client_id character varying,
    client_secret character varying,
    client_name character varying,
    redirect_uris text,
    response_types character varying,
    grant_types character varying,
    application_type character varying DEFAULT 'web'::character varying,
    contacts character varying,
    logo_uri character varying,
    client_uri character varying,
    policy_uri character varying,
    tos_uri character varying,
    sector_identifier_uri character varying,
    token_endpoint_auth_method character varying,
    jwks text,
    jwks_uri character varying,
    ppid boolean DEFAULT false,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.o_auth_applications OWNER TO diaspora;

--
-- Name: o_auth_applications_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.o_auth_applications_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.o_auth_applications_id_seq OWNER TO diaspora;

--
-- Name: o_auth_applications_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.o_auth_applications_id_seq OWNED BY public.o_auth_applications.id;


--
-- Name: o_embed_caches; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.o_embed_caches (
    id integer NOT NULL,
    url character varying(1024) NOT NULL,
    data text NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.o_embed_caches OWNER TO diaspora;

--
-- Name: o_embed_caches_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.o_embed_caches_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.o_embed_caches_id_seq OWNER TO diaspora;

--
-- Name: o_embed_caches_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.o_embed_caches_id_seq OWNED BY public.o_embed_caches.id;


--
-- Name: open_graph_caches; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.open_graph_caches (
    id integer NOT NULL,
    title character varying,
    ob_type character varying,
    image text,
    url text,
    description text,
    video_url text,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.open_graph_caches OWNER TO diaspora;

--
-- Name: open_graph_caches_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.open_graph_caches_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.open_graph_caches_id_seq OWNER TO diaspora;

--
-- Name: open_graph_caches_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.open_graph_caches_id_seq OWNED BY public.open_graph_caches.id;


--
-- Name: participations; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.participations (
    id integer NOT NULL,
    guid character varying,
    target_id integer,
    target_type character varying(60) NOT NULL,
    author_id integer,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    count integer DEFAULT 1 NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.participations OWNER TO diaspora;

--
-- Name: participations_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.participations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.participations_id_seq OWNER TO diaspora;

--
-- Name: participations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.participations_id_seq OWNED BY public.participations.id;


--
-- Name: people; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.people (
    id integer NOT NULL,
    guid character varying NOT NULL,
    diaspora_handle character varying NOT NULL,
    serialized_public_key text NOT NULL,
    owner_id integer,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    closed_account boolean DEFAULT false,
    fetch_status integer DEFAULT 0,
    pod_id integer,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.people OWNER TO diaspora;

--
-- Name: people_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.people_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.people_id_seq OWNER TO diaspora;

--
-- Name: people_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.people_id_seq OWNED BY public.people.id;


--
-- Name: photos; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.photos (
    id integer NOT NULL,
    author_id integer NOT NULL,
    public boolean DEFAULT false NOT NULL,
    guid character varying NOT NULL,
    pending boolean DEFAULT false NOT NULL,
    text text,
    remote_photo_path text,
    remote_photo_name character varying,
    random_string character varying,
    processed_image character varying,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    unprocessed_image character varying,
    status_message_guid character varying,
    comments_count integer,
    height integer,
    width integer,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.photos OWNER TO diaspora;

--
-- Name: photos_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.photos_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.photos_id_seq OWNER TO diaspora;

--
-- Name: photos_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.photos_id_seq OWNED BY public.photos.id;


--
-- Name: pods; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.pods (
    id integer NOT NULL,
    host character varying NOT NULL,
    ssl boolean,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    status integer DEFAULT 0,
    checked_at timestamp without time zone DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
    offline_since timestamp without time zone,
    response_time integer DEFAULT '-1'::integer,
    software character varying,
    error character varying,
    port integer,
    blocked boolean DEFAULT false,
    scheduled_check boolean DEFAULT false NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.pods OWNER TO diaspora;

--
-- Name: pods_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.pods_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.pods_id_seq OWNER TO diaspora;

--
-- Name: pods_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.pods_id_seq OWNED BY public.pods.id;


--
-- Name: poll_answers; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.poll_answers (
    id integer NOT NULL,
    answer character varying NOT NULL,
    poll_id integer NOT NULL,
    guid character varying,
    vote_count integer DEFAULT 0,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.poll_answers OWNER TO diaspora;

--
-- Name: poll_answers_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.poll_answers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.poll_answers_id_seq OWNER TO diaspora;

--
-- Name: poll_answers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.poll_answers_id_seq OWNED BY public.poll_answers.id;


--
-- Name: poll_participation_signatures; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.poll_participation_signatures (
    poll_participation_id integer NOT NULL,
    author_signature text NOT NULL,
    signature_order_id integer NOT NULL,
    additional_data text,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.poll_participation_signatures OWNER TO diaspora;

--
-- Name: poll_participations; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.poll_participations (
    id integer NOT NULL,
    poll_answer_id integer NOT NULL,
    author_id integer NOT NULL,
    poll_id integer NOT NULL,
    guid character varying,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.poll_participations OWNER TO diaspora;

--
-- Name: poll_participations_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.poll_participations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.poll_participations_id_seq OWNER TO diaspora;

--
-- Name: poll_participations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.poll_participations_id_seq OWNED BY public.poll_participations.id;


--
-- Name: polls; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.polls (
    id integer NOT NULL,
    question character varying NOT NULL,
    status_message_id integer NOT NULL,
    status boolean,
    guid character varying,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.polls OWNER TO diaspora;

--
-- Name: polls_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.polls_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.polls_id_seq OWNER TO diaspora;

--
-- Name: polls_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.polls_id_seq OWNED BY public.polls.id;


--
-- Name: posts; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.posts (
    id integer NOT NULL,
    author_id integer NOT NULL,
    public boolean DEFAULT false NOT NULL,
    guid character varying NOT NULL,
    type character varying(40) NOT NULL,
    text text,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    provider_display_name character varying,
    root_guid character varying,
    likes_count integer DEFAULT 0,
    comments_count integer DEFAULT 0,
    o_embed_cache_id integer,
    reshares_count integer DEFAULT 0,
    interacted_at timestamp without time zone,
    tweet_id character varying,
    open_graph_cache_id integer,
    tumblr_ids text,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.posts OWNER TO diaspora;

--
-- Name: posts_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.posts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.posts_id_seq OWNER TO diaspora;

--
-- Name: posts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.posts_id_seq OWNED BY public.posts.id;


--
-- Name: ppid; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.ppid (
    id integer NOT NULL,
    o_auth_application_id integer,
    user_id integer,
    guid character varying(32),
    identifier character varying,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.ppid OWNER TO diaspora;

--
-- Name: ppid_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.ppid_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ppid_id_seq OWNER TO diaspora;

--
-- Name: ppid_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.ppid_id_seq OWNED BY public.ppid.id;


--
-- Name: profiles; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.profiles (
    id integer NOT NULL,
    diaspora_handle character varying,
    first_name character varying(127),
    last_name character varying(127),
    image_url character varying,
    image_url_small character varying,
    image_url_medium character varying,
    birthday date,
    gender character varying,
    bio text,
    searchable boolean DEFAULT true NOT NULL,
    person_id integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    location character varying,
    full_name character varying(70),
    nsfw boolean DEFAULT false,
    public_details boolean DEFAULT false,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.profiles OWNER TO diaspora;

--
-- Name: profiles_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.profiles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.profiles_id_seq OWNER TO diaspora;

--
-- Name: profiles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.profiles_id_seq OWNED BY public.profiles.id;


--
-- Name: references; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public."references" (
    id bigint NOT NULL,
    source_id integer NOT NULL,
    source_type character varying(60) NOT NULL,
    target_id integer NOT NULL,
    target_type character varying(60) NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public."references" OWNER TO diaspora;

--
-- Name: references_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.references_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.references_id_seq OWNER TO diaspora;

--
-- Name: references_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.references_id_seq OWNED BY public."references".id;


--
-- Name: reports; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.reports (
    id integer NOT NULL,
    item_id integer NOT NULL,
    item_type character varying NOT NULL,
    reviewed boolean DEFAULT false,
    text text,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    user_id integer NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.reports OWNER TO diaspora;

--
-- Name: reports_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.reports_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.reports_id_seq OWNER TO diaspora;

--
-- Name: reports_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.reports_id_seq OWNED BY public.reports.id;


--
-- Name: roles; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.roles (
    id integer NOT NULL,
    person_id integer,
    name character varying,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.roles OWNER TO diaspora;

--
-- Name: roles_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.roles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.roles_id_seq OWNER TO diaspora;

--
-- Name: roles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.roles_id_seq OWNED BY public.roles.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.schema_migrations (
    version character varying NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.schema_migrations OWNER TO diaspora;

--
-- Name: services; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.services (
    id integer NOT NULL,
    type character varying(127) NOT NULL,
    user_id integer NOT NULL,
    uid character varying(127),
    access_token character varying,
    access_secret character varying,
    nickname character varying,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.services OWNER TO diaspora;

--
-- Name: services_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.services_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.services_id_seq OWNER TO diaspora;

--
-- Name: services_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.services_id_seq OWNED BY public.services.id;


--
-- Name: share_visibilities; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.share_visibilities (
    id integer NOT NULL,
    shareable_id integer NOT NULL,
    hidden boolean DEFAULT false NOT NULL,
    shareable_type character varying(60) DEFAULT 'Post'::character varying NOT NULL,
    user_id integer NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.share_visibilities OWNER TO diaspora;

--
-- Name: share_visibilities_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.share_visibilities_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.share_visibilities_id_seq OWNER TO diaspora;

--
-- Name: share_visibilities_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.share_visibilities_id_seq OWNED BY public.share_visibilities.id;


--
-- Name: signature_orders; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.signature_orders (
    id integer NOT NULL,
    "order" character varying NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.signature_orders OWNER TO diaspora;

--
-- Name: signature_orders_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.signature_orders_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.signature_orders_id_seq OWNER TO diaspora;

--
-- Name: signature_orders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.signature_orders_id_seq OWNED BY public.signature_orders.id;


--
-- Name: simple_captcha_data; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.simple_captcha_data (
    id integer NOT NULL,
    key character varying(40),
    value character varying(12),
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.simple_captcha_data OWNER TO diaspora;

--
-- Name: simple_captcha_data_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.simple_captcha_data_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.simple_captcha_data_id_seq OWNER TO diaspora;

--
-- Name: simple_captcha_data_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.simple_captcha_data_id_seq OWNED BY public.simple_captcha_data.id;


--
-- Name: tag_followings; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.tag_followings (
    id integer NOT NULL,
    tag_id integer NOT NULL,
    user_id integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.tag_followings OWNER TO diaspora;

--
-- Name: tag_followings_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.tag_followings_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.tag_followings_id_seq OWNER TO diaspora;

--
-- Name: tag_followings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.tag_followings_id_seq OWNED BY public.tag_followings.id;


--
-- Name: taggings; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.taggings (
    id integer NOT NULL,
    tag_id integer,
    taggable_id integer,
    taggable_type character varying(127),
    tagger_id integer,
    tagger_type character varying(127),
    context character varying(127),
    created_at timestamp without time zone,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.taggings OWNER TO diaspora;

--
-- Name: taggings_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.taggings_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.taggings_id_seq OWNER TO diaspora;

--
-- Name: taggings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.taggings_id_seq OWNED BY public.taggings.id;


--
-- Name: tags; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.tags (
    id integer NOT NULL,
    name character varying,
    taggings_count integer DEFAULT 0,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.tags OWNER TO diaspora;

--
-- Name: tags_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.tags_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.tags_id_seq OWNER TO diaspora;

--
-- Name: tags_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.tags_id_seq OWNED BY public.tags.id;


--
-- Name: user_preferences; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.user_preferences (
    id integer NOT NULL,
    email_type character varying,
    user_id integer,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.user_preferences OWNER TO diaspora;

--
-- Name: user_preferences_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.user_preferences_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.user_preferences_id_seq OWNER TO diaspora;

--
-- Name: user_preferences_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.user_preferences_id_seq OWNED BY public.user_preferences.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: diaspora
--

CREATE TABLE public.users (
    id integer NOT NULL,
    username character varying NOT NULL,
    serialized_private_key text,
    getting_started boolean DEFAULT true NOT NULL,
    disable_mail boolean DEFAULT false NOT NULL,
    language character varying,
    email character varying DEFAULT ''::character varying NOT NULL,
    encrypted_password character varying DEFAULT ''::character varying NOT NULL,
    reset_password_token character varying,
    remember_created_at timestamp without time zone,
    sign_in_count integer DEFAULT 0,
    current_sign_in_at timestamp without time zone,
    last_sign_in_at timestamp without time zone,
    current_sign_in_ip character varying,
    last_sign_in_ip character varying,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    invited_by_id integer,
    authentication_token character varying(30),
    unconfirmed_email character varying,
    confirm_email_token character varying(30),
    locked_at timestamp without time zone,
    show_community_spotlight_in_stream boolean DEFAULT true NOT NULL,
    auto_follow_back boolean DEFAULT false,
    auto_follow_back_aspect_id integer,
    hidden_shareables text,
    reset_password_sent_at timestamp without time zone,
    last_seen timestamp without time zone,
    remove_after timestamp without time zone,
    export character varying,
    exported_at timestamp without time zone,
    exporting boolean DEFAULT false,
    strip_exif boolean DEFAULT true,
    exported_photos_file character varying,
    exported_photos_at timestamp without time zone,
    exporting_photos boolean DEFAULT false,
    color_theme character varying,
    post_default_public boolean DEFAULT false,
    mark_as_delete boolean DEFAULT false
);


ALTER TABLE public.users OWNER TO diaspora;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: diaspora
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO diaspora;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: diaspora
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: account_deletions id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.account_deletions ALTER COLUMN id SET DEFAULT nextval('public.account_deletions_id_seq'::regclass);


--
-- Name: account_migrations id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.account_migrations ALTER COLUMN id SET DEFAULT nextval('public.account_migrations_id_seq'::regclass);


--
-- Name: aspect_memberships id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.aspect_memberships ALTER COLUMN id SET DEFAULT nextval('public.aspect_memberships_id_seq'::regclass);


--
-- Name: aspect_visibilities id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.aspect_visibilities ALTER COLUMN id SET DEFAULT nextval('public.aspect_visibilities_id_seq'::regclass);


--
-- Name: aspects id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.aspects ALTER COLUMN id SET DEFAULT nextval('public.aspects_id_seq'::regclass);


--
-- Name: authorizations id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.authorizations ALTER COLUMN id SET DEFAULT nextval('public.authorizations_id_seq'::regclass);


--
-- Name: blocks id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.blocks ALTER COLUMN id SET DEFAULT nextval('public.blocks_id_seq'::regclass);


--
-- Name: chat_contacts id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.chat_contacts ALTER COLUMN id SET DEFAULT nextval('public.chat_contacts_id_seq'::regclass);


--
-- Name: chat_fragments id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.chat_fragments ALTER COLUMN id SET DEFAULT nextval('public.chat_fragments_id_seq'::regclass);


--
-- Name: chat_offline_messages id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.chat_offline_messages ALTER COLUMN id SET DEFAULT nextval('public.chat_offline_messages_id_seq'::regclass);


--
-- Name: comments id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.comments ALTER COLUMN id SET DEFAULT nextval('public.comments_id_seq'::regclass);


--
-- Name: contacts id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.contacts ALTER COLUMN id SET DEFAULT nextval('public.contacts_id_seq'::regclass);


--
-- Name: conversation_visibilities id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.conversation_visibilities ALTER COLUMN id SET DEFAULT nextval('public.conversation_visibilities_id_seq'::regclass);


--
-- Name: conversations id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.conversations ALTER COLUMN id SET DEFAULT nextval('public.conversations_id_seq'::regclass);


--
-- Name: invitation_codes id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.invitation_codes ALTER COLUMN id SET DEFAULT nextval('public.invitation_codes_id_seq'::regclass);


--
-- Name: likes id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.likes ALTER COLUMN id SET DEFAULT nextval('public.likes_id_seq'::regclass);


--
-- Name: locations id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.locations ALTER COLUMN id SET DEFAULT nextval('public.locations_id_seq'::regclass);


--
-- Name: mentions id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.mentions ALTER COLUMN id SET DEFAULT nextval('public.mentions_id_seq'::regclass);


--
-- Name: messages id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.messages ALTER COLUMN id SET DEFAULT nextval('public.messages_id_seq'::regclass);


--
-- Name: notification_actors id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.notification_actors ALTER COLUMN id SET DEFAULT nextval('public.notification_actors_id_seq'::regclass);


--
-- Name: notifications id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.notifications ALTER COLUMN id SET DEFAULT nextval('public.notifications_id_seq'::regclass);


--
-- Name: o_auth_access_tokens id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.o_auth_access_tokens ALTER COLUMN id SET DEFAULT nextval('public.o_auth_access_tokens_id_seq'::regclass);


--
-- Name: o_auth_applications id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.o_auth_applications ALTER COLUMN id SET DEFAULT nextval('public.o_auth_applications_id_seq'::regclass);


--
-- Name: o_embed_caches id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.o_embed_caches ALTER COLUMN id SET DEFAULT nextval('public.o_embed_caches_id_seq'::regclass);


--
-- Name: open_graph_caches id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.open_graph_caches ALTER COLUMN id SET DEFAULT nextval('public.open_graph_caches_id_seq'::regclass);


--
-- Name: participations id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.participations ALTER COLUMN id SET DEFAULT nextval('public.participations_id_seq'::regclass);


--
-- Name: people id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.people ALTER COLUMN id SET DEFAULT nextval('public.people_id_seq'::regclass);


--
-- Name: photos id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.photos ALTER COLUMN id SET DEFAULT nextval('public.photos_id_seq'::regclass);


--
-- Name: pods id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.pods ALTER COLUMN id SET DEFAULT nextval('public.pods_id_seq'::regclass);


--
-- Name: poll_answers id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.poll_answers ALTER COLUMN id SET DEFAULT nextval('public.poll_answers_id_seq'::regclass);


--
-- Name: poll_participations id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.poll_participations ALTER COLUMN id SET DEFAULT nextval('public.poll_participations_id_seq'::regclass);


--
-- Name: polls id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.polls ALTER COLUMN id SET DEFAULT nextval('public.polls_id_seq'::regclass);


--
-- Name: posts id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.posts ALTER COLUMN id SET DEFAULT nextval('public.posts_id_seq'::regclass);


--
-- Name: ppid id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.ppid ALTER COLUMN id SET DEFAULT nextval('public.ppid_id_seq'::regclass);


--
-- Name: profiles id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.profiles ALTER COLUMN id SET DEFAULT nextval('public.profiles_id_seq'::regclass);


--
-- Name: references id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public."references" ALTER COLUMN id SET DEFAULT nextval('public.references_id_seq'::regclass);


--
-- Name: reports id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.reports ALTER COLUMN id SET DEFAULT nextval('public.reports_id_seq'::regclass);


--
-- Name: roles id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.roles ALTER COLUMN id SET DEFAULT nextval('public.roles_id_seq'::regclass);


--
-- Name: services id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.services ALTER COLUMN id SET DEFAULT nextval('public.services_id_seq'::regclass);


--
-- Name: share_visibilities id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.share_visibilities ALTER COLUMN id SET DEFAULT nextval('public.share_visibilities_id_seq'::regclass);


--
-- Name: signature_orders id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.signature_orders ALTER COLUMN id SET DEFAULT nextval('public.signature_orders_id_seq'::regclass);


--
-- Name: simple_captcha_data id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.simple_captcha_data ALTER COLUMN id SET DEFAULT nextval('public.simple_captcha_data_id_seq'::regclass);


--
-- Name: tag_followings id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.tag_followings ALTER COLUMN id SET DEFAULT nextval('public.tag_followings_id_seq'::regclass);


--
-- Name: taggings id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.taggings ALTER COLUMN id SET DEFAULT nextval('public.taggings_id_seq'::regclass);


--
-- Name: tags id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.tags ALTER COLUMN id SET DEFAULT nextval('public.tags_id_seq'::regclass);


--
-- Name: user_preferences id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.user_preferences ALTER COLUMN id SET DEFAULT nextval('public.user_preferences_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: account_deletions account_deletions_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.account_deletions
    ADD CONSTRAINT account_deletions_pkey PRIMARY KEY (id);


--
-- Name: account_migrations account_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.account_migrations
    ADD CONSTRAINT account_migrations_pkey PRIMARY KEY (id);


--
-- Name: ar_internal_metadata ar_internal_metadata_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.ar_internal_metadata
    ADD CONSTRAINT ar_internal_metadata_pkey PRIMARY KEY (key);


--
-- Name: aspect_memberships aspect_memberships_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.aspect_memberships
    ADD CONSTRAINT aspect_memberships_pkey PRIMARY KEY (id);


--
-- Name: aspect_visibilities aspect_visibilities_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.aspect_visibilities
    ADD CONSTRAINT aspect_visibilities_pkey PRIMARY KEY (id);


--
-- Name: aspects aspects_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.aspects
    ADD CONSTRAINT aspects_pkey PRIMARY KEY (id);


--
-- Name: authorizations authorizations_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.authorizations
    ADD CONSTRAINT authorizations_pkey PRIMARY KEY (id);


--
-- Name: blocks blocks_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.blocks
    ADD CONSTRAINT blocks_pkey PRIMARY KEY (id);


--
-- Name: chat_contacts chat_contacts_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.chat_contacts
    ADD CONSTRAINT chat_contacts_pkey PRIMARY KEY (id);


--
-- Name: chat_fragments chat_fragments_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.chat_fragments
    ADD CONSTRAINT chat_fragments_pkey PRIMARY KEY (id);


--
-- Name: chat_offline_messages chat_offline_messages_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.chat_offline_messages
    ADD CONSTRAINT chat_offline_messages_pkey PRIMARY KEY (id);


--
-- Name: comments comments_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comments_pkey PRIMARY KEY (id);


--
-- Name: contacts contacts_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.contacts
    ADD CONSTRAINT contacts_pkey PRIMARY KEY (id);


--
-- Name: conversation_visibilities conversation_visibilities_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.conversation_visibilities
    ADD CONSTRAINT conversation_visibilities_pkey PRIMARY KEY (id);


--
-- Name: conversations conversations_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.conversations
    ADD CONSTRAINT conversations_pkey PRIMARY KEY (id);


--
-- Name: invitation_codes invitation_codes_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.invitation_codes
    ADD CONSTRAINT invitation_codes_pkey PRIMARY KEY (id);


--
-- Name: likes likes_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.likes
    ADD CONSTRAINT likes_pkey PRIMARY KEY (id);


--
-- Name: locations locations_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.locations
    ADD CONSTRAINT locations_pkey PRIMARY KEY (id);


--
-- Name: mentions mentions_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.mentions
    ADD CONSTRAINT mentions_pkey PRIMARY KEY (id);


--
-- Name: messages messages_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.messages
    ADD CONSTRAINT messages_pkey PRIMARY KEY (id);


--
-- Name: notification_actors notification_actors_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.notification_actors
    ADD CONSTRAINT notification_actors_pkey PRIMARY KEY (id);


--
-- Name: notifications notifications_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_pkey PRIMARY KEY (id);


--
-- Name: o_auth_access_tokens o_auth_access_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.o_auth_access_tokens
    ADD CONSTRAINT o_auth_access_tokens_pkey PRIMARY KEY (id);


--
-- Name: o_auth_applications o_auth_applications_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.o_auth_applications
    ADD CONSTRAINT o_auth_applications_pkey PRIMARY KEY (id);


--
-- Name: o_embed_caches o_embed_caches_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.o_embed_caches
    ADD CONSTRAINT o_embed_caches_pkey PRIMARY KEY (id);


--
-- Name: open_graph_caches open_graph_caches_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.open_graph_caches
    ADD CONSTRAINT open_graph_caches_pkey PRIMARY KEY (id);


--
-- Name: participations participations_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.participations
    ADD CONSTRAINT participations_pkey PRIMARY KEY (id);


--
-- Name: people people_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.people
    ADD CONSTRAINT people_pkey PRIMARY KEY (id);


--
-- Name: photos photos_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.photos
    ADD CONSTRAINT photos_pkey PRIMARY KEY (id);


--
-- Name: pods pods_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.pods
    ADD CONSTRAINT pods_pkey PRIMARY KEY (id);


--
-- Name: poll_answers poll_answers_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.poll_answers
    ADD CONSTRAINT poll_answers_pkey PRIMARY KEY (id);


--
-- Name: poll_participations poll_participations_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.poll_participations
    ADD CONSTRAINT poll_participations_pkey PRIMARY KEY (id);


--
-- Name: polls polls_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.polls
    ADD CONSTRAINT polls_pkey PRIMARY KEY (id);


--
-- Name: posts posts_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_pkey PRIMARY KEY (id);


--
-- Name: ppid ppid_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.ppid
    ADD CONSTRAINT ppid_pkey PRIMARY KEY (id);


--
-- Name: profiles profiles_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.profiles
    ADD CONSTRAINT profiles_pkey PRIMARY KEY (id);


--
-- Name: references references_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public."references"
    ADD CONSTRAINT references_pkey PRIMARY KEY (id);


--
-- Name: reports reports_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.reports
    ADD CONSTRAINT reports_pkey PRIMARY KEY (id);


--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: services services_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.services
    ADD CONSTRAINT services_pkey PRIMARY KEY (id);


--
-- Name: share_visibilities share_visibilities_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.share_visibilities
    ADD CONSTRAINT share_visibilities_pkey PRIMARY KEY (id);


--
-- Name: signature_orders signature_orders_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.signature_orders
    ADD CONSTRAINT signature_orders_pkey PRIMARY KEY (id);


--
-- Name: simple_captcha_data simple_captcha_data_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.simple_captcha_data
    ADD CONSTRAINT simple_captcha_data_pkey PRIMARY KEY (id);


--
-- Name: tag_followings tag_followings_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.tag_followings
    ADD CONSTRAINT tag_followings_pkey PRIMARY KEY (id);


--
-- Name: taggings taggings_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.taggings
    ADD CONSTRAINT taggings_pkey PRIMARY KEY (id);


--
-- Name: tags tags_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.tags
    ADD CONSTRAINT tags_pkey PRIMARY KEY (id);


--
-- Name: user_preferences user_preferences_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.user_preferences
    ADD CONSTRAINT user_preferences_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: contacts_user_id_idx; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX contacts_user_id_idx ON public.contacts USING btree (user_id);


--
-- Name: conversations_author_id_fk; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX conversations_author_id_fk ON public.conversations USING btree (author_id);


--
-- Name: idx_key; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX idx_key ON public.simple_captcha_data USING btree (key);


--
-- Name: index_account_deletions_on_person_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_account_deletions_on_person_id ON public.account_deletions USING btree (person_id);


--
-- Name: index_account_migrations_on_old_person_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_account_migrations_on_old_person_id ON public.account_migrations USING btree (old_person_id);


--
-- Name: index_account_migrations_on_old_person_id_and_new_person_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_account_migrations_on_old_person_id_and_new_person_id ON public.account_migrations USING btree (old_person_id, new_person_id);


--
-- Name: index_aspect_memberships_on_aspect_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_aspect_memberships_on_aspect_id ON public.aspect_memberships USING btree (aspect_id);


--
-- Name: index_aspect_memberships_on_aspect_id_and_contact_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_aspect_memberships_on_aspect_id_and_contact_id ON public.aspect_memberships USING btree (aspect_id, contact_id);


--
-- Name: index_aspect_memberships_on_contact_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_aspect_memberships_on_contact_id ON public.aspect_memberships USING btree (contact_id);


--
-- Name: index_aspect_visibilities_on_aspect_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_aspect_visibilities_on_aspect_id ON public.aspect_visibilities USING btree (aspect_id);


--
-- Name: index_aspect_visibilities_on_shareable_and_aspect_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_aspect_visibilities_on_shareable_and_aspect_id ON public.aspect_visibilities USING btree (shareable_id, shareable_type, aspect_id);


--
-- Name: index_aspect_visibilities_on_shareable_id_and_shareable_type; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_aspect_visibilities_on_shareable_id_and_shareable_type ON public.aspect_visibilities USING btree (shareable_id, shareable_type);


--
-- Name: index_aspects_on_user_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_aspects_on_user_id ON public.aspects USING btree (user_id);


--
-- Name: index_aspects_on_user_id_and_name; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_aspects_on_user_id_and_name ON public.aspects USING btree (user_id, name);


--
-- Name: index_authorizations_on_o_auth_application_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_authorizations_on_o_auth_application_id ON public.authorizations USING btree (o_auth_application_id);


--
-- Name: index_authorizations_on_user_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_authorizations_on_user_id ON public.authorizations USING btree (user_id);


--
-- Name: index_blocks_on_user_id_and_person_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_blocks_on_user_id_and_person_id ON public.blocks USING btree (user_id, person_id);


--
-- Name: index_chat_contacts_on_user_id_and_jid; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_chat_contacts_on_user_id_and_jid ON public.chat_contacts USING btree (user_id, jid);


--
-- Name: index_chat_fragments_on_user_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_chat_fragments_on_user_id ON public.chat_fragments USING btree (user_id);


--
-- Name: index_comment_signatures_on_comment_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_comment_signatures_on_comment_id ON public.comment_signatures USING btree (comment_id);


--
-- Name: index_comments_on_commentable_id_and_commentable_type; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_comments_on_commentable_id_and_commentable_type ON public.comments USING btree (commentable_id, commentable_type);


--
-- Name: index_comments_on_guid; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_comments_on_guid ON public.comments USING btree (guid);


--
-- Name: index_comments_on_person_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_comments_on_person_id ON public.comments USING btree (author_id);


--
-- Name: index_contacts_on_person_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_contacts_on_person_id ON public.contacts USING btree (person_id);


--
-- Name: index_contacts_on_user_id_and_person_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_contacts_on_user_id_and_person_id ON public.contacts USING btree (user_id, person_id);


--
-- Name: index_conversation_visibilities_on_conversation_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_conversation_visibilities_on_conversation_id ON public.conversation_visibilities USING btree (conversation_id);


--
-- Name: index_conversation_visibilities_on_person_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_conversation_visibilities_on_person_id ON public.conversation_visibilities USING btree (person_id);


--
-- Name: index_conversation_visibilities_usefully; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_conversation_visibilities_usefully ON public.conversation_visibilities USING btree (conversation_id, person_id);


--
-- Name: index_conversations_on_guid; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_conversations_on_guid ON public.conversations USING btree (guid);


--
-- Name: index_like_signatures_on_like_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_like_signatures_on_like_id ON public.like_signatures USING btree (like_id);


--
-- Name: index_likes_on_guid; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_likes_on_guid ON public.likes USING btree (guid);


--
-- Name: index_likes_on_post_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_likes_on_post_id ON public.likes USING btree (target_id);


--
-- Name: index_likes_on_target_id_and_author_id_and_target_type; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_likes_on_target_id_and_author_id_and_target_type ON public.likes USING btree (target_id, author_id, target_type);


--
-- Name: index_locations_on_status_message_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_locations_on_status_message_id ON public.locations USING btree (status_message_id);


--
-- Name: index_mentions_on_mc_id_and_mc_type; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_mentions_on_mc_id_and_mc_type ON public.mentions USING btree (mentions_container_id, mentions_container_type);


--
-- Name: index_mentions_on_person_and_mc_id_and_mc_type; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_mentions_on_person_and_mc_id_and_mc_type ON public.mentions USING btree (person_id, mentions_container_id, mentions_container_type);


--
-- Name: index_mentions_on_person_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_mentions_on_person_id ON public.mentions USING btree (person_id);


--
-- Name: index_messages_on_author_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_messages_on_author_id ON public.messages USING btree (author_id);


--
-- Name: index_messages_on_guid; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_messages_on_guid ON public.messages USING btree (guid);


--
-- Name: index_notification_actors_on_notification_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_notification_actors_on_notification_id ON public.notification_actors USING btree (notification_id);


--
-- Name: index_notification_actors_on_notification_id_and_person_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_notification_actors_on_notification_id_and_person_id ON public.notification_actors USING btree (notification_id, person_id);


--
-- Name: index_notification_actors_on_person_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_notification_actors_on_person_id ON public.notification_actors USING btree (person_id);


--
-- Name: index_notifications_on_recipient_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_notifications_on_recipient_id ON public.notifications USING btree (recipient_id);


--
-- Name: index_notifications_on_target_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_notifications_on_target_id ON public.notifications USING btree (target_id);


--
-- Name: index_notifications_on_target_type_and_target_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_notifications_on_target_type_and_target_id ON public.notifications USING btree (target_type, target_id);


--
-- Name: index_o_auth_access_tokens_on_authorization_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_o_auth_access_tokens_on_authorization_id ON public.o_auth_access_tokens USING btree (authorization_id);


--
-- Name: index_o_auth_access_tokens_on_token; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_o_auth_access_tokens_on_token ON public.o_auth_access_tokens USING btree (token);


--
-- Name: index_o_auth_applications_on_client_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_o_auth_applications_on_client_id ON public.o_auth_applications USING btree (client_id);


--
-- Name: index_o_auth_applications_on_user_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_o_auth_applications_on_user_id ON public.o_auth_applications USING btree (user_id);


--
-- Name: index_o_embed_caches_on_url; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_o_embed_caches_on_url ON public.o_embed_caches USING btree (url);


--
-- Name: index_participations_on_author_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_participations_on_author_id ON public.participations USING btree (author_id);


--
-- Name: index_participations_on_guid; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_participations_on_guid ON public.participations USING btree (guid);


--
-- Name: index_participations_on_target_id_and_target_type_and_author_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_participations_on_target_id_and_target_type_and_author_id ON public.participations USING btree (target_id, target_type, author_id);


--
-- Name: index_people_on_diaspora_handle; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_people_on_diaspora_handle ON public.people USING btree (diaspora_handle);


--
-- Name: index_people_on_guid; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_people_on_guid ON public.people USING btree (guid);


--
-- Name: index_people_on_owner_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_people_on_owner_id ON public.people USING btree (owner_id);


--
-- Name: index_photos_on_author_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_photos_on_author_id ON public.photos USING btree (author_id);


--
-- Name: index_photos_on_guid; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_photos_on_guid ON public.photos USING btree (guid);


--
-- Name: index_photos_on_status_message_guid; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_photos_on_status_message_guid ON public.photos USING btree (status_message_guid);


--
-- Name: index_pods_on_checked_at; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_pods_on_checked_at ON public.pods USING btree (checked_at);


--
-- Name: index_pods_on_host_and_port; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_pods_on_host_and_port ON public.pods USING btree (host, port);


--
-- Name: index_pods_on_offline_since; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_pods_on_offline_since ON public.pods USING btree (offline_since);


--
-- Name: index_pods_on_status; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_pods_on_status ON public.pods USING btree (status);


--
-- Name: index_poll_answers_on_guid; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_poll_answers_on_guid ON public.poll_answers USING btree (guid);


--
-- Name: index_poll_answers_on_poll_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_poll_answers_on_poll_id ON public.poll_answers USING btree (poll_id);


--
-- Name: index_poll_participation_signatures_on_poll_participation_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_poll_participation_signatures_on_poll_participation_id ON public.poll_participation_signatures USING btree (poll_participation_id);


--
-- Name: index_poll_participations_on_guid; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_poll_participations_on_guid ON public.poll_participations USING btree (guid);


--
-- Name: index_poll_participations_on_poll_id_and_author_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_poll_participations_on_poll_id_and_author_id ON public.poll_participations USING btree (poll_id, author_id);


--
-- Name: index_polls_on_guid; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_polls_on_guid ON public.polls USING btree (guid);


--
-- Name: index_polls_on_status_message_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_polls_on_status_message_id ON public.polls USING btree (status_message_id);


--
-- Name: index_post_visibilities_on_post_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_post_visibilities_on_post_id ON public.share_visibilities USING btree (shareable_id);


--
-- Name: index_posts_on_author_id_and_root_guid; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_posts_on_author_id_and_root_guid ON public.posts USING btree (author_id, root_guid);


--
-- Name: index_posts_on_created_at_and_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_posts_on_created_at_and_id ON public.posts USING btree (created_at, id);


--
-- Name: index_posts_on_guid; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_posts_on_guid ON public.posts USING btree (guid);


--
-- Name: index_posts_on_id_and_type; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_posts_on_id_and_type ON public.posts USING btree (id, type);


--
-- Name: index_posts_on_person_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_posts_on_person_id ON public.posts USING btree (author_id);


--
-- Name: index_posts_on_root_guid; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_posts_on_root_guid ON public.posts USING btree (root_guid);


--
-- Name: index_ppid_on_o_auth_application_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_ppid_on_o_auth_application_id ON public.ppid USING btree (o_auth_application_id);


--
-- Name: index_ppid_on_user_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_ppid_on_user_id ON public.ppid USING btree (user_id);


--
-- Name: index_profiles_on_full_name; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_profiles_on_full_name ON public.profiles USING btree (full_name);


--
-- Name: index_profiles_on_full_name_and_searchable; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_profiles_on_full_name_and_searchable ON public.profiles USING btree (full_name, searchable);


--
-- Name: index_profiles_on_person_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_profiles_on_person_id ON public.profiles USING btree (person_id);


--
-- Name: index_references_on_source_and_target; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_references_on_source_and_target ON public."references" USING btree (source_id, source_type, target_id, target_type);


--
-- Name: index_references_on_source_id_and_source_type; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_references_on_source_id_and_source_type ON public."references" USING btree (source_id, source_type);


--
-- Name: index_reports_on_item_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_reports_on_item_id ON public.reports USING btree (item_id);


--
-- Name: index_roles_on_person_id_and_name; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_roles_on_person_id_and_name ON public.roles USING btree (person_id, name);


--
-- Name: index_services_on_type_and_uid; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_services_on_type_and_uid ON public.services USING btree (type, uid);


--
-- Name: index_services_on_user_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_services_on_user_id ON public.services USING btree (user_id);


--
-- Name: index_share_visibilities_on_user_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_share_visibilities_on_user_id ON public.share_visibilities USING btree (user_id);


--
-- Name: index_signature_orders_on_order; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_signature_orders_on_order ON public.signature_orders USING btree ("order");


--
-- Name: index_tag_followings_on_tag_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_tag_followings_on_tag_id ON public.tag_followings USING btree (tag_id);


--
-- Name: index_tag_followings_on_tag_id_and_user_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_tag_followings_on_tag_id_and_user_id ON public.tag_followings USING btree (tag_id, user_id);


--
-- Name: index_tag_followings_on_user_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_tag_followings_on_user_id ON public.tag_followings USING btree (user_id);


--
-- Name: index_taggings_on_created_at; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_taggings_on_created_at ON public.taggings USING btree (created_at);


--
-- Name: index_taggings_on_tag_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_taggings_on_tag_id ON public.taggings USING btree (tag_id);


--
-- Name: index_taggings_on_taggable_id_and_taggable_type_and_context; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_taggings_on_taggable_id_and_taggable_type_and_context ON public.taggings USING btree (taggable_id, taggable_type, context);


--
-- Name: index_taggings_uniquely; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_taggings_uniquely ON public.taggings USING btree (taggable_id, taggable_type, tag_id);


--
-- Name: index_tags_on_name; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_tags_on_name ON public.tags USING btree (name);


--
-- Name: index_user_preferences_on_user_id_and_email_type; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX index_user_preferences_on_user_id_and_email_type ON public.user_preferences USING btree (user_id, email_type);


--
-- Name: index_users_on_authentication_token; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_users_on_authentication_token ON public.users USING btree (authentication_token);


--
-- Name: index_users_on_email; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_users_on_email ON public.users USING btree (email);


--
-- Name: index_users_on_username; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX index_users_on_username ON public.users USING btree (username);


--
-- Name: likes_author_id_fk; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX likes_author_id_fk ON public.likes USING btree (author_id);


--
-- Name: messages_conversation_id_fk; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX messages_conversation_id_fk ON public.messages USING btree (conversation_id);


--
-- Name: shareable_and_hidden_and_user_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE INDEX shareable_and_hidden_and_user_id ON public.share_visibilities USING btree (shareable_id, shareable_type, hidden, user_id);


--
-- Name: shareable_and_user_id; Type: INDEX; Schema: public; Owner: diaspora
--

CREATE UNIQUE INDEX shareable_and_user_id ON public.share_visibilities USING btree (shareable_id, shareable_type, user_id);


--
-- Name: aspect_memberships aspect_memberships_aspect_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.aspect_memberships
    ADD CONSTRAINT aspect_memberships_aspect_id_fk FOREIGN KEY (aspect_id) REFERENCES public.aspects(id) ON DELETE CASCADE;


--
-- Name: aspect_memberships aspect_memberships_contact_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.aspect_memberships
    ADD CONSTRAINT aspect_memberships_contact_id_fk FOREIGN KEY (contact_id) REFERENCES public.contacts(id) ON DELETE CASCADE;


--
-- Name: aspect_visibilities aspect_visibilities_aspect_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.aspect_visibilities
    ADD CONSTRAINT aspect_visibilities_aspect_id_fk FOREIGN KEY (aspect_id) REFERENCES public.aspects(id) ON DELETE CASCADE;


--
-- Name: comment_signatures comment_signatures_comment_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.comment_signatures
    ADD CONSTRAINT comment_signatures_comment_id_fk FOREIGN KEY (comment_id) REFERENCES public.comments(id) ON DELETE CASCADE;


--
-- Name: comment_signatures comment_signatures_signature_orders_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.comment_signatures
    ADD CONSTRAINT comment_signatures_signature_orders_id_fk FOREIGN KEY (signature_order_id) REFERENCES public.signature_orders(id);


--
-- Name: comments comments_author_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comments_author_id_fk FOREIGN KEY (author_id) REFERENCES public.people(id) ON DELETE CASCADE;


--
-- Name: contacts contacts_person_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.contacts
    ADD CONSTRAINT contacts_person_id_fk FOREIGN KEY (person_id) REFERENCES public.people(id) ON DELETE CASCADE;


--
-- Name: conversation_visibilities conversation_visibilities_conversation_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.conversation_visibilities
    ADD CONSTRAINT conversation_visibilities_conversation_id_fk FOREIGN KEY (conversation_id) REFERENCES public.conversations(id) ON DELETE CASCADE;


--
-- Name: conversation_visibilities conversation_visibilities_person_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.conversation_visibilities
    ADD CONSTRAINT conversation_visibilities_person_id_fk FOREIGN KEY (person_id) REFERENCES public.people(id) ON DELETE CASCADE;


--
-- Name: conversations conversations_author_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.conversations
    ADD CONSTRAINT conversations_author_id_fk FOREIGN KEY (author_id) REFERENCES public.people(id) ON DELETE CASCADE;


--
-- Name: ppid fk_rails_150457f962; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.ppid
    ADD CONSTRAINT fk_rails_150457f962 FOREIGN KEY (o_auth_application_id) REFERENCES public.o_auth_applications(id);


--
-- Name: authorizations fk_rails_4ecef5b8c5; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.authorizations
    ADD CONSTRAINT fk_rails_4ecef5b8c5 FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: o_auth_access_tokens fk_rails_5debabcff3; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.o_auth_access_tokens
    ADD CONSTRAINT fk_rails_5debabcff3 FOREIGN KEY (authorization_id) REFERENCES public.authorizations(id);


--
-- Name: account_migrations fk_rails_610fe19943; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.account_migrations
    ADD CONSTRAINT fk_rails_610fe19943 FOREIGN KEY (new_person_id) REFERENCES public.people(id);


--
-- Name: o_auth_applications fk_rails_ad75323da2; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.o_auth_applications
    ADD CONSTRAINT fk_rails_ad75323da2 FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: account_migrations fk_rails_ddbe553eee; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.account_migrations
    ADD CONSTRAINT fk_rails_ddbe553eee FOREIGN KEY (old_person_id) REFERENCES public.people(id);


--
-- Name: authorizations fk_rails_e166644de5; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.authorizations
    ADD CONSTRAINT fk_rails_e166644de5 FOREIGN KEY (o_auth_application_id) REFERENCES public.o_auth_applications(id);


--
-- Name: ppid fk_rails_e6b8e5264f; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.ppid
    ADD CONSTRAINT fk_rails_e6b8e5264f FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: like_signatures like_signatures_like_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.like_signatures
    ADD CONSTRAINT like_signatures_like_id_fk FOREIGN KEY (like_id) REFERENCES public.likes(id) ON DELETE CASCADE;


--
-- Name: like_signatures like_signatures_signature_orders_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.like_signatures
    ADD CONSTRAINT like_signatures_signature_orders_id_fk FOREIGN KEY (signature_order_id) REFERENCES public.signature_orders(id);


--
-- Name: likes likes_author_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.likes
    ADD CONSTRAINT likes_author_id_fk FOREIGN KEY (author_id) REFERENCES public.people(id) ON DELETE CASCADE;


--
-- Name: messages messages_author_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.messages
    ADD CONSTRAINT messages_author_id_fk FOREIGN KEY (author_id) REFERENCES public.people(id) ON DELETE CASCADE;


--
-- Name: messages messages_conversation_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.messages
    ADD CONSTRAINT messages_conversation_id_fk FOREIGN KEY (conversation_id) REFERENCES public.conversations(id) ON DELETE CASCADE;


--
-- Name: notification_actors notification_actors_notification_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.notification_actors
    ADD CONSTRAINT notification_actors_notification_id_fk FOREIGN KEY (notification_id) REFERENCES public.notifications(id) ON DELETE CASCADE;


--
-- Name: people people_pod_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.people
    ADD CONSTRAINT people_pod_id_fk FOREIGN KEY (pod_id) REFERENCES public.pods(id) ON DELETE CASCADE;


--
-- Name: poll_participation_signatures poll_participation_signatures_poll_participation_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.poll_participation_signatures
    ADD CONSTRAINT poll_participation_signatures_poll_participation_id_fk FOREIGN KEY (poll_participation_id) REFERENCES public.poll_participations(id) ON DELETE CASCADE;


--
-- Name: poll_participation_signatures poll_participation_signatures_signature_orders_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.poll_participation_signatures
    ADD CONSTRAINT poll_participation_signatures_signature_orders_id_fk FOREIGN KEY (signature_order_id) REFERENCES public.signature_orders(id);


--
-- Name: posts posts_author_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_author_id_fk FOREIGN KEY (author_id) REFERENCES public.people(id) ON DELETE CASCADE;


--
-- Name: profiles profiles_person_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.profiles
    ADD CONSTRAINT profiles_person_id_fk FOREIGN KEY (person_id) REFERENCES public.people(id) ON DELETE CASCADE;


--
-- Name: services services_user_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.services
    ADD CONSTRAINT services_user_id_fk FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: share_visibilities share_visibilities_user_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: diaspora
--

ALTER TABLE ONLY public.share_visibilities
    ADD CONSTRAINT share_visibilities_user_id_fk FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: postgres
--

GRANT USAGE ON SCHEMA public TO cow1;


--
-- Name: TABLE account_deletions; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.account_deletions TO cow1;


--
-- Name: SEQUENCE account_deletions_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.account_deletions_id_seq TO cow1;


--
-- Name: TABLE account_migrations; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.account_migrations TO cow1;


--
-- Name: SEQUENCE account_migrations_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.account_migrations_id_seq TO cow1;


--
-- Name: TABLE ar_internal_metadata; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.ar_internal_metadata TO cow1;


--
-- Name: TABLE aspect_memberships; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.aspect_memberships TO cow1;


--
-- Name: SEQUENCE aspect_memberships_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.aspect_memberships_id_seq TO cow1;


--
-- Name: TABLE aspect_visibilities; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.aspect_visibilities TO cow1;


--
-- Name: SEQUENCE aspect_visibilities_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.aspect_visibilities_id_seq TO cow1;


--
-- Name: TABLE aspects; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.aspects TO cow1;


--
-- Name: SEQUENCE aspects_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.aspects_id_seq TO cow1;


--
-- Name: TABLE authorizations; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.authorizations TO cow1;


--
-- Name: SEQUENCE authorizations_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.authorizations_id_seq TO cow1;


--
-- Name: TABLE blocks; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.blocks TO cow1;


--
-- Name: SEQUENCE blocks_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.blocks_id_seq TO cow1;


--
-- Name: TABLE chat_contacts; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.chat_contacts TO cow1;


--
-- Name: SEQUENCE chat_contacts_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.chat_contacts_id_seq TO cow1;


--
-- Name: TABLE chat_fragments; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.chat_fragments TO cow1;


--
-- Name: SEQUENCE chat_fragments_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.chat_fragments_id_seq TO cow1;


--
-- Name: TABLE chat_offline_messages; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.chat_offline_messages TO cow1;


--
-- Name: SEQUENCE chat_offline_messages_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.chat_offline_messages_id_seq TO cow1;


--
-- Name: TABLE comment_signatures; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.comment_signatures TO cow1;


--
-- Name: TABLE comments; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.comments TO cow1;


--
-- Name: SEQUENCE comments_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.comments_id_seq TO cow1;


--
-- Name: TABLE contacts; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.contacts TO cow1;


--
-- Name: SEQUENCE contacts_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.contacts_id_seq TO cow1;


--
-- Name: TABLE conversation_visibilities; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.conversation_visibilities TO cow1;


--
-- Name: SEQUENCE conversation_visibilities_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.conversation_visibilities_id_seq TO cow1;


--
-- Name: TABLE conversations; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.conversations TO cow1;


--
-- Name: SEQUENCE conversations_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.conversations_id_seq TO cow1;


--
-- Name: TABLE invitation_codes; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.invitation_codes TO cow1;


--
-- Name: SEQUENCE invitation_codes_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.invitation_codes_id_seq TO cow1;


--
-- Name: TABLE like_signatures; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.like_signatures TO cow1;


--
-- Name: TABLE likes; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.likes TO cow1;


--
-- Name: SEQUENCE likes_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.likes_id_seq TO cow1;


--
-- Name: TABLE locations; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.locations TO cow1;


--
-- Name: SEQUENCE locations_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.locations_id_seq TO cow1;


--
-- Name: TABLE mentions; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.mentions TO cow1;


--
-- Name: SEQUENCE mentions_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.mentions_id_seq TO cow1;


--
-- Name: TABLE messages; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.messages TO cow1;


--
-- Name: SEQUENCE messages_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.messages_id_seq TO cow1;


--
-- Name: TABLE notification_actors; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.notification_actors TO cow1;


--
-- Name: SEQUENCE notification_actors_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.notification_actors_id_seq TO cow1;


--
-- Name: TABLE notifications; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.notifications TO cow1;


--
-- Name: SEQUENCE notifications_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.notifications_id_seq TO cow1;


--
-- Name: TABLE o_auth_access_tokens; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.o_auth_access_tokens TO cow1;


--
-- Name: SEQUENCE o_auth_access_tokens_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.o_auth_access_tokens_id_seq TO cow1;


--
-- Name: TABLE o_auth_applications; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.o_auth_applications TO cow1;


--
-- Name: SEQUENCE o_auth_applications_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.o_auth_applications_id_seq TO cow1;


--
-- Name: TABLE o_embed_caches; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.o_embed_caches TO cow1;


--
-- Name: SEQUENCE o_embed_caches_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.o_embed_caches_id_seq TO cow1;


--
-- Name: TABLE open_graph_caches; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.open_graph_caches TO cow1;


--
-- Name: SEQUENCE open_graph_caches_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.open_graph_caches_id_seq TO cow1;


--
-- Name: TABLE participations; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.participations TO cow1;


--
-- Name: SEQUENCE participations_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.participations_id_seq TO cow1;


--
-- Name: TABLE people; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.people TO cow1;


--
-- Name: SEQUENCE people_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.people_id_seq TO cow1;


--
-- Name: TABLE photos; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.photos TO cow1;


--
-- Name: SEQUENCE photos_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.photos_id_seq TO cow1;


--
-- Name: TABLE pods; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.pods TO cow1;


--
-- Name: SEQUENCE pods_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.pods_id_seq TO cow1;


--
-- Name: TABLE poll_answers; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.poll_answers TO cow1;


--
-- Name: SEQUENCE poll_answers_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.poll_answers_id_seq TO cow1;


--
-- Name: TABLE poll_participation_signatures; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.poll_participation_signatures TO cow1;


--
-- Name: TABLE poll_participations; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.poll_participations TO cow1;


--
-- Name: SEQUENCE poll_participations_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.poll_participations_id_seq TO cow1;


--
-- Name: TABLE polls; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.polls TO cow1;


--
-- Name: SEQUENCE polls_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.polls_id_seq TO cow1;


--
-- Name: TABLE posts; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.posts TO cow1;


--
-- Name: SEQUENCE posts_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.posts_id_seq TO cow1;


--
-- Name: TABLE ppid; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.ppid TO cow1;


--
-- Name: SEQUENCE ppid_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.ppid_id_seq TO cow1;


--
-- Name: TABLE profiles; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.profiles TO cow1;


--
-- Name: SEQUENCE profiles_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.profiles_id_seq TO cow1;


--
-- Name: TABLE "references"; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public."references" TO cow1;


--
-- Name: SEQUENCE references_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.references_id_seq TO cow1;


--
-- Name: TABLE reports; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.reports TO cow1;


--
-- Name: SEQUENCE reports_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.reports_id_seq TO cow1;


--
-- Name: TABLE roles; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.roles TO cow1;


--
-- Name: SEQUENCE roles_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.roles_id_seq TO cow1;


--
-- Name: TABLE schema_migrations; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.schema_migrations TO cow1;


--
-- Name: TABLE services; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.services TO cow1;


--
-- Name: SEQUENCE services_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.services_id_seq TO cow1;


--
-- Name: TABLE share_visibilities; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.share_visibilities TO cow1;


--
-- Name: SEQUENCE share_visibilities_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.share_visibilities_id_seq TO cow1;


--
-- Name: TABLE signature_orders; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.signature_orders TO cow1;


--
-- Name: SEQUENCE signature_orders_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.signature_orders_id_seq TO cow1;


--
-- Name: TABLE simple_captcha_data; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.simple_captcha_data TO cow1;


--
-- Name: SEQUENCE simple_captcha_data_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.simple_captcha_data_id_seq TO cow1;


--
-- Name: TABLE tag_followings; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.tag_followings TO cow1;


--
-- Name: SEQUENCE tag_followings_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.tag_followings_id_seq TO cow1;


--
-- Name: TABLE taggings; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.taggings TO cow1;


--
-- Name: SEQUENCE taggings_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.taggings_id_seq TO cow1;


--
-- Name: TABLE tags; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.tags TO cow1;


--
-- Name: SEQUENCE tags_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.tags_id_seq TO cow1;


--
-- Name: TABLE user_preferences; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.user_preferences TO cow1;


--
-- Name: SEQUENCE user_preferences_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.user_preferences_id_seq TO cow1;


--
-- Name: TABLE users; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON TABLE public.users TO cow1;


--
-- Name: SEQUENCE users_id_seq; Type: ACL; Schema: public; Owner: diaspora
--

GRANT ALL ON SEQUENCE public.users_id_seq TO cow1;


--
-- PostgreSQL database dump complete
--

