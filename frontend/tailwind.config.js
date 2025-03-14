/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // BQA colors
        primary: {
          DEFAULT: '#003366', // Navy blue primary color
          light: '#0055a4',
          dark: '#00254d',
        },
        secondary: {
          DEFAULT: '#e63946', // Red accent color
          light: '#ff4d5e',
          dark: '#c62b38',
        },
        accent: {
          DEFAULT: '#f0f4f8', // Light blue-gray accent
          dark: '#d0d8e0',
        },
        neutral: {
          100: '#f8f9fa',
          200: '#e9ecef',
          300: '#dee2e6',
          400: '#ced4da',
          500: '#adb5bd',
          600: '#6c757d',
          700: '#495057',
          800: '#343a40',
          900: '#212529',
        }
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', '-apple-system', 'sans-serif'],
      },
    },
  },
  plugins: [],
}