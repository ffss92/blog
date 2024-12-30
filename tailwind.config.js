/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["web/**/*.{tmpl,js}"],
  theme: {
    extend: {
      fontFamily: {
        sans: ["Open Sans"],
        heading: ["Raleway"],
        mono: ["JetBrainsMono"],
      },
      minHeight: {
        header: "calc(100dvh - 5rem)",
      },
    },
  },
  plugins: [require("@tailwindcss/typography")],
};
