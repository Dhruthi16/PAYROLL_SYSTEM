package models

type Payroll struct {
	ID      string  `bson:"_id,omitempty" json:"id"`
	EmpID   string  `bson:"emp_id" json:"emp_id"`
	EmpName string  `bson:"emp_name" json:"emp_name"`
	Salary  float64 `bson:"salary" json:"salary"`
	Month   string  `bson:"month" json:"month"`
}
