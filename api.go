package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/gorilla/securecookie"
)

//GetCode sets the cookie if none is found, returns json of salted hash and pw
func GetCode(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		var sc = securecookie.New(hashKey, blockKey)
		var pw = RandomString(10)
		var pwh = pw + config.Salt
		h := sha256.New()
		h.Write([]byte(pwh))
		hash := hex.EncodeToString(h.Sum(nil))

		token := map[string]string{
			"password": pw,
			"hash":     hash[:30],
		}

		if encoded, err := sc.Encode("token", token); err == nil {
			cookie := &http.Cookie{
				Name:  "token",
				Value: encoded,
				Path:  "/",
			}

			http.SetCookie(w, cookie)
		}

		output, _ := json.Marshal(token)
		fmt.Fprintln(w, string(output[:]))
	} else if token.Value == "" {
		var sc = securecookie.New(hashKey, blockKey)
		var pw = RandomString(10)
		var pwh = pw + config.Salt
		h := sha256.New()
		h.Write([]byte(pwh))
		hash := hex.EncodeToString(h.Sum(nil))

		token := map[string]string{
			"password": pw,
			"hash":     hash[:30],
		}

		if encoded, err := sc.Encode("token", token); err == nil {
			cookie := &http.Cookie{
				Name:  "token",
				Value: encoded,
				Path:  "/",
			}

			http.SetCookie(w, cookie)
		}

		output, _ := json.Marshal(token)
		fmt.Fprintln(w, string(output[:]))
	} else {
		var sc = securecookie.New(hashKey, blockKey)
		value := make(map[string]string)
		if err = sc.Decode("token", token.Value, &value); err == nil {
			output, _ := json.Marshal(value)
			fmt.Fprintln(w, string(output[:]))
		}
	}
}

//SetCode handles the hash and cookie if a user changes their password
func SetCode(w http.ResponseWriter, r *http.Request) {
	var jstoken Token
	_ = json.NewDecoder(r.Body).Decode(&jstoken)
	var sc = securecookie.New(hashKey, blockKey)
	var pw = jstoken.Password
	var pwh = pw + config.Salt
	h := sha256.New()
	h.Write([]byte(pwh))
	hash := hex.EncodeToString(h.Sum(nil))

	token := map[string]string{
		"password": pw,
		"hash":     hash[:30],
	}

	if encoded, err := sc.Encode("token", token); err == nil {
		cookie := &http.Cookie{
			Name:  "token",
			Value: encoded,
			Path:  "/",
		}

		http.SetCookie(w, cookie)
	}

	output, _ := json.Marshal(token)
	fmt.Fprintln(w, string(output[:]))
}

//DeleteCode expires the token cookie
func DeleteCode(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "token",
		Value:  "",
		MaxAge: 0,
		Path:   "/",
	}

	http.SetCookie(w, cookie)
}

//CreateBoard creates a new board
func CreateBoard(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	var sc = securecookie.New(hashKey, blockKey)
	value := make(map[string]string)
	if err = sc.Decode("token", token.Value, &value); err == nil {
		if IsAdmin(value["hash"]) == true {
			var board Board

			_ = json.NewDecoder(r.Body).Decode(&board)
			boardCollection := dbsession.Collection("boards")
			boardCollection.Insert(DBBoard{Title: board.Title, Description: board.Description, Position: board.Position})
		}
	}
}

//ListBoard provides a list of available boards
func ListBoard(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	var sc = securecookie.New(hashKey, blockKey)
	value := make(map[string]string)
	if err = sc.Decode("token", token.Value, &value); err == nil {
		var boards []Board
		var dbboards []DBBoard
		var err error

		q := dbsession.SelectFrom("boards").OrderBy("position")
		err = q.All(&dbboards)
		if err != nil {
			log.Fatalf("q.All(): %q\n", err)
		}

		for _, b := range dbboards {
			boards = append(boards, Board{b.ID, b.Title, b.Description, b.Position})
		}

		output, _ := json.Marshal(boards)
		fmt.Fprintln(w, string(output[:]))
	}
}

