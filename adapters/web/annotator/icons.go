package annotator

import (
	"fmt"
)

var TrashIcon = `
<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
  <path stroke-linecap="round" stroke-linejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
</svg>
`

var BadgeIcon = `<span class="w-fit inline-flex mx-1 my-1 overflow-hidden rounded-radius border border-secondary bg-surface text-xs font-medium text-secondary dark:border-secondary-dark dark:bg-surface-dark dark:text-secondary-dark">
    <span class="flex items-center gap-1 bg-secondary/10 px-2 py-1 dark:bg-secondary-dark/10">
		<a href="#" onclick="%v">
			<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" aria-hidden="true" stroke="currentColor" fill="none" stroke-width="1.4" class="size-4">
				<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
			</svg>
		</a>
		%v
    </span>
</span>`

var AddIcon = `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
  <line x1="12" y1="4" x2="12" y2="20"/>
  <line x1="4" y1="12" x2="20" y2="12"/>
</svg>`

var RectangleSelectorIcon = `<svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-square" aria-hidden="true">
	<rect width="18" height="18" x="3" y="3" rx="2"></rect>
</svg>`

var PolygonSelectorIcon = `<svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-triangle-right" aria-hidden="true">
<path d="M22 18a2 2 0 0 1-2 2H3c-1.1 0-1.3-.6-.4-1.3L20.4 4.3c.9-.7 1.6-.4 1.6.7Z">
</path>
</svg>`

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

var ColoredRectangleIcon = `<svg width="22" height="22" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <rect x="10" y="10" width="80" height="80" rx="10" ry="10" fill="%v" />
</svg>`
