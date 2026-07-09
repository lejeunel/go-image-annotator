package icons

import (
	_ "embed"
	"fmt"
)

//go:embed svg/trash.svg
var TrashIcon string

//go:embed svg/edit.svg
var EditIcon string

//go:embed svg/sun.svg
var SunIcon string

//go:embed svg/moon.svg
var MoonIcon string

var GitHubIcon = `
<svg xmlns="http://www.w3.org/2000/svg" fill="currentColor" class="size-5" viewBox="0 0 16 16">
    <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27s1.36.09 2 .27c1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.01 8.01 0 0 0 16 8c0-4.42-3.58-8-8-8"></path>
</svg>
`

//go:embed svg/user.svg
var UserIcon string

//go:embed svg/add-circle.svg
var AddCircleIcon string

//go:embed svg/add.svg
var AddIcon string

//go:embed svg/box.svg
var BoundingBoxIcon string

//go:embed svg/polygon.svg
var PolygonIcon string

func MakeColoredRectangleIcon(color string) string {
	return fmt.Sprintf(`<svg width="22" height="22" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <rect x="10" y="10" width="80" height="80" rx="10" ry="10" fill="%v" />
</svg>`, color)
}

func MakeColoredHexagonIcon(color string) string {
	return fmt.Sprintf(`<svg width="22" height="22" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <polygon points="50,10 84.64,30 84.64,70 50,90 15.36,70 15.36,30" fill="%v" />
</svg>`, color)
}
