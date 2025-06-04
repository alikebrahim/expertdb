# BQA Website CSS Analysis

## CSS Framework and Libraries Used

### 1. Main Stylesheets
- **Main Theme CSS**: `wp-content/themes/bqa/style.css` - Custom BQA theme styles
- **Bootstrap CSS**: `wp-content/themes/bqa/css/bootstrap.min.css` - Grid system and components
- **WordPress Core**: `wp-includes/css/dist/block-library/style.min.css` - Gutenberg blocks
- **Font Awesome**: Icon library for various icons
- **Cookie Notice**: `wp-content/plugins/cookie-notice/css/front.min.css`
- **WPML Language Switcher**: Multilingual support styles
- **Accessibility Plugin**: Screen reader support styles

### 2. Custom Fonts
- **Primary Font**: "Graphik" - Custom font family with multiple weights
  - Regular (400)
  - Medium (500) 
  - Semibold (600)
  - Bold (700)
  - Italic variants for each weight
- **Arabic Font**: "Graphik Arabic" - For Arabic language support
- **Fallback Font**: "Segoe UI" - System font fallback

### 3. Icon Fonts
- **Icomoon**: Custom icon font with BQA-specific icons
  - book-open, pin, envelope, calendar, home
  - backup, recovery, smartphone, support, maintenance
  - exam, reports, results, test, pencil, wrench, child, study

### 4. Key CSS Features

#### Color Scheme
- **Primary Green**: #397b26 (BQA brand color)
- **Secondary Red**: #e64125 (accent color)
- **Blue**: Various shades for service cards
- **Gray**: #494949 (text color)

#### Layout System
- **Bootstrap Grid**: 12-column responsive grid system
- **Flexbox**: Modern layout for components
- **CSS Grid**: Used for complex layouts

#### Component Styling
- **Nice Select**: Custom dropdown styling
- **Card Components**: Service cards with hover effects
- **Button Styles**: Consistent button design across site
- **Form Styling**: Contact Form 7 integration

#### Responsive Design
- **Mobile-first approach**
- **Breakpoints**: Standard Bootstrap breakpoints
- **Flexible typography**: Responsive font sizes
- **Touch-friendly**: Mobile interaction support

### 5. WordPress Integration
- **Gutenberg Blocks**: Full support for WordPress block editor
- **Widget Styling**: Custom widget area styles
- **Comment Forms**: Styled comment system
- **Search Forms**: Custom search functionality

### 6. Performance Optimizations
- **Minified CSS**: All stylesheets are minified
- **Font Display Swap**: Optimized font loading
- **CSS Caching**: Rocket CDN for fast delivery
- **Critical CSS**: Above-the-fold optimization

### 7. Accessibility Features
- **Screen Reader Support**: Dedicated accessibility styles
- **Focus States**: Keyboard navigation support
- **High Contrast**: Accessible color combinations
- **ARIA Labels**: Proper semantic markup

### 8. Third-party Integrations
- **Cookie Notice**: GDPR compliance styling
- **WPML**: Multilingual support
- **Contact Form 7**: Form styling
- **Chat Widget**: Customer support integration

## CSS Architecture

### File Structure
```
/wp-content/themes/bqa/
├── style.css (main theme styles)
├── css/
│   ├── bootstrap.min.css
│   └── [other component styles]
├── fonts/
│   ├── graphik-*.woff2
│   ├── graphik-arabic-*.woff2
│   └── icomoon.woff
└── js/
    └── [JavaScript files]
```

### Methodology
- **Component-based**: Modular CSS architecture
- **BEM-like naming**: Consistent class naming
- **Utility classes**: Helper classes for common styles
- **Custom properties**: CSS variables for theming

### Browser Support
- **Modern browsers**: Chrome, Firefox, Safari, Edge
- **IE11**: Limited support with fallbacks
- **Mobile browsers**: Full responsive support

