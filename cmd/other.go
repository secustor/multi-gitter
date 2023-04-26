package cmd

import (
	"os"

	"github.com/lindell/multi-gitter/internal/scm"
	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
)

func outputFlag() *flag.FlagSet {
	flags := flag.NewFlagSet("output", flag.ExitOnError)

	flags.StringP("output", "o", "-", `The file that the output of the script should be outputted to. "-" means stdout. If using 'print' you can also use templates e.g. '--output "./results/{{ .FullName }}.json"'`)

	return flags
}

func getToken(flag *flag.FlagSet) (string, error) {
	if OverrideVersionController != nil {
		return "", nil
	}

	token, _ := flag.GetString("token")

	if token == "" {
		if ght := os.Getenv("GITHUB_TOKEN"); ght != "" {
			token = ght
		} else if ght := os.Getenv("GITLAB_TOKEN"); ght != "" {
			token = ght
		} else if ght := os.Getenv("GITEA_TOKEN"); ght != "" {
			token = ght
		} else if ght := os.Getenv("BITBUCKET_SERVER_TOKEN"); ght != "" {
			token = ght
		}
	}

	if token == "" {
		return "", errors.New("either the --token flag or the GITHUB_TOKEN/GITLAB_TOKEN/GITEA_TOKEN/BITBUCKET_SERVER_TOKEN environment variable has to be set")
	}

	return token, nil
}

func getMergeTypes(flag *flag.FlagSet) ([]scm.MergeType, error) {
	mergeTypeStrs, _ := flag.GetStringSlice("merge-type") // Only used for the merge command

	// Convert all defined merge types (if any)
	var err error
	mergeTypes := make([]scm.MergeType, len(mergeTypeStrs))
	for i, mt := range mergeTypeStrs {
		mergeTypes[i], err = scm.ParseMergeType(mt)
		if err != nil {
			return nil, err
		}
	}

	return mergeTypes, nil
}
