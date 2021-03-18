/*-----------------------------------------------------------------------------
# Name:        Check-VAT
# Purpose:     Verify VAT status in mf.gov.pl database
#
# Author:      Rafal Wilk <rw@pcboot.pl>
#
# Created:     27-09-2020
# Modified:    18-03-2021
# Copyright:   (c) PcBoot 2021
# License:     BSD-new
-----------------------------------------------------------------------------*/

package main

import (
	"flag"
	"fmt"
	"os"

	vatpl "github.com/rwilk/go-vat-pl"
)

func main() {
	fmt.Println("Check-VAT command line VAT status verifier")
	fmt.Println("All rights reserved. (c) PcBoot 2021")
	fmt.Println()

	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Printf("Usage:\n\t%s <NIP> [<NIP> ...]\n\n", os.Args[0])
		os.Exit(2)
	}

	if len(flag.Args()) == 1 {
		// single NIP verification
		nip := flag.Arg(0)
		status, err := vatpl.VerifyByNIPRetry(nip)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Verification status of %s - %v\n", nip, status)
	} else {
		// bulk NIP verification
		nips := flag.Args()
		statuses, err := vatpl.VerifyByNIPBulkRetry(nips)
		if err != nil {
			panic(err)
		}

		for nip, status := range statuses {
			fmt.Printf("Verification status of %s - %v\n", nip, status)
		}
	}

}
