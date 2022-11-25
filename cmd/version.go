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

	"github.com/spf13/cobra"
)

var (
	GitCommit    string
	GitBranch    string
	GitTag       string
	GitDirty     string
	BuildDate    string
	BuildVersion string
	BuildHash    string
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Check the tool version",
	Long:  `Check the tool version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("GitCommit: ", GitCommit)
		fmt.Println("GitBranch: ", GitBranch)
		fmt.Println("GitTag: ", GitTag)
		fmt.Println("GitDirty: ", GitDirty)
		fmt.Println("BuildDate: ", BuildDate)
		fmt.Println("BuildVersion: ", BuildVersion)
		fmt.Println("BuildHash: ", BuildHash)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
