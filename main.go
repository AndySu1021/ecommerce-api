package main

import (
	"ecommerce-api/cmd"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"runtime"
)

var rootCmd = &cobra.Command{Use: "app"}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)

	rootCmd.AddCommand(cmd.NewMigrateCommand())
	rootCmd.AddCommand(cmd.NewServerCommand())
	rootCmd.AddCommand(cmd.NewAdminCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
