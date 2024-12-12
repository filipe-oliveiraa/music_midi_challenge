package logging

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"text/template"
	"time"
)

// CyclicFileWriter implements the io.Writer interface and wraps an underlying file.
// It ensures that the file never grows over a limit.
type CyclicFileWriter struct {
	mu        sync.Mutex
	writer    *os.File
	liveLog   string
	nextWrite uint64
	limit     uint64
	logStart  time.Time
	maxLogAge time.Duration

	archiveFilename *template.Template
}

// MakeCyclicFileWriter returns a writer that wraps a file to ensure it never grows too large
func MakeCyclicFileWriter(
	liveLogFilePath string,
	archiveFilePath string,
	sizeLimitBytes uint64,
	maxLogAge time.Duration,
) *CyclicFileWriter {
	var err error
	cyclic := CyclicFileWriter{
		mu:              sync.Mutex{},
		logStart:        time.Time{},
		writer:          nil,
		liveLog:         liveLogFilePath,
		nextWrite:       0,
		limit:           sizeLimitBytes,
		maxLogAge:       maxLogAge,
		archiveFilename: template.New("archiveFilename"),
	}

	cyclic.archiveFilename, err = cyclic.archiveFilename.Parse(archiveFilePath)
	if err != nil {
		panic(fmt.Sprintf("bad LogArchiveName: %s", err))
	}
	cyclic.logStart = time.Now()

	fs, err := os.Stat(liveLogFilePath)
	if err == nil {
		cyclic.nextWrite = uint64(fs.Size()) //nolint: gosec
	}

	writer, err := os.OpenFile(liveLogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		panic(fmt.Sprintf("CyclicFileWriter: cannot open log file %v", err))
	}
	cyclic.writer = writer
	return &cyclic
}

type archiveFilenameTemplateData struct {
	Year      string
	Month     string
	Day       string
	Hour      string
	Minute    string
	Second    string
	EndYear   string
	EndMonth  string
	EndDay    string
	EndHour   string
	EndMinute string
	EndSecond string
}

func (cyclic *CyclicFileWriter) getArchiveFilename(now time.Time) string {
	buf := strings.Builder{}
	_ = cyclic.archiveFilename.Execute(&buf, archiveFilenameTemplateData{
		fmt.Sprintf("%04d", cyclic.logStart.Year()),
		fmt.Sprintf("%02d", cyclic.logStart.Month()),
		fmt.Sprintf("%02d", cyclic.logStart.Day()),
		fmt.Sprintf("%02d", cyclic.logStart.Hour()),
		fmt.Sprintf("%02d", cyclic.logStart.Minute()),
		fmt.Sprintf("%02d", cyclic.logStart.Second()),
		fmt.Sprintf("%04d", now.Year()),
		fmt.Sprintf("%02d", now.Month()),
		fmt.Sprintf("%02d", now.Day()),
		fmt.Sprintf("%02d", now.Hour()),
		fmt.Sprintf("%02d", now.Minute()),
		fmt.Sprintf("%02d", now.Second()),
	})
	return buf.String()
}

func (cyclic *CyclicFileWriter) getArchiveGlob() string {
	buf := strings.Builder{}
	_ = cyclic.archiveFilename.Execute(&buf, archiveFilenameTemplateData{
		"*", "*", "*", "*", "*", "*",
		"*", "*", "*", "*", "*", "*",
	})
	return buf.String()
}

func procWait(cmd *exec.Cmd, cause string) {
	err := cmd.Wait()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", cause, err)
	}
}

// Write ensures the the underlying file can store an additional len(p) bytes.
// If there is not enough room left it seeks
// to the beginning of the file.
func (cyclic *CyclicFileWriter) Write(p []byte) (n int, err error) {
	cyclic.mu.Lock()
	defer cyclic.mu.Unlock()

	if uint64(len(p)) > cyclic.limit {
		return 0, dumpLongLine(p)
	}

	if cyclic.nextWrite+uint64(len(p)) > cyclic.limit {
		now := time.Now()
		// we don't have enough space to write the entry, so archive data
		cyclic.writer.Close()
		globPath := cyclic.getArchiveGlob()
		oldarchives, err := filepath.Glob(globPath)
		if err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%s: glob err: %s\n", globPath, err)
		} else if cyclic.maxLogAge != 0 {
			tooOld := now.Add(-cyclic.maxLogAge)
			for _, path := range oldarchives {
				finfo, err := os.Stat(path)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s: stat: %s\n", path, err)
					continue
				}
				if finfo.ModTime().Before(tooOld) {
					err = os.Remove(path)
					if err != nil {
						fmt.Fprintf(os.Stderr, "%s: rm: %s\n", path, err)
					}
				}
			}
		}
		archivePath := cyclic.getArchiveFilename(now)
		shouldGz := false
		shouldBz2 := false
		if strings.HasSuffix(archivePath, ".gz") {
			shouldGz = true
			archivePath = archivePath[:len(archivePath)-3]
		} else if strings.HasSuffix(archivePath, ".bz2") {
			shouldBz2 = true
			archivePath = archivePath[:len(archivePath)-4]
		}
		if err := os.Rename(cyclic.liveLog, archivePath); err != nil {
			panic(fmt.Sprintf("CyclicFileWriter: cannot archive full log %v", err))
		}
		if shouldGz {
			zipGz(archivePath)
		} else if shouldBz2 {
			zipBz2(archivePath)
		}
		cyclic.logStart = now
		cyclic.writer, err = os.OpenFile(cyclic.liveLog, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o666)
		if err != nil {
			panic(fmt.Sprintf("CyclicFileWriter: cannot open log file %v", err))
		}
		cyclic.nextWrite = 0
	}
	// write the data
	n, err = cyclic.writer.Write(p)
	cyclic.nextWrite += uint64(n) //nolint: gosec
	return
}

func zipGz(archivePath string) {
	cmd := exec.Command("gzip", archivePath)
	err := cmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: could not gzip: %s", archivePath, err)
	} else {
		go procWait(cmd, archivePath)
	}
}

func zipBz2(archivePath string) {
	cmd := exec.Command("bzip2", archivePath)
	err := cmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: could not bzip2: %s", archivePath, err)
	} else {
		go procWait(cmd, archivePath)
	}
}

func dumpLongLine(p []byte) error {
	// there's no hope for writing this entry to the log
	// for the large lines this is a clear indication something does wrong, dump data into stderr
	const minDebugLogLineSize = 10 * 1024 * 1024
	if len(p) >= minDebugLogLineSize {
		buf := make([]byte, 16*1024)
		stlen := runtime.Stack(buf, false)
		fmt.Fprintf(os.Stderr, "Attempt to write a large log line:\n%s\n", string(buf[:stlen]))
		fmt.Fprintf(os.Stderr, "The offending line:\n%s\n", string(p[:4096]))
	}

	return fmt.Errorf("CyclicFileWriter: input too long to write. Len = %v", len(p))
}

func (c *CyclicFileWriter) Sync() error {
	return c.writer.Sync()
}
