package cmd

import (
	"os"
	"path/filepath"

	"github.com/c3sr/evaluation"
	"github.com/spf13/cobra"
)

var accuracyCmd = &cobra.Command{
	Use: "accuracy",
	Aliases: []string{
		"top_accuracy",
	},
	Short: "Get accuracy summary from the database",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if databaseName == "" {
			databaseName = defaultDatabaseName[cmd.Name()]
		}
		err := rootSetup()
		if err != nil {
			return err
		}
		if modelName == "all" && outputFormat == "json" && outputFileName == "" {
			outputFileName = filepath.Join(mlArcWebAssetsPath, "accuracy")
		}
		if overwrite && isExists(outputFileName) {
			os.RemoveAll(outputFileName)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		run := func() error {
			accs, err := predictAccuracyInformationSummary()
			if err != nil {
				return err
			}

			writer := NewWriter(evaluation.SummaryModelAccuracyInformation{})
			defer writer.Close()

			for _, acc := range accs {
				writer.Row(acc)
			}

			return nil
		}
		return forallmodels(run)
	},
}

func predictAccuracyInformationSummary() (evaluation.SummaryModelAccuracyInformations, error) {
	evals, err := getEvaluations()
	if err != nil {
		return nil, err
	}

	accs, err := evals.PredictAccuracyInformationSummary(modelAccuracyCollection)
	if err != nil {
		return nil, err
	}

	return accs.Group()
}
