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
	PostgresDatabase = "rend_car"
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
	of_id, br_id, car_id := InsertDatabase(db)
	UpdateInfo(of_id, br_id, car_id, db)

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

func UpdateInfo(officeId, branchId, carId int, db *sql.DB) {
	car := &Offices{
		Name: "Rend_car",
		Branch: []Branches{
			{
				officeId: 2,
				Name:     "Shahar",
				Car: []Car{
					{
						BranchId: 2,
						Name:     "Gentra",
						Color:    "dark black",
						CostDay:  9500.22,
						Amount:   200,
						Customer: []Customer{
							{
								CarId:        2,
								Name:         "Abdukarim",
								Age:          28,
								Phone_number: "+99893 545 78 87",
								Address:      "Toshkent, Yunusobod",
							},
						},
					},
				},
				Address: []Address{
					{
						BranchId: 2,
						Street:   "Qodiriy street",
						City:     "Toshkent",
					},
				},
			},
		},
	}

	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Failed to begin transaction", err)
	}

	_, err = tx.Exec("UPDATE offices SET name=$1 WHERE id=$2", car.Name, officeId)
	if err != nil {
		tx.Rollback()
		fmt.Println("Failed to update offices", err)
	}

	for _, branch := range car.Branch {
		_, err := tx.Exec("UPDATE branches SET name=$1 WHERE office_id=$2", branch.Name, officeId)
		if err != nil {
			tx.Rollback()
			fmt.Println("Failed to Update branches: ", err)
		}

		for _, car := range branch.Car {
			_, err := tx.Exec("UPDATE cars SET name=$1,color=$2,cost_day=$3,amount=$4 WHERE branch_id=$5", car.Name, car.Color, car.CostDay, car.Amount, branchId)
			if err != nil {
				tx.Rollback()
				fmt.Println("Failed to update cars: ", err)
			}

			for _, customer := range car.Customer {
				_, err := tx.Exec("UPDATE customer SET name=$1,age=$2,phone_number=$3,address=$4 WHERE car_id=$5", customer.Name, customer.Age, customer.Phone_number, customer.Address, carId)
				if err != nil {
					tx.Rollback()
					fmt.Println("Failed to udate customer: ", err)
				}
			}
		}

		for _, address := range branch.Address {
			_, err := tx.Exec("UPDATE address SET street=$1,city=$2 WHERE branch_id=$3", address.Street, address.City, branchId)
			if err != nil {
				tx.Rollback()
				fmt.Println("Failed to update address:", err)
			}
		}
	}

}
