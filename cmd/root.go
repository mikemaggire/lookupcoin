/*
Copyright © 2022 Mike Maggire <mike.maggire@gmail.com>

*/
package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	coingecko "github.com/superoo7/go-gecko/v3"
	util "github.com/mikemaggire/gotermutil"

)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lookupcoin",
	Short: "Lookup coin in the CoinGecko database",
	Long: `Lookup coin by requesting the CoinGecko API and looking in the symbol, gcid and name fields.
	Returns one line per finding.

PARAMETER: the coin your're looking for.`,
	Run: func(cmd *cobra.Command, args []string) { 
		if len(args) == 0 || len(args[0]) == 0 {
			fmt.Println("Missing search parameter")
			return
		}
		search := args[0]

		fExactMatch, _ := cmd.Flags().GetBool("exact-match")

		lst, err := lookupCoin(search, fExactMatch) 
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%d coins match '%s'\n", len(lst), search) 
		
		sort.Strings(lst)
		for _, found := range lst {
			fmt.Println(found) 
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	rootCmd.Flags().BoolP("exact-match", "", false, "TRUE to lookup the exact string, FALSE to lookup the string everywhere in fields.")
}

// return a slice of founded coins, with formatted strings
func lookupCoin(search string, fExactMatch bool) ([]string, error) {

	// load coinGecko database
	httpClient := &http.Client{	Timeout: time.Second * 15 }
	cg := coingecko.NewClient(httpClient)
	gclist, err := cg.CoinsList()
	if err != nil {
		return nil, err
	}
	fmt.Printf("Lookup among %d coins in the CoinGecko database.\n", len(*gclist))

	// prepare for the lookup
	s := make([]string, 0)
	lsearch := strings.ToLower(search)

	// lookup
	for _, gc := range *gclist {
		fMatch := false
		if fExactMatch {
			if 	(strings.ToLower(gc.Symbol) == lsearch) || (strings.ToLower(gc.ID) == lsearch) {
				fMatch = true
			}
		} else { // !fExactMatch
			gcdata := gc.Symbol + "•" + gc.ID + "•" + gc.Name
			if strings.Contains(strings.ToLower(gcdata), lsearch) {
				fMatch = true
			}
		}
		// Format the found coin
		if fMatch {
			found := util.White(util.Tab(gc.Symbol, 15, false) + ": ")
			found += gc.ID + ", "
			found += gc.Name
			s = append(s, found)
		}
	}

	// That's it
	if len(s) == 0 {
		return nil, errors.New("no coin matches '"+ search + "'")
	}
	return s, nil
}

