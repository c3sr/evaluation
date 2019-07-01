// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

//go:generate go get github.com/cheekybits/genny

package evaluation

import (
	"errors"

	"github.com/spf13/cast"
	"github.com/uber/jaeger/model/json"
)

func getTagValueAsUint(span json.Span, key string) (uint, error) {
	var res uint
	if isZero(span) {
		return res, errors.New("nil span")
	}
	for _, tag := range span.Tags {
		if tag.Key == key {
			return cast.ToUintE(tag.Value)
		}
	}
	return res, errors.New("tag not found")
}

func mustGetTagValueAsUint(span json.Span, key string) uint {
	val, err := getTagValueAsUint(span, key)
	if err != nil {
		log.WithError(err).WithField("key", key).Fatal("failed to get tag")
	}
	return val
}

//go:generate go get github.com/cheekybits/genny

func getTagValueAsUint8(span json.Span, key string) (uint8, error) {
	var res uint8
	if isZero(span) {
		return res, errors.New("nil span")
	}
	for _, tag := range span.Tags {
		if tag.Key == key {
			return cast.ToUint8E(tag.Value)
		}
	}
	return res, errors.New("tag not found")
}

func mustGetTagValueAsUint8(span json.Span, key string) uint8 {
	val, err := getTagValueAsUint8(span, key)
	if err != nil {
		log.WithError(err).WithField("key", key).Fatal("failed to get tag")
	}
	return val
}

//go:generate go get github.com/cheekybits/genny

func getTagValueAsUint16(span json.Span, key string) (uint16, error) {
	var res uint16
	if isZero(span) {
		return res, errors.New("nil span")
	}
	for _, tag := range span.Tags {
		if tag.Key == key {
			return cast.ToUint16E(tag.Value)
		}
	}
	return res, errors.New("tag not found")
}

func mustGetTagValueAsUint16(span json.Span, key string) uint16 {
	val, err := getTagValueAsUint16(span, key)
	if err != nil {
		log.WithError(err).WithField("key", key).Fatal("failed to get tag")
	}
	return val
}

//go:generate go get github.com/cheekybits/genny

func getTagValueAsUint32(span json.Span, key string) (uint32, error) {
	var res uint32
	if isZero(span) {
		return res, errors.New("nil span")
	}
	for _, tag := range span.Tags {
		if tag.Key == key {
			return cast.ToUint32E(tag.Value)
		}
	}
	return res, errors.New("tag not found")
}

func mustGetTagValueAsUint32(span json.Span, key string) uint32 {
	val, err := getTagValueAsUint32(span, key)
	if err != nil {
		log.WithError(err).WithField("key", key).Fatal("failed to get tag")
	}
	return val
}

//go:generate go get github.com/cheekybits/genny

func getTagValueAsUint64(span json.Span, key string) (uint64, error) {
	var res uint64
	if isZero(span) {
		return res, errors.New("nil span")
	}
	for _, tag := range span.Tags {
		if tag.Key == key {
			return cast.ToUint64E(tag.Value)
		}
	}
	return res, errors.New("tag not found")
}

func mustGetTagValueAsUint64(span json.Span, key string) uint64 {
	val, err := getTagValueAsUint64(span, key)
	if err != nil {
		log.WithError(err).WithField("key", key).Fatal("failed to get tag")
	}
	return val
}

//go:generate go get github.com/cheekybits/genny

func getTagValueAsInt(span json.Span, key string) (int, error) {
	var res int
	if isZero(span) {
		return res, errors.New("nil span")
	}
	for _, tag := range span.Tags {
		if tag.Key == key {
			return cast.ToIntE(tag.Value)
		}
	}
	return res, errors.New("tag not found")
}

func mustGetTagValueAsInt(span json.Span, key string) int {
	val, err := getTagValueAsInt(span, key)
	if err != nil {
		log.WithError(err).WithField("key", key).Fatal("failed to get tag")
	}
	return val
}

//go:generate go get github.com/cheekybits/genny

