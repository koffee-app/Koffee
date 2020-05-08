-- INSERT INTO songs (albumid, name, duration, artists) VALUES (82, 'Great song3...', 13300, '{gabivlj}') 

-- SELECT JSON_AGG(result) AS artists, songs.name FROM (
--   SELECT 
--     profile.name,
--     profile.userid,
--     profile.username
--   FROM songs song, profiles profile WHERE 1=song.id AND profile.username=ANY(song.artists)
-- ) result, songs WHERE songs.albumid=82 GROUP BY songs.name

-- SELECT JSON_AGG(result) AS artists, json_agg(songs) AS songs,  albums.name FROM (
--   SELECT DISTINCT
--     profile.name,
--     profile.userid,
--     profile.username
--   FROM albums album, profiles profile WHERE 82=album.id AND profile.username=ANY(album.artists) GROUP BY profile.userid, profile.name, profile.username
-- ) result, (
--   SELECT JSON_AGG(result) AS artists, songs.name FROM (
--     SELECT DISTINCT
--       profile.name,
--       profile.userid,
--       profile.username
--     FROM songs song, profiles profile 
--     WHERE 1=song.id AND profile.username=ANY(song.artists)
--   )
--    result, songs WHERE songs.albumid=82 GROUP BY songs.name
-- ) songs, albums WHERE albums.id=82 GROUP BY albums.name;

select row_to_json(t)
from (
  select name,
  -- Fill songs
    (
      select array_to_json(array_agg(row_to_json(d)))
      from (
        select songs.name, songs.duration, songs.artists, 
        -- Fill artists
        (
          SELECT array_to_json(array_agg(row_to_json(profiles))) FROM profiles
                     WHERE profiles.username=ANY(songs.artists)
        ) as artists
        from songs
        where songs.albumid=albums.id
      ) d
    ) as songs,
    --  FIll profiles
    (  
      SELECT array_to_json(array_agg(row_to_json(profiles))) FROM profiles
      WHERE profiles.username=ANY(albums.artists) 
    ) as artists,
    (
      SELECT row_to_json(images) from images WHERE images.id=albums.id AND images.type='COVER_URL'
    ) as images
  from albums
) t