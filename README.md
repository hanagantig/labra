для работы goseract нужно отредактировать файлы
github.com/otiai10/gosseract/v2/preprocessflags_x.go

нужно чтобы содержимое было сдедующее

```go
package gosseract

// #cgo CXXFLAGS: -std=c++0x
// #cgo CPPFLAGS: -I/opt/homebrew/include
// #cgo CPPFLAGS: -Wno-unused-result
// #cgo LDFLAGS: -L/opt/homebrew/lib -lleptonica -ltesseract
import "C"
```

для openCV
brew install pkg-config opencv
