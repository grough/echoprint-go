package echo

import (
	"fmt"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

type Renderer struct {
	InputDecoder *wav.Decoder
	Tempo        int
	Bars         int
	Division     int
	OutputPath   string
}

func NewRenderer(inputPath string, outputPath string, tempo, bars, division int) (*Renderer, error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	decoder := wav.NewDecoder(f)
	if !decoder.IsValidFile() {
		return nil, fmt.Errorf("invalid WAV file: %s", inputPath)
	}
	return &Renderer{
		InputDecoder: decoder,
		Tempo:        tempo,
		Bars:         bars,
		Division:     division,
		OutputPath:   outputPath,
	}, nil
}

func (r *Renderer) Render() {
	inBuf, err := r.InputDecoder.FullPCMBuffer()
	if err != nil {
		fmt.Printf("Error reading input samples: %v\n", err)
		return
	}
	if inBuf == nil {
		fmt.Println("No audio data found in input file")
		return
	}

	framesPerBar := int(r.Tempo) * int(r.InputDecoder.SampleRate) / 60
	outputFrames := framesPerBar * r.Bars
	loopFrames := framesPerBar / r.Division

	input := inBuf.Data
	loop := make([]int, loopFrames*int(r.InputDecoder.NumChans))
	output := make([]int, outputFrames*int(r.InputDecoder.NumChans))
	loopIndex := 0

	for outIndex := 0; outIndex < len(output); outIndex++ {
		if outIndex < len(input) {
			loop[loopIndex] += input[outIndex]
		}
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
