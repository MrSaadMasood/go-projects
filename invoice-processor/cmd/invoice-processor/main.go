package main

import (
	"flag"
	"fmt"
	"invoice/internal/invoice"
	"invoice/internal/pkg/validate"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic Occured:", err)
		}
	}()

	dryRunPtr := flag.Bool("dry-run", false, "Skip output files generation")
	fromPtr := flag.String("from", "", "Filter invoices from provided dates")
	toPtr := flag.String("to", "", "Filter invoices to provided dates")
	pathPtr := flag.String("path", "data.json", "path to the invoices json files")

	flag.Parse()

	flags := validate.FlagsOptions{DryRun: *dryRunPtr, From: *fromPtr, To: *toPtr, Path: *pathPtr}

	vflags, err := validate.Flags(flags)
	if err != nil || vflags == nil {
		fmt.Println("Flags Validation Failed %w", err)
	}
	from := vflags.From
	to := vflags.To
	if (from == nil && to != nil) || (from != nil && to == nil) {
		fmt.Println("Date Filtering Required Both From And To To Exist Or Neither")
		return
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool)
	errChan := make(chan error)

	go func() {
		terminationSignal := <-sig
		fmt.Println("Received Termination Signal:", terminationSignal, "Terminating Processing")
		done <- true
	}()

	go invoice.Process(*vflags, done, errChan)

	select {
	case <-done:
		return
	case err := <-errChan:
		fmt.Println("Error Occured While Processing Invoices %w", err)
	}

}
