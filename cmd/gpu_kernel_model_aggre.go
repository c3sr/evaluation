package cmd

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rai-project/evaluation"
	"github.com/spf13/cobra"
)

var gpuKernelModelAggreCmd = &cobra.Command{
	Use:     model_aggre,
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

			cudaKernelInfos, err := evals.SummaryGPUKernelModelAggreInformations(performanceCollection)
			if err != nil {
				return err
			}

			if sortOutput || topKernels != -1 {
				sort.Sort(cudaKernelInfos)
				if topKernels != -1 {
					if topKernels >= len(cudaKernelInfos) {
						topKernels = len(cudaKernelInfos)
					}
					cudaKernelInfos = cudaKernelInfos[:topKernels]
				}
			}

			var writer *Writer
			if len(cudaKernelInfos) == 0 {
				writer = NewWriter(evaluation.SummaryGPUKernelInformation{})
				defer writer.Close()
			}
			writer = NewWriter(cudaKernelInfos[0])
			defer writer.Close()

			for _, elem := range cudaKernelInfos {
				writer.Row(elem)
			}

			return nil
		}

		return forallmodels(run)
	},
}