import { ro } from '@faker-js/faker/.';
import { Page } from '@playwright/test';

export async function login(page: Page) {
  await page.goto('/login?next=/admin');

  // Fill in login form
  await page.fill('input[name="email"]', 'admin@example.com');
  await page.fill('input[name="password"]', 'Password123!');

  // Submit form
  await page.click('button[type="submit"]');

  // Wait for successful login and redirect
  await page.waitForURL('/admin');
}

export async function logout(page: Page) {
  await page.goto('/logout');
  await page.waitForURL('/login');
}