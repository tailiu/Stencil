CREATE TABLE public.account_deletions (
    id SERIAL PRIMARY KEY,
    person_id integer,
    completed_at timestamp without time zone
);


CREATE TABLE public.account_migrations (
    id bigint PRIMARY KEY,
    old_person_id integer NOT NULL,
    new_person_id integer NOT NULL,
    completed_at timestamp without time zone
);

CREATE TABLE public.ar_internal_metadata (
    key character varying NOT NULL,
    value character varying,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

CREATE TABLE public.aspect_memberships (
    id SERIAL PRIMARY KEY,
    aspect_id integer NOT NULL,
    contact_id integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

CREATE TABLE public.aspect_visibilities (
    id SERIAL PRIMARY KEY,
    shareable_id integer NOT NULL,
    aspect_id integer NOT NULL,
    shareable_type character varying DEFAULT 'Post'::character varying NOT NULL
);

CREATE TABLE public.aspects (
    id SERIAL PRIMARY KEY,
    name character varying NOT NULL,
    user_id integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    order_id integer,
    chat_enabled boolean DEFAULT false,
    post_default boolean DEFAULT true
);


CREATE TABLE public.authorizations (
    id SERIAL PRIMARY KEY,
    user_id integer,
    o_auth_application_id integer,
    refresh_token character varying,
    code character varying,
    redirect_uri character varying,
    nonce character varying,
    scopes character varying,
    code_used boolean DEFAULT false,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

CREATE TABLE public.blocks (
    id SERIAL PRIMARY KEY,
    user_id integer,
    person_id integer
);

CREATE TABLE public.chat_contacts (
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL,
    jid character varying NOT NULL,
    name character varying(255),
    ask character varying(128),
    subscription character varying(128) NOT NULL
);

CREATE TABLE public.chat_fragments (
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL,
    root character varying(256) NOT NULL,
    namespace character varying(256) NOT NULL,
    xml text NOT NULL
);

CREATE TABLE public.chat_offline_messages (
    id SERIAL PRIMARY KEY,
    "from" character varying NOT NULL,
    "to" character varying NOT NULL,
    message text NOT NULL,
    created_at timestamp without time zone NOT NULL
);

CREATE TABLE public.comment_signatures (
    comment_id integer NOT NULL,
    author_signature text NOT NULL,
    signature_order_id integer NOT NULL,
    additional_data text
);

CREATE TABLE public.comments (
    id SERIAL PRIMARY KEY,
    text text NOT NULL,
    commentable_id integer NOT NULL,
    author_id integer NOT NULL,
    guid character varying NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    likes_count integer DEFAULT 0 NOT NULL,
    commentable_type character varying(60) DEFAULT 'Post'::character varying NOT NULL
);

CREATE TABLE public.contacts (
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL,
    person_id integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    sharing boolean DEFAULT false NOT NULL,
    receiving boolean DEFAULT false NOT NULL
);


CREATE TABLE public.conversation_visibilities (
    id SERIAL PRIMARY KEY,
    conversation_id integer NOT NULL,
    person_id integer NOT NULL,
    unread integer DEFAULT 0 NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

CREATE TABLE public.conversations (
    id SERIAL PRIMARY KEY,
    subject character varying,
    guid character varying NOT NULL,
    author_id integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

CREATE TABLE public.invitation_codes (
    id SERIAL PRIMARY KEY,
    token character varying,
    user_id integer,
    count integer,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);



CREATE TABLE public.like_signatures (
    like_id integer NOT NULL,
    author_signature text NOT NULL,
    signature_order_id integer NOT NULL,
    additional_data text
);


CREATE TABLE public.likes (
    id SERIAL PRIMARY KEY,
    positive boolean DEFAULT true,
    target_id integer,
    author_id integer,
    guid character varying,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    target_type character varying(60) NOT NULL
);

CREATE TABLE public.locations (
    id SERIAL PRIMARY KEY,
    address character varying,
    lat character varying,
    lng character varying,
    status_message_id integer,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

CREATE TABLE public.mentions (
    id SERIAL PRIMARY KEY,
    mentions_container_id integer NOT NULL,
    person_id integer NOT NULL,
    mentions_container_type character varying NOT NULL
);


CREATE TABLE public.messages (
    id SERIAL PRIMARY KEY,
    conversation_id integer NOT NULL,
    author_id integer NOT NULL,
    guid character varying NOT NULL,
    text text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

CREATE TABLE public.notification_actors (
    id SERIAL PRIMARY KEY,
    notification_id integer,
    person_id integer,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

CREATE TABLE public.notifications (
    id SERIAL PRIMARY KEY,
    target_type character varying,
    target_id integer,
    recipient_id integer NOT NULL,
    unread boolean DEFAULT true NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    type character varying
);


CREATE TABLE public.o_auth_access_tokens (
    id SERIAL PRIMARY KEY,
    authorization_id integer,
    token character varying,
    expires_at timestamp without time zone,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


CREATE TABLE public.o_auth_applications (
    id SERIAL PRIMARY KEY,
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
    updated_at timestamp without time zone NOT NULL
);

CREATE TABLE public.o_embed_caches (
    id SERIAL PRIMARY KEY,
    url character varying(1024) NOT NULL,
    data text NOT NULL
);

CREATE TABLE public.open_graph_caches (
    id SERIAL PRIMARY KEY,
    title character varying,
    ob_type character varying,
    image text,
    url text,
    description text,
    video_url text
);

CREATE TABLE public.participations (
    id SERIAL PRIMARY KEY,
    guid character varying,
    target_id integer,
    target_type character varying(60) NOT NULL,
    author_id integer,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    count integer DEFAULT 1 NOT NULL
);

CREATE TABLE public.people (
    id SERIAL PRIMARY KEY,
    guid character varying NOT NULL,
    diaspora_handle character varying NOT NULL,
    serialized_public_key text NOT NULL,
    owner_id integer,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    closed_account boolean DEFAULT false,
    fetch_status integer DEFAULT 0,
    pod_id integer
);


CREATE TABLE public.photos (
    id SERIAL PRIMARY KEY,
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
    width integer
);

CREATE TABLE public.pods (
    id SERIAL PRIMARY KEY,
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
    scheduled_check boolean DEFAULT false NOT NULL
);

CREATE TABLE public.poll_answers (
    id SERIAL PRIMARY KEY,
    answer character varying NOT NULL,
    poll_id integer NOT NULL,
    guid character varying,
    vote_count integer DEFAULT 0
);

CREATE TABLE public.poll_participation_signatures (
    poll_participation_id integer NOT NULL,
    author_signature text NOT NULL,
    signature_order_id integer NOT NULL,
    additional_data text
);

CREATE TABLE public.poll_participations (
    id SERIAL PRIMARY KEY,
    poll_answer_id integer NOT NULL,
    author_id integer NOT NULL,
    poll_id integer NOT NULL,
    guid character varying,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);


CREATE TABLE public.polls (
    id SERIAL PRIMARY KEY,
    question character varying NOT NULL,
    status_message_id integer NOT NULL,
    status boolean,
    guid character varying,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);


CREATE TABLE public.posts (
    id SERIAL PRIMARY KEY,
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
    tumblr_ids text
);

CREATE TABLE public.ppid (
    id SERIAL PRIMARY KEY,
    o_auth_application_id integer,
    user_id integer,
    guid character varying(32),
    identifier character varying
);

CREATE TABLE public.profiles (
    id SERIAL PRIMARY KEY,
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
    public_details boolean DEFAULT false
);


CREATE TABLE public."references" (
    id bigint PRIMARY KEY,
    source_id integer NOT NULL,
    source_type character varying(60) NOT NULL,
    target_id integer NOT NULL,
    target_type character varying(60) NOT NULL
);

CREATE TABLE public.reports (
    id SERIAL PRIMARY KEY,
    item_id integer NOT NULL,
    item_type character varying NOT NULL,
    reviewed boolean DEFAULT false,
    text text,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    user_id integer NOT NULL
);

CREATE TABLE public.roles (
    id SERIAL PRIMARY KEY,
    person_id integer,
    name character varying,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

CREATE TABLE public.schema_migrations (
    version character varying NOT NULL
);

CREATE TABLE public.services (
    id SERIAL PRIMARY KEY,
    type character varying(127) NOT NULL,
    user_id integer NOT NULL,
    uid character varying(127),
    access_token character varying,
    access_secret character varying,
    nickname character varying,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


CREATE TABLE public.share_visibilities (
    id SERIAL PRIMARY KEY,
    shareable_id integer NOT NULL,
    hidden boolean DEFAULT false NOT NULL,
    shareable_type character varying(60) DEFAULT 'Post'::character varying NOT NULL,
    user_id integer NOT NULL
);

CREATE TABLE public.signature_orders (
    id SERIAL PRIMARY KEY,
    "order" character varying NOT NULL
);

CREATE TABLE public.simple_captcha_data (
    id SERIAL PRIMARY KEY,
    key character varying(40),
    value character varying(12),
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

CREATE TABLE public.tag_followings (
    id SERIAL PRIMARY KEY,
    tag_id integer NOT NULL,
    user_id integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

CREATE TABLE public.taggings (
    id SERIAL PRIMARY KEY,
    tag_id integer,
    taggable_id integer,
    taggable_type character varying(127),
    tagger_id integer,
    tagger_type character varying(127),
    context character varying(127),
    created_at timestamp without time zone
);

CREATE TABLE public.tags (
    id SERIAL PRIMARY KEY,
    name character varying,
    taggings_count integer DEFAULT 0
);


CREATE TABLE public.user_preferences (
    id SERIAL PRIMARY KEY,
    email_type character varying,
    user_id integer,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


CREATE TABLE public.users (
    id SERIAL PRIMARY KEY,
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
    post_default_public boolean DEFAULT false
);

