package pastebox

import "embed"

// FS содержит HTML-шаблоны и статические файлы, вкомпилированные в бинарник.
//
//go:embed templates static
var FS embed.FS
