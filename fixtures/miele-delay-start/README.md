# Miele W1 Delay Start Calculator

A standalone web calculator to help determine when to start your Miele W1 washing machine so it finishes at your desired time.

## Features

- Calculate exact delay needed based on current time, program duration, and desired finish time
- Accounts for Miele W1-specific delay constraints:
  - 30-minute increments for delays up to 10 hours
  - 1-hour increments for delays from 10-24 hours
- Shows both exact mathematical result and closest W1-compatible setting
- Automatically handles next-day finish times
- Dark-themed, mobile-friendly interface

## Usage

### Quick Start

Just open `index.html` in your browser. No installation needed.

### Development

Install dependencies and run the development server:

```bash
npm install
npm run dev
```

This will start Vite's dev server at http://localhost:5173 with hot module replacement.

### Building for Production

Build the app for production:

```bash
npm run build
```

The optimized files will be output to the `dist/` directory.

### Preview Production Build

Preview the production build locally:

```bash
npm run preview
```

### Testing

Run the test suite:

```bash
npm test
```

Run tests with UI:

```bash
npm run test:ui
```

The test suite covers all core calculator functions including time parsing, formatting, and the W1-specific delay quantization logic.

## How It Works

1. Enter the current time (or click "Use now" to auto-fill)
2. Enter your washing program duration (e.g., 03:39 for 3 hours 39 minutes)
3. Enter when you want the wash to finish
4. Click "Calculate delay"

The calculator will show:
- **Exact math**: The precise delay and start time needed
- **Closest W1 delay setting**: The nearest delay your Miele W1 can actually be set to

## Miele W1 Delay Rules

Typical W1 behaviour (check your exact model manual):
- Delay start available from **30 minutes** up to **24 hours**
- Below 10 hours: steps of **30 minutes** (0:30, 1:00, 1:30, …)
- From 10-24 hours: steps of **1 hour** (10:00, 11:00, …, 24:00)

## Technical Details

- Pure HTML/CSS/JavaScript with ES modules
- Zero runtime dependencies
- Built with Vite for fast development and optimized production builds
- Tested with Vitest (13 unit tests)
- Hot module replacement in development mode
- Works offline once loaded

## License

MIT
