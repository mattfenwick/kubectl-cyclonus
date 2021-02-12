package cli

import (
	"github.com/mattfenwick/cyclonus/pkg/cli"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	//"k8s.io/cli-runtime/pkg/genericclioptions"
)

func doOrDie(err error) {
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func RunRootCommand() {
	command := cli.SetupAnalyzeCommand()
	command.Use = "cyclonus"
	doOrDie(errors.Wrapf(command.Execute(), "run root command"))
}

/*
type Config struct {
	LogLevel       string
	//KubeFlags      *genericclioptions.ConfigFlags
	AnalyzeArgs *cli.AnalyzeArgs
}

func setupRootCmd() *cobra.Command {
	args := &Config{}

	cmd := &cobra.Command{
		Use:   "cyclonus",
		Short: "",
		Long:  `.`,
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, as []string) error {
			// TODO detect whether this is running under kubectl or not, and modify help message accordingly
			//   see https://krew.sigs.k8s.io/docs/developer-guide/develop/best-practices/#auth-plugins
			//   if strings.HasPrefix(filepath.Base(os.Args[0]), "kubectl-") { }
			runRootCmd(args)
			return nil
		},
	}

	cmd.Flags().StringVar(&args.LogLevel, "v", "info", "log level")

	//args.KubeFlags = genericclioptions.NewConfigFlags(false)
	//args.KubeFlags.AddFlags(cmd.Flags())

	return cmd
}

func runRootCmd(args *Config) {
	level, err := log.ParseLevel(args.LogLevel)
	doOrDie(err)
	log.SetLevel(level)
}
 */