package cmd

import (
	"fmt"

	"github.com/dtan4/valec/aws"
	"github.com/dtan4/valec/secret"
	"github.com/dtan4/valec/util"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate SECRETDIR",
	Short: "Validate secrets in local files",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Please specify secret directory.")
		}
		dirname := args[0]

		files, err := util.ListYAMLFiles(dirname)
		if err != nil {
			return errors.Wrapf(err, "Failed to read directory. dirname=%s", dirname)
		}

		for _, file := range files {
			if err := validateFile(file); err != nil {
				return errors.Wrapf(err, "Failed to validate secrets. filename=%s", file)
			}
		}

		fmt.Println("All secrets are valid.")

		return nil
	},
}

func validateFile(filename string) error {
	fmt.Println(filename)

	secrets, err := secret.LoadFromYAML(filename)
	if err != nil {
		return errors.Wrapf(err, "Failed to load secrets. filename=%s", filename)
	}

	hasError := false
	red := color.New(color.FgRed)

	for _, secret := range secrets {
		if _, err := aws.KMS.DecryptBase64(secret.Key, secret.Value); err != nil {
			red.Printf("  Secret value is invalid. Please try `valec encrypt`. key=%s\n", secret.Key)
			hasError = true
		}
	}

	if hasError {
		return errors.New("Some secrets are invalid.")
	}

	return nil
}

func init() {
	RootCmd.AddCommand(validateCmd)
}
