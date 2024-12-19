import { ro } from '@faker-js/faker/.';
import { Page } from '@playwright/test';

export async function login(page: Page, serverNumber: number) {
  const rootURL = 'http://localhost:808' + serverNumber;
  await page.goto(rootURL + '/login?next=/admin');

  // Fill in login form
  await page.fill('input[name="email"]', 'admin@example.com');
  await page.fill('input[name="password"]', 'Password123!');

  // Submit form
  await page.click('button[type="submit"]');

  // Wait for successful login and redirect
  await page.waitForURL(rootURL + '/admin');
}

export async function logout(page: Page, serverNumber: number) {
  const rootURL = 'http://localhost:808' + serverNumber;
  await page.goto(rootURL + '/logout');
  await page.waitForURL(rootURL + '/');
}
