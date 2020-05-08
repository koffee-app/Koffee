-- json cant have multiple fields next to it, thats why the fields are built with froms always
-- When you wanna fill an attribute from another table you must open inside the from of the attribute a select with the desired data
-- 
select row_to_json(t)
from (
  -- Fill album
  select name,
  -- Fill songs (This should be a method like add_songs_query(album_id))
    (
      select array_to_json(array_agg(row_to_json(d)))
      from (
        select songs.name, songs.duration,  
        -- Fill artists field (This should be a method like add_artists(comparator))
        (
          SELECT array_to_json(array_agg(row_to_json(profiles))) 
          --  you can do directly this (but you cant specify fields)
          -- FROM profiles
          --            WHERE profiles.username=ANY(songs.artists)
          -- Or you can do this which you specify fields
          FROM (
            SELECT profiles.name, profiles.username FROM profiles WHERE profiles.username=ANY(songs.artists)
          ) profiles
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