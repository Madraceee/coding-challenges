package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	FieldColumnRange     string
	CharacterColumnRange string
	Delimiter            string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cut OPTION... [FILE]...",
	Short: "cut - remove sections from each line of files.",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		info, err := os.Stdin.Stat()
		if err != nil {
			fmt.Printf("Err: %s\n", err)
			os.Exit(1)
		}

		output := ""
		var reader *bufio.Reader
		if (info.Mode() & os.ModeCharDevice) == 0 {
			reader = bufio.NewReader(os.Stdin)
			res, err := GetCutData(reader, Delimiter)
			if err != nil {
				return err
			}

			for _, s := range *res {
				output += s + "\n"
			}

		} else {
			for _, file := range args {
				file, err := os.Open(file)
				if err != nil {
					return err
				}

				reader := bufio.NewReader(file)
				res, err := GetCutData(reader, Delimiter)
				if err != nil {
					return err
				}

				for _, s := range *res {
					output += s + "\n"
				}
			}
		}

		fmt.Fprintf(os.Stdout, "%s", output)

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&FieldColumnRange, "fields", "f", "", "Print the specified column(1-indexed)")
	rootCmd.Flags().StringVarP(&Delimiter, "delimiter", "d", "\t", "Delimiter required to perform split each row.")

	rootCmd.Flags().StringVarP(&CharacterColumnRange, "characters", "c", "", "Print the specified characters(1-indexed)")

	rootCmd.MarkFlagsMutuallyExclusive("characters", "fields")
	rootCmd.MarkFlagsMutuallyExclusive("characters", "delimiter")
}