//DeleteBoard deletes a board matching {id} in url
func DeleteBoard(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	var sc = securecookie.New(hashKey, blockKey)
	value := make(map[string]string)
	if err = sc.Decode("token", token.Value, &value); err == nil {
		if IsAdmin(value["hash"]) == true {
			var err error

			vars := mux.Vars(r)
			q := dbsession.DeleteFrom("boards").Where("id", vars["id"])
			_, err = q.Exec()

			if err != nil {
				log.Fatalf("q.Exec(): %q\n", err)
			}
		}
	}
}

//UpdateBoard updates the board with {id} in the url
func UpdateBoard(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	var sc = securecookie.New(hashKey, blockKey)
	value := make(map[string]string)
	if err = sc.Decode("token", token.Value, &value); err == nil {
		if IsAdmin(value["hash"]) == true {
			var board Board
			var err error

			vars := mux.Vars(r)
			_ = json.NewDecoder(r.Body).Decode(&board)
			q := dbsession.Update("boards").Set("title", board.Title, "description", board.Description, "position", board.Position).Where("id", vars["id"])
			_, err = q.Exec()

			if err != nil {
				log.Fatalf("q.Exec(): %q\n", err)
			}
		}
	}
}

//CreateThread creates a thread under board {id}
func CreateThread(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	var sc = securecookie.New(hashKey, blockKey)
	value := make(map[string]string)
	if err = sc.Decode("token", token.Value, &value); err == nil {
		var thread Thread

		vars := mux.Vars(r)
		_ = json.NewDecoder(r.Body).Decode(&thread)
		threadCollection := dbsession.Collection("threads")
		boardid, _ := strconv.ParseUint(vars["id"], 10, 64)
		threadCollection.Insert(DBThread{Board: uint16(boardid), Author: value["hash"], Title: thread.Title, Body: thread.Body, Created: time.Now(), LastModified: time.Now(), IP: r.RemoteAddr})
	}
}

//ListThread lists the threads of board {id} greater than {start} limit {amount}
func ListThread(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	var sc = securecookie.New(hashKey, blockKey)
	value := make(map[string]string)
	if err = sc.Decode("token", token.Value, &value); err == nil {
		var thread Thread
		var threads []Thread
		var dbthreads []DBThread
		var err error

		vars := mux.Vars(r)
		_ = json.NewDecoder(r.Body).Decode(&thread)
		boardid, _ := strconv.ParseUint(vars["id"], 10, 64)
		start, _ := strconv.ParseUint(vars["start"], 10, 64)
		amount, _ := strconv.Atoi(vars["amount"])

		q := dbsession.SelectFrom("threads").Where("board = ? AND id > ?", uint16(boardid), start).Limit(amount)
		err = q.All(&dbthreads)

		if err != nil {
			log.Fatalf("q.All(): %q\n", err)
		}

		for _, b := range dbthreads {
			if IsAdmin(value["hash"]) == true || IsMod(value["hash"]) == true {
				threads = append(threads, Thread{b.ID, b.Board, b.Author, b.Title, b.Body, uint64(b.Created.Unix()), uint64(b.LastModified.Unix()), b.IP, b.IsSticky, b.IsLocked})
			} else {
				threads = append(threads, Thread{b.ID, b.Board, b.Author, b.Title, b.Body, uint64(b.Created.Unix()), uint64(b.LastModified.Unix()), "", b.IsSticky, b.IsLocked})
			}
		}

		output, _ := json.Marshal(threads)
		fmt.Fprintln(w, string(output[:]))
	}
}

