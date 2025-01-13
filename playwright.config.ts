import { defineConfig, devices } from '@playwright/test';


export default defineConfig({
  globalSetup: require.resolve('./e2e/global.setup.ts'),
  timeout: 10000,
  testDir: './e2e/specs',
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: 2,
  reporter: 'html',
  use: {
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    baseURL: `http://localhost:8081`,
    timezoneId: 'Europe/Paris',
  },
  projects: [
    {
      name: 'setup user',
      testMatch: 'setup.spec.ts',
    },
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    
    },
  ],
  webServer: {
      command: [
        "CAPTAIN_DEBUG=1",
        `CAPTAIN_DB_PATH=./testdata/test.db`,
        `CAPTAIN_SERVER_PORT=8081`,
        "./dist/bin/captain run",
      ].join(' '),
      url: `http://localhost:8081`,
      reuseExistingServer: false,
      stdout: 'pipe',
      stderr: 'pipe',
    }
});
