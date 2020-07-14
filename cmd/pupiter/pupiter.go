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
	cmd := &cobra.Command{
		Use:   "pupiter [session]",
		Short: "Pupiter is an interpreter for Pulumi.",
		Long:  `Pupiter is an interpreter for Pulumi.`,
		Run: func(cmd *cobra.Command, args []string) {
			user, err := user.Current()
			if err != nil {
				panic(err)
			}
			fmt.Printf("Hello %s! Welcome to Pupiter!!\n", user.Username)
			fmt.Println("Define a pulumi resource...")
			fmt.Println("-----------------------RULES-----------------------")
			fmt.Println(`1. Add resources through individual var declarations (';' triggers evaluation):
			var x = new aws.s3.Bucket("mybucket");
			
			`)
			fmt.Println(`2. Read back variables:
			x;
			{...output...}

			`)
			fmt.Println(`3. List all variables in the current session:
			ls();
			x
			y
			z
			...
			
			`)
			fmt.Println(`4. Capture your complete program with eject:
			eject();

			var x = new aws.s3.Bucket("mybucket");
			var y = ...;
			var z = ...;
			...
			
			`)

			fmt.Println(`4. Scope traversal is allowed:
			var y = new aws.s3.BucketObject("y", { bucket: x.bucket });
			
			`)
			fmt.Println("Evaluation occurs after ';'")
			fmt.Println("Only simple statements in the form `var x = ...;` are supported")

			fmt.Println("-----------------------ENJOY-----------------------")

			rand.Seed(time.Now().UTC().UnixNano())

			session := petname.Generate(1 /*words*/, "" /*seperator*/)
			if len(args) > 0 {
				session = args[0]
			}
			repl.Start(os.Stdin, os.Stdout, session)
		},
	}

	cmd.AddCommand(NewNotebookCmd())

	return cmd
}
