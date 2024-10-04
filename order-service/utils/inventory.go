package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// UpdateProductInventory updates the product inventory by making a PUT request to the product service
func UpdateProductInventory(productName string, quantity int) error {
	url := fmt.Sprintf("http://localhost:8082/product/%s", productName)

	// Create the request body
	requestBody, err := json.Marshal(map[string]int{"quantity": quantity})
	if err != nil {
		return fmt.Errorf("error marshaling request body: %v", err)
	}

	// Create a new PUT request
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("error creating PUT request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending PUT request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	return nil
}
