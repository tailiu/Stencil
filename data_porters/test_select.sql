SELECT * FROM users
SELECT id, username FROM users
SELECT id, username, created_at, updated_at FROM users
SELECT * FROM people
SELECT id, diaspora_handle, guid FROM people
SELECT * FROM users JOIN people ON users.id = people.owner_id
SELECT id, guid, author_id, text FROM POSTS WHERE author_id = $1 AND id NOT IN (SELECT distinct(target_id) FROM likes UNION SELECT distinct(commentable_id) FROM comments) order by random()
SELECT id FROM contacts WHERE user_id = $1 AND person_id = $2
SELECT id FROM aspect_memberships WHERE aspect_id = $1 AND contact_id = $2 LIMIT 1
SELECT user_id, person_id, string_agg(aspect_id::text, ',') as aspects FROM ( SELECT users.id as user_id, people.id as person_id, aspects.id as aspect_id FROM users JOIN people ON users.id = people.owner_id JOIN aspects ON aspects.user_id = users.id ) tab GROUP BY user_id, person_id ORDER BY random()
SELECT user_id, person_id, string_agg(aspect_id::text, ',') as aspects FROM ( SELECT users.id as user_id, people.id as person_id, aspects.id as aspect_id FROM users JOIN people ON users.id = people.owner_id JOIN aspects ON aspects.user_id = users.id ) tab WHERE user_id NOT IN ( SELECT DISTINCT %s FROM %s ) GROUP BY user_id, person_id ORDER BY random()
SELECT users.id as user_id, people.id as person_id, string_agg(aspects.id::text, ',') as aspects, contacts.id as contact_id, am.aspect_id as contact_aspect FROM contacts JOIN aspect_memberships am on contacts.id = am.contact_id JOIN people on people.id = contacts.user_id JOIN users on users.id = people.owner_id JOIN aspects on aspects.user_id = users.id WHERE contacts.user_id = $1 AND contacts.sharing = true GROUP BY users.id, people.id, contacts.id, am.aspect_id
SELECT users.id as user_id, people.id as person_id FROM users JOIN people ON users.id = people.owner_id WHERE users.id NOT IN ($1) ORDER BY random() LIMIT 1
SELECT id FROM aspects WHERE user_id = $1
