package v1

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
)

func listFiles(directoryPath string) ([]string, error) {
	var files []string

	fileInfos, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		return nil, err
	}

	for _, fileInfo := range fileInfos {
		//iles = append(files, fileInfo.Name())
		files = append(files, directoryPath+fileInfo.Name())
	}

	return files, nil
}

func listFilesRecursive(rootPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		//fmt.Println(info.Name())
		// Skip directories
		if info.IsDir() {
			return nil
		}
		// Add the file path to the slice
		files = append(files, path)
		return nil
	})

	return files, err
}

func TestReadingLocalFiles(t *testing.T) {
	var (
		dir = "/Users/progers/baddat/loki_dev_006_index_19731/29/blooms/"
		//dir = "/data/blooms.old/bloom/loki_dev_006_index_19731/29/blooms/"
	)
	files, _ := listFilesRecursive(dir)
	for _, file := range files {
		cmd := exec.Command("mkdir", "/tmp/foo")
		_ = cmd.Run()
		fmt.Println(file)
		file, _ := os.Open(file)
		defer file.Close()
		reader := bufio.NewReader(file)
		UnTarGz("/tmp/foo", reader)
		r := NewDirectoryBlockReader("/tmp/foo")
		err := r.Init()
		require.NoError(t, err)

		_, err = r.Index()
		require.NoError(t, err)

		_, err = r.Blooms()
		require.NoError(t, err)

		block := NewBlock(r)
		blockQuerier := NewBlockQuerier(block)
		blockIters := NewPeekingIter[*SeriesWithBloom](blockQuerier)
		for blockIters.Next() {
			_ = blockIters.At()
		}
		blockQuerier.blooms.Next()

		cmd = exec.Command("rm", "-rf", "/tmp/foo")
		_ = cmd.Run()
	}
}

func TestReadingAllLocalFiles(t *testing.T) {
	var (
		dir = "/Users/progers/baddat/loki_dev_006_index_19731/29/blooms/"
	)
	cmd := exec.Command("mkdir", "/tmp/foo")
	_ = cmd.Run()
	files, _ := listFilesRecursive(dir)
	blockIters := make([]PeekingIterator[*SeriesWithBloom], len(files))
	for i, file := range files {
		tmpDirI := "/tmp/foo/" + strconv.Itoa(i)
		cmd := exec.Command("mkdir", tmpDirI)
		_ = cmd.Run()
		fmt.Println(file)
		file, _ := os.Open(file)
		defer file.Close()
		reader := bufio.NewReader(file)
		UnTarGz(tmpDirI, reader)
		r := NewDirectoryBlockReader(tmpDirI)
		err := r.Init()
		require.NoError(t, err)

		_, err = r.Index()
		require.NoError(t, err)

		_, err = r.Blooms()
		require.NoError(t, err)

		block := NewBlock(r)
		blockQuerier := NewBlockQuerier(block)
		blockIters[i] = NewPeekingIter[*SeriesWithBloom](blockQuerier)

	}
	heap := NewHeapIterForSeriesWithBloom(blockIters...)
	fmt.Printf("made heap iterator\n")
	_ = heap.Next()
	fmt.Println("Got here")
	cmd = exec.Command("rm", "-rf", "/tmp/foo")
	_ = cmd.Run()
}
