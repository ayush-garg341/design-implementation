package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Book struct {
	Genre           string
	AuthorLastName  string
	PublicationYear int
}

var genres = []string{
	"fiction", "history", "science", "fantasy", "biography",
}

var authors = []string{
	"Sharma", "Verma", "Smith", "Johnson", "Brown",
	"Taylor", "Anderson", "Thomas", "Jackson", "White",
}

func generateBooks(n int) []Book {
	rand.Seed(time.Now().UnixNano())

	books := make([]Book, n)
	for i := 0; i < n; i++ {
		books[i] = Book{
			Genre:           genres[rand.Intn(len(genres))],
			AuthorLastName:  authors[rand.Intn(len(authors))],
			PublicationYear: rand.Intn(50) + 1975,
		}
	}
	return books
}

func sortBooks(books []Book) {
	sort.Slice(books, func(i, j int) bool {
		if books[i].Genre != books[j].Genre {
			return books[i].Genre < books[j].Genre
		}
		if books[i].AuthorLastName != books[j].AuthorLastName {
			return books[i].AuthorLastName < books[j].AuthorLastName
		}
		return books[i].PublicationYear < books[j].PublicationYear
	})
}

func printBooks(books []Book, limit int) {
	for i := 0; i < limit && i < len(books); i++ {
		fmt.Printf("%-10s | %-10s | %d\n",
			books[i].Genre,
			books[i].AuthorLastName,
			books[i].PublicationYear,
		)
	}
}

func main() {
	books := generateBooks(100)

	fmt.Println("Before Sorting:")
	printBooks(books, 10)

	sortBooks(books)

	fmt.Println("\nAfter Sorting:")
	printBooks(books, 20)
}
