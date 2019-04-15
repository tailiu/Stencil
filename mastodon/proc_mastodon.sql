CREATE TABLE public.imports (
    id serial8 NOT NULL PRIMARY KEY,
    type integer NOT NULL,
    approved boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    data_file_name character varying,
    data_content_type character varying,
    data_file_size integer,
    data_updated_at timestamp without time zone,
    account_id bigint NOT NULL
);
CREATE TABLE public.tags (
    id serial8 NOT NULL PRIMARY KEY,
    name character varying DEFAULT ''::character varying NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);
CREATE TABLE public.web_push_subscriptions (
    id serial8 NOT NULL PRIMARY KEY,
    endpoint character varying NOT NULL,
    key_p256dh character varying NOT NULL,
    key_auth character varying NOT NULL,
    data json,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    access_token_id bigint,
    user_id bigint
);
CREATE TABLE public.invites (
    id serial8 NOT NULL PRIMARY KEY,
    user_id bigint NOT NULL,
    code character varying DEFAULT ''::character varying NOT NULL,
    expires_at timestamp without time zone,
    max_uses integer,
    uses integer DEFAULT 0 NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    autofollow boolean DEFAULT false NOT NULL
);
CREATE TABLE public.follow_requests (
    id serial8 NOT NULL PRIMARY KEY,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    account_id bigint NOT NULL,
    target_account_id bigint NOT NULL,
    show_reblogs boolean DEFAULT true NOT NULL,
    uri character varying
);
CREATE TABLE public.favourites (
    id serial8 NOT NULL PRIMARY KEY,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    account_id bigint NOT NULL,
    status_id bigint NOT NULL
);
CREATE TABLE public.conversation_mutes (
    id serial8 NOT NULL PRIMARY KEY,
    conversation_id bigint NOT NULL,
    account_id bigint NOT NULL
);
CREATE TABLE public.preview_cards (
    id serial8 NOT NULL PRIMARY KEY,
    url character varying DEFAULT ''::character varying NOT NULL,
    title character varying DEFAULT ''::character varying NOT NULL,
    description character varying DEFAULT ''::character varying NOT NULL,
    image_file_name character varying,
    image_content_type character varying,
    image_file_size integer,
    image_updated_at timestamp without time zone,
    type integer DEFAULT 0 NOT NULL,
    html text DEFAULT ''::text NOT NULL,
    author_name character varying DEFAULT ''::character varying NOT NULL,
    author_url character varying DEFAULT ''::character varying NOT NULL,
    provider_name character varying DEFAULT ''::character varying NOT NULL,
    provider_url character varying DEFAULT ''::character varying NOT NULL,
    width integer DEFAULT 0 NOT NULL,
    height integer DEFAULT 0 NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    embed_url character varying DEFAULT ''::character varying NOT NULL
);
CREATE TABLE public.accounts_tags (
    account_id bigint NOT NULL,
    tag_id bigint NOT NULL
);
CREATE TABLE public.account_domain_blocks (
    id serial8 NOT NULL PRIMARY KEY,
    domain character varying,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    account_id bigint
);
CREATE TABLE public.site_uploads (
    id serial8 NOT NULL PRIMARY KEY,
    var character varying DEFAULT ''::character varying NOT NULL,
    file_file_name character varying,
    file_content_type character varying,
    file_file_size integer,
    file_updated_at timestamp without time zone,
    meta json,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);
CREATE TABLE public.users (
    id serial8 NOT NULL PRIMARY KEY,
    email character varying DEFAULT ''::character varying NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    encrypted_password character varying DEFAULT ''::character varying NOT NULL,
    reset_password_token character varying,
    reset_password_sent_at timestamp without time zone,
    remember_created_at timestamp without time zone,
    sign_in_count integer DEFAULT 0 NOT NULL,
    current_sign_in_at timestamp without time zone,
    last_sign_in_at timestamp without time zone,
    current_sign_in_ip inet,
    last_sign_in_ip inet,
    admin boolean DEFAULT false NOT NULL,
    confirmation_token character varying,
    confirmed_at timestamp without time zone,
    confirmation_sent_at timestamp without time zone,
    unconfirmed_email character varying,
    locale character varying,
    encrypted_otp_secret character varying,
    encrypted_otp_secret_iv character varying,
    encrypted_otp_secret_salt character varying,
    consumed_timestep integer,
    otp_required_for_login boolean DEFAULT false NOT NULL,
    last_emailed_at timestamp without time zone,
    otp_backup_codes character varying[],
    filtered_languages character varying[] DEFAULT '{}'::character varying[] NOT NULL,
    account_id bigint NOT NULL,
    disabled boolean DEFAULT false NOT NULL,
    moderator boolean DEFAULT false NOT NULL,
    invite_id bigint,
    remember_token character varying,
    chosen_languages character varying[],
    created_by_application_id bigint
);
CREATE TABLE public.ar_internal_metadata (
    key character varying NOT NULL PRIMARY KEY,
    value character varying,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);
CREATE TABLE public.custom_filters (
    id serial8 NOT NULL PRIMARY KEY,
    account_id bigint,
    expires_at timestamp without time zone,
    phrase text DEFAULT ''::text NOT NULL,
    context character varying[] DEFAULT '{}'::character varying[] NOT NULL,
    irreversible boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    whole_word boolean DEFAULT true NOT NULL
);
CREATE TABLE public.admin_action_logs (
    id serial8 NOT NULL PRIMARY KEY,
    account_id bigint,
    action character varying DEFAULT ''::character varying NOT NULL,
    target_type character varying,
    target_id bigint,
    recorded_changes text DEFAULT ''::text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);
CREATE TABLE public.pghero_space_stats (
    id serial8 NOT NULL PRIMARY KEY,
    database text,
    schema text,
    relation text,
    size bigint,
    captured_at timestamp without time zone
);
CREATE TABLE public.reports (
    id serial8 NOT NULL PRIMARY KEY,
    status_ids bigint[] DEFAULT '{}'::bigint[] NOT NULL,
    comment text DEFAULT ''::text NOT NULL,
    action_taken boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    account_id bigint NOT NULL,
    action_taken_by_account_id bigint,
    target_account_id bigint NOT NULL,
    assigned_account_id bigint
);
CREATE TABLE public.settings (
    id serial8 NOT NULL PRIMARY KEY,
    var character varying NOT NULL,
    value text,
    thing_type character varying,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    thing_id bigint
);
CREATE TABLE public.identities (
    id serial8 NOT NULL PRIMARY KEY,
    provider character varying DEFAULT ''::character varying NOT NULL,
    uid character varying DEFAULT ''::character varying NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    user_id bigint
);
CREATE TABLE public.account_tag_stats (
    id serial8 NOT NULL PRIMARY KEY,
    tag_id bigint NOT NULL,
    accounts_count bigint DEFAULT 0 NOT NULL,
    hidden boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);
CREATE TABLE public.preview_cards_statuses (
    preview_card_id bigint NOT NULL,
    status_id bigint NOT NULL
);
CREATE TABLE public.account_stats (
    id serial8 NOT NULL PRIMARY KEY,
    account_id bigint NOT NULL,
    statuses_count bigint DEFAULT 0 NOT NULL,
    following_count bigint DEFAULT 0 NOT NULL,
    followers_count bigint DEFAULT 0 NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    last_status_at timestamp without time zone
);
CREATE TABLE public.status_pins (
    id serial8 NOT NULL PRIMARY KEY,
    account_id bigint NOT NULL,
    status_id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);
CREATE TABLE public.conversations (
    id serial8 NOT NULL PRIMARY KEY,
    uri character varying,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);
CREATE TABLE public.backups (
    id serial8 NOT NULL PRIMARY KEY,
    user_id bigint,
    dump_file_name character varying,
    dump_content_type character varying,
    dump_file_size integer,
    dump_updated_at timestamp without time zone,
    processed boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);
