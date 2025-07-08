package echo

import "fmt"

type Renderer struct {
	File     string
	Tempo    int16
	Bars     int16
	Division int16
}

func NewRenderer(file string, tempo int16, bars int16, division int16) *Renderer {
	return &Renderer{
		File:     file,
		Tempo:    tempo,
		Bars:     bars,
		Division: division,
	}
}

func (r *Renderer) Render() {
	fmt.Printf("Stub: rendering %s at %d BPM for %d bars (division %d)\n", r.File, r.Tempo, r.Bars, r.Division)
}
