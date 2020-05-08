--  Get all profiles with attached albums 

SELECT json_build_object(
  'id', p.userid,
  'artist', p.name,
  'albums', a
)
FROM  profiles p INNER JOIN albums a ON p.username=ANY(a.artists) 