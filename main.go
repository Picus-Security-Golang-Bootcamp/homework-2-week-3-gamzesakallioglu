package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func searchCommand(listOfBooks Books, searchText string) {

	searchResult := checkContains(listOfBooks, searchText)
	if len(searchResult) > 0 {
		fmt.Printf("--- SEARCH RESULTS (%d) ---", len(searchResult))
		fmt.Println()
		listCommand(searchResult)
	} else {
		fmt.Println("Sorry, we don't have what you look for.\nSearch something else or come later")
	}

}

// executes list command
func listCommand(listOfBooks Books) {
	i := 0
	for _, book := range listOfBooks {
		i++
		fmt.Print(i, ".")
		listOneBook(book)
	}
}

// execute delete   -- updates total stock as 0
func deleteCommand(listOfBooks Books, idToSearch int) Books {

	ifExists, i, b := getBookById(listOfBooks, idToSearch)
	if ifExists {
		b = setStock(b, b.totalStock)
		listOfBooks[i] = b
	} else {
		fmt.Println("Sorry, we don't have the book with given id")
	}
	return listOfBooks
}

// execute buy
func buyCommand(listOfBooks Books, idToSearch int, quantity int) Books {

	ifExists, i, b := getBookById(listOfBooks, idToSearch)
	if ifExists {
		// setStock - returns the updated book
		b = setStock(b, quantity)
	} else {
		fmt.Println("Sorry, we don't have the book with given Id")
	}
	listOfBooks[i] = b
	return listOfBooks
}

// execute the get command. Prints - the book if the book with given id exits - the message if the book Doesn't exists
func getCommand(listOfBooks Books, id int) {
	ifExists, _, b := getBookById(listOfBooks, id)
	if ifExists {
		listOneBook(b)
	} else {
		fmt.Println("Sorry, we don't have the book with id ", id)
	}
}

// To list the books properties
func listOneBook(book Book) {
	fmt.Print(" ", book.name, ", ", book.Author.name, " ", book.totalPages, " pages.")
	fmt.Print(" Price: ", book.price, " TL")
	fmt.Print(" Total Stock: ", book.totalStock)
	fmt.Println()
}

// returns search results for search command
func checkContains(listOfBooks Books, str string) Books {
	var listOfResults Books
	for _, book := range listOfBooks {
		if strings.Contains(strings.ToLower(book.name), strings.ToLower(str)) {
			listOfResults = append(listOfResults, book)
		}
	}
	return listOfResults
}

// returns false if user doesn't give any arguments
func checkIfOrderValid() bool {
	if len(os.Args) < 2 {
		return false
	}
	return true
}

// returns arguments slice from 1st index and arguments string
func readArguments() ([]string, string) {
	argumentsString := ""
	for _, v := range os.Args[1:] {
		argumentsString += v + " "
	}
	argumentsString = argumentsString[:len(argumentsString)-1]
	argumentsString = strings.ToLower(argumentsString)
	return os.Args[1:], argumentsString
}

// To see if arguments valid as quantitiy and data format
func checkIfArgumentsValid(argumentsStringArray []string, argumentNumber int) bool {
	if len(argumentsStringArray) > argumentNumber {
		return false
	}
	// If the user gives something other than a number it will return false (0 Ascii Number:0 9 Ascii Number: 57)
	for i := 1; i < argumentNumber; i++ {
		arg := argumentsStringArray[i]
		for _, v := range arg {
			if v < 48 || v > 57 {
				return false
			}
		}
	}
	return true
}

// if orders not valid print the message
func ordersNotValidM() {
	fmt.Println("--- The right usage of commands ---")
	fmt.Println("go run main.go list - to list all the books")
	fmt.Println("go run main.go get ID    - to list the book with given ID")
	fmt.Println("go run main.go delele ID - to delete the book with given ID")
	fmt.Println("go run main.go buy ID    - to buy given amount of the book with given ID")
	fmt.Println("go run main.go searcText - to search the books has given name")
}

