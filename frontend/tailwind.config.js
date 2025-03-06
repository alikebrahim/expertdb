/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.{js,ts,jsx,tsx,css}"],
  theme: {
    extend: {
      colors: {
        navy: "#133566",
        lightblue: "#1B4882",
        success: "#192012",
        warning: "#DC8335",
        error: "#FF4040",
        white: "#FFFFFF",
        gray: {
          50: "#F9FAFB",
          300: "#D1D5DB",
          400: "#9CA3AF",
          500: "#6B7280",
          600: "#4B5563",
          900: "#111827",
        },
      },
    },
  },
  plugins: [],
};
