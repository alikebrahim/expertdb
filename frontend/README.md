# ExpertDB Frontend

This is the frontend application for the ExpertDB system, built with Next.js, TypeScript, and shadcn/ui.

## Features

- Expert database search with filtering capabilities
- Expert request submission form
- Expert detail view
- Responsive design for desktop and mobile

## Technology Stack

- **Framework**: [Next.js](https://nextjs.org/) with TypeScript and App Router
- **Build Tool**: [Vite](https://vitejs.dev/) for fast development and optimized builds
- **UI Library**: [shadcn/ui](https://ui.shadcn.com/) for accessible, customizable components
- **State Management**: React Context API and React Query for data fetching
- **Form Validation**: [React Hook Form](https://react-hook-form.com/) with [Zod](https://github.com/colinhacks/zod)
- **HTTP Client**: [Axios](https://axios-http.com/) for API requests

## Getting Started

### Prerequisites

- Node.js 18+ and npm/yarn

### Installation

1. Clone the repository
2. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```
3. Install dependencies:
   ```bash
   npm install
   # or
   yarn install
   ```

### Development

To start the development server:

```bash
npm run dev
# or
yarn dev
```

The application will be available at [http://localhost:3000](http://localhost:3000).

### Building for Production

To create a production build:

```bash
npm run build
# or
yarn build
```

To start the production server:

```bash
npm start
# or
yarn start
```

## Project Structure

```
frontend/
├── app/                 # Next.js App Router pages
│   ├── expert/[id]/     # Expert detail page
│   ├── request/         # Expert request submission page
│   ├── search/          # Expert search page
│   ├── globals.css      # Global styles
│   ├── layout.tsx       # Root layout component
│   └── page.tsx         # Home page
├── components/          # Reusable UI components
│   ├── layout/          # Layout components (navbar, footer)
│   └── ui/              # UI components from shadcn/ui
├── lib/                 # Utility functions and API services
│   ├── api.ts           # API service for backend communication
│   └── utils.ts         # Utility functions
├── public/              # Static files
└── README.md            # This file
```

## Backend API Integration

The frontend communicates with the backend API through the API service defined in `lib/api.ts`. The API endpoints are proxied through Next.js to avoid CORS issues.

## Key Features Implementation

### Expert Request Form

- Located at `/request`
- Uses React Hook Form with Zod validation
- Submits request to the backend API endpoint

### Expert Search

- Located at `/search`
- Allows searching by name, area, role, and other filters
- Uses client-side data fetching with API service
- Displays results in a card-based layout

## Contributing

1. Create a feature branch (`git checkout -b feature/amazing-feature`)
2. Commit your changes (`git commit -m 'Add some amazing feature'`)
3. Push to the branch (`git push origin feature/amazing-feature`)
4. Open a Pull Request