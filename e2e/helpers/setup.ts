import { Page, expect } from '@playwright/test';

export async function setupAdmin(page: Page, serverNumber: number) {
  const rootURL = 'http://localhost:808' + serverNumber;
  await page.goto(rootURL + '/setup');

  // Setup initial configuration
  await page.fill('input[name="firstName"]', 'Admin');
  await page.fill('input[name="lastName"]', 'User');
  await page.fill('input[name="email"]', 'admin@example.com');
  await page.fill('input[name="password"]', 'Password123!');
  await page.click('button[type="submit"]');

  // Verify redirect to login
  await expect(page).toHaveURL(/.*\/login/);
}
