package main

import (
	"container/ring"

	"github.com/faiface/beep"
)

type RingBuffer struct {
	*ring.Ring
}

func NewBuffer(length int) *RingBuffer {
	r := &RingBuffer{ring.New(length)}
	r.zero()
	return r
}

func (r *RingBuffer) zero() {
	for i := 0; i < r.Len(); i++ {
		r.Value = float64(0)
		r.Ring = r.Next()
	}
}

func (r *RingBuffer) addSample(s float64) {
	r.Ring = r.Prev()
	r.Value = s
}

func (r *RingBuffer) bufferAsSlice() (s []float64) {
	r.Do(func(x interface{}) {
		s = append(s, x.(float64))
	})
	return s
}

// Filters
type Filter struct {
	*RingBuffer
	coeffs []float64
	gain   float64
}

func NewFilter(coeffs []float64, gain float64) *Filter {
	// some filter design tools produce a gain as well
	return &Filter{
		RingBuffer: NewBuffer(len(coeffs)),
		coeffs:     coeffs,
		gain:       gain,
	}
}

func (f *Filter) filterSample(s float64) (output float64) {
	f.addSample(s)
	rPtr := f.Ring
	for _, c := range f.coeffs {
		output += rPtr.Value.(float64) * c
		rPtr = rPtr.Next()
	}
	return output * f.gain
}

func (f *Filter) filterBuffer(buf []float64) []float64 {
	output := make([]float64, len(buf))
	for i, s := range buf {
		output[i] = f.filterSample(s)
	}
	return output
}

type FilterStreamer struct {
	*Filter
	Streamer beep.Streamer
}

func (fs *FilterStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	n, ok = fs.Streamer.Stream(samples)
	for i := range samples {
		mix := (samples[i][0] + samples[i][1]) / 2
		out := fs.filterSample(mix)
		samples[i][0], samples[i][1] = out, out
	}
	return n, ok
}

func (fs *FilterStreamer) Err() error {
	return fs.Streamer.Err()
}
