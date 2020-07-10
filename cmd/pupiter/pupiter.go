package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/evanboyle/pupiter/repl"
	"github.com/spf13/cobra"
)

func NewPupiterCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pupiter [session]",
		Short: "Pupiter is an interpreter for Pulumi.",
		Long:  `Pupiter is an interpreter for Pulumi.`,
		Run: func(cmd *cobra.Command, args []string) {

			// TODO use a randomly generated name that we print out for a "session ID"
			// Add a parameter to look up a session id ~/.pulumi/pupiter/<sessionID>

			user, err := user.Current()
			if err != nil {
				panic(err)
			}
			fmt.Printf("Hello %s! Welcome to pupiter!!\n", user.Username)
			fmt.Printf("Where engineers come...\n")
			fmt.Printf("...to get more educated about the cloud.\n")
			fmt.Printf("\n")
			fmt.Println("Define a pulumi resource...")
			fmt.Println("-----------------------RULES-----------------------")
			fmt.Println("Define pulumi resources using standard javascript.")
			fmt.Println("Evaluation occurs after ';'")
			fmt.Println("Pulumi resource declarations can refer to previous resource declarations.")
			fmt.Println("-----------------------ENJOY-----------------------")

			rand.Seed(time.Now().UTC().UnixNano())

			session := petname.Generate(1 /*words*/, "" /*seperator*/)
			if len(args) > 0 {
				session = args[0]
			}
			repl.Start(os.Stdin, os.Stdout, session)
		},
	}
}
