/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["web/**/*.{tmpl,js}"],
  theme: {
    extend: {
      fontFamily: {
        sans: ["Raleway"],
        heading: ["Raleway"],
      },
      minHeight: {
        header: "calc(100dvh - 5rem)",
      },
    },
  },
  plugins: [require("@tailwindcss/typography")],
};
