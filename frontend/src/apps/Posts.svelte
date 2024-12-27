<script lang="ts">
  import { Input, Label, Textarea, Button, Toggle, Select } from 'flowbite-svelte';
  import Tags from '../lib/Tags.svelte';
  import DateTimePicker from '../lib/DateTimePicker.svelte';

  let title: string = '';
  let slug: string = '';
  let tags: string[] = [];
  let excerpt: string = '';
  let content: string = '';
  let publish: string = 'immediately';
  let visible: boolean = false;


  const publishOptions = [
    {value: 'immediately', name: 'Immediately'},
    {value: 'scheduled', name: 'Scheduled'}
  ]

  function handleSubmit() {
    // Handle form submission
    console.log({
      title,
      slug,
      tags,
      excerpt,
      content,
      publish,
      visible
    });
  }
</script>

<div>
  <form on:submit|preventDefault={handleSubmit} class="space-y-6">
    <!-- Title -->
    <div>
      <Label for="title" class="block text-sm text-gray-700 font-bold">Title</Label>
      <Input
        type="text"
        id="title"
        bind:value={title}
        required
      />
    </div>

    <!-- Slug -->
    <div>
      <Label for="slug" class="block text-sm font-bold text-gray-700">Slug</Label>
      <Input
        type="text"
        id="slug"
        bind:value={slug}
        required
      />
    </div>

    <!-- Tags -->
    <div>
      <Label for="tags" class="block text-sm font-bold text-gray-700">Tags</Label>
      <Tags bind:tags={tags} />
    </div>

    <!-- Excerpt -->
    <div>
      <Label for="excerpt" class="block text-sm font-bold text-gray-700">Excerpt</Label>
      <Textarea
        id="excerpt"
        bind:value={excerpt}
        rows={4}
      ></Textarea>
    </div>

    <!-- Content -->
    <div>
      <Label for="content" class="block text-sm font-bold text-gray-700">Content</Label>
      <Textarea
        id="content"
        bind:value={content}
        rows={10}
      ></Textarea>
    </div>

    <div class="grid gap-4 sm:grid-cols-2 sm:gap-6">
        <!-- Visibility Toggle -->
        <div>
            <Label for="visible" class="block text-sm font-bold text-gray-700 mb-4">Visibility</Label>
            <Toggle bind:checked={visible}>
                <svelte:fragment slot="offLabel">Hidden</svelte:fragment>
                <span>Visible</span>
            </Toggle>
        </div>
        <!-- Publish Status -->
        <div>
            <Label for="publish" class="block text-sm font-bold text-gray-700">Publish</Label>
            <Select id="publish" bind:value={publish} items={publishOptions}></Select>

            <DateTimePicker />
        </div>
    </div>
    
    <!-- Submit Button -->
    <div class="flex justify-end">
      <Button
        type="submit"
        class="bg-indigo-600 py-2 px-4 text-sm font-bold text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
      >
        Save
      </Button>
    </div>
  </form>
</div>