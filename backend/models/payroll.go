package models

type Payroll struct {
	ID      string  `bson:"_id,omitempty"`
	EmpID   string  `bson:"emp_id"`
	EmpName string  `bson:"emp_name"`
	Salary  float64 `bson:"salary"`
	Month   string  `bson:"month"`
}

