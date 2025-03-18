package lg

type Pipe struct {
	Ser Serializer
	Wri Writer
}

func (p *Pipe) Handle(msg Message) {
	_ = p.Wri.Write(p.Ser.Serialize(msg))
}

func NewPipe(ser Serializer, wri Writer) *Pipe {
	return &Pipe{
		Ser: ser,
		Wri: wri,
	}
}

func DefaultConsolePipe(disableColors bool) *Pipe {
	return NewPipe(ConsoleSerializer{disableColors}, ConsoleWriter{})
}
