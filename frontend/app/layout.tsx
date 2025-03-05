import type { Metadata } from 'next'
import "./globals.css"
import { Inter, Rubik } from 'next/font/google'
import { Navbar } from '@/components/layout/navbar'
import { Footer } from '@/components/layout/footer'

// Using Rubik as a close approximation to Graphik
const rubik = Rubik({
  subsets: ['latin'],
  weight: ['300', '400', '500', '600', '700'],
  variable: '--font-rubik',
})

// Keep Inter as fallback
const inter = Inter({ 
  subsets: ['latin'],
  variable: '--font-inter',
})

export const metadata: Metadata = {
  title: 'BQA Expert Database',
  description: 'Bahrain Quality Assurance Authority Expert Database Management System',
  icons: {
    icon: '/images/logo/Icon Logo - Color.svg',
  },
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className={`${rubik.variable} ${inter.variable} font-sans`}>
        <div className="flex min-h-screen flex-col">
          <Navbar />
          <main className="flex-1">{children}</main>
          <Footer />
        </div>
      </body>
    </html>
  )
}