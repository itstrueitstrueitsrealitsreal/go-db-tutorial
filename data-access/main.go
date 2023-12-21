package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var db *sql.DB

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func main() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)
	// Hard code ID 2 to test the query
	alb, err := albumByID(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", alb)

	albID, err := addAlbum(Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID of added album: %v\n", albID)

	// testing delete
	id, err := deleteAlbum(3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album %v deleted.\n", id)

	// test view albums
	albums, err = viewAlbums()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Albums found: %v\n", albums)

	// test update album
	id, err = updateAlbum(1, Album{
		Title:  "My Beautiful Dark Twisted Fantasy",
		Artist: "Kanye West",
		Price:  49.99,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album %d updated.\n", id)
	alb, err = albumByID(1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album after update: %v\n", alb)

}

// albumsByArtist queries for albums that have the specified artist
// name.
func albumsByArtist(name string) ([]Album, error) {
	// An albums slice to hold data from returned rows.
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

// albumByID queries for the album with the specified ID.
func albumByID(id int64) (Album, error) {
	// An album to hold data from the returned row.
	var alb Album

	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}
	return alb, nil
}

// addAlbum adds the specified album to the database,
// returning the album ID of the new entry
func addAlbum(alb Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}

// deleteAlbum deletes the album with the specified ID,
// returning the album that was deleted
func deleteAlbum(id int64) (int64, error) {
	_, err := db.Exec("DELETE FROM album WHERE ID = ?", id)
	if err != nil {
		return 0, fmt.Errorf("deleteAlbum: %v", err)
	}
	return id, nil
}

// Implement viewAlbums which displays all albums currently in the collection
// viewAlbums returns all albums in the database.
func viewAlbums() ([]Album, error) {
	// An albums slice to hold data from returned rows.
	var albums []Album

	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		return nil, fmt.Errorf("viewAlbums: %v", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("viewAlbums: %v", err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("viewAlbums: %v", err)
	}
	return albums, nil
}

// updateAlbum updated the album with the specified ID with new attributes, returns the updated DB
func updateAlbum(id int64, album Album) (int64, error) {
	_, err := db.Query("UPDATE album SET title=?, artist=?, price=? WHERE id=?", album.Title, album.Artist,
		album.Price, id)
	if err != nil {
		return 0, fmt.Errorf("updateAlbum %d: %v", id, err)
	}

	return id, nil
}
