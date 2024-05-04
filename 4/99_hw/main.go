package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
)

type Row struct {
	ID        int    `xml:"id"`
	Age       int    `xml:"age"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Gender    string `xml:"gender"`
	About     string `xml:"about"`
}

func main() {
	query := "ab"
	order_field := "Id"
	order_by := -1
	limit := 10
	offset := 0

	file, err := os.Open("dataset.xml") // For read access.
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	input := bufio.NewReader(file)
	decoder := xml.NewDecoder(input)

	users := make([]User, 0, 64)
	var row Row
	i := 0

	for {
		t, tokenErr := decoder.Token()
		if tokenErr != nil && tokenErr != io.EOF {
			fmt.Println("error happend", tokenErr)
			break
		} else if tokenErr == io.EOF {
			break
		}
		if t == nil {
			fmt.Println("t is nil break")
		}

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "row" {
				decoder.DecodeElement(&row, &se)
			}

			if strings.Contains(strings.ToLower(row.FirstName), strings.ToLower(query)) ||
				strings.Contains(strings.ToLower(row.LastName), strings.ToLower(query)) ||
				strings.Contains(strings.ToLower(row.About), strings.ToLower(query)) {
				i++
				if i >= limit {
					break
				}
				if i > offset {
					users = append(users, User{
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

	if order_by != 0 {
		slices.SortFunc(users, func(a, b User) int {
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

	for i := 0; i < len(users); i++ {
		fmt.Println(users[i])
	}
}
