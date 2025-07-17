/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // BQA Design System Colors
        primary: {
          DEFAULT: '#397b26', // BQA Green - primary brand color, approvals, positive actions
          light: '#4a9632',
          dark: '#2d5f1e',
        },
        secondary: {
          DEFAULT: '#1c4679', // Deep Blue - navigation, headers, authority elements
          light: '#2557a0',
          dark: '#14335c',
        },
        accent: {
          DEFAULT: '#e64125', // Accent Red - call-to-action, warnings, rejections
          light: '#ff6b47',
          dark: '#c7341f',
        },
        highlight: {
          DEFAULT: '#e68835', // Orange - secondary actions and highlights
          light: '#ff9d52',
          dark: '#cc7429',
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
        sans: ['Graphik', 'Segoe UI', 'Helvetica Neue', 'Arial', 'sans-serif'],
      },
      fontSize: {
        'h1': '42px',
        'h2': '38px',
        'h3': '35px',
        'h4': '26px',
        'body': '16px',
      },
      screens: {
        'xs': '576px',
        'sm': '576px',
        'md': '768px',
        'lg': '992px',
        'xl': '1200px',
      },
      maxWidth: {
        'container': '1212px',
      },
    },
  },
  plugins: [],
}