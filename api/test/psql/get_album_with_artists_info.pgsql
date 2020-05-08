--  Get all profiles with attached albums 

SELECT JSON_AGG(result) AS json_array, albums.name FROM (
  SELECT 
    profile.name,
    profile.userid,
    profile.username
  FROM albums album, profiles profile WHERE 82=album.id AND profile.username=ANY(album.artists)
) result, (
  SELECT JSON_AGG(result) AS artists, songs.name FROM (
    SELECT 
      profile.name,
      profile.userid,
      profile.username
    FROM songs song, profiles profile 
    WHERE 1=song.id AND profile.username=ANY(song.artists)
  )
   result, songs WHERE songs.albumid=82 GROUP BY songs.name
) songs, albums WHERE albums.id=82 GROUP BY albums.name