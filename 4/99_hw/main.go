package main

import (
	"bufio"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"

	"hw4/user"
)

const token = "123"

type (
	fastUser = user.User
	Users    = user.Users
)

type Row struct {
	ID        int    `xml:"id"`
	Age       int    `xml:"age"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Gender    string `xml:"gender"`
	About     string `xml:"about"`
}

func runServer(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", SearchServer)

	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	fmt.Println("starting server at", addr)
	server.ListenAndServe()
}

func handleQuery(values url.Values) (string, string, int, int, int, error) {
	query := values.Get("query")
	order_field := values.Get("order_field")
	if order_field != "Id" && order_field != "Age" && order_field != "Name" {
		return "", "", 0, 0, 0, errors.New("bad_order_field")
	}
	order_by, err := strconv.Atoi(values.Get("order_by"))
	if err != nil {
		return "", "", 0, 0, 0, err
	}
	limit, err := strconv.Atoi(values.Get("limit"))
	if err != nil {
		return "", "", 0, 0, 0, err
	}
	offset, err := strconv.Atoi(values.Get("offset"))
	if err != nil {
		return "", "", 0, 0, 0, err
	}
	return query, order_field, order_by, limit, offset, err
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("AccessToken") != token {
		http.Error(w, "bad_token", http.StatusUnauthorized)
		return
	}

	query, order_field, order_by, limit, offset, err := handleQuery(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, err := os.Open("dataset.xml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	input := bufio.NewReader(file)
	decoder := xml.NewDecoder(input)

	users := make(Users, 0, 64)
	var row Row
	i := 0

	for {
		t, err := decoder.Token()
		if err != nil && err != io.EOF {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if err == io.EOF {
			break
		}
		if t == nil {
			fmt.Println("t is nil break")
		}

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "row" {
				decoder.DecodeElement(&row, &se)

				if strings.Contains(strings.ToLower(row.FirstName), strings.ToLower(query)) ||
					strings.Contains(strings.ToLower(row.LastName), strings.ToLower(query)) ||
					strings.Contains(strings.ToLower(row.About), strings.ToLower(query)) {
					i++
					if i > limit {
						break
					}
					if i > offset {
						users = append(users, fastUser{
							Id:     row.ID,
							Name:   row.FirstName + " " + row.LastName,
							Age:    row.Age,
							About:  row.About,
							Gender: row.Gender,
						})
					}

				}
			}
		}
	}

	if order_by != 0 {
		slices.SortFunc(users, func(a, b fastUser) int {
			switch order_field {
			case "Id":
				if a.Id < b.Id {
					return order_by
				} else {
					return -order_by
				}
			case "Age":
				if a.Age < b.Age {
					return order_by
				} else {
					return -order_by
				}
			default:
				if a.Name < b.Name {
					return order_by
				} else {
					return -order_by
				}
			}
		})
	}

	data, err := users.MarshalJSON()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(data)
}

func main() {
	runServer(":8080")
}
