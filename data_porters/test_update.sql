UPDATE users SET username = 'zt7@nyu.edu'
UPDATE users SET username = 'zt7@nyu.edu' WHERE id = 1
UPDATE users SET username = 'zt7@nyu.edu' WHERE username = 'zaintq@gmail.com'
UPDATE users SET username = 'zt7@nyu.edu' WHERE username = 'zaintq@gmail.com' AND id = 1
UPDATE posts SET updated_at = $1 WHERE posts.id = $2 
UPDATE posts SET comments_count = comments_count+1 WHERE posts.type IN ('StatusMessage') AND posts.id = $1
UPDATE posts SET updated_at = $1, interacted_at = $2 WHERE posts.id = $3 
UPDATE posts SET likes_count = likes_count+1 WHERE posts.type IN ('StatusMessage') AND posts.id = $1
UPDATE contacts SET sharing = $1, updated_at = $2 WHERE contacts.id = $3
UPDATE contacts SET receiving = $1, updated_at = $2 WHERE contacts.id = $3
UPDATE posts SET reshares_count = reshares_count+1 WHERE posts.type IN ('StatusMessage') AND posts.id = $1 
UPDATE conversations SET updated_at = $1 WHERE conversations.id = $2 
