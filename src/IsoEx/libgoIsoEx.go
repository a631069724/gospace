// +build linux,cgo windows,cgo

// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goIsoEx

/*
#cgo linux LDFLAGS: -fPIC -L${SRCDIR} -Wl,-rpath,${SRCDIR}  -lstdc++
#cgo linux CFLAGS: -fPIC -I${SRCDIR}
#cgo windows LDFLAGS: -fPIC -L${SRCDIR} -Wl,-rpath,${SRCDIR} ${SRCDIR}/libIsoEx.a
#cgo windows CFLAGS: -fPIC -I${SRCDIR} -DISLIB -DWIN32
*/
import "C"
