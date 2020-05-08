package db

// GetAlbumFullInformation Gets the full information of an album (album_id, published)
const GetAlbumFullInformation = `SELECT row_to_json(t) AS album
FROM
  (
		SELECT name, albums.description, albums.published, albums.uploaddate, albums.id,
     (SELECT array_to_json(array_agg(row_to_json(d)))
      FROM
        (SELECT songs.name,
								songs.duration,
								songs.albumid, 
								songs.id,
           (SELECT array_to_json(array_agg(row_to_json(profiles))) 
            FROM
              (SELECT profiles.name,
                      profiles.username
               FROM profiles
               WHERE profiles.username=ANY(songs.artists) ) profiles) AS artists
         FROM songs
         WHERE songs.albumid=albums.id ) d) AS songs,

     (SELECT array_to_json(array_agg(row_to_json(profiles)))
      FROM profiles
      WHERE profiles.username=ANY(albums.artists) ) AS artists,

		(SELECT array_to_json(array_agg(row_to_json(images))) FROM
		images WHERE images.id=albums.id AND 
		(images.type='cover_image' OR images.type='header_image') ) AS images
   FROM albums WHERE albums.id=$1 AND albums.published=$2) t
`

// GetAlbumsFullInformation Gets the full information of albums (album_ids, published)
const GetAlbumsFullInformation = `SELECT array_to_json(array_agg(row_to_json(t))) AS album
FROM
  (
		SELECT name, albums.description, albums.published, albums.uploaddate, albums.id,
     (SELECT array_to_json(array_agg(row_to_json(d)))
      FROM
        (SELECT songs.name,
								songs.duration,
								songs.albumid, 
								songs.id,
           (SELECT array_to_json(array_agg(row_to_json(profiles))) 
            FROM
              (SELECT profiles.name,
                      profiles.username
               FROM profiles
               WHERE profiles.username=ANY(songs.artists) ) profiles) AS artists
         FROM songs
         WHERE songs.albumid=albums.id ) d) AS songs,

     (SELECT array_to_json(array_agg(row_to_json(profiles)))
      FROM profiles
      WHERE profiles.username=ANY(albums.artists) ) AS artists,

		(SELECT array_to_json(array_agg(row_to_json(images))) FROM
		images WHERE images.id=albums.id AND 
		(images.type='cover_image' OR images.type='header_image') ) AS images
   FROM albums WHERE albums.id=ANY($1) AND albums.published=$2) t
`
