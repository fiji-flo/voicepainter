package main

import (
	"fmt"
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

func image(w http.ResponseWriter, r *http.Request) {
	var what string;

	what = r.FormValue("what");
	fmt.Println("what: " + what);

	if (what == "") {
		w.WriteHeader(http.StatusBadRequest);
		fmt.Fprintf(w, "usage: image?what=dog");
	}

	if (what == "cat") {
		fmt.Fprintf(w, "https://clipartfest.com/download/c67e552105c1c72d268d22a1ccff70cf7421fbe8.html");
	} else if (what == "dog") {
		fmt.Fprintf(w, "https://www.cesarsway.com/sites/newcesarsway/files/d6/images/features/2012/sept/Dyeing-Your-Dogs-Hair-Is-a-Bad-Idea.jpg");
	} else if (what == "house") {
		fmt.Fprintf(w, "http://images.clipartpanda.com/clipart-house-House-Clip-Art-87.jpg");
	}
}

func main() {
	fmt.Println("VRS")

	http.HandleFunc("/debug", debug)
	http.HandleFunc("/image", image)

	http.ListenAndServe(":8080", nil)
}
