package cmd

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/c3sr/evaluation"
	"github.com/spf13/cobra"
)

var gpuKernelModelAggreInfoCmd = &cobra.Command{
	Use:     "model_aggre_info",
	Aliases: []string{},
	Short:   "Get gpu information aggregated within the model from system library traces in a database. Specify model name as `all` to list information of all the models.",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if databaseName == "" {
			databaseName = defaultDatabaseName["cuda_kernel"]
		}
		err := rootSetup()
		if err != nil {
			return err
		}
		if modelName == "all" && outputFormat == "json" && outputFileName == "" {
			outputFileName = filepath.Join(mlArcWebAssetsPath, "cuda_kernel_launch")
		}
		if overwrite && isExists(outputFileName) {
			os.RemoveAll(outputFileName)
		}

		if kernelNameFilterString != "" {
			kernelNameFilterList = strings.Split(kernelNameFilterString, ",")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		run := func() error {
			evals, err := getEvaluations()
			if err != nil {
				return err
			}

			gpuKernelInfos, err := evals.SummaryGPUKernelModelAggreInformations(performanceCollection)
			if err != nil {
				return err
			}

			if sortOutput || topKernels != -1 {
				sort.Sort(gpuKernelInfos)
				if topKernels != -1 {
					if topKernels >= len(gpuKernelInfos) {
						topKernels = len(gpuKernelInfos)
					}
					gpuKernelInfos = gpuKernelInfos[:topKernels]
				}
			}

			var writer *Writer
			if len(gpuKernelInfos) == 0 {
				writer = NewWriter(evaluation.SummaryGPUKernelInformation{})
				defer writer.Close()
			}
			writer = NewWriter(gpuKernelInfos[0])
			defer writer.Close()

			for _, elem := range gpuKernelInfos {
				writer.Row(elem)
			}

			return nil
		}

		return forallmodels(run)
	},
}
