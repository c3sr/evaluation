package evaluation

import (
	json "encoding/json"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/k0kubun/pp"
	"github.com/rai-project/evaluation/writer"
	"github.com/rai-project/go-echarts/charts"
	"github.com/rai-project/tracer"
	"github.com/spf13/cast"
	model "github.com/uber/jaeger/model/json"
	db "upper.io/db.v3"
)

var (
	cntkLogMessageShown = false
)

//easyjson:json
type SummaryLayerInformation struct {
	SummaryBase        `json:",inline"`
	Index              int       `json:"index,omitempty"`
	Name               string    `json:"name,omitempty"`
	Type               string    `json:"type,omitempty"`
	StaticType         string    `json:"static_type,omitempty"`
	Shape              string    `json:"shap,omitempty"`
	Durations          []float64 `json:"durations,omitempty"`
	AllocatedBytes     []int64   `json:"allocated_bytes,omitempty"`
	HostTempMemSizes   []int64   `json:"host_temp_mem_sizes,omitempty"`
	DeviceTempMemSizes []int64   `json:"device_temp_mem_sizes,omitempty"`
}

type SummaryMeanLayerInformation struct {
	SummaryLayerInformation
}

//easyjson:json
type SummaryMeanLayerInformations []SummaryMeanLayerInformation

//easyjson:json
type SummaryLayerInformations []SummaryLayerInformation

func (SummaryLayerInformation) Header(iopts ...writer.Option) []string {
	extra := []string{
		"layer_index",
		"layer_name",
		"layer_type",
		"layer_shape",
		"layer_duration (us)",
		"layer_allocated_bytes",
		"layer_host_temp_mem_size",
		"layer_device_temp_mem_size",
	}
	opts := writer.NewOptions(iopts...)
	if opts.ShowSummaryBase {
		return append(SummaryBase{}.Header(iopts...), extra...)
	}
	return extra
}

func (s SummaryLayerInformation) Row(iopts ...writer.Option) []string {
	extra := []string{
		cast.ToString(s.Index),
		s.Name,
		s.Type,
		s.Shape,
		strings.Join(float64SliceToStringSlice(s.Durations), ","),
		strings.Join(int64SliceToStringSlice(s.AllocatedBytes), ","),
		strings.Join(int64SliceToStringSlice(s.HostTempMemSizes), ","),
		strings.Join(int64SliceToStringSlice(s.DeviceTempMemSizes), ","),
	}
	opts := writer.NewOptions(iopts...)
	if opts.ShowSummaryBase {
		return append(s.SummaryBase.Row(iopts...), extra...)
	}
	return extra
}

func (s SummaryMeanLayerInformation) Row(opts ...writer.Option) []string {
	return []string{
		cast.ToString(s.Index),
		s.Name,
		s.Type,
		s.Shape,
		cast.ToString(TrimmedMean(s.Durations, 0)),
		cast.ToString(TrimmedMean(convertInt64SliceToFloat64Slice(s.AllocatedBytes), 0)),
		cast.ToString(TrimmedMean(convertInt64SliceToFloat64Slice(s.HostTempMemSizes), 0)),
		cast.ToString(TrimmedMean(convertInt64SliceToFloat64Slice(s.DeviceTempMemSizes), 0)),
	}
}

