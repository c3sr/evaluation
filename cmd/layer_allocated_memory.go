package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/c3sr/evaluation"
	"github.com/spf13/cobra"
)

var layerAllocatedMemoryCmd = &cobra.Command{
	Use:     "allocated_memory",
	Aliases: []string{},
	Short:   "Get model layer memory information from framework traces in a database",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if databaseName == "" {
			databaseName = defaultDatabaseName["layer"]
		}
		err := rootSetup()
		if err != nil {
			return err
		}
		if modelName == "all" && outputFormat == "json" && outputFileName == "" {
			outputFileName = filepath.Join(mlArcWebAssetsPath, "layers")
		}
		if overwrite && isExists(outputFileName) {
			os.RemoveAll(outputFileName)
		}
		if plotPath == "" {
			plotPath = evaluation.TempFile("", "layer_allocated_memory_plot_*.html")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		run := func() error {
			evals, err := getEvaluations()
			if err != nil {
				return err
			}

			summary0, err := evals.SummaryLayerInformations(performanceCollection)
			if err != nil {
				return err
			}
			summary := evaluation.SummaryLayerAllocatedMemoryInformations(summary0)

			if sortOutput {
				sort.Slice(summary, func(ii, jj int) bool {
					return evaluation.TrimmedMeanInt64Slice(summary[ii].AllocatedBytes, evaluation.DefaultTrimmedMeanFraction) > evaluation.TrimmedMeanInt64Slice(summary[jj].AllocatedBytes, evaluation.DefaultTrimmedMeanFraction)
				})
			}

			if openPlot {
				return summary.OpenBarPlot()
			}

			if barPlot {
				err := summary.WriteBarPlot(plotPath)
				if err != nil {
					return err
				}
				fmt.Println("Created plot in " + plotPath)
				return nil
			}

			writer := NewWriter(evaluation.SummaryLayerInformation{})
			defer writer.Close()

			for _, v := range summary0 {
				writer.Row(v)
			}
			return nil
		}

		return forallmodels(run)
	},
}
