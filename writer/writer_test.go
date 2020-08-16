// File writer_test.go tests the writer.
package writer

import (
	"os"
	"path"
	"testing"

	"github.com/verless/verless/test"
)

const (
	testPath    = "../example"
	testOutPath = "../test-out-path"
)

func TestWriter_removeOutDirIfExists(t *testing.T) {

	tests := map[string]struct {
		// beforeTest is a callback which creates the folders/files
		// to test a specific testcase.
		beforeTest func()
		// cleanupTest is a callback which the folders/files created
		// by beforeTest and the test itself.
		cleanupTest   func()
		expectedError string
	}{
		"normal": {
			beforeTest:  func() {},
			cleanupTest: func() {},
		},
		"already exists": {
			beforeTest: func() {
				test.Ok(t, os.Mkdir(testOutPath, os.ModePerm))

				file, err := os.Create(path.Join(testOutPath, "anyFile.txt"))
				test.Ok(t, err)
				_ = file.Close()
			},
			cleanupTest: func() {
				err := os.RemoveAll(testOutPath)
				test.Ok(t, err)
			},
		},
		"already exists but without file": {
			beforeTest: func() {
				test.Ok(t, os.Mkdir(testOutPath, os.ModePerm))
			},
			cleanupTest: func() {
				err := os.RemoveAll(testOutPath)
				test.Ok(t, err)
			},
		},
	}

	for caseName, testCase := range tests {
		t.Logf("Testing '%s'", caseName)

		w := setupNewWriter(t)

		testCase.beforeTest()

		err := w.removeOutDirIfExists()

		if testCase.expectedError == "" {
			test.Ok(t, err)
		} else {
			test.Assert(t, err != nil && testCase.expectedError == err.Error(), "should error")
		}

		testCase.cleanupTest()
	}
}

func setupNewWriter(t testing.TB) *writer {
	w, err := New(testPath, testOutPath)
	if err != nil {
		t.Errorf("New should not throw an error: %v", err)
	}
	return w
}