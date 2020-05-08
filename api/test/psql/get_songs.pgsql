SELECT row_to_json(s) AS songs
FROM albums,

  (SELECT
     (SELECT array_to_json(array_agg(row_to_json(songs)))
      FROM
        (SELECT songs.name,
                songs.duration,
           (SELECT array_to_json(array_agg(row_to_json(profiles)))
            FROM
              (SELECT profiles.name,
                      profiles.username,
                      profiles.userid
               FROM profiles
               WHERE profiles.username=ANY(songs.artists) ) profiles) AS artists
         FROM songs
         WHERE songs.albumid=82 ) songs) AS songs) s
WHERE albums.id=82
  AND albums.published=FALSE