CREATE TABLE public.oauth_applications (
    id serial8 NOT NULL PRIMARY KEY,
    name character varying NOT NULL,
    uid character varying NOT NULL,
    secret character varying NOT NULL,
    redirect_uri text NOT NULL,
    scopes character varying DEFAULT ''::character varying NOT NULL,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    superapp boolean DEFAULT false NOT NULL,
    website character varying,
    owner_type character varying,
    owner_id bigint,
    confidential boolean DEFAULT true NOT NULL
);
CREATE TABLE public.statuses (
    id serial8 NOT NULL PRIMARY KEY,
    uri character varying,
    text text DEFAULT ''::text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    in_reply_to_id bigint,
    reblog_of_id bigint,
    url character varying,
    sensitive boolean DEFAULT false NOT NULL,
    visibility integer DEFAULT 0 NOT NULL,
    spoiler_text text DEFAULT ''::text NOT NULL,
    reply boolean DEFAULT false NOT NULL,
    language character varying,
    conversation_id bigint,
    local boolean,
    account_id bigint NOT NULL,
    application_id bigint,
    in_reply_to_account_id bigint
);
CREATE TABLE public.domain_blocks (
    id serial8 NOT NULL PRIMARY KEY,
    domain character varying DEFAULT ''::character varying NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    severity integer DEFAULT 0,
    reject_media boolean DEFAULT false NOT NULL,
    reject_reports boolean DEFAULT false NOT NULL
);
CREATE TABLE public.statuses_tags (
    status_id bigint NOT NULL,
    tag_id bigint NOT NULL
);
CREATE TABLE public.account_moderation_notes (
    id serial8 NOT NULL PRIMARY KEY,
    content text NOT NULL,
    account_id bigint NOT NULL,
    target_account_id bigint NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);
CREATE TABLE public.tombstones (
    id serial8 NOT NULL PRIMARY KEY,
    account_id bigint,
    uri character varying NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);
CREATE TABLE public.web_settings (
    id serial8 NOT NULL PRIMARY KEY,
    data json,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    user_id bigint NOT NULL
);
CREATE TABLE public.mutes (
    id serial8 NOT NULL PRIMARY KEY,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    account_id bigint NOT NULL,
    target_account_id bigint NOT NULL,
    hide_notifications boolean DEFAULT true NOT NULL
);
CREATE TABLE public.blocks (
    id serial8 NOT NULL PRIMARY KEY,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    account_id bigint NOT NULL,
    target_account_id bigint NOT NULL,
    uri character varying
);
CREATE TABLE public.follows (
    id serial8 NOT NULL PRIMARY KEY,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    account_id bigint NOT NULL,
    target_account_id bigint NOT NULL,
    show_reblogs boolean DEFAULT true NOT NULL,
    uri character varying
);
CREATE TABLE public.lists (
    id serial8 NOT NULL PRIMARY KEY,
    account_id bigint NOT NULL,
    title character varying DEFAULT ''::character varying NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);
CREATE TABLE public.stream_entries (
    id serial8 NOT NULL PRIMARY KEY,
    activity_id bigint,
    activity_type character varying,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    hidden boolean DEFAULT false NOT NULL,
    account_id bigint
);
CREATE TABLE public.notifications (
    id serial8 NOT NULL PRIMARY KEY,
    activity_id bigint NOT NULL,
    activity_type character varying NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    account_id bigint NOT NULL,
    from_account_id bigint NOT NULL
);
CREATE TABLE public.list_accounts (
    id serial8 NOT NULL PRIMARY KEY,
    list_id bigint NOT NULL,
    account_id bigint NOT NULL,
    follow_id bigint NOT NULL
);
CREATE TABLE public.mentions (
    id serial8 NOT NULL PRIMARY KEY,
    status_id bigint,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    account_id bigint,
    silent boolean DEFAULT false NOT NULL
);
CREATE TABLE public.media_attachments (
    id serial8 NOT NULL PRIMARY KEY,
    status_id bigint,
    file_file_name character varying,
    file_content_type character varying,
    file_file_size integer,
    file_updated_at timestamp without time zone,
    remote_url character varying DEFAULT ''::character varying NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    shortcode character varying,
    type integer DEFAULT 0 NOT NULL,
    file_meta json,
    account_id bigint,
    description text,
    scheduled_status_id bigint
);
CREATE TABLE public.account_conversations (
    id serial8 NOT NULL PRIMARY KEY,
    account_id bigint,
    conversation_id bigint,
    participant_account_ids bigint[] DEFAULT '{}'::bigint[] NOT NULL,
    status_ids bigint[] DEFAULT '{}'::bigint[] NOT NULL,
    last_status_id bigint,
    lock_version integer DEFAULT 0 NOT NULL,
    unread boolean DEFAULT false NOT NULL
);
CREATE TABLE public.account_warning_presets (
    id serial8 NOT NULL PRIMARY KEY,
    text text DEFAULT ''::text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);
CREATE TABLE public.account_warnings (
    id serial8 NOT NULL PRIMARY KEY,
    account_id bigint,
    target_account_id bigint,
    action integer DEFAULT 0 NOT NULL,
    text text DEFAULT ''::text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);
CREATE TABLE public.custom_emojis (
    id serial8 NOT NULL PRIMARY KEY,
    shortcode character varying DEFAULT ''::character varying NOT NULL,
    domain character varying,
    image_file_name character varying,
    image_content_type character varying,
    image_file_size integer,
    image_updated_at timestamp without time zone,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    disabled boolean DEFAULT false NOT NULL,
    uri character varying,
    image_remote_url character varying,
    visible_in_picker boolean DEFAULT true NOT NULL
);
CREATE TABLE public.relays (
    id serial8 NOT NULL PRIMARY KEY,
    inbox_url character varying DEFAULT ''::character varying NOT NULL,
    follow_activity_id character varying,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    state integer DEFAULT 0 NOT NULL
);
CREATE TABLE public.status_stats (
    id serial8 NOT NULL PRIMARY KEY,
    status_id bigint NOT NULL,
    replies_count bigint DEFAULT 0 NOT NULL,
    reblogs_count bigint DEFAULT 0 NOT NULL,
    favourites_count bigint DEFAULT 0 NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);
CREATE TABLE public.email_domain_blocks (
    id serial8 NOT NULL PRIMARY KEY,
    domain character varying DEFAULT ''::character varying NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);
CREATE TABLE public.account_pins (
    id serial8 NOT NULL PRIMARY KEY,
    account_id bigint,
    target_account_id bigint,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);
CREATE TABLE public.oauth_access_tokens (
    id serial8 NOT NULL PRIMARY KEY,
    token character varying NOT NULL,
    refresh_token character varying,
    expires_in integer,
    revoked_at timestamp without time zone,
    created_at timestamp without time zone NOT NULL,
    scopes character varying,
    application_id bigint,
    resource_owner_id bigint
);
CREATE TABLE public.report_notes (
    id serial8 NOT NULL PRIMARY KEY,
    content text NOT NULL,
    report_id bigint NOT NULL,
    account_id bigint NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);
