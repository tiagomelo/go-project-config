// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package cfg

// lineReader defines an interface for reading lines of text.
type lineReader interface {
	Scan() bool
	Text() string
	Err() error
}
