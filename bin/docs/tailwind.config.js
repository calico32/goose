const colors = require("tailwindcss/colors")

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.html"],
  theme: {
    extend: {
      colors: {
        primary: colors.sky,
      },
      fontFamily: {
        sans: [
          "Palanquin",
          "Inter",
          "IBM Plex Sans",
          "SF Sans",
          "system-ui",
          "-apple-system",
          "BlinkMacSystemFont",
          "Segoe UI",
          "Roboto",
          "Helvetica Neue",
          "Helvetica",
          "Arial",
          "sans-serif",
        ],
        mono: [
          "Monaspace Neon",
          "IBM Plex Mono",
          "SF Mono",
          "Fira Code",
          "Fira Mono",
          "Roboto Mono",
          "monospace",
        ],
      },
    },
  },
  plugins: [require("@tailwindcss/typography")],
}
