package echo

import (
	"fmt"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

type LoopRenderer struct {
	InputDecoder *wav.Decoder
	Tempo        float64
	Duration     float64
	Delay        float64
	OutputPath   string
}

func NewLoopRenderer(inputPath string, outputPath string, tempo, duration, delay float64) (*LoopRenderer, error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	decoder := wav.NewDecoder(f)
	if !decoder.IsValidFile() {
		return nil, fmt.Errorf("invalid WAV file: %s", inputPath)
	}
	return &LoopRenderer{
		InputDecoder: decoder,
		Tempo:        tempo,
		Duration:     duration,
		Delay:        delay,
		OutputPath:   outputPath,
	}, nil
}

func (r *LoopRenderer) Render() {
	inBuf, err := r.InputDecoder.FullPCMBuffer()
	if err != nil {
		fmt.Printf("Error reading input samples: %v\n", err)
		return
	}
	if inBuf == nil {
		fmt.Println("No audio data found in input file")
		return
	}

	framesPerBeat := float64(r.InputDecoder.SampleRate) / (r.Tempo / 60.0)
	outputFrames := int(framesPerBeat * r.Duration)
	loopFrames := int(framesPerBeat * r.Delay)

	input := inBuf.Data
	loop := make([]int, loopFrames*int(r.InputDecoder.NumChans))
	output := make([]int, outputFrames*int(r.InputDecoder.NumChans))
	loopIndex := 0

	for i := 0; i < len(input); i++ {
		loop[loopIndex] += input[i]
		loopIndex++
		if loopIndex == len(loop) {
			loopIndex = 0
		}
	}

	loopIndex = 0

	for outIndex := 0; outIndex < len(output); outIndex++ {
		output[outIndex] = loop[loopIndex]
		loopIndex++
		if loopIndex == len(loop) {
			loopIndex = 0
		}
	}

	out, err := os.Create(r.OutputPath)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer out.Close()

	outputEncoder := wav.NewEncoder(
		out,
		int(r.InputDecoder.SampleRate),
		int(r.InputDecoder.BitDepth),
		int(r.InputDecoder.NumChans),
		1, // PCM format
	)
	defer outputEncoder.Close()

	outBuf := &audio.IntBuffer{
		Data: output,
		Format: &audio.Format{
			NumChannels: int(r.InputDecoder.NumChans),
			SampleRate:  int(r.InputDecoder.SampleRate),
		},
		SourceBitDepth: int(r.InputDecoder.BitDepth),
	}

	if err := outputEncoder.Write(outBuf); err != nil {
		fmt.Printf("Error writing WAV data: %v\n", err)
		return
	}
	if err := outputEncoder.Close(); err != nil {
		fmt.Printf("Error closing encoder: %v\n", err)
		return
	}
}
