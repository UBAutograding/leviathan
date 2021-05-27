package chanwriter

type chanWriter struct {
	ch chan string
}

func NewChanWriter(ch chan string) *chanWriter {
	return &chanWriter{ch}
}

func (w *chanWriter) Write(p []byte) (int, error) {
	w.ch <- string(p)
	return len(p), nil
}
