// Command main renders terminal.svg: a green-phosphor terminal that boots and
// types a short session. Pure SMIL, no scripts, no web fonts, so it animates
// inside a GitHub README image. Regenerate: go run . > terminal.svg
package main

import (
	"fmt"
	"html"
	"strings"
)

type line struct {
	text  string
	typed bool // typed at the prompt, char by char; otherwise printed at once
}

var session = []line{
	{"FREDOS 1.0   640K OK   (C) 1986", false},
	{"", false},
	{"> whoami", true},
	{"fred :: systems programmer :: C . Go . Rust . asm", false},
	{"> uname -o", true},
	{"GNU/Linux", false},
	{"> cat ~/.plan", true},
	{"small binaries. fast hot loops. boring, correct software.", false},
	{"> ls ~/work", true},
	{"roxyapi/  sdk-go/  sdk-typescript/", false},
	{">", true},
}

const (
	width, height = 760, 340
	padX, padY    = 36, 54
	fontSize      = 16.0
	charW         = fontSize * 0.6 // monospace advance; pinned with textLength
	lineH         = 24.0
	typeRate      = 0.06 // seconds per typed character
	pause         = 0.45 // seconds between lines
)

func main() {
	var b strings.Builder
	p := func(f string, a ...any) { fmt.Fprintf(&b, f, a...) }

	p(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" role="img" aria-label="Fred, systems programmer: C, Go, Rust, assembly">`, width, height, width, height)
	p(`<defs>`)
	p(`<pattern id="scan" width="4" height="4" patternUnits="userSpaceOnUse"><rect width="4" height="2" fill="#000" opacity=".35"/></pattern>`)
	p(`<radialGradient id="vig" cx="50%%" cy="50%%" r="75%%"><stop offset="60%%" stop-color="#000" stop-opacity="0"/><stop offset="100%%" stop-color="#000" stop-opacity=".75"/></radialGradient>`)
	p(`<filter id="glow" x="-10%%" y="-10%%" width="120%%" height="120%%"><feGaussianBlur stdDeviation="1.6" result="b"/><feMerge><feMergeNode in="b"/><feMergeNode in="SourceGraphic"/></feMerge></filter>`)
	p(`</defs>`)
	p(`<rect width="%d" height="%d" rx="18" fill="#1a1a17"/>`, width, height)
	p(`<rect x="12" y="12" width="%d" height="%d" rx="12" fill="#041208"/>`, width-24, height-24)
	p(`<g filter="url(#glow)" fill="#41ff7d" font-family="'Courier New',Courier,ui-monospace,monospace" font-size="%g" xml:space="preserve">`, fontSize)
	p(`<animate attributeName="opacity" values="1;.93;1;.97;1" dur="0.18s" repeatCount="indefinite"/>`)

	t := 0.6 // boot delay
	var cursorX, cursorY float64
	for i, l := range session {
		y := padY + float64(i)*lineH
		n := len([]rune(l.text))
		w := float64(n) * charW
		id := fmt.Sprintf("c%d", i)
		if l.typed && n > 0 {
			// Discrete steps: one character per tick, like a real keystroke.
			var vals, keys []string
			for k := 0; k <= n; k++ {
				vals = append(vals, fmt.Sprintf("%.1f", float64(k)*charW))
				keys = append(keys, fmt.Sprintf("%.4f", float64(k)/float64(n+1)))
			}
			dur := float64(n+1) * typeRate
			p(`<clipPath id="%s"><rect x="%d" y="%.1f" height="%.1f" width="0"><animate attributeName="width" begin="%.2fs" dur="%.2fs" values="%s" keyTimes="%s" calcMode="discrete" fill="freeze"/></rect></clipPath>`,
				id, padX, y-fontSize, lineH, t, dur, strings.Join(vals, ";"), strings.Join(keys, ";"))
			t += dur
		} else {
			p(`<clipPath id="%s"><rect x="%d" y="%.1f" height="%.1f" width="0"><set attributeName="width" to="%.1f" begin="%.2fs" fill="freeze"/></rect></clipPath>`,
				id, padX, y-fontSize, lineH, w+charW, t)
		}
		if n > 0 {
			p(`<text x="%d" y="%.1f" textLength="%.1f" lengthAdjust="spacing" clip-path="url(#%s)">%s</text>`,
				padX, y, w, id, html.EscapeString(l.text))
		}
		cursorX, cursorY = padX+w+charW, y
		t += pause
	}
	// Block cursor parks after the last prompt and blinks forever.
	p(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" opacity="0"><set attributeName="opacity" to="1" begin="%.2fs" fill="freeze"/><animate attributeName="opacity" values="1;1;0;0" keyTimes="0;.5;.5;1" dur="1.06s" begin="%.2fs" repeatCount="indefinite"/></rect>`,
		cursorX-charW*0.5, cursorY-fontSize+2, charW, fontSize+2, t, t)
	p(`</g>`)
	p(`<rect x="12" y="12" width="%d" height="%d" rx="12" fill="url(#scan)"/>`, width-24, height-24)
	p(`<rect x="12" y="12" width="%d" height="%d" rx="12" fill="url(#vig)"/>`, width-24, height-24)
	p(`</svg>`)
	fmt.Println(b.String())
}
