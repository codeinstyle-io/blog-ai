import { defineConfig, devices } from '@playwright/test';


export default defineConfig({
  testDir: './e2e/specs',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    baseURL: 'http://localhost:8080',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
  webServer: {
    command: `mkdir -p ./dist && rm -f ./dist/test.db && CAPTAIN_DEBUG=1 CAPTAIN_DB_PATH=./dist/test.db CAPTAIN_SERVER_PORT=8081 make run`,
    url: 'http://localhost:8081',
    reuseExistingServer: !process.env.CI,
    stdout: 'pipe',
    stderr: 'pipe',
  },
});
