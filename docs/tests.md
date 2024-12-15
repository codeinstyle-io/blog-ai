# Testing Guide

This document explains how to run both unit tests and end-to-end (E2E) tests for the Captain project.

## Unit Tests

Unit tests are written in Go and can be run using the following commands:

```bash
# Run all unit tests
make test

# Run unit tests with coverage report
make test-coverage
```

The coverage report will be generated as `coverage.html` in the project root directory.

## End-to-End (E2E) Tests

E2E tests are written using Playwright and TypeScript. They test the entire application flow from the user's perspective.

### Prerequisites

1. Install Node.js (v16 or later)
2. Install dependencies:
```bash
npm install -D @playwright/test @faker-js/faker
npx playwright install
```

### Running E2E Tests

There are two ways to run the E2E tests:

1. Headless mode (CI-friendly):
```bash
make test-e2e
```

2. UI mode (interactive, good for development):
```bash
make test-e2e-ui
```

### Test Scenarios

The E2E tests cover the following scenarios:

1. Initial Setup
   - First start form completion
   - Admin login

2. Posts Management
   - Create post with specific date, excerpt, and tags
   - Verify post creation
   - Edit post and remove excerpt
   - Verify slug persistence
   - Create additional post with random data

3. Pages Management
   - Create new page
   - Verify page creation
   - Edit page
   - Verify slug persistence

4. Menu Management
   - Create menu item with custom URL
   - Create menu item linked to page
   - Verify menu items display and URLs

5. Tags Management
   - Verify post tags are listed
   - Test tag filtering
   - Verify public tag pages

6. Settings Management
   - Change timezone and verify post times
   - Modify posts per page setting
   - Update site title and subtitle
   - Verify changes on public site

### Test Reports

After running the tests:
- HTML report is available in `playwright-report/`
- Screenshots of failures are saved in `test-results/`
- Traces for debugging are available in UI mode

### Debugging

For detailed debugging:
1. Use UI mode: `make test-e2e-ui`
2. Step through tests
3. View network requests
4. Inspect DOM elements
5. Time-travel through test steps

### CI Integration

The tests are configured to run in CI environments with:
- Retries for flaky tests
- Parallel execution disabled
- Additional logging
- Failure screenshots
