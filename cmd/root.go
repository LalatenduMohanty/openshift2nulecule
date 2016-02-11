package cmd

import (
	"fmt"
	"os"

	"github.com/kadel/openshift2nulecule/o2n"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"

	cclientcmd "github.com/openshift/origin/pkg/cmd/util/clientcmd"
)

var OutputDir string
var Factory *cclientcmd.Factory

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "openshift2nulecule",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		jww.INFO.Printf("Staring")
		jww.DEBUG.Printf("args = %s", args)
		// TODO: check required arguments

		// export only k8s kinds, for now
		var kindsToExport []string
		kindsToExport = append(kindsToExport, "pods,replicationcontrollers,persistentvolumeclaims,services")
		err := o2n.Export(OutputDir, Factory, cmd, kindsToExport)
		if err != nil {
			jww.FATAL.Printf("ERROR: %s", err)
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	jww.SetStdoutThreshold(jww.LevelDebug)
	RootCmd.Flags().StringVar(&OutputDir, "output-dir", "", "Directory where new Nulecule app will be created.")

	f := cclientcmd.New(RootCmd.PersistentFlags())
	Factory = f
}
