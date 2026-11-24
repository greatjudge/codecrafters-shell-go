package liner

type Reader struct {
	readchan chan string
}

func (r *Reader) ReadLine() string {
	line := r.readchan
	return line
}
