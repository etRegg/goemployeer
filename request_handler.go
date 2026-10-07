package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 15 * time.Second,
}

// fetchData obtiene los empleados directamente del endpoint indicado,
// los convierte desde JSON, los procesa y finalmente los guarda en cache.
func fetchData(cache *CacheManager, endpoint string) error {
	resp, err := httpClient.Get(endpoint)
	if err != nil {
		return fmt.Errorf("HTTP GET: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf(
			"HTTP status %d: %s",
			resp.StatusCode,
			string(body),
		)
	}

	var employees []Employee

	if err := json.NewDecoder(resp.Body).Decode(&employees); err != nil {
		return fmt.Errorf("decode JSON: %w", err)
	}

	fmt.Printf("Endpoint %s returned %d employees\n", endpoint, len(employees))

	processedEmployees, err := processEmployees(employees)
	if err != nil {
		return fmt.Errorf("process employees: %w", err)
	}

	cache.Set(processedEmployees...)

	return nil
}
