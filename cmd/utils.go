package cmd

import (
	"errors"
	"go/build"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Unknwon/com"
	framework "github.com/c3sr/dlframework/framework/cmd"
	"github.com/c3sr/evaluation"
	udb "upper.io/db.v3"
)

func getEvaluations() (evaluation.Evaluations, error) {

	filter := udb.Cond{}
	if modelName != "" {
		filter["model.name"] = modelName
	}
	if modelVersion != "" {
		filter["model.version"] = modelVersion
	}
	if frameworkName != "" {
		filter["framework.name"] = frameworkName
	}
	if frameworkVersion != "" {
		filter["framework.version"] = frameworkVersion
	}
	if machineArchitecture != "" {
		filter["machinearchitecture"] = machineArchitecture
	}
	if hostName != "" {
		filter["hostname"] = hostName
	}
	if batchSize != 0 {
		filter["batch_size"] = batchSize
	}
	evals, err := evaluationCollection.Find(filter)
	if err != nil {
		return nil, err
	}

	if limit > 0 {
		// reverse list
		for i := len(evals)/2 - 1; i >= 0; i-- {
			opp := len(evals) - 1 - i
			evals[i], evals[opp] = evals[opp], evals[i]
		}
		evals = evals[:minInt(len(evals)-1, limit)]
	}

	return evaluation.Evaluations(evals), nil
}

func uptoIndex(arry []interface{}, idx int) int {
	if len(arry) <= idx {
		return len(arry) - 1
	}
	return idx
}

func maxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func minInt(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func getSrcPath(importPath string) (appPath string) {
	paths := com.GetGOPATHs()
	for _, p := range paths {
		d := filepath.Join(p, "src", importPath)
		if com.IsExist(d) {
			appPath = d
			break
		}
	}

	if len(appPath) == 0 {
		appPath = filepath.Join(goPath, "src", importPath)
	}

	return appPath
}

func isExists(s string) bool {
	return com.IsExist(s)
}

func getBuildFile() (string, error) {
	pkg, err := build.Default.ImportDir(sourcePath, build.ImportMode(0))
	if err == nil && pkg.IsCommand() {
		return pkg.SrcRoot, nil
	}

	mainPath := filepath.Join(sourcePath, "main.go")
	if com.IsFile(mainPath) {
		return mainPath, nil
	}

	return "", errors.New("unable to figure out what file to build")
}

// Random number state.
// We generate random temporary file names so that there's a good
// chance the file doesn't exist yet - keeps the number of tries in
// TempFile to a minimum.
var rand uint32
var randmu sync.Mutex

func reseed() uint32 {
	return uint32(time.Now().UnixNano() + int64(os.Getpid()))
}

func nextRandom() string {
	randmu.Lock()
	r := rand
	if r == 0 {
		r = reseed()
	}
	r = r*1664525 + 1013904223 // constants from Numerical Recipes
	rand = r
	randmu.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
}

// TempFile creates a new temporary file in the directory dir,
// opens the file for reading and writing, and returns the resulting *os.File.
// The filename is generated by taking pattern and adding a random
// string to the end. If pattern includes a "*", the random string
// replaces the last "*".
// If dir is the empty string, TempFile uses the default directory
// for temporary files (see os.TempDir).
// Multiple programs calling TempFile simultaneously
// will not choose the same file. The caller can use f.Name()
// to find the pathname of the file. It is the caller's responsibility
// to remove the file when no longer needed.
func TempFile(dir, pattern string) string {
	if dir == "" {
		dir = os.TempDir()
	}

	var prefix, suffix string
	if pos := strings.LastIndex(pattern, "*"); pos != -1 {
		prefix, suffix = pattern[:pos], pattern[pos+1:]
	} else {
		prefix = pattern
	}

	var name string

	nconflict := 0
	for i := 0; i < 10000; i++ {
		name = filepath.Join(dir, prefix+nextRandom()+suffix)
		if com.IsFile(name) {
			if nconflict++; nconflict > 10 {
				randmu.Lock()
				rand = reseed()
				randmu.Unlock()
			}
			continue
		}
		break
	}
	if !com.IsDir(filepath.Dir(name)) {
		os.MkdirAll(filepath.Dir(name), os.ModePerm)
	}
	return name
}

func forallmodels(run func() error) error {
	if modelName != "all" {
		return run()
	}

	outputDirectory := outputFileName
	if !com.IsDir(outputDirectory) {
		os.MkdirAll(outputDirectory, os.ModePerm)
	}
	for _, model := range framework.DefaultEvaulationModels {
		modelName, modelVersion = framework.ParseModelName(model)
		outputFileName = filepath.Join(outputDirectory, model+"."+outputFileExtension)
		run()
	}
	return nil
}
