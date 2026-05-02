/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/madraceee/coding-challenges/huffman-encoder-decoder/internal/cli/compress"
	"github.com/madraceee/coding-challenges/huffman-encoder-decoder/internal/cli/decompress"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "huffman-encoder-decoder command [INPUT FILE] [OUTPUT FILE]",
	Short: "Compress and Decompress files",
	Long: `Compress and Decompress files using huffman encoding decoding 
	technique.`,
}

var encodeCmd = &cobra.Command{
	Use:  "compress",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("Wrong input\nhuffman-encoder-decoder command [INPUT FILE] [OUTPUT FILE]")
		}
		return compress.Encode(args[0], args[1])
	},
}

var decodeCmd = &cobra.Command{
	Use:  "decompress",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("Wrong input\nhuffman-encoder-decoder command [INPUT FILE] [OUTPUT FILE]")
		}
		return decompress.Decode(args[0], args[1])
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(encodeCmd, decodeCmd)
}
