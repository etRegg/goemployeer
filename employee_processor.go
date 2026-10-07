package main

import "math/rand"

// processEmployees procesa los empleados recibidos desde las APIs.
func processEmployees(employees []Employee) ([]Employee, error) {
	for i := range employees {
		rate := rand.Float64() * 100
		employees[i].Salary = employees[i].Hours * rate
	}

	return employees, nil
}
