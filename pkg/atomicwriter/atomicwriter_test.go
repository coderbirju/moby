package atomicwriter

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

var testMode os.FileMode = 0o640

func init() {
	// Windows does not support full Linux file mode
	if runtime.GOOS == "windows" {
		testMode = 0o666
	}
}

func TestAtomicWriteToFile(t *testing.T) {
	tmpDir := t.TempDir()

	expected := []byte("barbaz")
	if err := WriteFile(filepath.Join(tmpDir, "foo"), expected, testMode); err != nil {
		t.Fatalf("Error writing to file: %v", err)
	}

	actual, err := os.ReadFile(filepath.Join(tmpDir, "foo"))
	if err != nil {
		t.Fatalf("Error reading from file: %v", err)
	}

	if !bytes.Equal(actual, expected) {
		t.Fatalf("Data mismatch, expected %q, got %q", expected, actual)
	}

	st, err := os.Stat(filepath.Join(tmpDir, "foo"))
	if err != nil {
		t.Fatalf("Error statting file: %v", err)
	}
	if expected := testMode; st.Mode() != expected {
		t.Fatalf("Mode mismatched, expected %o, got %o", expected, st.Mode())
	}
}

func TestAtomicWriteSetCommit(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.Mkdir(filepath.Join(tmpDir, "tmp"), 0o700); err != nil {
		t.Fatalf("Error creating tmp directory: %s", err)
	}

	targetDir := filepath.Join(tmpDir, "target")
	ws, err := NewWriteSet(filepath.Join(tmpDir, "tmp"))
	if err != nil {
		t.Fatalf("Error creating atomic write set: %s", err)
	}

	expected := []byte("barbaz")
	if err := ws.WriteFile("foo", expected, testMode); err != nil {
		t.Fatalf("Error writing to file: %v", err)
	}

	if _, err := os.ReadFile(filepath.Join(targetDir, "foo")); err == nil {
		t.Fatalf("Expected error reading file where should not exist")
	}

	if err := ws.Commit(targetDir); err != nil {
		t.Fatalf("Error committing file: %s", err)
	}

	actual, err := os.ReadFile(filepath.Join(targetDir, "foo"))
	if err != nil {
		t.Fatalf("Error reading from file: %v", err)
	}

	if !bytes.Equal(actual, expected) {
		t.Fatalf("Data mismatch, expected %q, got %q", expected, actual)
	}

	st, err := os.Stat(filepath.Join(targetDir, "foo"))
	if err != nil {
		t.Fatalf("Error statting file: %v", err)
	}
	if expected := testMode; st.Mode() != expected {
		t.Fatalf("Mode mismatched, expected %o, got %o", expected, st.Mode())
	}
}

func TestAtomicWriteSetCancel(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.Mkdir(filepath.Join(tmpDir, "tmp"), 0o700); err != nil {
		t.Fatalf("Error creating tmp directory: %s", err)
	}

	ws, err := NewWriteSet(filepath.Join(tmpDir, "tmp"))
	if err != nil {
		t.Fatalf("Error creating atomic write set: %s", err)
	}

	expected := []byte("barbaz")
	if err := ws.WriteFile("foo", expected, testMode); err != nil {
		t.Fatalf("Error writing to file: %v", err)
	}

	if err := ws.Cancel(); err != nil {
		t.Fatalf("Error committing file: %s", err)
	}

	if _, err := os.ReadFile(filepath.Join(tmpDir, "target", "foo")); err == nil {
		t.Fatalf("Expected error reading file where should not exist")
	} else if !os.IsNotExist(err) {
		t.Fatalf("Unexpected error reading file: %s", err)
	}
}
