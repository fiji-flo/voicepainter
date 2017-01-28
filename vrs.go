package main

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"net/http"
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

	if (what == "cat") {
		fmt.Fprintf(w, "https://clipartfest.com/download/c67e552105c1c72d268d22a1ccff70cf7421fbe8.html");
	} else if (what == "dog") {
		fmt.Fprintf(w, "https://www.cesarsway.com/sites/newcesarsway/files/d6/images/features/2012/sept/Dyeing-Your-Dogs-Hair-Is-a-Bad-Idea.jpg");
	} else if (what == "house") {
		fmt.Fprintf(w, "http://images.clipartpanda.com/clipart-house-House-Clip-Art-87.jpg");
	} else {
		fmt.Fprintf(w, getOpenclipart(what));
	}

}

func main() {
	fmt.Println("VRS")

	http.HandleFunc("/debug", debug)
	http.HandleFunc("/image", image)

	http.ListenAndServe(":8080", nil)
}
