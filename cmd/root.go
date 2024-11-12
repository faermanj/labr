/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// runShellCommand runs the shell script and logs its output in real-time
func runShellCommand(filePath string) error {
	// Create the command to run the shell script
	cmd := exec.Command("bash", filePath)

	// Create a pipe to capture the output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	// Create a pipe to capture the stderr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// Start the command execution
	if err := cmd.Start(); err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// Use a scanner to read the output line by line
	stdoutScanner := bufio.NewScanner(stdout)
	stderrScanner := bufio.NewScanner(stderr)

	go func() {
		defer wg.Done()
		for stdoutScanner.Scan() {
			log.Info(stdoutScanner.Text())
		}
		if err := stdoutScanner.Err(); err != nil {
			log.Error(err)
		}
	}()

	go func() {
		defer wg.Done()
		for stderrScanner.Scan() {
			log.Error(stderrScanner.Text())
		}
		if err := stderrScanner.Err(); err != nil {
			log.Error(err)
		}
	}()

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		return err
	}

	// Wait for the goroutines to finish
	wg.Wait()

	// Flush any remaining output
	if stdoutScanner.Scan() {
		log.Info(stdoutScanner.Text())
	}
	if stderrScanner.Scan() {
		log.Error(stderrScanner.Text())
	}

	return nil
}

// FindSHFiles takes a directory path and returns a sorted slice of files matching the criteria
func BuildPlan(dir string) ([]string, error) {
	var result []string

	// Regular expression to match files that are .sh files and start with 1 to 5 digits
	re := regexp.MustCompile(`^\d{1,5}_.*\.sh$`)

	// Walk through the directory recursively
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file matches the regex pattern and ends with .sh
		if re.MatchString(info.Name()) {
			result = append(result, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort the result alphabetically
	sort.Strings(result)

	return result, nil
}

func getCWD() string {
	dir, err := os.Getwd()
	if err != nil {
		log.WithFields(log.Fields{
			"dir": dir,
		}).Fatal("Error getting current working directory")
	}
	return dir
}

func labr(dir string) {
	cwd := getCWD()
	plan, err := BuildPlan(dir)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	log.WithFields(log.Fields{
		"dir":  dir,
		"cwd":  cwd,
		"plan": plan,
	}).Info("Running labr")
	// Iterate over the list of files
	for _, file := range plan {
		// Log the start of the command execution
		log.WithFields(log.Fields{
			"file": file,
		}).Info("Executing shell script")

		// Run the shell command and capture output in real-time
		if err := runShellCommand(file); err != nil {
			log.WithFields(log.Fields{
				"file": file,
			}).Error("Error executing script")
		} else {
			log.WithFields(log.Fields{
				"file": file,
			}).Debug("Script executed successfully")
		}
	}

}

func run(cmd *cobra.Command, args []string) {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}
	// check if the directory exists and is a directory
	dir_valid := false
	if _, err := os.Stat(dir); err == nil {
		dir_valid = true
	}
	if !dir_valid {
		log.WithFields(log.Fields{
			"dir": dir,
		}).Fatal("Directory not found")
		return
	}
	labr(dir)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "labr",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: run,
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
	// log.SetFormatter(&log.JSONFormatter{}) // JSON output
	// OR for a more readable format during development
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})
	log.SetLevel(log.InfoLevel)
	/*
		log.Trace("This is a trace message.")
		log.Debug("This is an debug message.")
		log.Info("This is an informational message.")
		log.Warn("This is a warning.")
		log.Error("This is an error message.")
		log.WithFields(log.Fields{
			"user":   "johndoe",
			"action": "login",
		}).Info("User action")
	*/
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.labr.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
