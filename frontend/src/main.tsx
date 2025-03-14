import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'

// Set up API mocking in development
if (import.meta.env.DEV) {
  // We can implement a mock service worker if needed later
  console.log('Running in development mode');
}

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
