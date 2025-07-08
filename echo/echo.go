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
}

func NewRenderer(filePath string, tempo, bars, division int) (*Renderer, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	decoder := wav.NewDecoder(f)
	if !decoder.IsValidFile() {
		return nil, fmt.Errorf("invalid WAV file: %s", filePath)
	}
	return &Renderer{
		InputDecoder: decoder,
		Tempo:        tempo,
		Bars:         bars,
		Division:     division,
	}, nil
}

func (r *Renderer) Render() {
	fmt.Printf("Rendering echo at %d BPM for %d bars (division %d)\n", r.Tempo, r.Bars, r.Division)

	samplesPerBar := int(r.Tempo) * int(r.InputDecoder.SampleRate) / 60
	samplesTotal := samplesPerBar * r.Bars * int(r.InputDecoder.NumChans)
	samplesLoop := samplesPerBar / r.Division * int(r.InputDecoder.NumChans)

	inBuf, err := r.InputDecoder.FullPCMBuffer()
	if err != nil {
		fmt.Printf("Error reading input samples: %v\n", err)
		return
	}
	if inBuf == nil {
		fmt.Println("No audio data found in input file")
		return
	}

	input := inBuf.Data
	loop := make([]int, samplesLoop)
	output := make([]int, samplesTotal)
	loopIndex := 0

	for renderIndex := 0; renderIndex < samplesTotal; renderIndex++ {
		if renderIndex < len(input) {
			loop[loopIndex] += input[renderIndex]
		}
		output[renderIndex] = loop[loopIndex]
		loopIndex++
		if loopIndex == len(loop) {
			loopIndex = 0
		}
	}

	out, err := os.Create("ignore/output.wav")
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
