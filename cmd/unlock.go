/*
Copyright © 2022 sam <contact@justsam.io>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log"

	cloudlock "github.com/jsam/cloudlock/pkg"
	"github.com/spf13/cobra"
)

// unlockCmd represents the unlock command
var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Release existing lock file.",
	Long: `Release existing lock file.
	
	Releases a lock file created by the lock command.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("requires a lock file name")
		}
		if len(args) > 1 {
			return fmt.Errorf("only one lock file is allowed")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		gcs := cloudlock.NewGCS(cloudlock.CloudlockBucket)
		_, err := gcs.Unlock(args[0])
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Lock is released")
	},
}

func init() {
	rootCmd.AddCommand(unlockCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// unlockCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// unlockCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
