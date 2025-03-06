/*
Copyright 2024 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package hostpath

import (
	"bytes"
	"io"
	"os"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"golang.org/x/net/context"
	"k8s.io/klog/v2"
)

// NOTE: This implementation of SnapshotMetadata service is used for demo and CI testing purpose only.
// This should not be used in production or as an example about how to write a real driver.

type fileBlockReader struct {
	base              *os.File
	target            *os.File
	offset            int64
	blockSize         int64
	blockMetadataType csi.BlockMetadataType
	maxResult         int32
}

func newFileBlockReader(
	basePath,
	targetPath string,
	startingOffset int64,
	blockSize int64,
	blockMetadataType csi.BlockMetadataType,
	maxResult int32,
) (*fileBlockReader, error) {
	base, target, err := openFiles(basePath, targetPath)
	if err != nil {
		return nil, err
	}

	return &fileBlockReader{
		base:              base,
		target:            target,
		offset:            startingOffset,
		blockSize:         blockSize,
		blockMetadataType: blockMetadataType,
		maxResult:         maxResult,
	}, nil
}

func (cb *fileBlockReader) seekToStartingOffset() error {
	if _, err := cb.target.Seek(cb.offset, io.SeekStart); err != nil {
		return err
	}
	if cb.base == nil {
		return nil
	}
	if _, err := cb.base.Seek(cb.offset, io.SeekStart); err != nil {
		return err
	}
	return nil
}

func (cb *fileBlockReader) Close() error {
	if cb.base != nil {
		if err := cb.base.Close(); err != nil {
			return err
		}
	}
	if cb.target != nil {
		if err := cb.target.Close(); err != nil {
			return err
		}
	}
	return nil
}

func openFiles(basePath, targetPath string) (base, target *os.File, err error) {
	target, err = os.Open(targetPath)
	if err != nil {
		return nil, nil, err
	}
	if basePath == "" {
		return nil, target, nil
	}
	base, err = os.Open(basePath)
	if err != nil {
		target.Close()
		return nil, nil, err
	}

	return base, target, nil
}

// getChangedBlockMetadata reads base and target files, compare block differences between them
// and returns list of changed block metadata. It reads all the blocks till it reaches EOF or size of changed block
// metadata list <= maxSize.
func (cb *fileBlockReader) getChangedBlockMetadata(ctx context.Context) ([]*csi.BlockMetadata, error) {
	if cb.base == nil {
		klog.V(4).Infof("finding allocated blocks by file: %s", cb.target.Name())
	} else {
		klog.V(4).Infof("finding changed blocks between two files: %s, %s", cb.base.Name(), cb.target.Name())
	}

	blockIndex := cb.offset / cb.blockSize
	sBuffer := make([]byte, cb.blockSize)
	tBuffer := make([]byte, cb.blockSize)
	zeroBlock := make([]byte, cb.blockSize)
	eofBaseFile, eofTargetFile := false, false

	changedBlocks := []*csi.BlockMetadata{}

	// Read blocks and compare them. Create the list of changed blocks metadata.
	// Once the number of blocks reaches to maxResult, return the result and
	// compute next batch of blocks.
	for int32(len(changedBlocks)) < cb.maxResult {
		select {
		case <-ctx.Done():
			klog.V(4).Infof("handling cancellation signal, closing goroutine")
			return nil, ctx.Err()
		default:
		}
		targetReadBytes, eofTarget, err := readFileBlock(cb.target, tBuffer, eofTargetFile)
		if err != nil {
			return nil, err
		}
		eofTargetFile = eofTarget

		// If base is nil, return blocks allocated by target file.
		if cb.base == nil {
			if eofTargetFile {
				return changedBlocks, io.EOF
			}
			// return only allocated blocks.
			if blockChanged(zeroBlock, tBuffer[:targetReadBytes]) {
				// if VARIABLE_LENGTH type is enabled, return blocks extend instead of individual blocks.
				blockMetadata := createBlockMetadata(blockIndex, cb.blockSize)
				if extendBlock(changedBlocks, csi.BlockMetadataType(cb.blockMetadataType), blockIndex, cb.blockSize) {
					changedBlocks[len(changedBlocks)-1].SizeBytes += cb.blockSize
					cb.offset += cb.blockSize
					blockIndex++
					continue
				}
				changedBlocks = append(changedBlocks, blockMetadata)
			}

			cb.offset += cb.blockSize
			blockIndex++
			continue
		}

		baseReadBytes, eofBase, err := readFileBlock(cb.base, sBuffer, eofBaseFile)
		if err != nil {
			return nil, err
		}
		eofBaseFile = eofBase

		// If both files have reached EOF, exit the loop.
		if eofBaseFile && eofTargetFile {
			klog.V(4).Infof("reached end of the files")
			return changedBlocks, io.EOF
		}

		// Compare the two blocks and add result.
		// Even if one of the file reaches to end, continue to add block metadata of other file.
		if blockChanged(sBuffer[:baseReadBytes], tBuffer[:targetReadBytes]) {
			blockMetadata := createBlockMetadata(blockIndex, cb.blockSize)
			// if VARIABLE_LEGTH type is enabled, check if blocks are adjacent,
			// extend the previous block if adjacent blocks found instead of adding new entry.
			if extendBlock(changedBlocks, csi.BlockMetadataType(cb.blockMetadataType), blockIndex, cb.blockSize) {
				changedBlocks[len(changedBlocks)-1].SizeBytes += cb.blockSize
				cb.offset += cb.blockSize
				blockIndex++
				continue
			}
			changedBlocks = append(changedBlocks, blockMetadata)
		}

		cb.offset += cb.blockSize
		blockIndex++
	}
	return changedBlocks, nil
}

// readFileBlock reads blocks from a file.
func readFileBlock(file *os.File, buffer []byte, eof bool) (int, bool, error) {
	if eof {
		return 0, true, nil
	}

	bytesRead, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return 0, false, err
	}

	return bytesRead, err == io.EOF, nil
}

func blockChanged(baseBlock, targetBlock []byte) bool {
	return !bytes.Equal(baseBlock, targetBlock)
}

func createBlockMetadata(blockIndex, blockSize int64) *csi.BlockMetadata {
	return &csi.BlockMetadata{
		ByteOffset: blockIndex * blockSize,
		SizeBytes:  blockSize,
	}
}

func extendBlock(changedBlocks []*csi.BlockMetadata, metadataType csi.BlockMetadataType, blockIndex, blockSize int64) bool {
	blockMetadata := createBlockMetadata(blockIndex, blockSize)
	// if VARIABLE_LEGTH type is enabled, check if blocks are adjacent,
	// extend the previous block if adjacent blocks found instead of adding new entry.
	if len(changedBlocks) < 1 {
		return false
	}
	lastBlock := changedBlocks[len(changedBlocks)-1]
	if blockMetadata.ByteOffset == lastBlock.ByteOffset+lastBlock.SizeBytes &&
		metadataType == csi.BlockMetadataType_VARIABLE_LENGTH {
		return true
	}
	return false
}
