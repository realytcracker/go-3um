package main

import "time"

//Config is loaded from config.json
type Config struct {
	MySQLHost     string `json:"MySQLHost"`
	MySQLPort     string `json:"MySQLPort"`
	MySQLUser     string `json:"MySQLUser"`
	MySQLPassword string `json:"MySQLPassword"`
	MySQLDatabase string `json:"MySQLDatabase"`
	HashKey       string `json:"HashKey"`
	BlockKey      string `json:"BlockKey"`
	Salt          string `json:"Salt"`
}

//Token is the stuff dumped to and from the /api/user/code endpoint,
//but i also used map[string]string like a numbnuts
type Token struct {
	Password string `json:"password"`
	Hash     string `json:"hash"`
}

//Board is the json object for a board
type Board struct {
	ID          uint16 `json:"id,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Position    uint16 `json:"position"`
}

//DBBoard is the database object for a board
type DBBoard struct {
	ID          uint16 `db:"id,omitempty"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Position    uint16 `db:"position"`
}

//Thread is the json object for a thread
type Thread struct {
	ID           uint64 `json:"id,omitempty"`
	Board        uint16 `json:"board"`
	Author       string `json:"author"`
	Title        string `json:"title"`
	Body         string `json:"body"`
	Created      uint64 `json:"created"`
	LastModified uint64 `json:"last_modified"`
	IP           string `json:"ip,omitempty"`
	IsSticky     bool   `json:"is_sticky"`
	IsLocked     bool   `json:"is_locked"`
}

//ThreadFull is the json object for a thread with posts
type ThreadFull struct {
	ID           uint64 `json:"id,omitempty"`
	Board        uint16 `json:"board"`
	Author       string `json:"author"`
	Title        string `json:"title"`
	Body         string `json:"body"`
	Created      uint64 `json:"created"`
	LastModified uint64 `json:"last_modified"`
	IP           string `json:"ip,omitempty"`
	Posts        []Post `json:"posts"`
	IsSticky     bool   `json:"is_sticky"`
	IsLocked     bool   `json:"is_locked"`
}

//DBThread is the database object for a thread
type DBThread struct {
	ID           uint64    `db:"id,omitempty"`
	Board        uint16    `db:"board"`
	Author       string    `db:"author"`
	Title        string    `db:"title"`
	Body         string    `db:"body"`
	Created      time.Time `db:"created"`
	LastModified time.Time `db:"last_modified"`
	IP           string    `db:"ip,omitempty"`
	IsSticky     bool      `db:"is_sticky"`
	IsLocked     bool      `db:"is_locked"`
}

//Post is the json object for a post
type Post struct {
	ID           uint64 `json:"id,omitempty"`
	Thread       uint64 `json:"thread"`
	Author       string `json:"author"`
	Body         string `json:"body"`
	Created      uint64 `json:"created"`
	LastModified uint64 `json:"last_modified"`
	IP           string `json:"ip,omitempty"`
}

//DBPost is the database object for a post
type DBPost struct {
	ID           uint64    `db:"id,omitempty"`
	Thread       uint64    `db:"thread"`
	Author       string    `db:"author"`
	Body         string    `db:"body"`
	Created      time.Time `db:"created"`
	LastModified time.Time `db:"last_modified"`
	IP           string    `db:"ip,omitempty"`
}

//User is the json object for a user
type User struct {
	ID      uint64 `json:"id,omitempty"`
	Hash    string `json:"hash"`
	IsAdmin bool   `json:"is_admin"`
	IsMod   bool   `json:"is_mod"`
}

//DBUser is the database object for a user
type DBUser struct {
	ID      uint64 `db:"id,omitempty"`
	Hash    string `db:"hash"`
	IsAdmin bool   `db:"is_admin"`
	IsMod   bool   `db:"is_mod"`
}
