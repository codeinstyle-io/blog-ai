<script>
import Posts from './apps/Posts.svelte'
import Pages from './apps/Pages.svelte'
</script>

<div class="max-w-5xl mx-auto p-6">
  <h1 class="text-2xl font-bold mb-4 text-gray-900">Posts</h1>
  <Posts
    title="My First Post"
    slug="my-first-post"
    excerpt="This is my first post."
    content="Hello, world!"
    tags={["tag1", "tag2", "tag3"]}
    visible={true}
    onSubmit={(data) => {
      console.log('Submitted post data:', data)
    }}
  />
</div>

<hr>

<div class="max-w-5xl mx-auto p-6">
  <h1 class="text-2xl font-bold mb-4 text-gray-900">Pages</h1>
  <Pages
    title="My First Page"
    slug="my-first-page"
    contentType="markdown"
    content="Hello, world!"
    visible={true}
    onSubmit={(data) => {
      console.log('Submitted page data:', data)
    }}
  />
</div>