//GetThread lists the posts of thread {id} greater than {start} limit {amount}
func GetThread(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	var sc = securecookie.New(hashKey, blockKey)
	value := make(map[string]string)
	if err = sc.Decode("token", token.Value, &value); err == nil {
		var thread ThreadFull
		var dbthread DBThread
		var posts []Post
		var dbposts []DBPost
		var err error

		vars := mux.Vars(r)
		_ = json.NewDecoder(r.Body).Decode(&thread)
		threadid, _ := strconv.ParseUint(vars["id"], 10, 64)
		start, _ := strconv.ParseUint(vars["start"], 10, 64)
		amount, _ := strconv.Atoi(vars["amount"])

		q := dbsession.SelectFrom("threads").Where("id", threadid)
		err = q.One(&dbthread)

		if err != nil {
			log.Fatalf("q.One(): %q\n", err)
		}

		q2 := dbsession.SelectFrom("posts").Where("thread = ? AND id > ?", threadid, start).Limit(amount)
		err = q2.All(&dbposts)

		if err != nil {
			log.Fatalf("q2.All(): %q\n", err)
		}

		for _, b := range dbposts {
			if IsAdmin(value["hash"]) == true || IsMod(value["hash"]) == true {
				posts = append(posts, Post{b.ID, b.Thread, b.Author, b.Body, uint64(b.Created.Unix()), uint64(b.LastModified.Unix()), b.IP})
			} else {
				posts = append(posts, Post{b.ID, b.Thread, b.Author, b.Body, uint64(b.Created.Unix()), uint64(b.LastModified.Unix()), ""})
			}
		}

		thread = ThreadFull{dbthread.ID, dbthread.Board, dbthread.Author, dbthread.Title, dbthread.Body, uint64(dbthread.Created.Unix()), uint64(dbthread.LastModified.Unix()), "", posts, dbthread.IsSticky, dbthread.IsLocked}

		output, _ := json.Marshal(thread)
		fmt.Fprintln(w, string(output[:]))
	}
}

//DeleteThread deletes thread {id} and all related posts - very messy
func DeleteThread(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	var sc = securecookie.New(hashKey, blockKey)
	value := make(map[string]string)
	if err = sc.Decode("token", token.Value, &value); err == nil {
		if IsAdmin(value["hash"]) == true || IsMod(value["hash"]) == true {
			var err error

			vars := mux.Vars(r)
			q := dbsession.DeleteFrom("threads").Where("id", vars["id"])
			_, err = q.Exec()
			if err != nil {
				log.Fatalf("q.Exec(): %q\n", err)
			}
			q2 := dbsession.DeleteFrom("posts").Where("thread", vars["id"])
			_, err = q2.Exec()
			if err != nil {
				log.Fatalf("q2.Exec(): %q\n", err)
			}
		}
	}
}

//UpdateThread changes thread data
func UpdateThread(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	var sc = securecookie.New(hashKey, blockKey)
	value := make(map[string]string)
	if err = sc.Decode("token", token.Value, &value); err == nil {
		vars := mux.Vars(r)
		threadid, _ := strconv.ParseUint(vars["id"], 10, 64)
		if IsAdmin(value["hash"]) == true || IsMod(value["hash"]) == true || IsThreadOwner(value["hash"], threadid) == true {
			var thread Thread
			var err error

			_ = json.NewDecoder(r.Body).Decode(&thread)
			q := dbsession.Update("threads").Set("title", thread.Title, "body", thread.Body, "last_modified", time.Now(), "is_sticky", thread.IsSticky, "is_locked", thread.IsLocked).Where("id", threadid)
			_, err = q.Exec()

			if err != nil {
				log.Fatalf("q.Exec(): %q\n", err)
			}
		}
	}
}

