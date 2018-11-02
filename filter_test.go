package main

import (
	"reflect"
	"testing"
)

func TestRingBuffer(t *testing.T) {
	testCases := []struct {
		name         string
		length       int
		samplesToAdd []float64
		expectedBuf  []float64
	}{
		{"initialize to zeros", 2, []float64{}, []float64{0, 0}},
		{"add samples", 2, []float64{0, 1, 2, 3}, []float64{3, 2}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// initialize buffer
			buffer := NewBuffer(tc.length)
			for _, s := range tc.samplesToAdd {
				buffer.addSample(s)
			}

			got := buffer.bufferAsSlice()

			if !reflect.DeepEqual(tc.expectedBuf, got) {
				t.Errorf("wanted %v, got %v", got, tc.expectedBuf)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	testCases := []struct {
		name          string
		coeffs        []float64
		gain          float64
		inputSamples  []float64
		outputSamples []float64
	}{
		{"all zeros", []float64{0, 0}, 1, []float64{1, 100, -24}, []float64{0, 0, 0}},
		{"ones", []float64{1, 1}, 1, []float64{1, 2, 3}, []float64{1, 3, 5}},
		{"ones with half gain", []float64{1, 1}, 0.5, []float64{1, 2, 3}, []float64{0.5, 1.5, 2.5}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// initialize filter
			filter := NewFilter(tc.coeffs, tc.gain)
			got := filter.filterBuffer(tc.inputSamples)

			if !reflect.DeepEqual(tc.outputSamples, got) {
				t.Errorf("wanted %v, got %v", tc.outputSamples, got)
			}
		})
	}
}
