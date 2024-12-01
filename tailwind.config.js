/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["web/**/*.{tmpl,js}"],
  theme: {
    extend: {
      fontFamily: {
        sans: ["Open Sans"],
        heading: ["Raleway"],
      },
    },
  },
  plugins: [],
};
