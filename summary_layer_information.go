package evaluation

import (
	"errors"
	"strings"

	"github.com/rai-project/config"
	"github.com/rai-project/tracer"
	"github.com/spf13/cast"
	model "github.com/uber/jaeger/model/json"
	db "upper.io/db.v3"
)

type LayerInformation struct {
	Name      string    `json:"name,omitempty"`
	Durations []float64 `json:"durations,omitempty"`
}

type LayerInformations []LayerInformation

type SummaryLayerInformation struct {
	SummaryBase         `json:",inline"`
	MachineArchitecture string            `json:"machine_architecture,omitempty"`
	UsingGPU            bool              `json:"using_gpu,omitempty"`
	BatchSize           int               `json:"batch_size,omitempty"`
	HostName            string            `json:"host_name,omitempty"`
	LayerInformations   LayerInformations `json:"layer_informations,omitempty"`
}

type SummaryLayerInformations []SummaryLayerInformation

func (LayerInformation) Header() []string {
	return []string{
		"name",
		"durations",
	}
}

func (info LayerInformation) Row() []string {
	return []string{
		info.Name,
		strings.Join(float64SliceToStringSlice(info.Durations), "\t"),
	}
}

func (LayerInformations) Header() []string {
	return LayerInformation{}.Header()
}

func (s LayerInformations) Rows() [][]string {
	rows := [][]string{}
	for _, e := range s {
		rows = append(rows, e.Row())
	}
	return rows
}

func (SummaryLayerInformation) Header() []string {
	extra := []string{
		"machine_architecture",
		"using_gpu",
		"batch_size",
		"hostname",
		"layer_informations",
	}
	return append(SummaryBase{}.Header(), extra...)
}

func (s SummaryLayerInformation) Row() []string {
	infos := []string{}
	for _, row := range s.LayerInformations.Rows() {
		infos = append(infos, strings.Join(row, ":"))
	}
	extra := []string{
		s.MachineArchitecture,
		cast.ToString(s.UsingGPU),
		cast.ToString(s.BatchSize),
		s.HostName,
		strings.Join(infos, ";"),
	}
	return append(s.SummaryBase.Row(), extra...)
}

func (SummaryLayerInformations) Header() []string {
	return SummaryLayerInformation{}.Header()
}

func (s SummaryLayerInformations) Rows() [][]string {
	rows := [][]string{}
	for _, e := range s {
		rows = append(rows, e.Row())
	}
	return rows
}

type layerInformationMap map[string]LayerInformation

func (p Performance) LayerInformationSummary(e Evaluation) (*SummaryLayerInformation, error) {
	sspans := getSpanLayersFromSpans(p.Spans())
	numSSpans := len(sspans)

	summary := &SummaryLayerInformation{
		SummaryBase:         e.summaryBase(),
		MachineArchitecture: e.MachineArchitecture,
		UsingGPU:            e.UsingGPU,
		BatchSize:           e.BatchSize,
		HostName:            e.Hostname,
		LayerInformations:   LayerInformations{},
	}
	if numSSpans == 0 {
		return summary, nil
	}

	infosFullMap := make([]layerInformationMap, numSSpans)
	for ii, spans := range sspans {
		if infosFullMap[ii] == nil {
			infosFullMap[ii] = layerInformationMap{}
		}
		for _, span := range spans {
			opName := strings.ToLower(span.OperationName)
			if _, ok := infosFullMap[ii][opName]; !ok {
				infosFullMap[ii][opName] = LayerInformation{
					Name:      span.OperationName,
					Durations: []float64{},
				}
			}
			info := infosFullMap[ii][opName]
			info.Durations = append(info.Durations, cast.ToFloat64(span.Duration))
			infosFullMap[ii][opName] = info
		}
	}

	infoMap := layerInformationMap{}
	for _, span := range sspans[0] {
		opName := strings.ToLower(span.OperationName)
		if _, ok := infoMap[opName]; !ok {
			infoMap[opName] = LayerInformation{
				Name:      span.OperationName,
				Durations: []float64{},
			}
		}
		info := infoMap[opName]
		allDurations := [][]float64{}
		for ii := range sspans {
			allDurations = append(allDurations, infosFullMap[ii][opName].Durations)
		}
		transposedDurations := transpose(allDurations)
		durations := []float64{}
		for _, tr := range transposedDurations {
			durations = append(durations, trimmedMean(tr, DefaultTrimmedMeanFraction))
		}
		info.Durations = durations
		infoMap[opName] = info
	}

	infos := []LayerInformation{}
	for _, v := range infoMap {
		infos = append(infos, v)
	}

	summary.LayerInformations = infos
	return summary, nil
}

func (e Evaluation) LayerInformationSummary(perfCol *PerformanceCollection) (*SummaryLayerInformation, error) {
	perfs, err := perfCol.Find(db.Cond{"_id": e.PerformanceID})
	if err != nil {
		return nil, err
	}
	if len(perfs) != 1 {
		return nil, errors.New("expecting on performance output")
	}
	perf := perfs[0]
	return perf.LayerInformationSummary(e)
}

