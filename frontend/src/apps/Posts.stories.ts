import type { Meta, StoryObj } from "@storybook/svelte";
import Posts from "./Posts.svelte";

const meta = {
  title: "Apps/Posts",
  component: Posts,
  tags: ["autodocs"],
  argTypes: {
    onSubmit: { action: "submitted" },
  },
} satisfies Meta<Posts>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {},
};

export const WithPrefilledData: Story = {
  args: {
    title: "Sample Blog Post",
    slug: "sample-blog-post",
    tags: ["svelte", "typescript", "web"],
    excerpt: "This is a sample blog post excerpt.",
    content: "This is the main content of the blog post.",
    publish: "draft",
    visible: true,
  },
};
