/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./app/annotator/*.go",
            "./app/annotator/templates/*.html",
            "./domain/images/list_viewer.go",
            "./domain/locations/view_all.go",
            "./domain/annotation_profiles/view_all.go",
            "./generic/*.go"],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/typography')
  ],
}
