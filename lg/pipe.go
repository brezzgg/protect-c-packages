package lg

import "fmt"

type (
	Pipe struct {
		ser Serializer
		wri Writer
	}

	PipeOption func(*Pipe)
)

func (p *Pipe) Handle(msg Message) {
	err := p.wri.Write(p.ser.Serialize(msg))
	if err != nil {
		fmt.Printf("lg: writer error: %s\n", err.Error())
	}
}

func NewPipe(opts ...PipeOption) *Pipe {
	p := &Pipe{
		wri: NewConsoleWriter(),
		ser: NewConsoleSerializer(),
	}
	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *Pipe) Flush() {
	p.wri.Flush()
}

func (p *Pipe) Close() {
	p.wri.Close()
}

func WithSerializer(ser Serializer) PipeOption {
	return func(p *Pipe) {
		p.ser = ser
	}
}

func WithWriter(writer Writer) PipeOption {
	return func(p *Pipe) {
		p.wri = writer
	}
}

func AsDefaultConsole() PipeOption {
	return func(p *Pipe) {
		p.ser = NewConsoleSerializer()
		p.wri = NewConsoleWriter()
	}
}

func AsDefaultFile(filename string) PipeOption {
	return func(p *Pipe) {
		p.ser = NewJSONSerializer()
		p.wri = NewFileWriter(filename)
	}
}
