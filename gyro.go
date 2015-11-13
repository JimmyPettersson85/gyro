package gyro

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

const (
	flags            = os.O_WRONLY | os.O_APPEND | os.O_CREATE
	fileMode         = 0644
	defaultLayout    = "2006-01-02T15"
	defaultExtension = "log"
)

// Logger holds all the state for rotating logs
type Logger struct {
	// path to where files are written
	path string

	// placeholders for formatting filenames
	prefix    string
	suffix    string
	separator string
	layout    string
	extension string

	// what time to use when formatting filenames
	timeFn func() time.Time

	// holds the pre-formatted string for filenames
	format string

	// protects against concurrent writes
	mu *sync.Mutex
}

// New returns a new rotating logger with the default values.
func New(path string) (*Logger, error) {
	l := &Logger{
		path:      path,
		mu:        &sync.Mutex{},
		layout:    defaultLayout,
		extension: defaultExtension,
		timeFn:    func() time.Time { return time.Now().UTC() },
	}

	// Make sure we have write permissions in l.path
	if err := l.canWrite(); err != nil {
		return nil, err
	}

	l.buildFormatString()

	return l, nil
}

// SetPrefix sets filename prefix
func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
	l.buildFormatString()
}

// SetSuffix sets filename suffix
func (l *Logger) SetSuffix(suffix string) {
	l.suffix = suffix
	l.buildFormatString()
}

// SetSeparator sets the field separator
func (l *Logger) SetSeparator(sep string) {
	l.separator = sep
	l.buildFormatString()
}

// SetLayout sets the layout passed to time.Format()
// This controls how often the logs are rotated depending
// on the highest resolution element in the layout string.
func (l *Logger) SetLayout(layout string) {
	l.layout = layout
}

// SetExtension sets the filename extension
func (l *Logger) SetExtension(ext string) {
	l.extension = ext
	l.buildFormatString()
}

// SetTimeFunction sets the function that returns the time
// passed to time.Format()
func (l *Logger) SetTimeFunction(f func() time.Time) {
	l.timeFn = f
}

// FileName returns the current filename
func (l *Logger) FileName() string {
	return fmt.Sprintf(l.format, l.timeFn().Format(l.layout))
}

// Write writes the byte data to the log file
func (l *Logger) Write(data []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	f, err := os.OpenFile(path.Join(l.path, l.FileName()), flags, fileMode)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	n, err := f.Write(data)
	if n != len(data) {
		return n, fmt.Errorf("Didnt write all data. Wrote %d out of %d bytes", n, len(data))
	}

	return n, err
}

// WriteString writes the data string to the log file
func (l *Logger) WriteString(data string) error {
	_, err := l.Write([]byte(data))
	return err
}

// String returns a debug string of the logger
func (l *Logger) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("gyro.Logger\n")
	buffer.WriteString(fmt.Sprintf("  path: %s\n", l.path))
	buffer.WriteString(fmt.Sprintf("  prefix: %q\n", l.prefix))
	buffer.WriteString(fmt.Sprintf("  suffix: %q\n", l.suffix))
	buffer.WriteString(fmt.Sprintf("  separator: %q\n", l.separator))
	buffer.WriteString(fmt.Sprintf("  extension: %q\n", l.extension))
	buffer.WriteString(fmt.Sprintf("  layout: %s\n", l.layout))
	buffer.WriteString(fmt.Sprintf("  format: %s\n", l.format))
	buffer.WriteString(fmt.Sprintf("  current filename: %s\n", l.FileName()))

	return strings.TrimSpace(buffer.String())
}

// canWrite tests if we have permissions to create/write files in l.path
func (l *Logger) canWrite() error {
	f, err := ioutil.TempFile(l.path, "gyro.Logger.")
	if err != nil {
		return err
	}

	f.Close()
	return os.Remove(f.Name())
}

// buildFormatString precalculates the format string for the filenames
// so we dont have to do the string interpolation on each call
func (l *Logger) buildFormatString() {
	p, s := len(l.prefix) > 0, len(l.suffix) > 0

	if p && s {
		l.format = strings.TrimSpace(fmt.Sprintf("%s%s%s%s%s.%s", l.prefix, l.separator, "%s", l.separator, l.suffix, l.extension))
	} else if p && !s {
		l.format = strings.TrimSpace(fmt.Sprintf("%s%s%s.%s", l.prefix, l.separator, "%s", l.extension))
	} else if !p && s {
		l.format = strings.TrimSpace(fmt.Sprintf("%s%s%s.%s", "%s", l.separator, l.suffix, l.extension))
	} else {
		l.format = strings.TrimSpace(fmt.Sprintf("%s.%s", "%s", l.extension))
	}

	if len(l.extension) == 0 {
		l.format = l.format[:len(l.format)-1]
	}
}
