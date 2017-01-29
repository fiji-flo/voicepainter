package main

import (
	"fmt"
	"log"
	"strings"
	"encoding/json"
	"io/ioutil"
	"database/sql"
	"net/url"
	"net/http"
	_ "github.com/mattn/go-sqlite3"
)

func debug(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nrequest: ", r.RequestURI)

	fmt.Println("vals: ")
	r.ParseForm()

	for i, _ := range r.Form {
		fmt.Println("\t" + i + ": " + r.FormValue(i))
	}

	fmt.Println("From: " + r.RemoteAddr)
}

func getOpenclipart(s string) string {
	var u *url.URL;

	type Svg struct {
		Url string
	}
	type Payload struct {
		Title string
		Svg Svg
	}
	type ApiResponse struct {
		Msg string
		Payload []Payload
	}

	var ar ApiResponse;

	u, err := url.Parse("https://openclipart.org/search/json/");
	if (err != nil) {
		fmt.Println("Parsing url failed\n");
		return "";
	}

	q := u.Query();
	q.Set("query", s);
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if (err != nil || resp.StatusCode != 200) {
		fmt.Println("GET failed\n");
		return "";
	}

	body, err := ioutil.ReadAll(resp.Body);
	resp.Body.Close();
	if (err != nil) {
		fmt.Println("Read failed\n");
		return "";
	}

	err = json.Unmarshal(body, &ar);
	if (err != nil) {
		fmt.Println("unmarshal failed\n");
		return "";
	}

	return ar.Payload[0].Svg.Url;
}

func image(w http.ResponseWriter, r *http.Request) {
	var what string;

	what = r.FormValue("what");
	fmt.Println("what: " + what);

	if (what == "") {
		w.WriteHeader(http.StatusBadRequest);
		fmt.Fprintf(w, "usage: image?what=dog");
		return;
	}

	fmt.Fprintf(w, getOpenclipart(what));
}

func check_noun(s string, db *sql.DB) string {
	rows, err := db.Query("SELECT noun FROM nouns WHERE noun = ?", s);
	if (err != nil) {
		log.Fatal("select failed");
	}

	for rows.Next() {
		var result string;
		err = rows.Scan(&result);
		if (result == s) {
			return result
		}
	}

	return "";
}

func find_noun(s []string, pos int, db *sql.DB) int {

	// circle around the word until we find a noun
	for i := 1; i < len(s); i++ {
		var ppos int;

		ppos = pos - i;
		if (ppos >= 0 && ppos < len(s) && s[ppos] != "") {
			if (check_noun(s[ppos], db) != "") {
				return ppos;
			}
		}
		ppos = pos + i;
		if (ppos >= 0 && ppos < len(s) && s[ppos] != "") {
			if (check_noun(s[ppos], db) != "") {
				return ppos;
			}
		}
	}

	return -1;
}

func image_nlp(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var what string;

	type ApiResponse struct {
		Left string
		Right string
		Background string
		Middle string
	}
	var ar ApiResponse;

	what = r.FormValue("what");

	if (what == "") {
		w.WriteHeader(http.StatusBadRequest);
		fmt.Fprintf(w, "usage: image2?what=draw%20a%20house%20in%20the%20background");
		return;
	}

	var words []string;
	words = strings.Split(what, " ");

	for index, elem := range words {
		var left int;
		var right int;
		var middle int;
		var background int;
		if (elem == "left" || elem == "Left") {
			left = find_noun(words, index, db);
			ar.Left = getOpenclipart(words[left]);
			words[left] = "";
			words[index] = "";
		} else if (elem == "right" || elem == "Right") {
			right = find_noun(words, index, db);
			ar.Right = getOpenclipart(words[right]);
			words[right] = "";
			words[index] = "";
		} else if (elem == "middle" || elem == "middle" || elem == "mid" || elem == "Mid") {
			middle = find_noun(words, index, db);
			ar.Middle = getOpenclipart(words[middle]);
			words[middle] = "";
			words[index] = "";
		} else if (elem == "background" || elem == "Background") {
			background = find_noun(words, index, db);
			ar.Background = getOpenclipart(words[background]);
			words[background] = "";
			words[index] = "";
		}
		//fmt.Printf("index %d: %s\n", index, words[index]);
	}

	var resp []byte;
	resp, _ = json.Marshal(ar);
	fmt.Fprintf(w, string(resp));
}

func staticFiles(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/"+r.URL.Path[1:])
}

func main() {
	fmt.Println("VRS")

	db, err := sql.Open("sqlite3", "./words.db")
	if (err != nil) {
		log.Fatal("could not open words.db");
	}

	rows, err := db.Query("SELECT * FROM adjectives LIMIT 10")
	defer rows.Close()
	for rows.Next() {
		var a string;
		err = rows.Scan(&a)
		if (err != nil) {
			log.Fatal("rows.Scan failed");
		}
	}


	http.HandleFunc("/debug", debug)
	http.HandleFunc("/static/", staticFiles)
	http.HandleFunc("/image", image)
	http.HandleFunc("/image2", func(w http.ResponseWriter, r *http.Request){image_nlp(w, r, db)})

	http.ListenAndServe(":8080", nil)
}
