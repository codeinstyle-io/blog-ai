import { Page } from '@playwright/test';

interface CreatePostOptions {
  title: string;
  content: string;
  visible: boolean;
  publishedAt: string;
}


export async function fillDateTimeField(page: Page, date: string) {
  await page.evaluate((date) => {
    const input = document.querySelector('input[name="datetime"]') as HTMLInputElement;
    if (!input) {
      return;
    }
    input.value = date;
  }, date);
}

export async function setVisibility(page: Page, visible: boolean) {
  await page.evaluate(() => {
    const input = document.querySelector('input[name="visible"]') as HTMLInputElement;
    if (!input) {
      return;
    }
    input.checked = visible;
  });
}

export async function createPost(page: Page, options: CreatePostOptions) {
  // Navigate to create post page
  await page.goto('/admin/posts/create');

  // Fill in the form
  await page.fill('input[name="title"]', options.title);
  await page.fill('textarea[name="content"]', options.content);

  // Set visibility
  if (options.visible) {
    await setVisibility(page, true);
  }

  // Set publish type to "scheduled"
  await page.selectOption('select[name="publish"]', 'scheduled');

  // Set published date
  await fillDateTimeField(page, options.publishedAt);

  // Submit the form
  await page.click('button[type="submit"]');

  // Wait for navigation after submission
  await page.waitForURL(/\/admin\/posts/);
}
