# NOTES:
- After recent changes, there has been an error build_err below for details

# ERRORS:

- build_err:
Error: ./app/panel/page.tsx
Error:   [31m×[0m Expression expected
    ╭─[[36;1;4m/home/alikebrahim/dev/expertdb_new/frontend/app/panel/page.tsx[0m:65:1]
 [2m62[0m │   };
 [2m63[0m │   
 [2m64[0m │   return (
 [2m65[0m │     <>
    · [35;1m     ─[0m
 [2m66[0m │       <Navbar />
 [2m67[0m │       <div className="container py-10">
 [2m68[0m │         <div className="mb-6">
    ╰────
  [31m×[0m Expected ',', got 'className'
    ╭─[[36;1;4m/home/alikebrahim/dev/expertdb_new/frontend/app/panel/page.tsx[0m:67:1]
 [2m64[0m │   return (
 [2m65[0m │     <>
 [2m66[0m │       <Navbar />
 [2m67[0m │       <div className="container py-10">
    · [35;1m           ─────────[0m
 [2m68[0m │         <div className="mb-6">
 [2m69[0m │           <h1 className="text-3xl font-bold">AI Expert Panel Suggestion</h1>
 [2m70[0m │           <p className="text-muted-foreground mt-2">
    ╰────

Caused by:
    Syntax Error
    at BuildError (webpack-internal:///(pages-dir-browser)/./node_modules/next/dist/client/components/react-dev-overlay/ui/container/build-error.js:43:41)
    at renderWithHooks (webpack-internal:///(pages-dir-browser)/./node_modules/react-dom/cjs/react-dom.development.js:15486:18)
    at updateFunctionComponent (webpack-internal:///(pages-dir-browser)/./node_modules/react-dom/cjs/react-dom.development.js:19619:24)
    at beginWork (webpack-internal:///(pages-dir-browser)/./node_modules/react-dom/cjs/react-dom.development.js:21635:16)
    at beginWork$1 (webpack-internal:///(pages-dir-browser)/./node_modules/react-dom/cjs/react-dom.development.js:27460:14)
    at performUnitOfWork (webpack-internal:///(pages-dir-browser)/./node_modules/react-dom/cjs/react-dom.development.js:26591:12)
    at workLoopSync (webpack-internal:///(pages-dir-browser)/./node_modules/react-dom/cjs/react-dom.development.js:26500:5)
    at renderRootSync (webpack-internal:///(pages-dir-browser)/./node_modules/react-dom/cjs/react-dom.development.js:26468:7)
    at performConcurrentWorkOnRoot (webpack-internal:///(pages-dir-browser)/./node_modules/react-dom/cjs/react-dom.development.js:25772:74)
    at workLoop (webpack-internal:///(pages-dir-browser)/./node_modules/scheduler/cjs/scheduler.development.js:266:34)
    at flushWork (webpack-internal:///(pages-dir-browser)/./node_modules/scheduler/cjs/scheduler.development.js:239:14)
    at MessagePort.performWorkUntilDeadline (webpack-internal:///(pages-dir-browser)/./node_modules/scheduler/cjs/scheduler.development.js:533:21)
