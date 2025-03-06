# shadcn/ui Integration Issues Report

## Summary
We encountered issues attempting to integrate shadcn/ui components into the ExpertDB frontend project. The installation process stalled during dependency installation despite proper configuration setup.

## Details of the Issue

### Error Description
- Initial error: "No import alias found in your tsconfig.json file"
- After fixing the import alias: Installation process stalls during the dependency installation phase
- Installation appears to detect React 19 and prompts for handling peer dependency issues (with options to use `--force` or `--legacy-peer-deps`), but doesn't proceed after selection

### Attempted Steps
1. Ran `npx shadcn@latest init` - Failed with import alias error
2. Updated tsconfig.json to include proper path aliasing
3. Created components.json configuration file
4. Created utils.ts for shadcn integration
5. Installed required dependencies (clsx, tailwind-merge)
6. Attempted to add components with `npx shadcn@latest add button input table select`
7. Tried individual component installation with `npx shadcn@latest add button`

## Root Cause Analysis

The issue appears to be related to React 19 compatibility with shadcn/ui. The shadcn CLI detects React 19 and warns about possible peer dependency issues. When prompted to use `--force` or `--legacy-peer-deps`, the installation process stalls instead of proceeding.

Additionally, the initial error with import alias detection appears to be because shadcn looks specifically at the root tsconfig.json rather than referenced configurations like tsconfig.app.json.

## Configuration Files

### tsconfig.json (Updated)
```json
{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "react-jsx",
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"]
    }
  },
  "include": ["src"],
  "references": [
    { "path": "./tsconfig.app.json" },
    { "path": "./tsconfig.node.json" }
  ]
}
```

### tsconfig.app.json
```json
{
  "compilerOptions": {
    "tsBuildInfoFile": "./node_modules/.tmp/tsconfig.app.tsbuildinfo",
    "target": "ES2020",
    "useDefineForClassFields": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,

    /* Bundler mode */
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "isolatedModules": true,
    "moduleDetection": "force",
    "noEmit": true,
    "jsx": "react-jsx",
    "baseUrl": ".",
    "paths": {
      "@/*": ["src/*"]
    },
    /* Linting */
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "noUncheckedSideEffectImports": true
  },
  "include": ["src"]
}
```

### vite.config.ts
```typescript
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { resolve } from "path";

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": resolve(__dirname, "./src"),
    },
  },
});
```

### components.json (Created for shadcn)
```json
{
  "$schema": "https://ui.shadcn.com/schema.json",
  "style": "default",
  "rsc": false,
  "tsx": true,
  "tailwind": {
    "config": "tailwind.config.js",
    "css": "src/index.css",
    "baseColor": "neutral",
    "cssVariables": true
  },
  "aliases": {
    "components": "@/components",
    "utils": "@/lib/utils"
  }
}
```

### src/lib/utils.ts (Created for shadcn)
```typescript
import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
```

## Recommendations

1. **Alternative Approach**: Consider implementing the Search functionality without relying on shadcn/ui initially. Use basic HTML elements with Tailwind classes as an interim solution.

2. **React Version**: Check if downgrading React from version 19 to 18 is feasible, as shadcn/ui might have better compatibility with React 18.

3. **Manual Component Creation**: Instead of using the shadcn CLI, consider manually copying the component source code from the shadcn website and adapting it to the project's needs.

4. **Package Manager Switch**: Try using Yarn or pnpm instead of npm, which might handle the peer dependencies differently.

5. **Build Tool Configuration**: Verify if there are any specific Vite configurations needed for React 19 + shadcn integration.

## Next Steps

The most pragmatic approach may be to implement the Search page functionality using standard HTML elements and Tailwind CSS classes for now. This allows us to make progress on the functional requirements while the shadcn/ui integration issues are resolved.

If component consistency is crucial, we should investigate downgrading React to version 18 or explore the direct installation of the underlying libraries that shadcn/ui wraps (like Radix UI).