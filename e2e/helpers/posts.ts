import { Page } from '@playwright/test';

interface CreatePostOptions {
  title: string;
  content: string;
  visible: boolean;
  publishedAt: string;
}

export async function createPost(page: Page, options: CreatePostOptions) {
  // Navigate to create post page
  await page.goto('/admin/posts/create');

  // Fill in the form
  await page.fill('input[name="title"]', options.title);
  await page.fill('textarea[name="content"]', options.content);

  // Set visibility
  if (options.visible) {
    await page.click('.toggle-switch');
  }

  // Set publish type to "scheduled"
  await page.selectOption('select[id="publishType"]', 'scheduled');

  // Set published date
  const [first, second] = options.publishedAt.split(':');

  await page.fill('input[name="publishedAt"]', `${first}:${second}`);

  // Submit the form
  await page.click('button[type="submit"]');

  // Wait for navigation after submission
  await page.waitForURL(/\/admin\/posts/);
}
