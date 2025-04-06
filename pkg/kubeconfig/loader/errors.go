package loader

import "errors"

var ErrNoKubeDir = errors.New("cannot find .kube directory")
var ErrReading = errors.New("read error")
var ErrParsing = errors.New("parse error")
