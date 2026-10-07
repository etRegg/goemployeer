package main

import "sync"

type CacheManager struct {
	mu        sync.RWMutex
	employees map[int]Employee
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		employees: make(map[int]Employee),
	}
}

func (cache *CacheManager) Get(id int) (Employee, bool) {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	employee, found := cache.employees[id]
	return employee, found
}

func (cache *CacheManager) Set(employees ...Employee) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	for _, employee := range employees {
		cache.employees[employee.ID] = employee
	}
}

func (cache *CacheManager) GetProcessedEmployees() []Employee {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	processedEmployees := make([]Employee, 0, len(cache.employees))

	for _, employee := range cache.employees {
		processedEmployees = append(processedEmployees, employee)
	}

	return processedEmployees
}
