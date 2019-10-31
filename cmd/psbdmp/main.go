package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/traviscampbell/psbdmp"
)

func main() {
	dlFlag := flag.String("dl", "", "ID of dump to be downloaded")
	
	domainFlag := flag.String("domain", "", "domain to search dumps for")
	emailFlag := flag.String("email", "", "email to search the dumps for")
	sinceFlag := flag.Int("since", 0, "number of days ago to start getting all dumps from")
	searchFlag := flag.String("search", "", "keyword to search the dumps for")

	fetchDumpFlag := flag.Bool("fetch", false, "fetch each dump found (instead of just getting a list of IDs)")
	outdirFlag := flag.String("out", "/tmp", "directory where the fetched dumps should be written to disk")
	
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Printf("[!] Gotta choose and option...\nEither dl a specific dump or search them for something.\n\n")
		flag.Usage()
		os.Exit(1)
	}

	pd := psbdmp.NewDumpClient()

	if *dlFlag != "" {
		content, err := pd.GetDumpContent(*dlFlag)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(content)
		os.Exit(0)
	}

	var (
		err error
		dumps []psbdmp.Dump
	)

	if *domainFlag != "" {
		dumps, err = pd.SearchByDomain(*domainFlag)
	}

	if *emailFlag != "" {
		dumps, err = pd.SearchByEmail(*emailFlag)
	}

	if *searchFlag != "" {
		dumps, err = pd.Search(*searchFlag)
	}

	if *sinceFlag != 0 {
		daysAgo := *sinceFlag
		if daysAgo < 0 {
			daysAgo = daysAgo * -1
		}
		dumps, err = pd.GetByDate(time.Now(), time.Now().AddDate(0, 0, -daysAgo))
	}

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("[+]", len(dumps), "dumps found!")

	if *fetchDumpFlag {
		fmt.Println("[+] fetching the dumps and writing to disk...")
		for _, d := range dumps {
			content, err := pd.GetDumpContent(d.ID)
			if err != nil {
				log.Println(d.ID, err)
				continue
			}
			
			fname := filepath.Join(*outdirFlag, d.ID)
			if err := ioutil.WriteFile(fname, []byte(content), 0666); err != nil {
				log.Println(err)
			}
		}
	} else {
		fmt.Println("[+] dump IDs matching the query:")
		for _, d := range dumps {
			fmt.Println(d.ID)
		}
	}
}
