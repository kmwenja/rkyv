package rkyv

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
)

// MAGIC is the identifier for an .rkyv file
var MAGIC = []byte("RKYV")

// Rkyv is a helper struct for constructing, interogating
// and editing .rkyv files
type Rkyv struct {
	Version     string      `json:"version"`
	UUID        string      `json:"uuid"`
	DateCreated time.Time   `json:"date_created"`
	DateUpdated time.Time   `json:"date_updated"`
	Files       []*rkyvFile `json:"files"`
	Tags        []string    `json:"tags"`
	SearchTerms []string    `json:"search"`
}

type rkyvFile struct {
	Name string `json:"name"`
	data []byte
	Type string `json:"type"`
	Size int    `json:"size"`
	Hash []byte `json:"hash"`
}

// OpenFile opens a .rkyv file and populates a Rkyv struct
func OpenFile(filename string) (*Rkyv, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file `%s`: %v", filename, err)
	}
	defer f.Close()

	// make sure it's an rkyv file (check for RKYV)
	readBuf := make([]byte, len(MAGIC))
	n, err := f.Read(readBuf)
	if err != nil {
		return nil, fmt.Errorf("could not read file for RKYV `%s`: %v", filename, err)
	}
	if n != len(MAGIC) {
		return nil, fmt.Errorf("could not read file for RKYV `%s`: %d bytes read", filename, n)
	}

	// read the meta pointer
	ptrBuf := make([]byte, 8)
	n, err = f.Read(ptrBuf)
	if err != nil {
		return nil, fmt.Errorf("could not read file for meta pointer `%s`: %v", filename, err)
	}
	if n != len(ptrBuf) {
		return nil, fmt.Errorf("could not read file for meta pointer `%s`: %d bytes read", filename, n)
	}
	ptr := binary.LittleEndian.Uint64(ptrBuf)

	// read meta
	_, err = f.Seek(int64(ptr), 1)
	if err != nil {
		return nil, fmt.Errorf("could not seek to meta `%s`: %v", filename, err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("could not read file for meta `%s`: %v", filename, err)
	}
	var r Rkyv
	err = json.Unmarshal(data, &r)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal meta json `%s`: %v", filename, err)
	}

	return &r, nil
}

// AddFile adds the filename and its data to the list
// of files in the .rkyv file. It also performs some
// prep work i.e. determine filetype, size and hash
func (r *Rkyv) AddFile(name string, data []byte) {
	ct := http.DetectContentType(data)

	hasher := sha256.New()
	hasher.Write(data)
	hash := hasher.Sum(nil)

	f := &rkyvFile{
		Name: name,
		data: data,
		Type: ct,
		Size: len(data),
		Hash: hash,
	}

	r.Files = append(r.Files, f)
}

// Filename returns the .rkyv filename
func (r *Rkyv) Filename() string {
	return fmt.Sprintf("%s.rkyv", r.UUID)
}

// Flush flushes the contents of the struct into the .rkyv file
func (r *Rkyv) Flush() error {
	// prep meta
	if r.Version == "" {
		r.Version = "1.0"
	}

	if r.UUID == "" {
		r.UUID = uuid.NewV4().String()
	}

	nilTime := time.Time{}
	if r.DateCreated == nilTime {
		r.DateCreated = time.Now()
	}
	r.DateUpdated = time.Now()

	// open file for writing
	f, err := os.OpenFile(r.Filename(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0775)
	if err != nil {
		return fmt.Errorf("could not open rkyv file to write `%s`: %v", r.Filename(), err)
	}
	defer f.Close()

	// write an identifying string first
	n, err := f.Write(MAGIC)
	if err != nil {
		return fmt.Errorf("could not write RKYV to rkyv file: %v", err)
	}
	if n != len(MAGIC) {
		return fmt.Errorf("could not write RKYV to rkyv file: %d bytes written", n)
	}

	// write the meta pointer
	var ptr uint64
	for _, rf := range r.Files {
		ptr += uint64(rf.Size)
	}
	ptrBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(ptrBuf, ptr)

	n, err = f.Write(ptrBuf)
	if err != nil {
		return fmt.Errorf("could not write pointer to rkyv file: %v", err)
	}
	if n != len(ptrBuf) {
		return fmt.Errorf("could not write pointer to rkyv file: %d bytes written", n)
	}

	// write the files next
	for _, rf := range r.Files {
		n, err := f.Write(rf.data)
		if err != nil {
			return fmt.Errorf("could not write file `%s` to rkyv file: %v", rf.Name, err)
		}
		if n != rf.Size {
			return fmt.Errorf("could not write file `%s` to rkyv file: %d bytes written", rf.Name, n)
		}
	}

	// write meta next
	data, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("could not marshal meta to json: %v", err)
	}

	n, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("could not write meta to rkyv file: %v", err)
	}
	if n != len(data) {
		return fmt.Errorf("could not write meta to rkyv file: %d bytes written", n)
	}

	return nil
}
