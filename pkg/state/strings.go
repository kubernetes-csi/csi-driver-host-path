/*
Copyright 2021 The Kubernetes Authors.

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

package state

// Strings is an ordered set of strings with helper functions for
// adding, searching and removing entries.
type Strings []string

// Add appends at the end.
func (s *Strings) Add(str string) {
	*s = append(*s, str)
}

// Has checks whether the string is already present.
func (s *Strings) Has(str string) bool {
	for _, str2 := range *s {
		if str == str2 {
			return true
		}
	}
	return false
}

// Empty returns true if the list is empty.
func (s *Strings) Empty() bool {
	return len(*s) == 0
}

// Remove removes the first occurence of the string, if present.
func (s *Strings) Remove(str string) {
	for i, str2 := range *s {
		if str == str2 {
			*s = append((*s)[:i], (*s)[i+1:]...)
			return
		}
	}
}
