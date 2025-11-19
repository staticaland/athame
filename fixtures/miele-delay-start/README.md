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

### Development Server

If you want to run a local development server:

```bash
npm run dev
```

This will start a server at http://localhost:8080 and open it in your browser.

### Production Server

To run without opening the browser:

```bash
npm start
```

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

- Pure HTML/CSS/JavaScript - no build step required
- Zero dependencies for the app itself
- Uses `http-server` for development convenience (optional)
- Works offline once loaded

## License

MIT
