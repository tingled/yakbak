package main

import (
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

func main() {
	f, err := os.Open("sample.wav")
	if err != nil {
		log.Fatal(err)
	}
	s, format, _ := wav.Decode(f)

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan struct{})

	// raised cosine lowpass filter, corner freq 2kHz
	// generated from https://www-users.cs.york.ac.uk/~fisher/cgi-bin/mkfscript
	filterCoeffs := []float64{
		+0.0001830294, +0.0003699462, +0.0005836596, +0.0007831327,
		+0.0009243916, +0.0009683493, +0.0008884547, +0.0006767556,
		+0.0003471317, -0.0000651488, -0.0005078017, -0.0009193970,
		-0.0012396205, -0.0014202847, -0.0014351547, -0.0012866806,
		-0.0010080975, -0.0006600359, -0.0003216861, -0.0000775473,
		-0.0000016939, -0.0001421492, -0.0005082200, -0.0010634394,
		-0.0017260602, -0.0023779246, -0.0028811413, -0.0031005435,
		-0.0029286343, -0.0023088683, -0.0012528903, +0.0001521593,
		+0.0017489329, +0.0033285777, +0.0046614493, +0.0055361538,
		+0.0058012634, +0.0054027577, +0.0044099178, +0.0030232319,
		+0.0015598644, +0.0004152511, +0.0000030509, +0.0006795489,
		+0.0026620886, +0.0059536379, +0.0102866705, +0.0150988129,
		+0.0195500612, +0.0225869312, +0.0230530991, +0.0198395420,
		+0.0120607292, -0.0007620787, -0.0185330470, -0.0404533068,
		-0.0649539024, -0.0897075180, -0.1117289738, -0.1275650961,
		-0.1335641754, -0.1262050063, -0.1024567182, -0.0601344186,
		+0.0017870031, +0.0829371583, +0.1814232082, +0.2938330192,
		+0.4153775254, +0.5401675018, +0.6616058940, +0.7728642018,
		+0.8674016331, +0.9394800102, +0.9846264573, +0.9999999445,
		+0.9846264573, +0.9394800102, +0.8674016331, +0.7728642018,
		+0.6616058940, +0.5401675018, +0.4153775254, +0.2938330192,
		+0.1814232082, +0.0829371583, +0.0017870031, -0.0601344186,
		-0.1024567182, -0.1262050063, -0.1335641754, -0.1275650961,
		-0.1117289738, -0.0897075180, -0.0649539024, -0.0404533068,
		-0.0185330470, -0.0007620787, +0.0120607292, +0.0198395420,
		+0.0230530991, +0.0225869312, +0.0195500612, +0.0150988129,
		+0.0102866705, +0.0059536379, +0.0026620886, +0.0006795489,
		+0.0000030509, +0.0004152511, +0.0015598644, +0.0030232319,
		+0.0044099178, +0.0054027577, +0.0058012634, +0.0055361538,
		+0.0046614493, +0.0033285777, +0.0017489329, +0.0001521593,
		-0.0012528903, -0.0023088683, -0.0029286343, -0.0031005435,
		-0.0028811413, -0.0023779246, -0.0017260602, -0.0010634394,
		-0.0005082200, -0.0001421492, -0.0000016939, -0.0000775473,
		-0.0003216861, -0.0006600359, -0.0010080975, -0.0012866806,
		-0.0014351547, -0.0014202847, -0.0012396205, -0.0009193970,
		-0.0005078017, -0.0000651488, +0.0003471317, +0.0006767556,
		+0.0008884547, +0.0009683493, +0.0009243916, +0.0007831327,
		+0.0005836596, +0.0003699462, +0.0001830294,
	}
	gain := 1 / 1.11e1

	orig, clone := beep.Dup(s)
	filter := NewFilter(filterCoeffs, gain)
	fs := &FilterStreamer{filter, clone}

	speaker.Play(
		beep.Seq(orig, fs, beep.Callback(func() {
			close(done)
		})))

	<-done
}