func summaryLayerInformations(es Evaluations, spans Spans) (SummaryLayerInformations, error) {
	summary := SummaryLayerInformations{}
	if len(es) == 0 {
		return summary, errors.New("no evaluation is found in the database")
	}

	cPredictSpans := spans.FilterByOperationNameAndEvalTraceLevel("c_predict", tracer.FRAMEWORK_TRACE.String())
	groupedLayerSpans, err := getGroupedLayerSpansFromSpans(cPredictSpans, spans)
	if err != nil {
		return summary, err
	}
	numGroups := len(groupedLayerSpans)
	if numGroups == 0 {
		return summary, errors.New("no group of spans is found")
	}

	groupedLayerInfos := make([][]SummaryLayerInformation, numGroups)
	for ii, spans := range groupedLayerSpans {
		if groupedLayerInfos[ii] == nil {
			groupedLayerInfos[ii] = []SummaryLayerInformation{}
		}
		for _, span := range spans {
			idx, err := getTagValueAsString(span, "layer_sequence_index")
			if err != nil || idx == "" {
				return summary, errors.New("cannot find tag layer_sequence_index")
			}
			shape, _ := getTagValueAsString(span, "shape")
			staticType, _ := getTagValueAsString(span, "static_type")
			allocation := getAllocationBytes(span)
			hostTempMemSize, _ := getTagValueAsString(span, "temp_memory_size")
			deviceTempMemSize, _ := getTagValueAsString(span, "device_temp_memory_size")
			layerInfo := SummaryLayerInformation{
				Index:      cast.ToInt(idx),
				Name:       span.OperationName,
				Type:       getOpName(span),
				StaticType: staticType,
				Shape:      shape,
				Durations: []float64{
					cast.ToFloat64(span.Duration),
				},
				AllocatedBytes: []int64{
					cast.ToInt64(allocation),
				},
				HostTempMemSizes: []int64{
					cast.ToInt64(hostTempMemSize),
				},
				DeviceTempMemSizes: []int64{
					cast.ToInt64(deviceTempMemSize),
				},
			}
			groupedLayerInfos[ii] = append(groupedLayerInfos[ii], layerInfo)
		}
	}

	for ii, span := range groupedLayerSpans[0] {
		durations := []float64{}
		allocations := []int64{}
		hostTempMems := []int64{}
		deviceTempMems := []int64{}
		idx, err := getTagValueAsString(span, "layer_sequence_index")
		if err != nil || idx == "" {
			return summary, errors.New("cannot find tag layer_sequence_index")
		}
		for _, infos := range groupedLayerInfos {
			if len(infos) <= ii {
				continue
			}
			durationToAppend := []float64{}
			allocationToAppend := []int64{}
			hostTemMemToAppend := []int64{}
			deviceTemMemToAppend := []int64{}
			for _, info := range infos {
				if info.Index == cast.ToInt(idx) && info.Name == span.OperationName {
					durationToAppend = append(durationToAppend, info.Durations...)
					allocationToAppend = append(allocationToAppend, info.AllocatedBytes...)
					hostTemMemToAppend = append(hostTemMemToAppend, info.HostTempMemSizes...)
					deviceTemMemToAppend = append(deviceTemMemToAppend, info.DeviceTempMemSizes...)
				}
			}
			durations = append(durations, durationToAppend...)
			allocations = append(allocations, allocationToAppend...)
			hostTempMems = append(hostTempMems, hostTemMemToAppend...)
			deviceTempMems = append(deviceTempMems, deviceTemMemToAppend...)
		}
		shape, _ := getTagValueAsString(span, "shape")
		staticType, _ := getTagValueAsString(span, "static_type")
		summary = append(summary,
			SummaryLayerInformation{
				SummaryBase:        es[0].summaryBase(),
				Index:              cast.ToInt(idx),
				Name:               span.OperationName,
				Type:               getOpName(span),
				StaticType:         staticType,
				Shape:              shape,
				Durations:          durations,
				AllocatedBytes:     allocations,
				HostTempMemSizes:   hostTempMems,
				DeviceTempMemSizes: deviceTempMems,
			})
	}

	return summary, nil
}

func (es Evaluations) SummaryLayerInformations(perfCol *PerformanceCollection) (SummaryLayerInformations, error) {
	summary := SummaryLayerInformations{}
	spans := []model.Span{}
	for _, e := range es {
		foundPerfs, err := perfCol.Find(db.Cond{"_id": e.PerformanceID})
		if err != nil {
			return summary, err
		}
		if len(foundPerfs) != 1 {
			return summary, errors.New("no performance is found for the evaluation")
		}
		perf := foundPerfs[0]
		spans = append(spans, perf.Spans()...)
	}
	if len(spans) == 0 {
		return summary, errors.New("no span is found for the evaluation")
	}

	return summaryLayerInformations(es, spans)
}

func sortByLayerIndex(spans Spans) {
	sort.Slice(spans, func(ii, jj int) bool {
		li, foundI := spanTagValue(spans[ii], "layer_sequence_index")
		if !foundI {
			return false
		}
		lj, foundJ := spanTagValue(spans[jj], "layer_sequence_index")
		if !foundJ {
			return true
		}

		return cast.ToInt64(li) < cast.ToInt64(lj)
	})
}

func getGroupedLayerSpansFromSpans(cPredictSpans Spans, spans Spans) ([]Spans, error) {
	groupedSpans, err := getGroupedSpansFromSpans(cPredictSpans, spans)
	if err != nil {
		return nil, err
	}
	numPredictSpans := len(groupedSpans)

	groupedLayerSpans := make([]Spans, numPredictSpans)
	for ii, grsp := range groupedSpans {
		if len(grsp) == 0 {
			continue
		}

		groupedLayerSpans[ii] = Spans{}
		for _, sp := range grsp {
			traceLevel, err := getTagValueAsString(sp, "trace_level")
			if err != nil || traceLevel == "" {
				continue
			}
			if tracer.LevelFromName(traceLevel) != tracer.FRAMEWORK_TRACE {
				continue
			}
			groupedLayerSpans[ii] = append(groupedLayerSpans[ii], sp)
		}

		sortByLayerIndex(groupedLayerSpans[ii])
	}

	return groupedLayerSpans, nil
}

