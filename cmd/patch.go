/*
Copyright © 2022 Aria Taylor <ari@aricodes.net>

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
	"os"
	"path"

	"gips/ips"

	"github.com/aricodes-oss/std"
	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
)

const IPS_EXTENSION = ".ips"

var log = std.Logger

func getFileSet(args []string) (rom, patch string) {
	for _, fileName := range args {
		if path.Ext(fileName) == IPS_EXTENSION {
			patch = fileName
		} else {
			rom = fileName
		}
	}

	return
}

var patchCmd = &cobra.Command{
	Use:   "patch [inputFile.bin] [ipsPatch.ips]",
	Short: "Takes a ROM and an IPS patch and applies it",
	Long: `Applies an IPS patch and produces a patched file.
The arguments can be given in any order, so long as the IPS file has a .ips extension.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		romFile, patchFile := getFileSet(args)
		outFile, _ := cmd.Flags().GetString("output")

		if outFile == "" {
			outFile = "patched." + romFile
		}

		log.Debugf("Loading %v as patch file", aurora.Bold(patchFile))
		patcher := ips.LoadPatchFile(patchFile)
		err := patcher.LoadRecords()
		if err != nil {
			log.Fatalf("Failed loading records with error: %v", err)
		}

		log.Debugf("Loading %v as ROM file", aurora.Bold(romFile))
		data, err := os.ReadFile(romFile)
		if err != nil {
			log.Fatalf("Failed loading ROM file with error: %v")
		}

		log.Infof("Patching now!")
		data = patcher.Patch(data)
		os.WriteFile(outFile, data, 0755)

		log.Infof("Finished writing output file %v, enjoy!", aurora.Bold(aurora.BrightGreen(outFile)))
	},
}

func init() {
	patchCmd.PersistentFlags().StringP("output", "o", "", "output file name")
	rootCmd.AddCommand(patchCmd)
}
