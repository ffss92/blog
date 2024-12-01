/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["web/**/*.{tmpl,js}"],
  theme: {
    extend: {
      fontFamily: {
        sans: ["Raleway"],
        heading: ["Raleway"],
      },
    },
  },
  plugins: [require("@tailwindcss/typography")],
};
