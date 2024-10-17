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

func (hp *hostPath) getAllocatedBlockMetadata(ctx context.Context, filePath string, startingOffset, blockSize int64, maxResult int32, allocBlocksChan chan<- []*csi.BlockMetadata) error {
	klog.V(4).Infof("finding allocated blocks in the file: %s", filePath)
	defer close(allocBlocksChan)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Seek(startingOffset, 0); err != nil {
		return err
	}

	return hp.compareBlocks(ctx, nil, file, startingOffset, blockSize, maxResult, allocBlocksChan)
}

func (hp *hostPath) getChangedBlockMetadata(ctx context.Context, sourcePath, targetPath string, startingOffset, blockSize int64, maxResult int32, changedBlocksChan chan<- []*csi.BlockMetadata) error {
	klog.V(4).Infof("finding changed blocks between two files: %s, %s", sourcePath, targetPath)
	defer close(changedBlocksChan)

	source, target, err := openFiles(sourcePath, targetPath)
	if err != nil {
		return err
	}
	defer source.Close()
	defer target.Close()

	if err := seekToOffset(source, target, startingOffset); err != nil {
		return err
	}

	return hp.compareBlocks(ctx, source, target, startingOffset, blockSize, maxResult, changedBlocksChan)
}

func openFiles(sourcePath, targetPath string) (source, target *os.File, err error) {
	source, err = os.Open(sourcePath)
	if err != nil {
		return nil, nil, err
	}

	target, err = os.Open(targetPath)
	if err != nil {
		source.Close()
		return nil, nil, err
	}

	return source, target, nil
}

func seekToOffset(source, target *os.File, startingOffset int64) error {
	if _, err := source.Seek(startingOffset, 0); err != nil {
		return err
	}
	if _, err := target.Seek(startingOffset, 0); err != nil {
		return err
	}
	return nil
}

// Compare blocks from source and target, and send changed blocks to channel.
// If source if nil, returns blocks allocated by target.
func (hp *hostPath) compareBlocks(ctx context.Context, source, target *os.File, startingOffset, blockSize int64, maxResult int32, changedBlocksChan chan<- []*csi.BlockMetadata) error {
	blockIndex := startingOffset / blockSize
	sBuffer := make([]byte, blockSize)
	tBuffer := make([]byte, blockSize)
	eofSourceFile, eofTargetFile := false, false

	for {
		changedBlocks := []*csi.BlockMetadata{}

		// Read blocks and compare them. Create the list of changed blocks metadata.
		// Once the number of blocks reaches to maxResult, return the result and
		// compute next batch of blocks.
		for int32(len(changedBlocks)) < maxResult {
			select {
			case <-ctx.Done():
				klog.V(4).Infof("handling cancellation signal, closing goroutine")
				return nil
			default:
				targetReadBytes, eofTarget, err := readFileBlock(target, tBuffer, eofTargetFile)
				if err != nil {
					return err
				}
				eofTargetFile = eofTarget

				if source == nil {
					// If source is nil, return blocks allocated by target file.
					if eofTargetFile {
						if len(changedBlocks) != 0 {
							changedBlocksChan <- changedBlocks
						}
						return nil
					}
					changedBlocks = append(changedBlocks, createBlockMetadata(blockIndex, blockSize))
					blockIndex++
					continue
				}

				sourceReadBytes, eofSource, err := readFileBlock(source, sBuffer, eofSourceFile)
				if err != nil {
					return err
				}
				eofSourceFile = eofSource

				// If both files have reached EOF, exit the loop.
				if eofSourceFile && eofTargetFile {
					klog.V(4).Infof("reached end of the files")
					if len(changedBlocks) != 0 {
						changedBlocksChan <- changedBlocks
					}
					return nil
				}

				// Compare the two blocks and add result.
				// Even if one of the file reaches to end, continue to add block metadata of other file.
				if blockChanged(sBuffer[:sourceReadBytes], tBuffer[:targetReadBytes]) {
					// TODO: Support for VARIABLE sized block metadata
					changedBlocks = append(changedBlocks, createBlockMetadata(blockIndex, blockSize))
				}

				blockIndex++
			}
		}

		if len(changedBlocks) > 0 {
			changedBlocksChan <- changedBlocks
		}
	}
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

func blockChanged(sourceBlock, targetBlock []byte) bool {
	return !bytes.Equal(sourceBlock, targetBlock)
}

func createBlockMetadata(blockIndex, blockSize int64) *csi.BlockMetadata {
	return &csi.BlockMetadata{
		ByteOffset: blockIndex * blockSize,
		SizeBytes:  blockSize,
	}
}
