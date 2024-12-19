import { defineConfig, devices } from '@playwright/test';


let portPrefix: string = "808";

function buildCommand(serverNumber: number) {
  const command = [
    "CAPTAIN_DEBUG=1",
    `CAPTAIN_DB_PATH=./testdata/test-${serverNumber}.db`,
    `CAPTAIN_SERVER_PORT=808${serverNumber}`,
    "./dist/bin/captain run",
  ].join(' ');

  console.log(command);

  return command;
}


export default defineConfig({
  timeout: 10000,
  testDir: './e2e/specs',
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: 1,
  reporter: 'html',
  use: {
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    baseURL: `http://localhost:${portPrefix}1`,
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
  webServer: [
    {
      command: buildCommand(1),
      url: `http://localhost:${portPrefix}1`,
      reuseExistingServer: true,
      stdout: 'pipe',
      stderr: 'pipe',
    },
    {
      command: buildCommand(2),
      url: `http://localhost:${portPrefix}2`,
      reuseExistingServer: true,
      stdout: 'pipe',
      stderr: 'pipe',
    },
],
});