CREATE TABLE public.oauth_access_grants (
    id serial8 NOT NULL PRIMARY KEY,
    token character varying NOT NULL,
    expires_in integer NOT NULL,
    redirect_uri text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    revoked_at timestamp without time zone,
    scopes character varying,
    application_id bigint NOT NULL,
    resource_owner_id bigint NOT NULL
);
CREATE TABLE public.accounts (
    id serial8 NOT NULL PRIMARY KEY,
    username character varying DEFAULT ''::character varying NOT NULL,
    domain character varying,
    secret character varying DEFAULT ''::character varying NOT NULL,
    private_key text,
    public_key text DEFAULT ''::text NOT NULL,
    remote_url character varying DEFAULT ''::character varying NOT NULL,
    salmon_url character varying DEFAULT ''::character varying NOT NULL,
    hub_url character varying DEFAULT ''::character varying NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    note text DEFAULT ''::text NOT NULL,
    display_name character varying DEFAULT ''::character varying NOT NULL,
    uri character varying DEFAULT ''::character varying NOT NULL,
    url character varying,
    avatar_file_name character varying,
    avatar_content_type character varying,
    avatar_file_size integer,
    avatar_updated_at timestamp without time zone,
    header_file_name character varying,
    header_content_type character varying,
    header_file_size integer,
    header_updated_at timestamp without time zone,
    avatar_remote_url character varying,
    subscription_expires_at timestamp without time zone,
    silenced boolean DEFAULT false NOT NULL,
    suspended boolean DEFAULT false NOT NULL,
    locked boolean DEFAULT false NOT NULL,
    header_remote_url character varying DEFAULT ''::character varying NOT NULL,
    last_webfingered_at timestamp without time zone,
    inbox_url character varying DEFAULT ''::character varying NOT NULL,
    outbox_url character varying DEFAULT ''::character varying NOT NULL,
    shared_inbox_url character varying DEFAULT ''::character varying NOT NULL,
    followers_url character varying DEFAULT ''::character varying NOT NULL,
    protocol integer DEFAULT 0 NOT NULL,
    memorial boolean DEFAULT false NOT NULL,
    moved_to_account_id bigint,
    featured_collection_url character varying,
    fields jsonb,
    actor_type character varying,
    discoverable boolean,
    also_known_as character varying[]
);
CREATE TABLE public.session_activations (
    id serial8 NOT NULL PRIMARY KEY,
    session_id character varying NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    user_agent character varying DEFAULT ''::character varying NOT NULL,
    ip inet,
    access_token_id bigint,
    user_id bigint NOT NULL,
    web_push_subscription_id bigint
);
CREATE TABLE public.subscriptions (
    id serial8 NOT NULL PRIMARY KEY,
    callback_url character varying DEFAULT ''::character varying NOT NULL,
    secret character varying,
    expires_at timestamp without time zone,
    confirmed boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    last_successful_delivery_at timestamp without time zone,
    domain character varying,
    account_id bigint NOT NULL
);
CREATE TABLE public.scheduled_statuses (
    id serial8 NOT NULL PRIMARY KEY,
    account_id bigint,
    scheduled_at timestamp without time zone,
    params jsonb
);
CREATE TABLE public.schema_migrations (
    version character varying NOT NULL PRIMARY KEY
);
COPY public.account_conversations (id, account_id, conversation_id, participant_account_ids, status_ids, last_status_id, lock_version, unread) FROM stdin;
\.
COPY public.account_domain_blocks (id, domain, created_at, updated_at, account_id) FROM stdin;
\.
COPY public.account_moderation_notes (id, content, account_id, target_account_id, created_at, updated_at) FROM stdin;
\.
COPY public.account_pins (id, account_id, target_account_id, created_at, updated_at) FROM stdin;
\.
COPY public.account_stats (id, account_id, statuses_count, following_count, followers_count, created_at, updated_at, last_status_at) FROM stdin;
1	1	0	0	0	2019-02-06 11:03:11.679828	2019-02-06 11:03:11.679828	\N
2	2	0	0	0	2019-02-06 11:23:01.692073	2019-02-06 11:23:01.692073	\N
\.
COPY public.account_tag_stats (id, tag_id, accounts_count, hidden, created_at, updated_at) FROM stdin;
\.
COPY public.account_warning_presets (id, text, created_at, updated_at) FROM stdin;
\.
COPY public.account_warnings (id, account_id, target_account_id, action, text, created_at, updated_at) FROM stdin;
\.
COPY public.accounts (id, username, domain, secret, private_key, public_key, remote_url, salmon_url, hub_url, created_at, updated_at, note, display_name, uri, url, avatar_file_name, avatar_content_type, avatar_file_size, avatar_updated_at, header_file_name, header_content_type, header_file_size, header_updated_at, avatar_remote_url, subscription_expires_at, silenced, suspended, locked, header_remote_url, last_webfingered_at, inbox_url, outbox_url, shared_inbox_url, followers_url, protocol, memorial, moved_to_account_id, featured_collection_url, fields, actor_type, discoverable, also_known_as) FROM stdin;
1	admin	\N		-----BEGIN RSA PRIVATE KEY-----\nMIIEpQIBAAKCAQEAqn7P+Ml6xwCh2sT02YLhEGZCbopsaPFb2eBPBv/vZmT4PKk7\nSXicqJPQshcGIp9pu/oHioHLsYoThUzprhG5RhRKaZAheiZETgU9F4Cc4ySVFx26\nIC340YD/v1LQlaiZGYJ7dPesT6tq+B48V3HVyD9fNuFqelqlUql2MlxurhyXV7uI\n1SjyRQSl65nXoMprNKNnP0OcpGxOU1yZJzdziexo6qG6V+q7cV1VdeG7Y9VSpFUB\nbRh3ZQhWUiR+/Hmxh0HzAdogg/RF2FSsr6p+NX/u1sLSfN7OKDkSZYm9dc6x9Knh\nFbpz7vz2dsodWPQLeixA1RkuH4iZ791WxBlfwQIDAQABAoIBAQCaqEgNfO6jwD4S\nDiGxgViZoLlYPrbSh0ZzoFbvmZBXiPXpSPYf0ooBHXztX5dQJt0qCEd46/6TQRYu\nEDPVk/xFxrgtg/HqNPY28+eT/zXRkeiwPGYPNMSFfwf/TKcravHeQw+sbdLfvjZd\ndkf0Zq2vZVUAmoAVF07qahBu5Iv94TRY0n/nLYJG9AxEjVPlt8Vy3fb5sIIYzMO1\nwP3+BTpRHp0QmcQf5iH6ytK7jMNIokotzm5sQJBvMR2fY3d5BK8QevqWFvlnazox\nC4bn+Dd+DN+dYekL8QUJjXvqIsFYxZUk5RwL56anurVWB6zm593jsIVkdhl++0Vo\nUHsSb2HBAoGBANjkv72le/knx0Hqn6Ework2MTXs3KkfpxiVbzW8Gbzje9AiI0xT\nrM1Sj5eTPUTZwAQhjYM5sX4MmRGRcJ6FyOMpXI2IIIK4Gqjxvtrcvey94p9HkvDZ\nzQrvvZMLSr37VrtQlEpsJLg04i+bWsnq3OkhDowi+XIrewLLCvyZqckJAoGBAMk8\ncRZlRxLAXK26Vwxf9N5MPGrYg0LeLTCZGhMzxqdBAvkHyHI0Im5jySYonKgW2dAX\nzT2jnRMFJ4lXBOp/0+fb3ZJb3QY1Eu0K3t7TmGfXnZivslligHMVeH1edpocuQE8\nilWhRVnkev+yLDWma8GlOgBpeMPrYrxHMKb8Eab5AoGBAKVKxxlneTBrcT456Ud1\ngj12IFDBX3UAK17f223vGQpLrzryGUZ86k9boRTZ4DKNY/mB/I/KMwsl3K130oTs\n3ijIh8FQwb39QkwIV/QBkDhQidnrOP+WbN3t0OK0E1Tvq6x6/1gsTFuZ6dpwIeOJ\nuqtsRuLjcIjivA9n38qb7LnJAoGADxLM8b2CVmA8UPMNRCsH34LcX7B6HI8h2WsO\nbfPJ5ItVGqw/knZfQd+NmKMgIOMdS54MzJot0Nfo/zuabapHiC2K6kShSK6/DSxs\nR0qYNucKsf4vIMzlDDnGfbWOsrqGDRao0gMze6lGoVKKRzaBCc9DifZcimheS/YV\nKdzlwZkCgYEAqLe/eqm6Vs9/MTvoB+DUsdt6p8ET5xgfa1eqnFEPShkk39CNc56U\nHUtEjUu4kiHBlCk5poDfk1DMkh4UPJPI/GyI3oRydypKTuApJXWkYlFYjEaa6vLM\nlGmWl60h6FgDoJH/rATdX+1bUhnSltDYQQOvLg6Tpzt6iVdoP42u52c=\n-----END RSA PRIVATE KEY-----\n	-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqn7P+Ml6xwCh2sT02YLh\nEGZCbopsaPFb2eBPBv/vZmT4PKk7SXicqJPQshcGIp9pu/oHioHLsYoThUzprhG5\nRhRKaZAheiZETgU9F4Cc4ySVFx26IC340YD/v1LQlaiZGYJ7dPesT6tq+B48V3HV\nyD9fNuFqelqlUql2MlxurhyXV7uI1SjyRQSl65nXoMprNKNnP0OcpGxOU1yZJzdz\niexo6qG6V+q7cV1VdeG7Y9VSpFUBbRh3ZQhWUiR+/Hmxh0HzAdogg/RF2FSsr6p+\nNX/u1sLSfN7OKDkSZYm9dc6x9KnhFbpz7vz2dsodWPQLeixA1RkuH4iZ791WxBlf\nwQIDAQAB\n-----END PUBLIC KEY-----\n				2019-02-06 11:03:11.576318	2019-02-06 11:03:11.576318				\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	f	f	f		\N					0	f	\N	\N	\N	\N	\N	\N
2	tailiu	\N		-----BEGIN RSA PRIVATE KEY-----\nMIIEpQIBAAKCAQEAq7+iw5pfySDgHrRnBWqLOfkLQAe0AmO8dH27Jl/squPruOMu\nycHpwtxjjb7AqlFzLjJTiERcPGY5qyyYa/MPJo29ydPR3v52UyZyDMgG21OOBGPs\n93z6a8bQorWSpcCmQuLcpXfDHlMHp6OOcp5s38qnAYXQdkjBbLX+n1vYz8h7Wq1r\nqgSXt3WSPjfbwI4NYV2oiI/axoPA7RxTVGPmEogqQEE7ruF0FAUOLS8Vn2cqKZ1i\n1dSWCoixtmNyqYm3ZsTVjNjp5Dyey+GSi08bJXpBfCWMDF5Q+RsXrdyqzt0TN8bR\nvuRUqsR8nrTiSikv6d7gjSUd1F2B2xdQizwVqwIDAQABAoIBAHfnkICy9BB87TtC\ng3IakhzlK9+GATxx2Q4TAWenLJmaCeXIZc/hm4u5RZ+d/vBpcNpdtEe8QhDA5Z4F\nwlwLExa9ejS/txPR31Xpz1HxDChvSwTxpmyaSlKYOOx/i2RA/VJRA+5ZkFVJblyx\nKTAIPsZ2uuWrQIB0BuYYoS0seE+G6R7JnYwQ0pfc//hqzCY+y3rbL1otLeuHmH+W\nvua4vV/aCBD03ViIRa+9dMFsZEdDi0m1+VO6OBbpX69v/A8TBqgRiX9R4Sb4kT4O\nX4Xp8UChFikJJ0AZm411SU8BHo3l8yNSiyr4n9L60yllaec/HDopiu6VCj6doM9k\nJL46zeECgYEA0n540g/H6j8lwQoc4RDU9PlfJWmwvI/HPipRlhvogCE+DN595L6z\noBIU+c12SZExPpBs1ndFgPecoSR9vJw2fOfVYtru27E3tXIgSC3y9pK126EVN1oS\nDFskaSRAswgYJeKbya4Oc5TbRDIKcVQfBf1lWEgYd8h096TEnq2Wj6MCgYEA0ODY\nZsRPdmBbDZym4H7OFv+KWEETx2gnxPp05VGBaCenu6yNvInBrEJ3O6WZRpdK40bp\nYi8GcFjAocU6nVaGnaK1pxVD2tTGtvrL5rax5gzOdBmqgkdWIOzfbJG4RUNv9Bh/\nBeLpZt4Z+cTUrWKI89s7V4B1cOIoe/6t+/z1olkCgYEAkfv8t1MShzc8a+EjnkQa\nLbw1bLEcTeo5eLfI1Z6NZS+o5Sv5jAdmdIGV4pnIi8USrh1kHmmh3ovcKTYxrfl5\nIK94opLMTbletYxtLyIO+0tMrQHOwRDKq58aZYErDf9zH/NFsF3yz95RI77A11BM\nI89V1iBKN+jilk3Dv3kMjpkCgYEAowCD/5Z8yE0jYTj5RUHPhEUA6iRG0hsmxeIJ\nrRbw3J3tmFhs90+tUsc/ks2FEoBoUXp6EEPQS4YHNXbbagMm5AcgqPXAURowxIRs\n8Gtr4rHlvtZ0qFwRC3quVGRXH74jtKIVJjvQlpUGQlLnATNe2qYf5gX6IBBtNW4m\nyfm6mmkCgYEAhOdGYtrpB9htvF8yhIqGJzwjuOwhj6aWwXqt7LS7VAxA9zIee2An\nNBKoj9TP/PTZUK4V/pQJwfHVct+SGCiuESm3w2CYhBu0S2smI+PCUZAuqvBYuaX1\nDNki1WTQbMePBQmwXTpqVVG3dzG2MtIrnY2WvFAC8oswJEjyjO+1Xo0=\n-----END RSA PRIVATE KEY-----\n	-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAq7+iw5pfySDgHrRnBWqL\nOfkLQAe0AmO8dH27Jl/squPruOMuycHpwtxjjb7AqlFzLjJTiERcPGY5qyyYa/MP\nJo29ydPR3v52UyZyDMgG21OOBGPs93z6a8bQorWSpcCmQuLcpXfDHlMHp6OOcp5s\n38qnAYXQdkjBbLX+n1vYz8h7Wq1rqgSXt3WSPjfbwI4NYV2oiI/axoPA7RxTVGPm\nEogqQEE7ruF0FAUOLS8Vn2cqKZ1i1dSWCoixtmNyqYm3ZsTVjNjp5Dyey+GSi08b\nJXpBfCWMDF5Q+RsXrdyqzt0TN8bRvuRUqsR8nrTiSikv6d7gjSUd1F2B2xdQizwV\nqwIDAQAB\n-----END PUBLIC KEY-----\n				2019-02-06 11:23:01.609905	2019-02-06 11:23:01.609905				\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	\N	f	f	f		\N					0	f	\N	\N	\N	\N	\N	\N
\.
COPY public.accounts_tags (account_id, tag_id) FROM stdin;
\.
COPY public.admin_action_logs (id, account_id, action, target_type, target_id, recorded_changes, created_at, updated_at) FROM stdin;
\.
COPY public.ar_internal_metadata (key, value, created_at, updated_at) FROM stdin;
environment	development	2019-02-06 11:02:58.134562	2019-02-06 11:02:58.134562
\.
COPY public.backups (id, user_id, dump_file_name, dump_content_type, dump_file_size, dump_updated_at, processed, created_at, updated_at) FROM stdin;
\.
COPY public.blocks (id, created_at, updated_at, account_id, target_account_id, uri) FROM stdin;
\.
COPY public.conversation_mutes (id, conversation_id, account_id) FROM stdin;
\.
COPY public.conversations (id, uri, created_at, updated_at) FROM stdin;
\.
COPY public.custom_emojis (id, shortcode, domain, image_file_name, image_content_type, image_file_size, image_updated_at, created_at, updated_at, disabled, uri, image_remote_url, visible_in_picker) FROM stdin;
\.
COPY public.custom_filters (id, account_id, expires_at, phrase, context, irreversible, created_at, updated_at, whole_word) FROM stdin;
\.
COPY public.domain_blocks (id, domain, created_at, updated_at, severity, reject_media, reject_reports) FROM stdin;
\.
COPY public.email_domain_blocks (id, domain, created_at, updated_at) FROM stdin;
\.
COPY public.favourites (id, created_at, updated_at, account_id, status_id) FROM stdin;
\.
COPY public.follow_requests (id, created_at, updated_at, account_id, target_account_id, show_reblogs, uri) FROM stdin;
\.
COPY public.follows (id, created_at, updated_at, account_id, target_account_id, show_reblogs, uri) FROM stdin;
\.
COPY public.identities (id, provider, uid, created_at, updated_at, user_id) FROM stdin;
\.
COPY public.imports (id, type, approved, created_at, updated_at, data_file_name, data_content_type, data_file_size, data_updated_at, account_id) FROM stdin;
\.
COPY public.invites (id, user_id, code, expires_at, max_uses, uses, created_at, updated_at, autofollow) FROM stdin;
\.
COPY public.list_accounts (id, list_id, account_id, follow_id) FROM stdin;
\.
COPY public.lists (id, account_id, title, created_at, updated_at) FROM stdin;
\.
COPY public.media_attachments (id, status_id, file_file_name, file_content_type, file_file_size, file_updated_at, remote_url, created_at, updated_at, shortcode, type, file_meta, account_id, description, scheduled_status_id) FROM stdin;
\.
COPY public.mentions (id, status_id, created_at, updated_at, account_id, silent) FROM stdin;
\.
COPY public.mutes (id, created_at, updated_at, account_id, target_account_id, hide_notifications) FROM stdin;
\.
COPY public.notifications (id, activity_id, activity_type, created_at, updated_at, account_id, from_account_id) FROM stdin;
\.
COPY public.oauth_access_grants (id, token, expires_in, redirect_uri, created_at, revoked_at, scopes, application_id, resource_owner_id) FROM stdin;
\.
COPY public.oauth_access_tokens (id, token, refresh_token, expires_in, revoked_at, created_at, scopes, application_id, resource_owner_id) FROM stdin;
\.
COPY public.oauth_applications (id, name, uid, secret, redirect_uri, scopes, created_at, updated_at, superapp, website, owner_type, owner_id, confidential) FROM stdin;
1	Web	03026bcbc926a65379e545719e987222284c7bc3f95584de4b9afc9feb30a6d2	f821fde75bf0d1e1b8e90473b990b108a402cbcd18969d7c13ba9b3e230358f5	urn:ietf:wg:oauth:2.0:oob	read write follow	2019-02-06 11:03:11.507833	2019-02-06 11:03:11.507833	t	\N	\N	\N	t
\.
COPY public.pghero_space_stats (id, database, schema, relation, size, captured_at) FROM stdin;
\.
COPY public.preview_cards (id, url, title, description, image_file_name, image_content_type, image_file_size, image_updated_at, type, html, author_name, author_url, provider_name, provider_url, width, height, created_at, updated_at, embed_url) FROM stdin;
\.
COPY public.preview_cards_statuses (preview_card_id, status_id) FROM stdin;
\.
COPY public.relays (id, inbox_url, follow_activity_id, created_at, updated_at, state) FROM stdin;
\.
COPY public.report_notes (id, content, report_id, account_id, created_at, updated_at) FROM stdin;
\.
COPY public.reports (id, status_ids, comment, action_taken, created_at, updated_at, account_id, action_taken_by_account_id, target_account_id, assigned_account_id) FROM stdin;
\.
COPY public.scheduled_statuses (id, account_id, scheduled_at, params) FROM stdin;
\.
COPY public.schema_migrations (version) FROM stdin;
20190117114553
20160220211917
20161222201034
20180814171349
20170127165745
20170322162804
20181204193439
20181024224956
20170520145338
20170112154826
20161027172456
20170418160728
20181007025445
20170506235850
20171201000000
20181219235220
20171212195226
20160223165723
20171028221157
20161203164520
20171125031751
20171010023049
20171010025614
20180609104432
20170114203041
20161205214545
20171130000000
20180514140000
20170217012631
20161128103007
20171226094803
20160221003621
20170317193015
20180402040909
20170403172249
20171020084748
20170330164118
20161009120834
20170901141119
20170718211102
20161119211120
20171107143332
20181203003808
20161222204147
20170205175257
20171005171936
20170720000000
20170125145934
20161104173623
20181116173541
20160322193748
20170623152212
20171125190735
20180707154237
20180812173710
20170624134742
20160826155805
20170917153509
20170424003227
20170424112722
20161006213403
20180711152640
20170405112956
20170409170753
20160926213048
20171125185353
20170928082043
20170123203248
20160314164231
20170920032311
20170901142658
20180808175627
20170610000000
20181226021420
20160221003140
20171125024930
20160222143943
20181026034033
20190103124754
20170507000211
20180617162849
20161130185319
20160223164502
20180615122121
20171114080328
20161130142058
20180211015820
20160325130944
20160312193225
20170913000752
20170918125918
20170425131920
20171116161857
20181127130500
20170604144747
20171122120436
20160919221059
20170114194937
20181203021853
20160222122600
20171119172437
20180206000000
20160227230233
20170303212857
20181213185533
20170214110202
20170301222600
20181204215309
20170601210557
20180109143959
20160223162837
20161122163057
20180410204633
20170927215609
20161202132159
20161003142332
20180514130000
20180608213548
20170508230434
20180510214435
20161105130633
20180628181026
20170713175513
20170414080609
20180416210259
20170609145826
20160220174730
20170713190709
20170711225116
20161116162355
20181207011115
20170823162448
20170427011934
20161003145426
20180106000232
20170318214217
20171107143624
20170330163835
20170105224407
20171006142024
20181213184704
20160306172223
20170109120109
20171005102658
20160920003904
20170920024819
20170414132105
20180510230049
20170905165803
20170209184350
20170507141759
20181010141500
20161221152630
20170824103029
20170330021336
20171114231651
20160223171800
20170829215220
20171109012327
20180528141303
20180616192031
20180204034416
20170425202925
20180812123222
20160223165855
20160905150353
20170322021028
20160224223247
20170606113804
20171118012443
20181017170937
20160305115639
20170714184731
20170516072309
20170304202101
20170625140443
20170406215816
20170924022025
20180820232245
20180402031200
20170322143850
20190103124649
20170905044538
20180310000000
20171129172043
20170713112503
20170716191202
20180929222014
20161123093447
20170123162658
20180304013859
20170129000348
20180812162710
20160316103650
20181018205649
20181116165755
20180506221944
20170423005413
20170119214911
20180813113448
20181116184611
\.
COPY public.session_activations (id, session_id, created_at, updated_at, user_agent, ip, access_token_id, user_id, web_push_subscription_id) FROM stdin;
\.
COPY public.settings (id, var, value, thing_type, created_at, updated_at, thing_id) FROM stdin;
\.
COPY public.site_uploads (id, var, file_file_name, file_content_type, file_file_size, file_updated_at, meta, created_at, updated_at) FROM stdin;
\.
COPY public.status_pins (id, account_id, status_id, created_at, updated_at) FROM stdin;
\.
COPY public.status_stats (id, status_id, replies_count, reblogs_count, favourites_count, created_at, updated_at) FROM stdin;
\.
COPY public.statuses (id, uri, text, created_at, updated_at, in_reply_to_id, reblog_of_id, url, sensitive, visibility, spoiler_text, reply, language, conversation_id, local, account_id, application_id, in_reply_to_account_id) FROM stdin;
\.
COPY public.statuses_tags (status_id, tag_id) FROM stdin;
\.
COPY public.stream_entries (id, activity_id, activity_type, created_at, updated_at, hidden, account_id) FROM stdin;
\.
COPY public.subscriptions (id, callback_url, secret, expires_at, confirmed, created_at, updated_at, last_successful_delivery_at, domain, account_id) FROM stdin;
\.
COPY public.tags (id, name, created_at, updated_at) FROM stdin;
\.
COPY public.tombstones (id, account_id, uri, created_at, updated_at) FROM stdin;
\.
COPY public.users (id, email, created_at, updated_at, encrypted_password, reset_password_token, reset_password_sent_at, remember_created_at, sign_in_count, current_sign_in_at, last_sign_in_at, current_sign_in_ip, last_sign_in_ip, admin, confirmation_token, confirmed_at, confirmation_sent_at, unconfirmed_email, locale, encrypted_otp_secret, encrypted_otp_secret_iv, encrypted_otp_secret_salt, consumed_timestep, otp_required_for_login, last_emailed_at, otp_backup_codes, filtered_languages, account_id, disabled, moderator, invite_id, remember_token, chosen_languages, created_by_application_id) FROM stdin;
2	tai.liu@nyu.edu	2019-02-06 11:23:01.67976	2019-02-06 11:23:01.67976	$2a$10$zJp1RkXUY079LXfEPfNmgOVa9CEzdAmaq/CkCnZbwOHdONBUvb7TK	\N	\N	\N	0	\N	\N	\N	\N	f	XFtUso4-EpjTt-6xzXx1	\N	2019-02-06 11:23:01.679862	\N	en	\N	\N	\N	\N	f	\N	\N	{}	2	f	f	\N	\N	\N	\N
1	admin@localhost:3000	2019-02-06 11:03:11.81281	2019-02-07 12:07:58.355697	$2a$10$gkI.zeitRHBK8r0AZbsdM.cwEfHydKIHaRo93zYodUs1meLczwHXm	\N	\N	\N	2	2019-02-07 12:07:03.693727	2019-02-06 11:16:16.493899	127.0.0.1	127.0.0.1	t	\N	2019-02-06 11:03:11.697185	\N	\N	\N	\N	\N	\N	\N	f	\N	\N	{}	1	f	f	\N	\N	\N	\N
\.
COPY public.web_push_subscriptions (id, endpoint, key_p256dh, key_auth, data, created_at, updated_at, access_token_id, user_id) FROM stdin;
\.
COPY public.web_settings (id, data, created_at, updated_at, user_id) FROM stdin;
1	{"onboarded":false,"notifications":{"alerts":{"follow":true,"favourite":true,"reblog":true,"mention":true},"quickFilter":{"active":"all","show":true,"advanced":false},"shows":{"follow":true,"favourite":true,"reblog":true,"mention":true},"sounds":{"follow":true,"favourite":true,"reblog":true,"mention":true}},"public":{"regex":{"body":""}},"direct":{"regex":{"body":""}},"community":{"regex":{"body":""}},"skinTone":1,"trends":{"show":true},"columns":[{"id":"COMPOSE","uuid":"8f39d61f-cd43-419c-b52a-e34c0e12bf42","params":{}},{"id":"HOME","uuid":"66a9929d-e900-4d2f-a6d2-a3b5c9ce8b61","params":{}},{"id":"NOTIFICATIONS","uuid":"ab738ed7-73ff-4630-bff4-f26710bd15ab","params":{}}],"introductionVersion":20181216044202,"home":{"shows":{"reblog":true,"reply":true},"regex":{"body":""}}}	2019-02-06 11:16:34.324871	2019-02-06 11:16:34.324871	1
\.
CREATE UNIQUE INDEX account_activity ON public.notifications USING btree (account_id, activity_id, activity_type);
CREATE INDEX hashtag_search_index ON public.tags USING btree (name, text_pattern_ops);
CREATE INDEX index_account_conversations_on_account_id ON public.account_conversations USING btree (account_id);
CREATE INDEX index_account_conversations_on_conversation_id ON public.account_conversations USING btree (conversation_id);
CREATE UNIQUE INDEX index_account_domain_blocks_on_account_id_and_domain ON public.account_domain_blocks USING btree (account_id, domain);
CREATE INDEX index_account_moderation_notes_on_account_id ON public.account_moderation_notes USING btree (account_id);
CREATE INDEX index_account_moderation_notes_on_target_account_id ON public.account_moderation_notes USING btree (target_account_id);
CREATE INDEX index_account_pins_on_account_id ON public.account_pins USING btree (account_id);
CREATE UNIQUE INDEX index_account_pins_on_account_id_and_target_account_id ON public.account_pins USING btree (account_id, target_account_id);
CREATE INDEX index_account_pins_on_target_account_id ON public.account_pins USING btree (target_account_id);
CREATE UNIQUE INDEX index_account_stats_on_account_id ON public.account_stats USING btree (account_id);
CREATE UNIQUE INDEX index_account_tag_stats_on_tag_id ON public.account_tag_stats USING btree (tag_id);
CREATE INDEX index_account_warnings_on_account_id ON public.account_warnings USING btree (account_id);
CREATE INDEX index_account_warnings_on_target_account_id ON public.account_warnings USING btree (target_account_id);
CREATE INDEX index_accounts_on_moved_to_account_id ON public.accounts USING btree (moved_to_account_id);
CREATE INDEX index_accounts_on_uri ON public.accounts USING btree (uri);
CREATE INDEX index_accounts_on_url ON public.accounts USING btree (url);
CREATE UNIQUE INDEX index_accounts_on_username_and_domain_lower ON public.accounts USING btree (username, domain);
CREATE INDEX index_accounts_tags_on_account_id_and_tag_id ON public.accounts_tags USING btree (account_id, tag_id);
CREATE UNIQUE INDEX index_accounts_tags_on_tag_id_and_account_id ON public.accounts_tags USING btree (tag_id, account_id);
CREATE INDEX index_admin_action_logs_on_account_id ON public.admin_action_logs USING btree (account_id);
CREATE INDEX index_admin_action_logs_on_target_type_and_target_id ON public.admin_action_logs USING btree (target_type, target_id);
CREATE UNIQUE INDEX index_blocks_on_account_id_and_target_account_id ON public.blocks USING btree (account_id, target_account_id);
CREATE INDEX index_blocks_on_target_account_id ON public.blocks USING btree (target_account_id);
CREATE UNIQUE INDEX index_conversation_mutes_on_account_id_and_conversation_id ON public.conversation_mutes USING btree (account_id, conversation_id);
CREATE UNIQUE INDEX index_conversations_on_uri ON public.conversations USING btree (uri);
CREATE UNIQUE INDEX index_custom_emojis_on_shortcode_and_domain ON public.custom_emojis USING btree (shortcode, domain);
CREATE INDEX index_custom_filters_on_account_id ON public.custom_filters USING btree (account_id);
CREATE UNIQUE INDEX index_domain_blocks_on_domain ON public.domain_blocks USING btree (domain);
CREATE UNIQUE INDEX index_email_domain_blocks_on_domain ON public.email_domain_blocks USING btree (domain);
CREATE INDEX index_favourites_on_account_id_and_id ON public.favourites USING btree (account_id, id);
CREATE UNIQUE INDEX index_favourites_on_account_id_and_status_id ON public.favourites USING btree (account_id, status_id);
CREATE INDEX index_favourites_on_status_id ON public.favourites USING btree (status_id);
CREATE UNIQUE INDEX index_follow_requests_on_account_id_and_target_account_id ON public.follow_requests USING btree (account_id, target_account_id);
CREATE UNIQUE INDEX index_follows_on_account_id_and_target_account_id ON public.follows USING btree (account_id, target_account_id);
CREATE INDEX index_follows_on_target_account_id ON public.follows USING btree (target_account_id);
CREATE INDEX index_identities_on_user_id ON public.identities USING btree (user_id);
CREATE UNIQUE INDEX index_invites_on_code ON public.invites USING btree (code);
CREATE INDEX index_invites_on_user_id ON public.invites USING btree (user_id);
CREATE UNIQUE INDEX index_list_accounts_on_account_id_and_list_id ON public.list_accounts USING btree (account_id, list_id);
CREATE INDEX index_list_accounts_on_follow_id ON public.list_accounts USING btree (follow_id);
CREATE INDEX index_list_accounts_on_list_id_and_account_id ON public.list_accounts USING btree (list_id, account_id);
CREATE INDEX index_lists_on_account_id ON public.lists USING btree (account_id);
CREATE INDEX index_media_attachments_on_account_id ON public.media_attachments USING btree (account_id);
CREATE INDEX index_media_attachments_on_scheduled_status_id ON public.media_attachments USING btree (scheduled_status_id);
CREATE UNIQUE INDEX index_media_attachments_on_shortcode ON public.media_attachments USING btree (shortcode);
CREATE INDEX index_media_attachments_on_status_id ON public.media_attachments USING btree (status_id);
CREATE UNIQUE INDEX index_mentions_on_account_id_and_status_id ON public.mentions USING btree (account_id, status_id);
CREATE INDEX index_mentions_on_status_id ON public.mentions USING btree (status_id);
CREATE UNIQUE INDEX index_mutes_on_account_id_and_target_account_id ON public.mutes USING btree (account_id, target_account_id);
CREATE INDEX index_mutes_on_target_account_id ON public.mutes USING btree (target_account_id);
CREATE INDEX index_notifications_on_account_id_and_id ON public.notifications USING btree (account_id, id DESC);
CREATE INDEX index_notifications_on_activity_id_and_activity_type ON public.notifications USING btree (activity_id, activity_type);
CREATE INDEX index_notifications_on_from_account_id ON public.notifications USING btree (from_account_id);
CREATE INDEX index_oauth_access_grants_on_resource_owner_id ON public.oauth_access_grants USING btree (resource_owner_id);
CREATE UNIQUE INDEX index_oauth_access_grants_on_token ON public.oauth_access_grants USING btree (token);
CREATE UNIQUE INDEX index_oauth_access_tokens_on_refresh_token ON public.oauth_access_tokens USING btree (refresh_token);
CREATE INDEX index_oauth_access_tokens_on_resource_owner_id ON public.oauth_access_tokens USING btree (resource_owner_id);
CREATE UNIQUE INDEX index_oauth_access_tokens_on_token ON public.oauth_access_tokens USING btree (token);
CREATE INDEX index_oauth_applications_on_owner_id_and_owner_type ON public.oauth_applications USING btree (owner_id, owner_type);
CREATE UNIQUE INDEX index_oauth_applications_on_uid ON public.oauth_applications USING btree (uid);
CREATE INDEX index_pghero_space_stats_on_database_and_captured_at ON public.pghero_space_stats USING btree (database, captured_at);
CREATE UNIQUE INDEX index_preview_cards_on_url ON public.preview_cards USING btree (url);
CREATE INDEX index_preview_cards_statuses_on_status_id_and_preview_card_id ON public.preview_cards_statuses USING btree (status_id, preview_card_id);
CREATE INDEX index_report_notes_on_account_id ON public.report_notes USING btree (account_id);
CREATE INDEX index_report_notes_on_report_id ON public.report_notes USING btree (report_id);
CREATE INDEX index_reports_on_account_id ON public.reports USING btree (account_id);
CREATE INDEX index_reports_on_target_account_id ON public.reports USING btree (target_account_id);
CREATE INDEX index_scheduled_statuses_on_account_id ON public.scheduled_statuses USING btree (account_id);
CREATE INDEX index_scheduled_statuses_on_scheduled_at ON public.scheduled_statuses USING btree (scheduled_at);
CREATE INDEX index_session_activations_on_access_token_id ON public.session_activations USING btree (access_token_id);
CREATE UNIQUE INDEX index_session_activations_on_session_id ON public.session_activations USING btree (session_id);
CREATE INDEX index_session_activations_on_user_id ON public.session_activations USING btree (user_id);
CREATE UNIQUE INDEX index_settings_on_thing_type_and_thing_id_and_var ON public.settings USING btree (thing_type, thing_id, var);
CREATE UNIQUE INDEX index_site_uploads_on_var ON public.site_uploads USING btree (var);
CREATE UNIQUE INDEX index_status_pins_on_account_id_and_status_id ON public.status_pins USING btree (account_id, status_id);
CREATE UNIQUE INDEX index_status_stats_on_status_id ON public.status_stats USING btree (status_id);
CREATE INDEX index_statuses_20180106 ON public.statuses USING btree (account_id, id DESC, visibility, updated_at);
CREATE INDEX index_statuses_on_in_reply_to_account_id ON public.statuses USING btree (in_reply_to_account_id);
CREATE INDEX index_statuses_on_in_reply_to_id ON public.statuses USING btree (in_reply_to_id);
CREATE INDEX index_statuses_on_reblog_of_id_and_account_id ON public.statuses USING btree (reblog_of_id, account_id);
CREATE UNIQUE INDEX index_statuses_on_uri ON public.statuses USING btree (uri);
CREATE INDEX index_statuses_tags_on_status_id ON public.statuses_tags USING btree (status_id);
CREATE UNIQUE INDEX index_statuses_tags_on_tag_id_and_status_id ON public.statuses_tags USING btree (tag_id, status_id);
CREATE INDEX index_stream_entries_on_account_id_and_activity_type_and_id ON public.stream_entries USING btree (account_id, activity_type, id);
CREATE INDEX index_stream_entries_on_activity_id_and_activity_type ON public.stream_entries USING btree (activity_id, activity_type);
CREATE UNIQUE INDEX index_subscriptions_on_account_id_and_callback_url ON public.subscriptions USING btree (account_id, callback_url);
CREATE UNIQUE INDEX index_tags_on_name ON public.tags USING btree (name);
CREATE INDEX index_tombstones_on_account_id ON public.tombstones USING btree (account_id);
CREATE INDEX index_tombstones_on_uri ON public.tombstones USING btree (uri);
CREATE UNIQUE INDEX index_unique_conversations ON public.account_conversations USING btree (account_id, conversation_id, participant_account_ids);
CREATE INDEX index_users_on_account_id ON public.users USING btree (account_id);
CREATE UNIQUE INDEX index_users_on_confirmation_token ON public.users USING btree (confirmation_token);
CREATE INDEX index_users_on_created_by_application_id ON public.users USING btree (created_by_application_id);
CREATE UNIQUE INDEX index_users_on_email ON public.users USING btree (email);
CREATE UNIQUE INDEX index_users_on_reset_password_token ON public.users USING btree (reset_password_token);
CREATE INDEX index_web_push_subscriptions_on_access_token_id ON public.web_push_subscriptions USING btree (access_token_id);
CREATE INDEX index_web_push_subscriptions_on_user_id ON public.web_push_subscriptions USING btree (user_id);
CREATE UNIQUE INDEX index_web_settings_on_user_id ON public.web_settings USING btree (user_id);
CREATE INDEX search_index ON public.accounts USING gin ((((setweight(to_tsvector('simple'::regconfig, (display_name)::text), 'A'::"char") || setweight(to_tsvector('simple'::regconfig, (username)::text), 'B'::"char")) || setweight(to_tsvector('simple'::regconfig, (COALESCE(domain, ''::character varying))::text), 'C'::"char"))));
ALTER TABLE ONLY public.web_settings    ADD CONSTRAINT fk_11910667b2 FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.account_domain_blocks    ADD CONSTRAINT fk_206c6029bd FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.conversation_mutes    ADD CONSTRAINT fk_225b4212bb FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.statuses_tags    ADD CONSTRAINT fk_3081861e21 FOREIGN KEY (tag_id) REFERENCES public.tags(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.follows    ADD CONSTRAINT fk_32ed1b5560 FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.oauth_access_grants    ADD CONSTRAINT fk_34d54b0a33 FOREIGN KEY (application_id) REFERENCES public.oauth_applications(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.blocks    ADD CONSTRAINT fk_4269e03e65 FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.reports    ADD CONSTRAINT fk_4b81f7522c FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.users    ADD CONSTRAINT fk_50500f500d FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.stream_entries    ADD CONSTRAINT fk_5659b17554 FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.favourites    ADD CONSTRAINT fk_5eb6c2b873 FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.oauth_access_grants    ADD CONSTRAINT fk_63b044929b FOREIGN KEY (resource_owner_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.imports    ADD CONSTRAINT fk_6db1b6e408 FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.follows    ADD CONSTRAINT fk_745ca29eac FOREIGN KEY (target_account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.follow_requests    ADD CONSTRAINT fk_76d644b0e7 FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.follow_requests    ADD CONSTRAINT fk_9291ec025d FOREIGN KEY (target_account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.blocks    ADD CONSTRAINT fk_9571bfabc1 FOREIGN KEY (target_account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.session_activations    ADD CONSTRAINT fk_957e5bda89 FOREIGN KEY (access_token_id) REFERENCES public.oauth_access_tokens(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.media_attachments    ADD CONSTRAINT fk_96dd81e81b FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.mentions    ADD CONSTRAINT fk_970d43f9d1 FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.subscriptions    ADD CONSTRAINT fk_9847d1cbb5 FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.statuses    ADD CONSTRAINT fk_9bda1543f7 FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.oauth_applications    ADD CONSTRAINT fk_b0988c7c0a FOREIGN KEY (owner_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.favourites    ADD CONSTRAINT fk_b0e856845e FOREIGN KEY (status_id) REFERENCES public.statuses(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.mutes    ADD CONSTRAINT fk_b8d8daf315 FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.reports    ADD CONSTRAINT fk_bca45b75fd FOREIGN KEY (action_taken_by_account_id) REFERENCES public.accounts(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.identities    ADD CONSTRAINT fk_bea040f377 FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.notifications    ADD CONSTRAINT fk_c141c8ee55 FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.statuses    ADD CONSTRAINT fk_c7fa917661 FOREIGN KEY (in_reply_to_account_id) REFERENCES public.accounts(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.status_pins    ADD CONSTRAINT fk_d4cb435b62 FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.session_activations    ADD CONSTRAINT fk_e5fda67334 FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.oauth_access_tokens    ADD CONSTRAINT fk_e84df68546 FOREIGN KEY (resource_owner_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.reports    ADD CONSTRAINT fk_eb37af34f0 FOREIGN KEY (target_account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.mutes    ADD CONSTRAINT fk_eecff219ea FOREIGN KEY (target_account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.oauth_access_tokens    ADD CONSTRAINT fk_f5fc4c1ee3 FOREIGN KEY (application_id) REFERENCES public.oauth_applications(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.notifications    ADD CONSTRAINT fk_fbd6b0bf9e FOREIGN KEY (from_account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.backups    ADD CONSTRAINT fk_rails_096669d221 FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.account_conversations    ADD CONSTRAINT fk_rails_1491654f9f FOREIGN KEY (conversation_id) REFERENCES public.conversations(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.account_tag_stats    ADD CONSTRAINT fk_rails_1fa34bab2d FOREIGN KEY (tag_id) REFERENCES public.tags(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.account_stats    ADD CONSTRAINT fk_rails_215bb31ff1 FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.accounts    ADD CONSTRAINT fk_rails_2320833084 FOREIGN KEY (moved_to_account_id) REFERENCES public.accounts(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.scheduled_statuses    ADD CONSTRAINT fk_rails_23bd9018f9 FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.statuses    ADD CONSTRAINT fk_rails_256483a9ab FOREIGN KEY (reblog_of_id) REFERENCES public.statuses(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.media_attachments    ADD CONSTRAINT fk_rails_31fc5aeef1 FOREIGN KEY (scheduled_status_id) REFERENCES public.scheduled_statuses(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.lists    ADD CONSTRAINT fk_rails_3853b78dac FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.media_attachments    ADD CONSTRAINT fk_rails_3ec0cfdd70 FOREIGN KEY (status_id) REFERENCES public.statuses(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.account_moderation_notes    ADD CONSTRAINT fk_rails_3f8b75089b FOREIGN KEY (account_id) REFERENCES public.accounts(id);
ALTER TABLE ONLY public.list_accounts    ADD CONSTRAINT fk_rails_40f9cc29f1 FOREIGN KEY (follow_id) REFERENCES public.follows(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.status_stats    ADD CONSTRAINT fk_rails_4a247aac42 FOREIGN KEY (status_id) REFERENCES public.statuses(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.reports    ADD CONSTRAINT fk_rails_4e7a498fb4 FOREIGN KEY (assigned_account_id) REFERENCES public.accounts(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.mentions    ADD CONSTRAINT fk_rails_59edbe2887 FOREIGN KEY (status_id) REFERENCES public.statuses(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.conversation_mutes    ADD CONSTRAINT fk_rails_5ab139311f FOREIGN KEY (conversation_id) REFERENCES public.conversations(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.status_pins    ADD CONSTRAINT fk_rails_65c05552f1 FOREIGN KEY (status_id) REFERENCES public.statuses(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.account_conversations    ADD CONSTRAINT fk_rails_6f5278b6e9 FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.web_push_subscriptions    ADD CONSTRAINT fk_rails_751a9f390b FOREIGN KEY (access_token_id) REFERENCES public.oauth_access_tokens(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.report_notes    ADD CONSTRAINT fk_rails_7fa83a61eb FOREIGN KEY (report_id) REFERENCES public.reports(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.list_accounts    ADD CONSTRAINT fk_rails_85fee9d6ab FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.custom_filters    ADD CONSTRAINT fk_rails_8b8d786993 FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.users    ADD CONSTRAINT fk_rails_8fb2a43e88 FOREIGN KEY (invite_id) REFERENCES public.invites(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.statuses    ADD CONSTRAINT fk_rails_94a6f70399 FOREIGN KEY (in_reply_to_id) REFERENCES public.statuses(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.account_pins    ADD CONSTRAINT fk_rails_a176e26c37 FOREIGN KEY (target_account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.account_warnings    ADD CONSTRAINT fk_rails_a65a1bf71b FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.admin_action_logs    ADD CONSTRAINT fk_rails_a7667297fa FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.account_warnings    ADD CONSTRAINT fk_rails_a7ebbb1e37 FOREIGN KEY (target_account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.web_push_subscriptions    ADD CONSTRAINT fk_rails_b006f28dac FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.report_notes    ADD CONSTRAINT fk_rails_cae66353f3 FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.account_pins    ADD CONSTRAINT fk_rails_d44979e5dd FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.account_moderation_notes    ADD CONSTRAINT fk_rails_dd62ed5ac3 FOREIGN KEY (target_account_id) REFERENCES public.accounts(id);
ALTER TABLE ONLY public.statuses_tags    ADD CONSTRAINT fk_rails_df0fe11427 FOREIGN KEY (status_id) REFERENCES public.statuses(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.list_accounts    ADD CONSTRAINT fk_rails_e54e356c88 FOREIGN KEY (list_id) REFERENCES public.lists(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.users    ADD CONSTRAINT fk_rails_ecc9536e7c FOREIGN KEY (created_by_application_id) REFERENCES public.oauth_applications(id) ON DELETE SET NULL;
ALTER TABLE ONLY public.tombstones    ADD CONSTRAINT fk_rails_f95b861449 FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;
ALTER TABLE ONLY public.invites    ADD CONSTRAINT fk_rails_ff69dbb2ac FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;
