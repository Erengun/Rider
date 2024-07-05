/*
Copyright © 2024 NAME HERE hi@erengun.dev
*/

package cmd

import (
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type PlainFormatter struct {
}

func (f *PlainFormatter) Format(entry *log.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("%s\n", entry.Message)), nil
}
func toggleDebug(*cobra.Command, []string) {
	if Debug {
		log.Info("Debug logs enabled")
		log.SetLevel(log.DebugLevel)
		log.SetFormatter(&log.TextFormatter{})
	} else {
		plainFormatter := new(PlainFormatter)
		log.SetFormatter(plainFormatter)
	}
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init command is used to start the process of changing app name, package name, and other configurations.",
	Long: `
init command is used to start the process of changing app name, package name, and other configurations.
It uses the information provided by the user to update the configurations and rename the files accordingly.
Write 'template init' in the terminal to start the process.
	`,
	PreRun: toggleDebug,
	Run: func(cmd *cobra.Command, args []string) {
		// Clear the terminal
		_, stdErr := os.Stdout.WriteString("\x1b[3;J\x1b[H\x1b[2J")
		if stdErr != nil {
			return
		}

		// Check if the project is a Flutter project
		if _, err := os.Stat("pubspec.yaml"); os.IsNotExist(err) {
			// If not a Flutter project, show error message and exit
			// log.Debug("Not a Flutter project")
			log.Fatal("Not a Flutter project. Please run this command in a Flutter project directory.")
			return
		}
		log.Debug("Flutter project detected. Running flutter pub get to get the dependencies.")

		log.Info("Welcome to the template CLI. This command will help you update the app name, package name, app icon, and splash screen.")

		stopChan := make(chan struct{})
		go LoadingSpinner(stopChan, "Getting Dependencies...")
		// If Flutter project, run flutter pub get to get the dependencies
		log.Debug("Running flutter pub get to get the dependencies.")
		_, err := exec.Command("flutter", "pub", "get").Output()
		if err != nil {
			log.Fatal("Failed to run flutter pub get:", err)
		}
		// Run flutter pub add rename_app
		log.Debug("Adding rename_app package to the project.")
		_, err = exec.Command("flutter", "pub", "add", "rename_app").Output()
		if err != nil {
			log.Fatal("Failed to add rename_app package:", err)
		}

		log.Debug("Adding change_app_package_name package to the project.")
		_, err = exec.Command("flutter", "pub", "add", "change_app_package_name", "--dev").Output()
		if err != nil {
			log.Fatal("Failed to add change_app_package_name package:", err)
		}

		close(stopChan)

	Actions:
		action := promptGetSelect(
			promptContent{
				label:   "What do you want to update?",
				choices: []string{"App Name", "Package Name", "App Icon", "Done"},
			},
		)
		switch action {
		case "App Name":
			updateAppName()
			goto Actions
		case "Package Name":
			updatePackageName()
			goto Actions
		case "App Icon":
			updateAppIcon()
			goto Actions
		case "Done":
			break
		}

		cleanUp()

	},
}

func cleanUp() {
	stopChan := make(chan struct{})
	go LoadingSpinner(stopChan, "Cleaning up...")
	// Run flutter pub get
	log.Debug("Removing change_app_package_name package from the project.")
	_, err := exec.Command("flutter", "pub", "remove", "change_app_package_name").Output()
	if err != nil {
		log.Debug("Failed to remove change_app_package_name package:", err)
		return
	}
	log.Debug("Removing rename_app package from the project.")
	_, err = exec.Command("flutter", "pub", "remove", "rename_app").Output()
	if err != nil {
		log.Debug("Failed to remove rename_app package:", err)
		return
	}
	close(stopChan)
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
