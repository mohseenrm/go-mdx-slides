package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-ini/ini"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

/* FLAGS */
var cfgFile string
var devPort string
var email string
var phoneNumber string
var redisPort string
var userName string

/* Go Template strings */
var t *template.Template

const templates = `
{{ define "celery.url" }}redis://hostname:6379/{{.}}{{end}}
{{ define "cache_dir" }}/var/cache/xxx-{{.}}/data/{{end}}
{{ define "upload_dir" }}/var/spool/xxx-{{.}}/uploads/{{end}}
`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "create-ini",
	Short: "Test Description",
	Long: `[DEBUG] A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Print: " + strings.Join(args, " "))

		/* Get Flag Values */
		devPort, _ := cmd.Flags().GetString("devPort")
		email, _ := cmd.Flags().GetString("email")
		phoneNumber, _ := cmd.Flags().GetString("phoneNumber")
		redisPort, _ := cmd.Flags().GetString("redisPort")
		userName, _ := cmd.Flags().GetString("userName")

		outputFileName := fmt.Sprintf("%s.ini", userName)

		/* Read Initial Config */
		currentDirectory, _ := filepath.Abs(filepath.Dir(os.Args[1]))
		templatePath := path.Join(currentDirectory, "template", "new_dev.ini")
		outputPath := path.Join(currentDirectory, "output", outputFileName)

		cfg, err := ini.Load(templatePath)
		if err != nil {
			fmt.Printf("[ERROR] reading %v template file\n\n%v", templatePath, err)
			os.Exit(1)
		}

		/* Sections */
		appMain := cfg.Section("app:main")
		defaultSection := cfg.Section("DEFAULT")
		serverMain := cfg.Section("server:main")

		/* Set new values */
		appMain.Key("celery.prefix").SetValue(userName)
		appMain.Key("celery.db").SetValue(redisPort)
		appMain.Key("twilio.debug_phone").SetValue(phoneNumber)
		defaultSection.Key("email_to").SetValue(email)
		defaultSection.Key("mail.debug_email").SetValue(email)
		serverMain.Key("port").SetValue(devPort)

		/* Go Template Objects */
		var tplOutputBuf bytes.Buffer
		templateObj, err := template.New("newDevTemplate").Parse(templates)

		if err := templateObj.ExecuteTemplate(
			&tplOutputBuf,
			"cache_dir",
			userName); err != nil {
			panic(err)
		}
		appMain.Key("cache_dir").SetValue(tplOutputBuf.String())
		tplOutputBuf.Reset()

		if err := templateObj.ExecuteTemplate(
			&tplOutputBuf,
			"upload_dir",
			userName); err != nil {
			panic(err)
		}
		appMain.Key("upload_dir").SetValue(tplOutputBuf.String())
		tplOutputBuf.Reset()

		if err := templateObj.ExecuteTemplate(
			&tplOutputBuf,
			"celery.url",
			redisPort); err != nil {
			panic(err)
		}
		appMain.Key("celery.url").SetValue(tplOutputBuf.String())
		tplOutputBuf.Reset()

		cfg.SaveTo(outputPath)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.create-ini.yaml)")
	rootCmd.PersistentFlags().StringVar(&devPort, "devPort", "00000", "dev port for app")
	rootCmd.PersistentFlags().StringVar(&email, "email", "", "email for rerouting emails")
	rootCmd.PersistentFlags().StringVar(&phoneNumber, "phoneNumber", "", "phone number for debugging sms")
	rootCmd.PersistentFlags().StringVar(&redisPort, "redisPort", "0", "redis slot for dev")
	rootCmd.PersistentFlags().StringVar(&userName, "userName", "", "username for dev")

	rootCmd.MarkPersistentFlagRequired("devPort")
	rootCmd.MarkPersistentFlagRequired("email")
	rootCmd.MarkPersistentFlagRequired("phoneNumber")
	rootCmd.MarkPersistentFlagRequired("redisPort")
	rootCmd.MarkPersistentFlagRequired("userName")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".create-ini" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".create-ini")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}