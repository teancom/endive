package endive

// GenericBook interface for Books
type GenericBook interface {
	FullPath() string
	ShortString() string
}

/*
// GenericBooks interface for slices of Books
type GenericBooks interface {
	FindByID(string) (*GenericBook, error)
	FindByHash(string) (*GenericBook, error)
	FindByMetadata(string) (*GenericBook, error)
	FindByFullPath(string) (*GenericBook, error)
	Books() []GenericBook
}
*/

// Indexer provides an interface ofr indexing books.
type Indexer interface {
	SetPath(path string)
	Rebuild(all []GenericBook) error
	Update(new map[string]GenericBook, mod map[string]GenericBook, del map[string]GenericBook) error
	Query(query string) ([]string, error)
	Count() uint64
}

// TODO: Check(all []GenericBook) error : checks all books have an entry + no extraneous ones

/*
// Database is the interface for loading/saving Book information
type Database interface {
	SetPath(string)
	Add(*GenericBook) error
	Remove(id int) error
	Save() error
	Load() error
	Check() error
}
*/

// UserInterface deals with user input, output and logging.
type UserInterface interface {
	// input
	YesOrNo(string) bool
	Choose(string, string, string, string) (string, error)
	UpdateValues(string, string, []string) ([]string, error)
	// output
	Title(string, ...interface{})
	SubTitle(string, ...interface{})
	SubPart(string, ...interface{})
	Choice(string, ...interface{})
	Display(string)
	// log
	InitLogger(string) error
	CloseLog()
	Error(string)
	Errorf(string, ...interface{})
	Warning(string)
	Warningf(string, ...interface{})
	Info(string)
	Infof(string, ...interface{})
	Debug(string)
	Debugf(string, ...interface{})
}