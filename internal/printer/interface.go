package printer

type Printer interface {
	Add(item any)
	Close()
}
