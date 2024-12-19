import { test, expect } from '@playwright/test';
import { login, logout } from '../helpers/auth';
import { createPost } from '../helpers/posts';

test.describe('Post Visibility', () => {

  const SERVER_NUMBER = 2;
  const now = new Date('2024-12-18T22:30:12+01:00'); // Fixed time from context
  const tomorrow = new Date(now);
  tomorrow.setDate(tomorrow.getDate() + 1);

  test('should handle draft and scheduled posts correctly', async ({ page }) => {
    // Setup admin account
    await logout(page);

    // Create a draft post
    await login(page);
    await createPost(page, {
      title: 'Draft Post',
      content: 'This is a draft post',
      visible: false,
      publishedAt: now.toISOString(),
    });

    // Create a scheduled post
    await createPost(page, {
      title: 'Scheduled Post',
      content: 'This is a scheduled post',
      visible: true,
      publishedAt: tomorrow.toISOString(),
    });

    // Verify posts are visible and properly marked when logged in
    await page.goto('/');

    // Check draft post
    const draftPost = page.locator('article:has-text("Draft Post")');
    await expect(draftPost).toBeVisible();
    await expect(draftPost).toHaveClass(/draft-post/);
    await expect(draftPost.locator('.draft-indicator')).toBeVisible();
    await expect(draftPost.locator('.edit-link')).toBeVisible();

    // Check scheduled post
    const scheduledPost = page.locator('article:has-text("Scheduled Post")');
    await expect(scheduledPost).toBeVisible();
    await expect(scheduledPost).toHaveClass(/scheduled-post/);
    await expect(scheduledPost.locator('.scheduled-indicator')).toBeVisible();
    await expect(scheduledPost.locator('.edit-link')).toBeVisible();

    // Logout and verify posts are not visible to anonymous users
    await page.goto('/logout');
    await page.goto('/');

    // Verify posts are not visible
    await expect(page.locator('article:has-text("Draft Post")')).not.toBeVisible();
    await expect(page.locator('article:has-text("Scheduled Post")')).not.toBeVisible();
  });
});
