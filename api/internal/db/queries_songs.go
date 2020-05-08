package db

// GetSongsByAlbumIDQuery gets the songs $1 - albumID, $2 - published
const GetSongsByAlbumIDQuery = `
SELECT row_to_json(s) AS songs
FROM albums,
  (SELECT
     (SELECT array_to_json(array_agg(row_to_json(songs)))
      FROM
        (SELECT songs.name,
								songs.duration,
								songs.id,
								songs.albumid,
           (SELECT array_to_json(array_agg(row_to_json(profiles)))
            FROM
              (SELECT profiles.name,
											profiles.username,
											profiles.userid
               FROM profiles
               WHERE profiles.username=ANY(songs.artists) ) profiles) AS artists
         FROM songs
         WHERE songs.albumid=$1 ) songs) AS songs) s
WHERE albums.id=$1
  AND albums.published=$2
`