func spanIsCUPTI(span model.Span) bool {
	for _, tag := range span.Tags {
		key := strings.ToLower(tag.Key)
		switch key {
		case "cupti_domain", "cupti_callback_id":
			return true
		}
	}
	return false
}

func spanTagExists(span model.Span, key string) bool {
	for _, tag := range span.Tags {
		key0 := strings.ToLower(tag.Key)
		if key0 == key {
			return true
		}
	}
	return false
}

func spanTagValue(span model.Span, key string) (interface{}, bool) {
	for _, tag := range span.Tags {
		key0 := strings.ToLower(tag.Key)
		if key0 == key {
			return tag.Value, true
		}
	}
	return nil, false
}

func spanTagEquals(span model.Span, key string, value string) bool {
	for _, tag := range span.Tags {
		key0 := strings.ToLower(tag.Key)
		if key0 == key {
			return strings.ToLower(cast.ToString(tag.Value)) == value
		}
	}
	return false
}

func spanComponentIs(span model.Span, name string) bool {
	for _, tag := range span.Tags {
		key := strings.ToLower(tag.Key)
		switch key {
		case "component":
			return strings.ToLower(cast.ToString(tag.Value)) == name
		}
	}
	return false
}

func selectTensorflowLayerSpans(spans Spans) Spans {
	res := []model.Span{}
	for _, span := range spans {
		if spanIsCUPTI(span) {
			continue
		}
		if !spanComponentIs(span, config.App.Name) {
			continue
		}
		if !spanTagExists(span, "thread_id") {
			continue
		}
		if !spanTagExists(span, "timeline_label") {
			continue
		}
		res = append(res, span)
	}
	return nil
}
func selectMXNetLayerSpans(spans Spans) Spans {
	res := []model.Span{}
	for _, span := range spans {
		if spanIsCUPTI(span) {
			continue
		}
		if !spanComponentIs(span, config.App.Name) {
			continue
		}
		if !spanTagExists(span, "thread_id") {
			continue
		}
		if !spanTagExists(span, "process_id") {
			continue
		}
		if spanTagEquals(span, "metadata", "") {
			res = append(res, span)
			continue
		}
	}
	return nil
}
func selectCaffeLayerSpans(spans Spans) Spans {
	return selectCaffe2LayerSpans(spans)
}
func selectCaffe2LayerSpans(spans Spans) Spans {
	res := []model.Span{}
	for _, span := range spans {
		if spanIsCUPTI(span) {
			continue
		}
		if !spanComponentIs(span, config.App.Name) {
			continue
		}
		if !spanTagExists(span, "metadata") {
			continue
		}
		if !spanTagExists(span, "thread_id") {
			continue
		}
		if spanTagEquals(span, "process_id", "0") {
			res = append(res, span)
			continue
		}
	}
	return nil
}
func selectCNTKLayerSpans(spans Spans) Spans {
	log.WithField("function", "selectCNTKLayerSpans").Error("layer information is not currently supported by cntk")
	return Spans{}
}
func selectTensorRTLayerSpans(spans Spans) Spans {
	return selectCaffe2LayerSpans(spans)
}

func getSpanLayersFromSpans(spans Spans) []Spans {
	predictSpans := spans.FilterByOperationName("Predict")
	predictIndexOf := func(span model.Span) int {
		for ii, predict := range predictSpans {
			if span.ParentSpanID == predict.SpanID {
				return ii
			}
		}
		return -1
	}
	groupedSpans := make([]Spans, len(predictSpans))
	for _, span := range spans {
		idx := predictIndexOf(span)
		if idx == -1 {
			continue
		}
		groupedSpans[idx] = append(groupedSpans[idx], span)
	}
	groupedLayerSpans := make([]Spans, len(predictSpans))
	for ii, grp := range groupedSpans {
		groupedLayerSpans[ii] = Spans{}
		if len(grp) == 0 {
			continue
		}
		predict := predictSpans[ii]
		traceLevel0, ok := spanTagValue(predict, "trace_level")
		if !ok {
			continue
		}
		traceLevel, ok := traceLevel0.(string)
		if !ok {
			continue
		}
		if traceLevel == "" {
			continue
		}
		if tracer.LevelFromName(traceLevel) < tracer.FRAMEWORK_TRACE {
			continue
		}
		frameworkName := strings.ToLower(frameworkNameOfSpan(predict))
		switch frameworkName {
		case "tensorflow":
			groupedLayerSpans[ii] = selectTensorflowLayerSpans(grp)
		case "mxnet":
			groupedLayerSpans[ii] = selectMXNetLayerSpans(grp)
		case "caffe":
			groupedLayerSpans[ii] = selectCaffeLayerSpans(grp)
		case "caffe2":
			groupedLayerSpans[ii] = selectCaffe2LayerSpans(grp)
		case "cntk":
			groupedLayerSpans[ii] = selectCNTKLayerSpans(grp)
		case "tensorrt":
			groupedLayerSpans[ii] = selectTensorRTLayerSpans(grp)
		}
	}
	return groupedLayerSpans
}

func frameworkNameOfSpan(predictSpan model.Span) string {
	tagName := "framework_name"
	for _, tag := range predictSpan.Tags {
		if tag.Key == tagName {
			return cast.ToString(tag.Value)
		}
	}
	return ""
}
