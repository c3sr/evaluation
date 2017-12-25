package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/rai-project/evaluation"
)

var accuracymd = &cobra.Command{
	Use: "accuracy",
	Aliases: []string{
		"top_accuracy",
	},
	Short: "Get accuracy summary from CarML",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if databaseName == "" {
			databaseName = defaultDatabaseName[cmd.Name()]
		}
		rootSetup()
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
			accs, err := predictDurationInformationSummary()
			if err != nil {
				return err
			}

			writer := NewWriter(evaluation.SummaryPredictDurationInformation{})
			defer writer.Close()

			for _, acc := range accs {
				writer.Row(acc)
			}

			return nil
		}
		return forallmodels(run)
	},
}

func predictAccuracyInformationSummary() (evaluation.SummaryPredictDurationInformations, error) {
	evals, err := getEvaluations()
	if err != nil {
		return nil, err
	}
	return evals.PredictDurationInformationSummary(performanceCollection)
}