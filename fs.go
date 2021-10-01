package discordfs

import (
	"bytes"
	"io/fs"

	dg "github.com/bwmarrin/discordgo"
)

// DiscStorage contains the necessary information to connect to the cloud channel
type DiscStorage struct {
	session   *dg.Session
	channelId string
}

// NewStorage builds a new DiscStorage
func NewStorage(s *dg.Session, channelId string) DiscStorage {
	return DiscStorage{
		session:   s,
		channelId: channelId,
	}
}

// DiscFile is a lazy file descriptor. It downloads the file it points
// to only when Read() is called
type DiscFile struct {
	storage DiscStorage
	info    FileInfo

	// to make the file lazy, it only downloads from the server when it's read
	contents      bytes.Buffer
	hasDownloaded bool
}

func (st DiscStorage) Open(filename string) (fs.File, error) {
	return st.GetFile(filename)
}

func (st DiscStorage) Stat(filename string) (fs.FileInfo, error) {
	file, err := st.GetFile(filename)
	if err != nil {
		return nil, err
	}
	return file.Stat()
}

func (st DiscStorage) ReadDir(dir string) ([]fs.DirEntry, error) {
	files, err := st.ListFiles(CleanPath(dir))
	if err != nil {
		return nil, err
	}

	output := make([]fs.DirEntry, len(files))
	for i, file := range files {
		output[i] = DiscFolder{file}
	}

	return output, nil
}

func (df DiscFile) ConcreteStat() (FileInfo, error) { return df.info, nil }
func (df DiscFile) Stat() (fs.FileInfo, error)      { return df.ConcreteStat() }
func (df *DiscFile) Close() error {
	df.contents = bytes.Buffer{}
	return nil
}

func (df *DiscFile) Read(input []byte) (int, error) {
	if !df.hasDownloaded {
		err := df.storage.Receive(&df.contents, df.info.name)
		if err != nil {
			return 0, err
		}

		df.hasDownloaded = true
	}

	// todo: download only parts of a file?
	return df.contents.Read(input)
}

type DiscFolder struct{ FileInfo }

func (df DiscFolder) Type() fs.FileMode          { return fs.ModeDir | df.Mode() }
func (df DiscFolder) Info() (fs.FileInfo, error) { return df, nil }
