package invoice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"invoice/internal/pkg/validate"
	"os"
	"path"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
)

type invoice struct {
	ClientId    string `json:"client_id" validate:"required"`
	Date        string `json:"date" validate:"required"`
	HoursWorked int    `json:"hours_worked" validate:"required"`
	RatePerHour int    `json:"rate_per_hour" validate:"required"`
	Project     string `json:"project"`
	Description string `json:"description"`
}

type Compiled struct {
	Total_hours  int    `json:"total_hours"`
	Total_amount int    `json:"total_amount"`
	Client_id    string `json:"client_id"`
}

func Process(flags validate.ValidatedFlags, done chan bool, errChan chan error) {

	handleError := func(err error) {
		errChan <- err
	}

	invoices, err := unmarshall(flags.Path)
	fmt.Print(err)
	if err != nil {
		handleError(err)
		return
	}

	cm, err := clientMap(invoices, flags)
	if err != nil {
		handleError(err)
		return
	}

	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()

	result := make(chan Compiled, len(cm))
	doneChannel := make(chan bool)

	for k, v := range cm {
		wg.Add(1)
		go func() {
			defer wg.Done()
			totalSum := 0
			totalHours := 0
			for _, invoice := range v {
				totalHours += invoice.HoursWorked
				totalSum += invoice.HoursWorked * invoice.RatePerHour
			}
			result <- Compiled{Client_id: k, Total_hours: totalHours, Total_amount: totalSum}

		}()
	}

	go func() {
		wg.Wait()
		defer close(result)
		doneChannel <- true
	}()

	select {
	case <-doneChannel:
		err := results(result, flags)
		if err != nil {
			handleError(err)
			return
		}
		done <- true
	case <-ctx.Done():
		close(result)
		fmt.Println("Processing timedout")
		done <- true
	}

}

func clientMap(invoices []invoice, flags validate.ValidatedFlags) (map[string][]invoice, error) {
	clientInvoicesMap := map[string][]invoice{}
	applyDateFilter := flags.From != nil && flags.To != nil

	for _, v := range invoices {
		if applyDateFilter == true {
			date, err := validate.Date(v.Date)
			if err != nil {
				return nil, fmt.Errorf("Incorrect invoice date %w", err)
			}

			wihtinRange := date.After(*flags.From) && date.Before(*flags.To)
			if wihtinRange == false {
				continue
			}
		}
		clientInvoicesMap[v.ClientId] = append(clientInvoicesMap[v.ClientId], v)
	}
	return clientInvoicesMap, nil
}

func results(r <-chan Compiled, flags validate.ValidatedFlags) error {
	compiled := make([]Compiled, 0)

	for v := range r {
		compiled = append(compiled, v)
	}

	if flags.DryRun == true {
		fmt.Println(compiled)
		fmt.Println("Dry Run Completed")
		return nil
	}

	err := store(compiled)
	if err != nil {
		return fmt.Errorf("Error occured while compiling results %w", err)
	}
	return nil
}

func store(data []Compiled) error {
	j, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Error occured while marshelling %w", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Failed to get working directory %w", err)
	}

	finalFile := "invoices.json"
	storagePath := path.Join(wd, finalFile)
	_, err = os.Stat(storagePath)
	fileExists := !errors.Is(err, os.ErrNotExist)
	if fileExists == true {
		err := os.Remove(storagePath)
		if err != nil {
			return fmt.Errorf("Failed to removed Already Existing File %w", err)
		}
	}
	err = os.WriteFile(storagePath, j, 0644)
	if err != nil {
		return fmt.Errorf("error occured while storing results %w", err)
	}

	fmt.Println("Invoices Processed Successfully")
	return nil

}

func unmarshall(path string) ([]invoice, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed reading the invoices file %w", err)
	}

	var invoices []invoice

	err = json.Unmarshal(data, &invoices)
	if err != nil {
		return nil, fmt.Errorf("Error occured while unmarshelling invoices %w", err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Var(invoices, "required,min=1,dive"); err != nil {

		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			return nil, fmt.Errorf("InvalidValidationError %w", err)
		}

		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return nil, fmt.Errorf("ValidationErrors %w", err)
		}

		return nil, fmt.Errorf("Error occured while validatng %w", err)

	}

	return invoices, nil

}