// returns if the book with given id exists, the index and the book
func getBookById(listOfBooks Books, idToSearch int) (bool, int, Book) {
	index := 0
	bookUserLookedFor := Book{}
	for i, book := range listOfBooks {
		if book.id == idToSearch {
			index = i
			bookUserLookedFor = book
			return true, index, bookUserLookedFor
		}
	}
	return false, index, bookUserLookedFor
}

// set the given book's stock - returns the updated book
func setStock(book Book, quantity int) Book {
	if book.totalStock < quantity {
		fmt.Println("We have only", book.totalStock, "amount of the book.")
		book.totalStock = 0
	} else {
		book.totalStock -= quantity
	}
	fmt.Println("--- Updated Book ---")
	listOneBook(book)
	return book
}

type Author struct {
	id   int
	name string
}

type Book struct {
	id         int
	name       string
	totalPages int
	totalStock int
	price      float32
	stockCode  string
	ISBN       string
	Author
}
type Books []Book

func addNewBookToArchive(listOfBooks Books, id int, name string, totalPages int, totalStock int, price float32, stockCode string, ISBN string, author Author) Books {

	listOfBooks = append(listOfBooks, Book{id, name, totalPages, totalStock, price, stockCode, ISBN, author})
	return listOfBooks
}

func main() {

	var listOfBooks Books

	listOfBooks = addNewBookToArchive(listOfBooks, 1, "Lord of The Rings: The Fellowship of the Ring", 479, 30, 30, "056SBF", "234532", Author{2, "J.R.R Tolkien"})
	listOfBooks = addNewBookToArchive(listOfBooks, 2, "Lord of The Rings: The Two Towers", 415, 15, 30, "056SBF", "234532", Author{2, "J.R.R Tolkien"})
	listOfBooks = addNewBookToArchive(listOfBooks, 3, "Lord of The Rings: The Return of the King", 347, 10, 30, "056SBF", "234532", Author{2, "J.R.R Tolkien"})
	listOfBooks = addNewBookToArchive(listOfBooks, 4, "The Great Gatsby", 208, 40, 30, "056SBF", "234532", Author{1, "F. Scott Fitzgerald"})

	// check if any argument given. If not print the right usage
	orderValid := checkIfOrderValid()
	if orderValid {
		// returns arguments slice from 1st index and arguments in the lower case string format
		argumentsStringArray, _ := readArguments()
		switch argumentsStringArray[0] {
		case "get":
			// First, check if argument number is correct and arguments are int
			if checkIfArgumentsValid(argumentsStringArray, 2) {
				id, _ := strconv.Atoi(argumentsStringArray[1])
				getCommand(listOfBooks, id)
			} else {
				fmt.Println("The right usage of get command is: get IdNumber")
			}

		case "delete":
			// First, check if argument number is correct and arguments are int
			if checkIfArgumentsValid(argumentsStringArray, 2) {
				id, _ := strconv.Atoi(argumentsStringArray[1])
				listOfBooks = deleteCommand(listOfBooks, id)
			} else {
				fmt.Println("The right usage of delete command is: delete IdNumber")
			}

		case "buy":
			// First, check if argument number is correct and arguments are int
			if checkIfArgumentsValid(argumentsStringArray, 3) {
				id, _ := strconv.Atoi(argumentsStringArray[1])
				quantity, _ := strconv.Atoi(argumentsStringArray[2])
				listOfBooks = buyCommand(listOfBooks, id, quantity)
			} else {
				fmt.Println("The right usage of buy command is: buy IdNumber quantity")
			}

		case "list":
			fmt.Println("--- LIST OF ALL BOOKS ---")
			listCommand(listOfBooks)

		case "search":
			wordToSearch := ""
			for _, v := range argumentsStringArray[1:] {
				wordToSearch = wordToSearch + v + " "
			}
			searchCommand(listOfBooks, wordToSearch)
		}
	} else {
		ordersNotValidM()
	}
}