func (s SummaryLayerInformations) GetLayerInfoByName(name string) SummaryLayerInformation {
	for _, info := range s {
		if info.Name == name {
			return info
		}
	}
	return SummaryLayerInformation{}
}

func (o SummaryLayerInformations) BarPlot(title string) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.TitleOpts{Title: title},
		charts.ToolboxOpts{Show: true},
	)
	bar = o.BarPlotAdd(bar)
	return bar
}

func (o SummaryLayerInformations) BarPlotAdd(bar *charts.Bar) *charts.Bar {
	timeUnit := time.Microsecond
	labels := []string{}
	for _, elem := range o {
		labels = append(labels, elem.Name)
	}
	bar.AddXAxis(labels)

	durations := make([]time.Duration, len(o))
	for ii, elem := range o {
		val := TrimmedMean(elem.Durations, 0)
		durations[ii] = time.Duration(val)
	}
	bar.AddYAxis("", durations)
	bar.SetSeriesOptions(charts.LabelTextOpts{Show: false})
	bar.SetGlobalOptions(
		charts.XAxisOpts{Name: "Layer Index"},
		charts.YAxisOpts{Name: "Latency(" + unitName(timeUnit) + ")"},
	)
	return bar
}

func (o SummaryLayerInformations) BoxPlot(title string) *charts.BoxPlot {
	box := charts.NewBoxPlot()
	box.SetGlobalOptions(
		charts.TitleOpts{Title: title},
		charts.ToolboxOpts{Show: true},
	)
	box = o.BoxPlotAdd(box)
	return box
}

func (o SummaryLayerInformations) BoxPlotAdd(box *charts.BoxPlot) *charts.BoxPlot {
	timeUnit := time.Microsecond

	isPrivate := func(info SummaryLayerInformation) bool {
		return strings.HasPrefix(info.Name, "_")
	}

	labels := []string{}
	for _, elem := range o {
		if isPrivate(elem) {
			continue
		}
		labels = append(labels, elem.Name)
	}
	box.AddXAxis(labels)

	durations := make([][]time.Duration, 0, len(o))
	for _, elem := range o {
		if isPrivate(elem) {
			continue
		}
		ts := make([]time.Duration, len(elem.Durations))
		for jj, t := range elem.Durations {
			ts[jj] = time.Duration(t)
		}
		durations = append(durations, prepareBoxplotData(ts))
	}
	if false {
		pp.Println(len(labels))
		pp.Println(len(durations))
		pp.Println(len(durations[0]))
	}
	box.AddYAxis("", durations)
	box.SetSeriesOptions(charts.LabelTextOpts{Show: false})
	jsLabelsBts, _ := json.Marshal(labels)
	jsFun := `function (name, index) {
    var labels = ` + strings.Replace(string(jsLabelsBts), `"`, "'", -1) + `;
    return labels.indexOf(name);
  }`
	box.SetGlobalOptions(
		charts.XAxisOpts{
			Name:      "Layer Name",
			Type:      "category",
			AxisLabel: charts.LabelTextOpts{Show: true, Rotate: 45, Formatter: charts.FuncOpts(jsFun)},
			SplitLine: charts.SplitLineOpts{Show: false},
			SplitArea: charts.SplitAreaOpts{Show: true},
		},
		charts.YAxisOpts{
			Name: "Latency(" + unitName(timeUnit) + ")",
			Type: "value",
			// NameRotate: 90,
			AxisLabel: charts.LabelTextOpts{Formatter: "{value}" + unitName(timeUnit)},
			SplitArea: charts.SplitAreaOpts{Show: true},
			Mix:       0,
		},
		charts.DataZoomOpts{
			Type:       "slider",
			XAxisIndex: []int{0},
			Start:      0,
			End:        float32(len(labels)),
		},
	)
	return box
}

func prepareBoxplotData(ds []time.Duration) []time.Duration {
	min := durationMin(ds)
	q1 := durationPercentile(ds, 25)
	q2 := durationPercentile(ds, 50)
	q3 := durationPercentile(ds, 75)
	max := durationMax(ds)
	return []time.Duration{min, q1, q2, q3, max}
}

func (o SummaryLayerInformations) Name() string {
	if len(o) == 0 {
		return ""
	}
	return o[0].ModelName + " Layer Latency"
}

func (o SummaryLayerInformations) WriteBarPlot(path string) error {
	return writeBarPlot(o, path)
}

func (o SummaryLayerInformations) WriteBoxPlot(path string) error {
	return writeBoxPlot(o, path)
}

func (o SummaryLayerInformations) OpenBarPlot() error {
	return openBarPlot(o)
}

func (o SummaryLayerInformations) OpenBoxPlot() error {
	return openBoxPlot(o)
}