func getTagValueAsInt8(span json.Span, key string) (int8, error) {
	var res int8
	if isZero(span) {
		return res, errors.New("nil span")
	}
	for _, tag := range span.Tags {
		if tag.Key == key {
			return cast.ToInt8E(tag.Value)
		}
	}
	return res, errors.New("tag not found")
}

func mustGetTagValueAsInt8(span json.Span, key string) int8 {
	val, err := getTagValueAsInt8(span, key)
	if err != nil {
		log.WithError(err).WithField("key", key).Fatal("failed to get tag")
	}
	return val
}

//go:generate go get github.com/cheekybits/genny

func getTagValueAsInt16(span json.Span, key string) (int16, error) {
	var res int16
	if isZero(span) {
		return res, errors.New("nil span")
	}
	for _, tag := range span.Tags {
		if tag.Key == key {
			return cast.ToInt16E(tag.Value)
		}
	}
	return res, errors.New("tag not found")
}

func mustGetTagValueAsInt16(span json.Span, key string) int16 {
	val, err := getTagValueAsInt16(span, key)
	if err != nil {
		log.WithError(err).WithField("key", key).Fatal("failed to get tag")
	}
	return val
}

//go:generate go get github.com/cheekybits/genny

func getTagValueAsInt32(span json.Span, key string) (int32, error) {
	var res int32
	if isZero(span) {
		return res, errors.New("nil span")
	}
	for _, tag := range span.Tags {
		if tag.Key == key {
			return cast.ToInt32E(tag.Value)
		}
	}
	return res, errors.New("tag not found")
}

func mustGetTagValueAsInt32(span json.Span, key string) int32 {
	val, err := getTagValueAsInt32(span, key)
	if err != nil {
		log.WithError(err).WithField("key", key).Fatal("failed to get tag")
	}
	return val
}

//go:generate go get github.com/cheekybits/genny

func getTagValueAsInt64(span json.Span, key string) (int64, error) {
	var res int64
	if isZero(span) {
		return res, errors.New("nil span")
	}
	for _, tag := range span.Tags {
		if tag.Key == key {
			return cast.ToInt64E(tag.Value)
		}
	}
	return res, errors.New("tag not found")
}

func mustGetTagValueAsInt64(span json.Span, key string) int64 {
	val, err := getTagValueAsInt64(span, key)
	if err != nil {
		log.WithError(err).WithField("key", key).Fatal("failed to get tag")
	}
	return val
}

//go:generate go get github.com/cheekybits/genny

func getTagValueAsFloat32(span json.Span, key string) (float32, error) {
	var res float32
	if isZero(span) {
		return res, errors.New("nil span")
	}
	for _, tag := range span.Tags {
		if tag.Key == key {
			return cast.ToFloat32E(tag.Value)
		}
	}
	return res, errors.New("tag not found")
}

func mustGetTagValueAsFloat32(span json.Span, key string) float32 {
	val, err := getTagValueAsFloat32(span, key)
	if err != nil {
		log.WithError(err).WithField("key", key).Fatal("failed to get tag")
	}
	return val
}

//go:generate go get github.com/cheekybits/genny

func getTagValueAsFloat64(span json.Span, key string) (float64, error) {
	var res float64
	if isZero(span) {
		return res, errors.New("nil span")
	}
	for _, tag := range span.Tags {
		if tag.Key == key {
			return cast.ToFloat64E(tag.Value)
		}
	}
	return res, errors.New("tag not found")
}

func mustGetTagValueAsFloat64(span json.Span, key string) float64 {
	val, err := getTagValueAsFloat64(span, key)
	if err != nil {
		log.WithError(err).WithField("key", key).Fatal("failed to get tag")
	}
	return val
}

//go:generate go get github.com/cheekybits/genny

func getTagValueAsString(span json.Span, key string) (string, error) {
	var res string
	if isZero(span) {
		return res, errors.New("nil span")
	}
	for _, tag := range span.Tags {
		if tag.Key == key {
			return cast.ToStringE(tag.Value)
		}
	}
	return res, errors.New("tag not found")
}

func mustGetTagValueAsString(span json.Span, key string) string {
	val, err := getTagValueAsString(span, key)
	if err != nil {
		log.WithError(err).WithField("key", key).Fatal("failed to get tag")
	}
	return val
}