//CreatePost creates a post under thread {id}
func CreatePost(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	var sc = securecookie.New(hashKey, blockKey)
	value := make(map[string]string)
	if err = sc.Decode("token", token.Value, &value); err == nil {
		var post Post
		var err error

		vars := mux.Vars(r)
		_ = json.NewDecoder(r.Body).Decode(&post)
		postCollection := dbsession.Collection("posts")
		threadid, _ := strconv.ParseUint(vars["id"], 10, 64)
		postCollection.Insert(DBPost{Thread: threadid, Author: value["hash"], Body: post.Body, Created: time.Now(), LastModified: time.Now(), IP: r.RemoteAddr})
		q := dbsession.Update("threads").Set("last_modified", time.Now()).Where("id", threadid)
		_, err = q.Exec()

		if err != nil {
			log.Fatalf("q.Exec(): %q\n", err)
		}
	}
}

//GetPost gets post {id} and dumps the data on it
func GetPost(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	var sc = securecookie.New(hashKey, blockKey)
	value := make(map[string]string)
	if err = sc.Decode("token", token.Value, &value); err == nil {
		vars := mux.Vars(r)
		postid, _ := strconv.ParseUint(vars["id"], 10, 64)
		var post Post
		var dbpost DBPost
		var err error

		q := dbsession.SelectFrom("posts").Where("id", postid)
		err = q.One(&dbpost)
		if err != nil {
			log.Fatalf("q.One(): %q\n", err)
		}
		post = Post{dbpost.ID, dbpost.Thread, dbpost.Author, dbpost.Body, uint64(dbpost.Created.Unix()), uint64(dbpost.LastModified.Unix()), ""}

		output, _ := json.Marshal(post)
		fmt.Fprintln(w, string(output[:]))
	}
}

//DeletePost deletes a post {id} if admin, mod, or userhash is author of post
func DeletePost(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	var sc = securecookie.New(hashKey, blockKey)
	value := make(map[string]string)
	if err = sc.Decode("token", token.Value, &value); err == nil {
		vars := mux.Vars(r)
		postid, _ := strconv.ParseUint(vars["id"], 10, 64)
		if IsAdmin(value["hash"]) == true || IsMod(value["hash"]) == true || IsPostOwner(value["hash"], postid) == true {
			var err error

			q := dbsession.DeleteFrom("posts").Where("id", postid)
			_, err = q.Exec()
			if err != nil {
				log.Fatalf("q.Exec(): %q\n", err)
			}
		}
	}
}

//UpdatePost changes post body
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	var sc = securecookie.New(hashKey, blockKey)
	value := make(map[string]string)
	if err = sc.Decode("token", token.Value, &value); err == nil {
		vars := mux.Vars(r)
		postid, _ := strconv.ParseUint(vars["id"], 10, 64)
		if IsAdmin(value["hash"]) == true || IsMod(value["hash"]) == true || IsPostOwner(value["hash"], postid) == true {
			var post Post
			var err error

			_ = json.NewDecoder(r.Body).Decode(&post)
			q := dbsession.Update("posts").Set("body", post.Body, "last_modified", time.Now()).Where("id", postid)
			_, err = q.Exec()

			if err != nil {
				log.Fatalf("q.Exec(): %q\n", err)
			}
		}
	}
}

//SetupBBS initializes the first admin
func SetupBBS(w http.ResponseWriter, r *http.Request) {
	if IsSetup() == false {
		var sc = securecookie.New(hashKey, blockKey)
		var pw = RandomString(10)
		var pwh = pw + config.Salt
		h := sha256.New()
		h.Write([]byte(pwh))
		hash := hex.EncodeToString(h.Sum(nil))

		token := map[string]string{
			"password": pw,
			"hash":     hash[:30],
		}

		if encoded, err := sc.Encode("token", token); err == nil {
			cookie := &http.Cookie{
				Name:  "token",
				Value: encoded,
				Path:  "/",
			}

			http.SetCookie(w, cookie)
		}

		userCollection := dbsession.Collection("users")
		userCollection.Insert(DBUser{Hash: token["hash"], IsAdmin: true, IsMod: true})

		output, _ := json.Marshal(token)
		fmt.Fprintln(w, string(output[:]))
	}
}
