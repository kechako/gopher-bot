package plugin

import (
	"bufio"
	"io"
	"strings"
	"sync"

	"github.com/mattn/go-runewidth"
)

var defaultIndentSpacer = indentSpacer{}

type indentSpacer struct {
	buf []byte
	mux sync.RWMutex
}

func (s *indentSpacer) WriteTo(w io.Writer, count int) (n int, err error) {
	if s.count() < count {
		s.grow(count)
	}

	s.mux.RLock()
	n, err = w.Write(s.buf[:count])
	s.mux.RUnlock()

	return
}

func (s *indentSpacer) WriteToString(w io.StringWriter, count int) (n int, err error) {
	if s.count() < count {
		s.grow(count)
	}

	s.mux.RLock()
	n, err = w.WriteString(string(s.buf[:count]))
	s.mux.RUnlock()

	return
}

func (s *indentSpacer) count() int {
	s.mux.RLock()
	count := len(s.buf)
	s.mux.RUnlock()
	return count
}

func (s *indentSpacer) grow(count int) {
	s.mux.Lock()
	if len(s.buf) < count {
		s.buf = make([]byte, count)
		for i := 0; i < len(s.buf); i++ {
			s.buf[i] = byte(' ')
		}
	}
	s.mux.Unlock()
}

var DefaultHelpFormatter = &HelpFormatter{
	NameSuffix:    ": ",
	CommandSuffix: ": ",
	Indent:        4,
}

type HelpFormatter struct {
	NameSuffix    string
	CommandSuffix string
	Indent        int
}

func (f *HelpFormatter) Format(help *Help) string {
	var doc strings.Builder

	f.format(&doc, help)

	return doc.String()
}

func (f *HelpFormatter) format(w io.StringWriter, help *Help) {
	w.WriteString(help.Name)
	w.WriteString(f.NameSuffix)

	descIndent := runewidth.StringWidth(help.Name) + runewidth.StringWidth(f.NameSuffix)
	desc := bufio.NewScanner(strings.NewReader(help.Description))
	if desc.Scan() {
		// first line
		w.WriteString(desc.Text())
	}
	for desc.Scan() {
		w.WriteString("\n")
		writeSpaces(w, descIndent)
		w.WriteString(desc.Text())
	}

	var maxCmdWidth int
	cmdSuffixWidth := runewidth.StringWidth(f.CommandSuffix)
	for _, cmd := range help.Commands {
		cmdWidth := runewidth.StringWidth(cmd.Command) + cmdSuffixWidth
		if cmdWidth > maxCmdWidth {
			maxCmdWidth = cmdWidth
		}
	}

	cmdDescIndent := f.Indent + maxCmdWidth
	for _, cmd := range help.Commands {
		w.WriteString("\n")
		writeSpaces(w, f.Indent)
		w.WriteString(runewidth.FillRight(cmd.Command+f.CommandSuffix, maxCmdWidth))

		cmdDesc := bufio.NewScanner(strings.NewReader(cmd.Description))
		if cmdDesc.Scan() {
			// first line
			w.WriteString(cmdDesc.Text())
		}
		for cmdDesc.Scan() {
			w.WriteString("\n")
			writeSpaces(w, cmdDescIndent)
			w.WriteString(cmdDesc.Text())
		}
	}
}

func writeSpaces(w io.StringWriter, width int) {
	if ww, ok := w.(io.Writer); ok {
		defaultIndentSpacer.WriteTo(ww, width)
	} else {
		defaultIndentSpacer.WriteToString(w, width)
	}
}
