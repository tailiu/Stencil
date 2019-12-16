-- Drop table

-- DROP TABLE public.account_deletions

CREATE TABLE public.account_deletions (
	id serial NOT NULL,
	person_id int4 NULL,
	completed_at timestamp NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT account_deletions_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX index_account_deletions_on_person_id ON public.account_deletions USING btree (person_id);

-- Drop table

-- DROP TABLE public.account_migrations

CREATE TABLE public.account_migrations (
	id bigserial NOT NULL,
	old_person_id int4 NOT NULL,
	new_person_id int4 NOT NULL,
	completed_at timestamp NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT account_migrations_pkey PRIMARY KEY (id),
	CONSTRAINT fk_rails_610fe19943 FOREIGN KEY (new_person_id) REFERENCES people(id),
	CONSTRAINT fk_rails_ddbe553eee FOREIGN KEY (old_person_id) REFERENCES people(id)
);
CREATE UNIQUE INDEX index_account_migrations_on_old_person_id ON public.account_migrations USING btree (old_person_id);
CREATE UNIQUE INDEX index_account_migrations_on_old_person_id_and_new_person_id ON public.account_migrations USING btree (old_person_id, new_person_id);

-- Drop table

-- DROP TABLE public.ar_internal_metadata

CREATE TABLE public.ar_internal_metadata (
	"key" varchar NOT NULL,
	value varchar NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT ar_internal_metadata_pkey PRIMARY KEY (key)
);

-- Drop table

-- DROP TABLE public.aspect_memberships

CREATE TABLE public.aspect_memberships (
	id serial NOT NULL,
	aspect_id int4 NOT NULL,
	contact_id int4 NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT aspect_memberships_pkey PRIMARY KEY (id),
	CONSTRAINT aspect_memberships_aspect_id_fk FOREIGN KEY (aspect_id) REFERENCES aspects(id) ON DELETE CASCADE,
	CONSTRAINT aspect_memberships_contact_id_fk FOREIGN KEY (contact_id) REFERENCES contacts(id) ON DELETE CASCADE
);
CREATE INDEX index_aspect_memberships_on_aspect_id ON public.aspect_memberships USING btree (aspect_id);
CREATE UNIQUE INDEX index_aspect_memberships_on_aspect_id_and_contact_id ON public.aspect_memberships USING btree (aspect_id, contact_id);
CREATE INDEX index_aspect_memberships_on_contact_id ON public.aspect_memberships USING btree (contact_id);

-- Drop table

-- DROP TABLE public.aspect_visibilities

CREATE TABLE public.aspect_visibilities (
	id serial NOT NULL,
	shareable_id int4 NOT NULL,
	aspect_id int4 NOT NULL,
	shareable_type varchar NOT NULL DEFAULT 'Post'::character varying,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT aspect_visibilities_pkey PRIMARY KEY (id),
	CONSTRAINT aspect_visibilities_aspect_id_fk FOREIGN KEY (aspect_id) REFERENCES aspects(id) ON DELETE CASCADE
);
CREATE INDEX index_aspect_visibilities_on_aspect_id ON public.aspect_visibilities USING btree (aspect_id);
CREATE UNIQUE INDEX index_aspect_visibilities_on_shareable_and_aspect_id ON public.aspect_visibilities USING btree (shareable_id, shareable_type, aspect_id);
CREATE INDEX index_aspect_visibilities_on_shareable_id_and_shareable_type ON public.aspect_visibilities USING btree (shareable_id, shareable_type);

-- Drop table

-- DROP TABLE public.aspects

CREATE TABLE public.aspects (
	id serial NOT NULL,
	name varchar NOT NULL,
	user_id int4 NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	order_id int4 NULL,
	chat_enabled bool NULL DEFAULT false,
	post_default bool NULL DEFAULT true,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT aspects_pkey PRIMARY KEY (id)
);
CREATE INDEX index_aspects_on_user_id ON public.aspects USING btree (user_id);
CREATE UNIQUE INDEX index_aspects_on_user_id_and_name ON public.aspects USING btree (user_id, name);

-- Drop table

-- DROP TABLE public.authorizations

CREATE TABLE public.authorizations (
	id serial NOT NULL,
	user_id int4 NULL,
	o_auth_application_id int4 NULL,
	refresh_token varchar NULL,
	code varchar NULL,
	redirect_uri varchar NULL,
	nonce varchar NULL,
	scopes varchar NULL,
	code_used bool NULL DEFAULT false,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT authorizations_pkey PRIMARY KEY (id),
	CONSTRAINT fk_rails_4ecef5b8c5 FOREIGN KEY (user_id) REFERENCES users(id),
	CONSTRAINT fk_rails_e166644de5 FOREIGN KEY (o_auth_application_id) REFERENCES o_auth_applications(id)
);
CREATE INDEX index_authorizations_on_o_auth_application_id ON public.authorizations USING btree (o_auth_application_id);
CREATE INDEX index_authorizations_on_user_id ON public.authorizations USING btree (user_id);

-- Drop table

-- DROP TABLE public.blocks

CREATE TABLE public.blocks (
	id serial NOT NULL,
	user_id int4 NULL,
	person_id int4 NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT blocks_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX index_blocks_on_user_id_and_person_id ON public.blocks USING btree (user_id, person_id);

-- Drop table

-- DROP TABLE public.chat_contacts

CREATE TABLE public.chat_contacts (
	id serial NOT NULL,
	user_id int4 NOT NULL,
	jid varchar NOT NULL,
	name varchar(255) NULL,
	ask varchar(128) NULL,
	"subscription" varchar(128) NOT NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT chat_contacts_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX index_chat_contacts_on_user_id_and_jid ON public.chat_contacts USING btree (user_id, jid);

-- Drop table

-- DROP TABLE public.chat_fragments

CREATE TABLE public.chat_fragments (
	id serial NOT NULL,
	user_id int4 NOT NULL,
	root varchar(256) NOT NULL,
	namespace varchar(256) NOT NULL,
	"xml" text NOT NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT chat_fragments_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX index_chat_fragments_on_user_id ON public.chat_fragments USING btree (user_id);

-- Drop table

-- DROP TABLE public.chat_offline_messages

CREATE TABLE public.chat_offline_messages (
	id serial NOT NULL,
	"from" varchar NOT NULL,
	"to" varchar NOT NULL,
	message text NOT NULL,
	created_at timestamp NOT NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT chat_offline_messages_pkey PRIMARY KEY (id)
);

-- Drop table

-- DROP TABLE public.comment_signatures

CREATE TABLE public.comment_signatures (
	comment_id int4 NOT NULL,
	author_signature text NOT NULL,
	signature_order_id int4 NOT NULL,
	additional_data text NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT comment_signatures_comment_id_fk FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
	CONSTRAINT comment_signatures_signature_orders_id_fk FOREIGN KEY (signature_order_id) REFERENCES signature_orders(id)
);
CREATE UNIQUE INDEX index_comment_signatures_on_comment_id ON public.comment_signatures USING btree (comment_id);

-- Drop table

-- DROP TABLE public."comments"

CREATE TABLE public."comments" (
	id serial NOT NULL,
	"text" text NOT NULL,
	commentable_id int4 NOT NULL,
	author_id int4 NOT NULL,
	guid varchar NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	likes_count int4 NOT NULL DEFAULT 0,
	commentable_type varchar(60) NOT NULL DEFAULT 'Post'::character varying,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT comments_pkey PRIMARY KEY (id),
	CONSTRAINT comments_author_id_fk FOREIGN KEY (author_id) REFERENCES people(id) ON DELETE CASCADE
);
CREATE INDEX index_comments_on_commentable_id_and_commentable_type ON public.comments USING btree (commentable_id, commentable_type);
CREATE UNIQUE INDEX index_comments_on_guid ON public.comments USING btree (guid);
CREATE INDEX index_comments_on_person_id ON public.comments USING btree (author_id);

-- Drop table

-- DROP TABLE public.contacts

CREATE TABLE public.contacts (
	id serial NOT NULL,
	user_id int4 NOT NULL,
	person_id int4 NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	sharing bool NOT NULL DEFAULT false,
	receiving bool NOT NULL DEFAULT false,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT contacts_pkey PRIMARY KEY (id),
	CONSTRAINT contacts_person_id_fk FOREIGN KEY (person_id) REFERENCES people(id) ON DELETE CASCADE
);
CREATE INDEX contacts_user_id_idx ON public.contacts USING btree (user_id);
CREATE INDEX index_contacts_on_person_id ON public.contacts USING btree (person_id);
CREATE UNIQUE INDEX index_contacts_on_user_id_and_person_id ON public.contacts USING btree (user_id, person_id);

-- Drop table

-- DROP TABLE public.conversation_visibilities

CREATE TABLE public.conversation_visibilities (
	id serial NOT NULL,
	conversation_id int4 NOT NULL,
	person_id int4 NOT NULL,
	unread int4 NOT NULL DEFAULT 0,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT conversation_visibilities_pkey PRIMARY KEY (id),
	CONSTRAINT conversation_visibilities_conversation_id_fk FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
	CONSTRAINT conversation_visibilities_person_id_fk FOREIGN KEY (person_id) REFERENCES people(id) ON DELETE CASCADE
);
CREATE INDEX index_conversation_visibilities_on_conversation_id ON public.conversation_visibilities USING btree (conversation_id);
CREATE INDEX index_conversation_visibilities_on_person_id ON public.conversation_visibilities USING btree (person_id);
CREATE UNIQUE INDEX index_conversation_visibilities_usefully ON public.conversation_visibilities USING btree (conversation_id, person_id);

-- Drop table

-- DROP TABLE public.conversations

CREATE TABLE public.conversations (
	id serial NOT NULL,
	subject varchar NULL,
	guid varchar NOT NULL,
	author_id int4 NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT conversations_pkey PRIMARY KEY (id),
	CONSTRAINT conversations_author_id_fk FOREIGN KEY (author_id) REFERENCES people(id) ON DELETE CASCADE
);
CREATE INDEX conversations_author_id_fk ON public.conversations USING btree (author_id);
CREATE UNIQUE INDEX index_conversations_on_guid ON public.conversations USING btree (guid);

-- Drop table

-- DROP TABLE public.invitation_codes

CREATE TABLE public.invitation_codes (
	id serial NOT NULL,
	token varchar NULL,
	user_id int4 NULL,
	count int4 NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT invitation_codes_pkey PRIMARY KEY (id)
);

-- Drop table

-- DROP TABLE public.like_signatures

CREATE TABLE public.like_signatures (
	like_id int4 NOT NULL,
	author_signature text NOT NULL,
	signature_order_id int4 NOT NULL,
	additional_data text NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT like_signatures_like_id_fk FOREIGN KEY (like_id) REFERENCES likes(id) ON DELETE CASCADE,
	CONSTRAINT like_signatures_signature_orders_id_fk FOREIGN KEY (signature_order_id) REFERENCES signature_orders(id)
);
CREATE UNIQUE INDEX index_like_signatures_on_like_id ON public.like_signatures USING btree (like_id);

-- Drop table

-- DROP TABLE public.likes

CREATE TABLE public.likes (
	id serial NOT NULL,
	positive bool NULL DEFAULT true,
	target_id int4 NULL,
	author_id int4 NULL,
	guid varchar NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	target_type varchar(60) NOT NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT likes_pkey PRIMARY KEY (id),
	CONSTRAINT likes_author_id_fk FOREIGN KEY (author_id) REFERENCES people(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX index_likes_on_guid ON public.likes USING btree (guid);
CREATE INDEX index_likes_on_post_id ON public.likes USING btree (target_id);
CREATE UNIQUE INDEX index_likes_on_target_id_and_author_id_and_target_type ON public.likes USING btree (target_id, author_id, target_type);
CREATE INDEX likes_author_id_fk ON public.likes USING btree (author_id);

-- Drop table

-- DROP TABLE public.locations

CREATE TABLE public.locations (
	id serial NOT NULL,
	address varchar NULL,
	lat varchar NULL,
	lng varchar NULL,
	status_message_id int4 NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT locations_pkey PRIMARY KEY (id)
);
CREATE INDEX index_locations_on_status_message_id ON public.locations USING btree (status_message_id);

-- Drop table

-- DROP TABLE public.mentions

CREATE TABLE public.mentions (
	id serial NOT NULL,
	mentions_container_id int4 NOT NULL,
	person_id int4 NOT NULL,
	mentions_container_type varchar NOT NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT mentions_pkey PRIMARY KEY (id)
);
CREATE INDEX index_mentions_on_mc_id_and_mc_type ON public.mentions USING btree (mentions_container_id, mentions_container_type);
CREATE UNIQUE INDEX index_mentions_on_person_and_mc_id_and_mc_type ON public.mentions USING btree (person_id, mentions_container_id, mentions_container_type);
CREATE INDEX index_mentions_on_person_id ON public.mentions USING btree (person_id);

-- Drop table

-- DROP TABLE public.messages

CREATE TABLE public.messages (
	id serial NOT NULL,
	conversation_id int4 NOT NULL,
	author_id int4 NOT NULL,
	guid varchar NOT NULL,
	"text" text NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT messages_pkey PRIMARY KEY (id),
	CONSTRAINT messages_author_id_fk FOREIGN KEY (author_id) REFERENCES people(id) ON DELETE CASCADE,
	CONSTRAINT messages_conversation_id_fk FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE
);
CREATE INDEX index_messages_on_author_id ON public.messages USING btree (author_id);
CREATE UNIQUE INDEX index_messages_on_guid ON public.messages USING btree (guid);
CREATE INDEX messages_conversation_id_fk ON public.messages USING btree (conversation_id);

-- Drop table

-- DROP TABLE public.notification_actors

CREATE TABLE public.notification_actors (
	id serial NOT NULL,
	notification_id int4 NULL,
	person_id int4 NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT notification_actors_pkey PRIMARY KEY (id),
	CONSTRAINT notification_actors_notification_id_fk FOREIGN KEY (notification_id) REFERENCES notifications(id) ON DELETE CASCADE
);
CREATE INDEX index_notification_actors_on_notification_id ON public.notification_actors USING btree (notification_id);
CREATE UNIQUE INDEX index_notification_actors_on_notification_id_and_person_id ON public.notification_actors USING btree (notification_id, person_id);
CREATE INDEX index_notification_actors_on_person_id ON public.notification_actors USING btree (person_id);

-- Drop table

-- DROP TABLE public.notifications

CREATE TABLE public.notifications (
	id serial NOT NULL,
	target_type varchar NULL,
	target_id int4 NULL,
	recipient_id int4 NOT NULL,
	unread bool NOT NULL DEFAULT true,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	"type" varchar NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT notifications_pkey PRIMARY KEY (id)
);
CREATE INDEX index_notifications_on_recipient_id ON public.notifications USING btree (recipient_id);
CREATE INDEX index_notifications_on_target_id ON public.notifications USING btree (target_id);
CREATE INDEX index_notifications_on_target_type_and_target_id ON public.notifications USING btree (target_type, target_id);

-- Drop table

-- DROP TABLE public.o_auth_access_tokens

CREATE TABLE public.o_auth_access_tokens (
	id serial NOT NULL,
	authorization_id int4 NULL,
	token varchar NULL,
	expires_at timestamp NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT o_auth_access_tokens_pkey PRIMARY KEY (id),
	CONSTRAINT fk_rails_5debabcff3 FOREIGN KEY (authorization_id) REFERENCES authorizations(id)
);
CREATE INDEX index_o_auth_access_tokens_on_authorization_id ON public.o_auth_access_tokens USING btree (authorization_id);
CREATE UNIQUE INDEX index_o_auth_access_tokens_on_token ON public.o_auth_access_tokens USING btree (token);

-- Drop table

-- DROP TABLE public.o_auth_applications

CREATE TABLE public.o_auth_applications (
	id serial NOT NULL,
	user_id int4 NULL,
	client_id varchar NULL,
	client_secret varchar NULL,
	client_name varchar NULL,
	redirect_uris text NULL,
	response_types varchar NULL,
	grant_types varchar NULL,
	application_type varchar NULL DEFAULT 'web'::character varying,
	contacts varchar NULL,
	logo_uri varchar NULL,
	client_uri varchar NULL,
	policy_uri varchar NULL,
	tos_uri varchar NULL,
	sector_identifier_uri varchar NULL,
	token_endpoint_auth_method varchar NULL,
	jwks text NULL,
	jwks_uri varchar NULL,
	ppid bool NULL DEFAULT false,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT o_auth_applications_pkey PRIMARY KEY (id),
	CONSTRAINT fk_rails_ad75323da2 FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE UNIQUE INDEX index_o_auth_applications_on_client_id ON public.o_auth_applications USING btree (client_id);
CREATE INDEX index_o_auth_applications_on_user_id ON public.o_auth_applications USING btree (user_id);

-- Drop table

-- DROP TABLE public.o_embed_caches

CREATE TABLE public.o_embed_caches (
	id serial NOT NULL,
	url varchar(1024) NOT NULL,
	"data" text NOT NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT o_embed_caches_pkey PRIMARY KEY (id)
);
CREATE INDEX index_o_embed_caches_on_url ON public.o_embed_caches USING btree (url);

-- Drop table

-- DROP TABLE public.open_graph_caches

CREATE TABLE public.open_graph_caches (
	id serial NOT NULL,
	title varchar NULL,
	ob_type varchar NULL,
	image text NULL,
	url text NULL,
	description text NULL,
	video_url text NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT open_graph_caches_pkey PRIMARY KEY (id)
);

-- Drop table

-- DROP TABLE public.participations

CREATE TABLE public.participations (
	id serial NOT NULL,
	guid varchar NULL,
	target_id int4 NULL,
	target_type varchar(60) NOT NULL,
	author_id int4 NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	count int4 NOT NULL DEFAULT 1,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT participations_pkey PRIMARY KEY (id)
);
CREATE INDEX index_participations_on_author_id ON public.participations USING btree (author_id);
CREATE INDEX index_participations_on_guid ON public.participations USING btree (guid);
CREATE UNIQUE INDEX index_participations_on_target_id_and_target_type_and_author_id ON public.participations USING btree (target_id, target_type, author_id);

-- Drop table

-- DROP TABLE public.people

CREATE TABLE public.people (
	id serial NOT NULL,
	guid varchar NOT NULL,
	diaspora_handle varchar NOT NULL,
	serialized_public_key text NOT NULL,
	owner_id int4 NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	closed_account bool NULL DEFAULT false,
	fetch_status int4 NULL DEFAULT 0,
	pod_id int4 NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT people_pkey PRIMARY KEY (id),
	CONSTRAINT people_pod_id_fk FOREIGN KEY (pod_id) REFERENCES pods(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX index_people_on_diaspora_handle ON public.people USING btree (diaspora_handle);
CREATE UNIQUE INDEX index_people_on_guid ON public.people USING btree (guid);
CREATE UNIQUE INDEX index_people_on_owner_id ON public.people USING btree (owner_id);

-- Drop table

-- DROP TABLE public.photos

CREATE TABLE public.photos (
	id serial NOT NULL,
	author_id int4 NOT NULL,
	public bool NOT NULL DEFAULT false,
	guid varchar NOT NULL,
	pending bool NOT NULL DEFAULT false,
	"text" text NULL,
	remote_photo_path text NULL,
	remote_photo_name varchar NULL,
	random_string varchar NULL,
	processed_image varchar NULL,
	created_at timestamp NULL,
	updated_at timestamp NULL,
	unprocessed_image varchar NULL,
	status_message_guid varchar NULL,
	comments_count int4 NULL,
	height int4 NULL,
	width int4 NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT photos_pkey PRIMARY KEY (id)
);
CREATE INDEX index_photos_on_author_id ON public.photos USING btree (author_id);
CREATE UNIQUE INDEX index_photos_on_guid ON public.photos USING btree (guid);
CREATE INDEX index_photos_on_status_message_guid ON public.photos USING btree (status_message_guid);

-- Drop table

-- DROP TABLE public.pods

CREATE TABLE public.pods (
	id serial NOT NULL,
	host varchar NOT NULL,
	ssl bool NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	status int4 NULL DEFAULT 0,
	checked_at timestamp NULL DEFAULT '1970-01-01 00:00:00'::timestamp without time zone,
	offline_since timestamp NULL,
	response_time int4 NULL DEFAULT '-1'::integer,
	software varchar NULL,
	error varchar NULL,
	port int4 NULL,
	blocked bool NULL DEFAULT false,
	scheduled_check bool NOT NULL DEFAULT false,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT pods_pkey PRIMARY KEY (id)
);
CREATE INDEX index_pods_on_checked_at ON public.pods USING btree (checked_at);
CREATE UNIQUE INDEX index_pods_on_host_and_port ON public.pods USING btree (host, port);
CREATE INDEX index_pods_on_offline_since ON public.pods USING btree (offline_since);
CREATE INDEX index_pods_on_status ON public.pods USING btree (status);

-- Drop table

-- DROP TABLE public.poll_answers

CREATE TABLE public.poll_answers (
	id serial NOT NULL,
	answer varchar NOT NULL,
	poll_id int4 NOT NULL,
	guid varchar NULL,
	vote_count int4 NULL DEFAULT 0,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT poll_answers_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX index_poll_answers_on_guid ON public.poll_answers USING btree (guid);
CREATE INDEX index_poll_answers_on_poll_id ON public.poll_answers USING btree (poll_id);

-- Drop table

-- DROP TABLE public.poll_participation_signatures

CREATE TABLE public.poll_participation_signatures (
	poll_participation_id int4 NOT NULL,
	author_signature text NOT NULL,
	signature_order_id int4 NOT NULL,
	additional_data text NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT poll_participation_signatures_poll_participation_id_fk FOREIGN KEY (poll_participation_id) REFERENCES poll_participations(id) ON DELETE CASCADE,
	CONSTRAINT poll_participation_signatures_signature_orders_id_fk FOREIGN KEY (signature_order_id) REFERENCES signature_orders(id)
);
CREATE UNIQUE INDEX index_poll_participation_signatures_on_poll_participation_id ON public.poll_participation_signatures USING btree (poll_participation_id);

-- Drop table

-- DROP TABLE public.poll_participations

CREATE TABLE public.poll_participations (
	id serial NOT NULL,
	poll_answer_id int4 NOT NULL,
	author_id int4 NOT NULL,
	poll_id int4 NOT NULL,
	guid varchar NULL,
	created_at timestamp NULL,
	updated_at timestamp NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT poll_participations_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX index_poll_participations_on_guid ON public.poll_participations USING btree (guid);
CREATE UNIQUE INDEX index_poll_participations_on_poll_id_and_author_id ON public.poll_participations USING btree (poll_id, author_id);

-- Drop table

-- DROP TABLE public.polls

CREATE TABLE public.polls (
	id serial NOT NULL,
	question varchar NOT NULL,
	status_message_id int4 NOT NULL,
	status bool NULL,
	guid varchar NULL,
	created_at timestamp NULL,
	updated_at timestamp NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT polls_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX index_polls_on_guid ON public.polls USING btree (guid);
CREATE INDEX index_polls_on_status_message_id ON public.polls USING btree (status_message_id);

-- Drop table

-- DROP TABLE public.posts

CREATE TABLE public.posts (
	id serial NOT NULL,
	author_id int4 NOT NULL,
	public bool NOT NULL DEFAULT false,
	guid varchar NOT NULL,
	"type" varchar(40) NOT NULL,
	"text" text NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	provider_display_name varchar NULL,
	root_guid varchar NULL,
	likes_count int4 NULL DEFAULT 0,
	comments_count int4 NULL DEFAULT 0,
	o_embed_cache_id int4 NULL,
	reshares_count int4 NULL DEFAULT 0,
	interacted_at timestamp NULL,
	tweet_id varchar NULL,
	open_graph_cache_id int4 NULL,
	tumblr_ids text NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT posts_pkey PRIMARY KEY (id),
	CONSTRAINT posts_author_id_fk FOREIGN KEY (author_id) REFERENCES people(id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX index_posts_on_author_id_and_root_guid ON public.posts USING btree (author_id, root_guid);
CREATE INDEX index_posts_on_created_at_and_id ON public.posts USING btree (created_at, id);
CREATE UNIQUE INDEX index_posts_on_guid ON public.posts USING btree (guid);
CREATE INDEX index_posts_on_id_and_type ON public.posts USING btree (id, type);
CREATE INDEX index_posts_on_person_id ON public.posts USING btree (author_id);
CREATE INDEX index_posts_on_root_guid ON public.posts USING btree (root_guid);

-- Drop table

-- DROP TABLE public.ppid

CREATE TABLE public.ppid (
	id serial NOT NULL,
	o_auth_application_id int4 NULL,
	user_id int4 NULL,
	guid varchar(32) NULL,
	identifier varchar NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT ppid_pkey PRIMARY KEY (id),
	CONSTRAINT fk_rails_150457f962 FOREIGN KEY (o_auth_application_id) REFERENCES o_auth_applications(id),
	CONSTRAINT fk_rails_e6b8e5264f FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE INDEX index_ppid_on_o_auth_application_id ON public.ppid USING btree (o_auth_application_id);
CREATE INDEX index_ppid_on_user_id ON public.ppid USING btree (user_id);

-- Drop table

-- DROP TABLE public.profiles

CREATE TABLE public.profiles (
	id serial NOT NULL,
	diaspora_handle varchar NULL,
	first_name varchar(127) NULL,
	last_name varchar(127) NULL,
	image_url varchar NULL,
	image_url_small varchar NULL,
	image_url_medium varchar NULL,
	birthday date NULL,
	gender varchar NULL,
	bio text NULL,
	searchable bool NOT NULL DEFAULT true,
	person_id int4 NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	"location" varchar NULL,
	full_name varchar(70) NULL,
	nsfw bool NULL DEFAULT false,
	public_details bool NULL DEFAULT false,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT profiles_pkey PRIMARY KEY (id),
	CONSTRAINT profiles_person_id_fk FOREIGN KEY (person_id) REFERENCES people(id) ON DELETE CASCADE
);
CREATE INDEX index_profiles_on_full_name ON public.profiles USING btree (full_name);
CREATE INDEX index_profiles_on_full_name_and_searchable ON public.profiles USING btree (full_name, searchable);
CREATE INDEX index_profiles_on_person_id ON public.profiles USING btree (person_id);

-- Drop table

-- DROP TABLE public."references"

CREATE TABLE public."references" (
	id bigserial NOT NULL,
	source_id int4 NOT NULL,
	source_type varchar(60) NOT NULL,
	target_id int4 NOT NULL,
	target_type varchar(60) NOT NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT references_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX index_references_on_source_and_target ON public."references" USING btree (source_id, source_type, target_id, target_type);
CREATE INDEX index_references_on_source_id_and_source_type ON public."references" USING btree (source_id, source_type);

-- Drop table

-- DROP TABLE public.reports

CREATE TABLE public.reports (
	id serial NOT NULL,
	item_id int4 NOT NULL,
	item_type varchar NOT NULL,
	reviewed bool NULL DEFAULT false,
	"text" text NULL,
	created_at timestamp NULL,
	updated_at timestamp NULL,
	user_id int4 NOT NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT reports_pkey PRIMARY KEY (id)
);
CREATE INDEX index_reports_on_item_id ON public.reports USING btree (item_id);

-- Drop table

-- DROP TABLE public.roles

CREATE TABLE public.roles (
	id serial NOT NULL,
	person_id int4 NULL,
	name varchar NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT roles_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX index_roles_on_person_id_and_name ON public.roles USING btree (person_id, name);

-- Drop table

-- DROP TABLE public.schema_migrations

CREATE TABLE public.schema_migrations (
	"version" varchar NOT NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT schema_migrations_pkey PRIMARY KEY (version)
);

-- Drop table

-- DROP TABLE public.services

CREATE TABLE public.services (
	id serial NOT NULL,
	"type" varchar(127) NOT NULL,
	user_id int4 NOT NULL,
	uid varchar(127) NULL,
	access_token varchar NULL,
	access_secret varchar NULL,
	nickname varchar NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT services_pkey PRIMARY KEY (id),
	CONSTRAINT services_user_id_fk FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX index_services_on_type_and_uid ON public.services USING btree (type, uid);
CREATE INDEX index_services_on_user_id ON public.services USING btree (user_id);

-- Drop table

-- DROP TABLE public.share_visibilities

CREATE TABLE public.share_visibilities (
	id serial NOT NULL,
	shareable_id int4 NOT NULL,
	hidden bool NOT NULL DEFAULT false,
	shareable_type varchar(60) NOT NULL DEFAULT 'Post'::character varying,
	user_id int4 NOT NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT share_visibilities_pkey PRIMARY KEY (id),
	CONSTRAINT share_visibilities_user_id_fk FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX index_post_visibilities_on_post_id ON public.share_visibilities USING btree (shareable_id);
CREATE INDEX index_share_visibilities_on_user_id ON public.share_visibilities USING btree (user_id);
CREATE INDEX shareable_and_hidden_and_user_id ON public.share_visibilities USING btree (shareable_id, shareable_type, hidden, user_id);
CREATE UNIQUE INDEX shareable_and_user_id ON public.share_visibilities USING btree (shareable_id, shareable_type, user_id);

-- Drop table

-- DROP TABLE public.signature_orders

CREATE TABLE public.signature_orders (
	id serial NOT NULL,
	"order" varchar NOT NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT signature_orders_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX index_signature_orders_on_order ON public.signature_orders USING btree ("order");

-- Drop table

-- DROP TABLE public.simple_captcha_data

CREATE TABLE public.simple_captcha_data (
	id serial NOT NULL,
	"key" varchar(40) NULL,
	value varchar(12) NULL,
	created_at timestamp NULL,
	updated_at timestamp NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT simple_captcha_data_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_key ON public.simple_captcha_data USING btree (key);

-- Drop table

-- DROP TABLE public.tag_followings

CREATE TABLE public.tag_followings (
	id serial NOT NULL,
	tag_id int4 NOT NULL,
	user_id int4 NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT tag_followings_pkey PRIMARY KEY (id)
);
CREATE INDEX index_tag_followings_on_tag_id ON public.tag_followings USING btree (tag_id);
CREATE UNIQUE INDEX index_tag_followings_on_tag_id_and_user_id ON public.tag_followings USING btree (tag_id, user_id);
CREATE INDEX index_tag_followings_on_user_id ON public.tag_followings USING btree (user_id);

-- Drop table

-- DROP TABLE public.taggings

CREATE TABLE public.taggings (
	id serial NOT NULL,
	tag_id int4 NULL,
	taggable_id int4 NULL,
	taggable_type varchar(127) NULL,
	tagger_id int4 NULL,
	tagger_type varchar(127) NULL,
	context varchar(127) NULL,
	created_at timestamp NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT taggings_pkey PRIMARY KEY (id)
);
CREATE INDEX index_taggings_on_created_at ON public.taggings USING btree (created_at);
CREATE INDEX index_taggings_on_tag_id ON public.taggings USING btree (tag_id);
CREATE INDEX index_taggings_on_taggable_id_and_taggable_type_and_context ON public.taggings USING btree (taggable_id, taggable_type, context);
CREATE UNIQUE INDEX index_taggings_uniquely ON public.taggings USING btree (taggable_id, taggable_type, tag_id);

-- Drop table

-- DROP TABLE public.tags

CREATE TABLE public.tags (
	id serial NOT NULL,
	name varchar NULL,
	taggings_count int4 NULL DEFAULT 0,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT tags_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX index_tags_on_name ON public.tags USING btree (name);

-- Drop table

-- DROP TABLE public.user_preferences

CREATE TABLE public.user_preferences (
	id serial NOT NULL,
	email_type varchar NULL,
	user_id int4 NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT user_preferences_pkey PRIMARY KEY (id)
);
CREATE INDEX index_user_preferences_on_user_id_and_email_type ON public.user_preferences USING btree (user_id, email_type);

-- Drop table

-- DROP TABLE public.users

CREATE TABLE public.users (
	id serial NOT NULL,
	username varchar NOT NULL,
	serialized_private_key text NULL,
	getting_started bool NOT NULL DEFAULT true,
	disable_mail bool NOT NULL DEFAULT false,
	"language" varchar NULL,
	email varchar NOT NULL DEFAULT ''::character varying,
	encrypted_password varchar NOT NULL DEFAULT ''::character varying,
	reset_password_token varchar NULL,
	remember_created_at timestamp NULL,
	sign_in_count int4 NULL DEFAULT 0,
	current_sign_in_at timestamp NULL,
	last_sign_in_at timestamp NULL,
	current_sign_in_ip varchar NULL,
	last_sign_in_ip varchar NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	invited_by_id int4 NULL,
	authentication_token varchar(30) NULL,
	unconfirmed_email varchar NULL,
	confirm_email_token varchar(30) NULL,
	locked_at timestamp NULL,
	show_community_spotlight_in_stream bool NOT NULL DEFAULT true,
	auto_follow_back bool NULL DEFAULT false,
	auto_follow_back_aspect_id int4 NULL,
	hidden_shareables text NULL,
	reset_password_sent_at timestamp NULL,
	last_seen timestamp NULL,
	remove_after timestamp NULL,
	export varchar NULL,
	exported_at timestamp NULL,
	exporting bool NULL DEFAULT false,
	strip_exif bool NULL DEFAULT true,
	exported_photos_file varchar NULL,
	exported_photos_at timestamp NULL,
	exporting_photos bool NULL DEFAULT false,
	color_theme varchar NULL,
	post_default_public bool NULL DEFAULT false,
	mark_as_delete bool NULL DEFAULT false,
	CONSTRAINT users_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX index_users_on_authentication_token ON public.users USING btree (authentication_token);
CREATE UNIQUE INDEX index_users_on_email ON public.users USING btree (email);
CREATE UNIQUE INDEX index_users_on_username ON public.users USING btree (username);
