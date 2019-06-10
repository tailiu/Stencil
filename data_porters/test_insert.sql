INSERT INTO users (id,username,email,encrypted_password) VALUES (1,'zain','zaintq@gmail.com','091289yhewidbx')
INSERT INTO users (username, serialized_private_key, language, email, encrypted_password, created_at, updated_at, color_theme, last_seen, sign_in_count, current_sign_in_ip, last_sign_in_ip, current_sign_in_at, last_sign_in_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id
INSERT INTO people (guid,diaspora_handle,serialized_public_key,owner_id,created_at,updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
INSERT INTO profiles (person_id, created_at, updated_at, full_name) VALUES ($1, $2, $3, $4)
INSERT INTO aspects (name, user_id, created_at, updated_at, order_id) VALUES ($1, $2, $3, $4, $5)  RETURNING id
INSERT INTO posts (author_id, guid, type, text, created_at, updated_at, interacted_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
INSERT INTO aspect_visibilities (shareable_id, aspect_id) VALUES ($1, $2)
INSERT INTO share_visibilities (shareable_id, user_id) VALUES ($1, $2)
INSERT INTO comments (text,commentable_id,author_id,guid,created_at,updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
INSERT INTO notifications (target_type,target_id,recipient_id,created_at,updated_at,type) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
INSERT INTO notification_actors (notification_id,person_id,created_at,updated_at) VALUES ($1, $2, $3, $4)
INSERT INTO likes (target_id, author_id, guid, created_at, updated_at, target_type) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
INSERT INTO notifications (target_type,target_id,recipient_id,created_at, updated_at, type) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
INSERT INTO notification_actors (notification_id,person_id,created_at,updated_at) VALUES ($1, $2, $3, $4)
INSERT INTO participations (guid,target_id,target_type,author_id,created_at,updated_at) VALUES ($1, $2, $3, $4, $5, $6)
INSERT INTO aspect_memberships (aspect_id,contact_id,created_at,updated_at) VALUES ($1, $2, $3, $4)
INSERT INTO contacts (user_id,person_id,created_at,updated_at,receiving) VALUES ($1, $2, $3, $4, $5) RETURNING id
INSERT INTO contacts (user_id,person_id,created_at,updated_at,sharing) VALUES ($1, $2, $3, $4, $5)
INSERT INTO aspect_memberships (aspect_id,contact_id,created_at,updated_at) VALUES ($1, $2, $3, $4)
INSERT INTO notifications (target_type,target_id,recipient_id,created_at,updated_at,type) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
INSERT INTO notification_actors (notification_id,person_id,created_at,updated_at) VALUES ($1, $2, $3, $4) RETURNING id
INSERT INTO posts (author_id, public, guid, type, text, created_at, updated_at, root_guid, interacted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id
INSERT INTO participations (guid, target_id, target_type, author_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)
INSERT INTO participations (guid, target_id, target_type, author_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)
INSERT INTO notifications (target_type, target_id, recipient_id, created_at, updated_at, type) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
INSERT INTO notification_actors (notification_id, person_id, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id
INSERT INTO conversations (subject,guid,author_id,created_at,updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id
INSERT INTO conversation_visibilities (conversation_id,person_id,created_at,updated_at) VALUES ($1, $2, $3, $4), ($5, $6, $7, $8)
INSERT INTO messages (conversation_id,author_id,guid,text,created_at,updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
