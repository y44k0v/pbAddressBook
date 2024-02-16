package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	addressbookpb "pbAddressBook/proto"

	"time"

	// Deprecated but if not used proto>Marshal doesn't work
	"github.com/golang/protobuf/proto"

	// recommended but proto.Marshal doesn't work
	// TODO figure out how to use it
	//"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Creates an address book with one person
func createAB() *addressbookpb.AddressBook {

	timestamp := getTimeStamp()
	// person 1
	person1 := &addressbookpb.Person{
		Name:  "Jack Black",
		Id:    666123,
		Email: "jb@gmail.com",
		Phones: []*addressbookpb.Person_PhoneNumber{
			{
				Number: "666-123-5567",
				Type:   addressbookpb.Person_HOME,
			},
			{
				Number: "666-123-5568",
				Type:   addressbookpb.Person_MOBILE,
			},
			{
				Number: "666-123-8887",
				Type:   addressbookpb.Person_WORK,
			},
		},
		LastUpdated: timestamp,
	}

	ab := addressbookpb.AddressBook{
		People: []*addressbookpb.Person{
			person1,
		},
	}
	return &ab

}

// display persons detail in separated lines
func printPerson(p *addressbookpb.Person) {

	fmt.Println("Name:", p.GetName())
	fmt.Println("Id:", p.GetId())
	fmt.Println("Email:", p.GetEmail())
	for i, phone := range p.GetPhones() {
		fmt.Printf("Number %d: %s, Type: %s\n", i+1, phone.GetNumber(), phone.GetType())
	}
	fmt.Println("Last Updated:", p.GetLastUpdated())

}

// generated timestampt to update record TODO change format
func getTimeStamp() *timestamppb.Timestamp {
	timestamp := timestamppb.Timestamp{
		Seconds: time.Now().Unix(),
		Nanos:   int32(time.Nanosecond),
	}
	return &timestamp

}

// sets phonve type according to string in record
func setPhoneType(record string) addressbookpb.Person_PhoneType {
	if record == "home" {
		return addressbookpb.Person_HOME
	} else if record == "mobile" {
		return addressbookpb.Person_MOBILE
	} else {
		return addressbookpb.Person_WORK
	}
}

func main() {

	// creating a person
	p := addressbookpb.Person{
		Name:  "Kyle Gas",
		Id:    666177,
		Email: "kg@gmail.com",
		Phones: []*addressbookpb.Person_PhoneNumber{
			{
				Number: "555-123-5567",
				Type:   addressbookpb.Person_HOME,
			},
			{
				Number: "555-123-5568",
				Type:   addressbookpb.Person_MOBILE,
			},
		},
		LastUpdated: getTimeStamp(),
	}

	printPerson(&p)

	// 2 seconds delay to force change in timestamp
	time.Sleep(2000000000)

	// creating book with one person
	myAB := createAB()
	fmt.Println("\n", *myAB)

	// adding an extra person to the book
	myAB.People = append(myAB.People, &p)
	fmt.Println("\n", *myAB, "\n")

	// printing persons' names and emails
	for _, p := range myAB.People {
		fmt.Println(p.GetName(), p.GetEmail())
	}

	// reading CSV and adding the content to the book
	// Open the CSV file
	file, err := os.Open("data/random_data.csv")
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read the CSV headers
	headers, err := reader.Read()
	if err != nil {
		fmt.Println("Error reading CSV headers:", err)
		return
	}

	fmt.Printf("\nCSV headers:\n%s\n\n", headers)

	// read records from csv
	for {
		record, err := reader.Read()
		if err != nil {
			break // End of file
		}
		// converting record for ID to int32 from string
		var id int32
		fmt.Sscan(record[0], &id)

		// Create a new person
		person := addressbookpb.Person{
			Name:  record[1],
			Id:    id,
			Email: record[2],
			Phones: []*addressbookpb.Person_PhoneNumber{
				{
					Number: record[3],
					Type:   setPhoneType(record[6]),
				},
				{
					Number: record[4],
					Type:   setPhoneType(record[7]),
				},
				{
					Number: record[5],
					Type:   setPhoneType(record[8]),
				},
			},
			LastUpdated: getTimeStamp(),
		}

		// Append person to the book
		myAB.People = append(myAB.People, &person)

	}

	fmt.Println("Person in myAB:")
	printPerson(myAB.People[55])

	// write book to disk as binary file
	out, err := proto.Marshal(myAB)
	if err != nil {
		log.Fatalln("Failed to encode address book:", err)
	}
	if err := os.WriteFile("data/addressBook.bin", out, 0644); err != nil {
		log.Fatalln("Failed to write address book:", err)
	}

	// Read the existing address book.
	in, err := os.ReadFile("data/addressBook.bin")
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}

	book := &addressbookpb.AddressBook{}
	if err := proto.Unmarshal(in, book); err != nil {
		log.Fatalln("Failed to parse address book:", err)
	}

	fmt.Println("\nPerson in book read from binary file:")
	printPerson(book.People[76])

}
