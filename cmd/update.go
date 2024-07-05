/*
Copyright © 2024 Mikael Luca Pinheiro Costa <mikaellpc4@gmail.com>

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
	"context"
	"fmt"
	"time"

	"github.com/GoCEP/api/cep/repository"
	"github.com/GoCEP/api/cep/repository/implementations"
	"github.com/GoCEP/api/cep/services"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update databases",
	Long:  `Updates any existent repository with the newest data on OpenCEP`,
	Run: func(cmd *cobra.Command, args []string) {
    start := time.Now()
		// sqliteRepo := implementations.NewSqliteCepRepo()
		firebirdRepo := implementations.NewFirebirdCepRepo()

		repos := []repository.CepRepositary{firebirdRepo}
		service := services.NewCepService(repos)

		context := context.Background()
		err := service.UpdateData(context)
		if err != nil {
			fmt.Printf("\na error ocurred while updating the data: %s\n", err)
		}
    finishedIn := time.Since(start)
    fmt.Printf("finished in: %s", finishedIn.String())
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
