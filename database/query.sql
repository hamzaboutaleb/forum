SELECT postId FROM post_reactions WHERE userId = 1;
-- SELECT * FROM tags;
-- SELECT * FROM post_tags;

-- SELECT postId from post_tags WHERE tagId = (SELECT id FROM tags WHERE name="rzarza");

SELECT p.id, p.title, u.username, COALESCE(SUM(pr.isLike), 0) AS likeCount FROM posts p
LEFT JOIN users u ON p.userId = u.id
LEFT JOIN post_reactions pr ON p.id = pr.postId
GROUP BY p.id 
HAVING p.id IN (SELECT postId from post_tags WHERE tagId = (SELECT id FROM tags WHERE name="rzarza"))
--AND p.id IN (SELECT id FROM posts WHERE userId = 1)
AND p.id IN (SELECT postId FROM post_reactions WHERE userId = 1); --LIKED POSTS

