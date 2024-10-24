SELECT c.id ,c.postId ,c.userId, c.comment ,c.createdAt ,u.username , 
(SELECT count(*) from comment_reactions WHERE isLike=1 AND commentId=c.id ) likes,
(SELECT count(*) from comment_reactions WHERE isLike=-1 AND commentId=c.id ) dislike
 FROM comments c 
LEFT JOIN comment_reactions l ON c.id = l.commentId 
LEFT JOIN users u ON c.userId = u.id WHERE c.postId = 6 GROUP BY c.id HAVING count(c.id) > 0 ORDER BY c.createdAt desc