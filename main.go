package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Employee struct {
	ID          int     `db:"id" json:"id"`
	FirstName   string  `db:"first_name" json:"first_name"`
	LastName    string  `db:"last_name" json:"last_name"`
	Email       string  `db:"email" json:"email"`
	PhoneNumber string  `db:"phone_number" json:"phone_number"`
	HireDate    string  `db:"hire_date" json:"hire_date"`
	JobTitle    string  `db:"job_title" json:"job_title"`
	Salary      float64 `db:"salary" json:"salary"`
}

var db *sqlx.DB

func initDB() {
	var err error
	connStr := "user=postgres password=1234 dbname=company sslmode=disable"
	db, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		panic(err)
	}
}

func main() {
	initDB()

	r := gin.Default()

	r.GET("/employees", getEmployees)
	r.GET("/employees/:id", getEmployeeByID)
	r.POST("/employees", createEmployee)
	r.PUT("/employees/:id", updateEmployee)
	r.DELETE("/employees/:id", deleteEmployee)

	r.Run(":8080")
}

func getEmployees(c *gin.Context) {
	var employees []Employee // Slices of employees
	err := db.Select(&employees, "SELECT * FROM employees")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employees)
}

func getEmployeeByID(c *gin.Context) {
	id := c.Param("id")
	var employee Employee
	err := db.Get(&employee, "SELECT * FROM employees WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}
	c.JSON(http.StatusOK, employee)
}

func createEmployee(c *gin.Context) {
	var employee Employee

	// Bind the incoming JSON payload from the request body to the employee struct.
	// It parses the JSON and maps it to the struct fields.
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Inserts a new employee record into the employees table using named parameters that correspond to the struct fields.
	_, err := db.NamedExec(`INSERT INTO employees (first_name, last_name, email, phone_number, hire_date, job_title, salary) VALUES (:first_name, :last_name, :email, :phone_number, :hire_date, :job_title, :salary)`, &employee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, employee)
}

func updateEmployee(c *gin.Context) {
	id := c.Param("id")
	numID, _ := strconv.Atoi(id)
	var employee Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	employee.ID = numID
	_, err := db.NamedExec(`UPDATE employees SET first_name=:first_name, last_name=:last_name, email=:email, phone_number=:phone_number, hire_date=:hire_date, job_title=:job_title, salary=:salary WHERE id=:id`, &employee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employee)
}

func deleteEmployee(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM employees WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Employee deleted"})
}
