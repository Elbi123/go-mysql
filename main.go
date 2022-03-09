package main

import (
	"fmt"
	"os"
	"log"
	"database/sql"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Album struct {
	Id int64
	Title string
	Artist string
	Price float32
}

func albumsByArtist(name string) ([]Album, error) {

	var albums []Album

	// query first
	rows, err := db.Query("SELECT * FROM album WHERE artist = ? ", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	// close the reading with defer
	defer rows.Close()

	// loop through every row
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.Id, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	return albums, nil

}

func albumById(id int64) (Album, error) {
	var  album Album

	row := db.QueryRow("SELECT * FROM album where id = ? ", id)

	//defer row.Close();

	if err := row.Scan(&album.Id, &album.Title, &album.Artist, &album.Price); err != nil {
		if err == sql.ErrNoRows {
			return album, fmt.Errorf("albumsById %d: no such album", id)
		}
		return album, fmt.Errorf("albumsById %d: %v", id, err)
	}

	return album, nil
}

func addAlbum(album Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album(title, artist, price) VALUES(?, ?, ?)", album.Title, album.Artist, album.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	return id, nil
}

func main() {
	fmt.Println("working with dbs")
	db_config := mysql.Config{
		User: os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		AllowNativePasswords: true,
		Net:    "tcp",
	        Addr:   "localhost:3306",
		DBName: "recordings",
	}

	var err error
	db, err = sql.Open("mysql", db_config.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping();
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Database connected")

	albums, err := albumsByArtist("John Coltrane")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Albums found: %v\n", albums)

	album, err := albumById(1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Album found: %v\n", album)

	albumId, err := addAlbum(Album{
		Title: "The newly added album",
		Artist: "The New Artist",
		Price: 45.2,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID of added album: %v\n", albumId)
}
