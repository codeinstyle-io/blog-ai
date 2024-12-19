import { chromium, type FullConfig } from '@playwright/test';

async function globalSetup(config: FullConfig) {
  console.log('creating admin user...');

  const { baseURL } = config.projects[0].use;
  const browser = await chromium.launch();
  const page = await browser.newPage();
  await page.goto(baseURL! + '/setup');

  // Verify setup page
  if (!page.url().includes('setup')) {
    await browser.close();
    return;
  }

  // We are still on the setup page, proceed with initial configuration
  await page.fill('input[name="firstName"]', 'Admin');
  await page.fill('input[name="lastName"]', 'User');
  await page.fill('input[name="email"]', 'admin@example.com');
  await page.fill('input[name="password"]', 'Password123!');
  await page.click('button[type="submit"]');

  await browser.close();
}

export default globalSetup;