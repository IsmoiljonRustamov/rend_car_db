package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	PostgresHost     = "localhost"
	PostgresPort     = 5432
	PostgresUser     = "ismoiljon12"
	PostgresPassword = "12"
	PostgresDatabase = "rend_migration_db"
)

type Response struct {
	responseBranch []*Branches
	responseOffice []*Offices
	responseCar    []*Car
}

type Offices struct {
	Id     int
	Name   string
	Branch []Branches
}

type Branches struct {
	Id       int
	officeId int
	Name     string
	Car      []Car
	Address  []Address
}

type Car struct {
	Id       int
	BranchId int
	Name     string
	Color    string
	CostDay  float64
	Amount   int
	Customer []Customer
}

type Address struct {
	Id       int
	BranchId int
	Street   string
	City     string
}

type Customer struct {
	Id           int
	CarId        int
	Name         string
	Age          int
	Phone_number string
	Address      string
}

func main() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		PostgresHost, PostgresPort, PostgresUser, PostgresPassword, PostgresDatabase)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println("Failed to connected database: ", err)
	} else {
		log.Println("Connected")
	}

	defer db.Close()
	InsertDatabase(db)

}

func InsertDatabase(db *sql.DB) (int, int, int) {
	car := &Offices{
		Name: "Rental",
		Branch: []Branches{
			{
				officeId: 1,
				Name:     "Toshkent",
				Car: []Car{
					{
						BranchId: 1,
						Name:     "Malibu 2",
						Color:    "Black",
						CostDay:  1200.22,
						Amount:   100,
						Customer: []Customer{
							{
								CarId:        1,
								Name:         "Aziz",
								Age:          23,
								Phone_number: "+99896 558 78 87",
								Address:      "Toshkent, Chilonzor-19",
							},
						},
					},
				},
				Address: []Address{
					{
						BranchId: 1,
						Street:   "Shuhrat street",
						City:     "Toshkent",
					},
				},
			},
		},
	}
	tx, err := db.Begin()
	if err != nil {
		log.Println("Failed to begin:", err)
	}
	var OfficeId, branchId, carId int
	err = tx.QueryRow("INSERT INTO offices(name) VALUES($1) RETURNING id", car.Name).Scan(&OfficeId)
	if err != nil {
		tx.Rollback()
		log.Println("Failed to insert offices:", err)
	}

	for _, branch := range car.Branch {
		err = tx.QueryRow("INSERT INTO branches(office_id,name) VALUES($1,$2) RETURNING id", OfficeId, branch.Name).Scan(&branchId)
		if err != nil {
			tx.Rollback()
			log.Println("Failed to insert branches:", err)
		}
		for _, car := range branch.Car {
			err = tx.QueryRow("INSERT INTO cars(branch_id,name,color,cost_day,amount) VALUES($1,$2,$3,$4,$5) RETURNING id", branchId, car.Name, car.Color, car.CostDay, car.Amount).Scan(&carId)
			if err != nil {
				tx.Rollback()
				log.Println("Failed to insert cars: ", err)
			}
			for _, customer := range car.Customer {
				_, err = tx.Exec("INSERT INTO customer(car_id,name,age,phone_number,address) VALUES($1,$2,$3,$4,$5)", carId, customer.Name, customer.Age, customer.Phone_number, customer.Address)
				if err != nil {
					tx.Rollback()
					log.Println("Failed to insert customers: ", err)
				}
			}
		}
		for _, addres := range branch.Address {
			_, err = tx.Exec("INSERT INTO address(branch_id,street,city) VALUES($1,$2,$3)", branchId, addres.Street, addres.City)
			if err != nil {
				tx.Rollback()
				log.Println("Failed to insert address: ", err)
			}
		}

	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Println("Failed to commit: ", err)
	}
	return OfficeId, branchId, carId

}
