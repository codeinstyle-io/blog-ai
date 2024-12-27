<script lang="ts">
    import { onMount } from 'svelte';
    import { Alert, Input, Label, Textarea, Toggle, Select } from 'flowbite-svelte';

    import { type Pages } from '../utils/types/pages';
    import { slugify } from '../utils/text';
    import SubmitButton from '../lib/SubmitButton.svelte';
  
    let protectedSlug = $state(false);
    let protectedSlugViolation = $state(false);
  
    let {
      title = '', 
      content = '', 
      visible = false, 
      slug = '',
      savingState = 'draft',
      contentType = '',
      onSubmit = (data: any) => {}
   }: Pages = $props();
  
    onMount(() => {
      if (slug !== '') {
          protectedSlug = true;
      }
    });

    const contentTypeOptions = [
        {value: 'html', name: 'HTML'},
        {value: 'markdown', name: 'Markdown'}
    ]
  
    function updateSlug() {
      if (protectedSlug) {
          return;
      }
      slug = slugify(title);
    }
  
    function onSlugChange(e: Event) {
      const target = e.target as HTMLInputElement;
      protectedSlugViolation = (target.value !== slug) && protectedSlug;
    }
  
    function handleSubmit(e: Event) {
      e.preventDefault();
  
      onSubmit({
          title,
          content,
          visible,
          slug
      });
    }
  </script>
  
  <div>
    <form onsubmit={handleSubmit} class="space-y-6">
      <!-- Title -->
      <div>
        <Label for="title" class="block text-sm font-bold text-gray-700 mb-2">Title</Label>
        <Input
          type="text"
          id="title"
          bind:value={title}
          onkeyup={updateSlug}
          required
        />
      </div>
  
      <!-- Slug -->
      <div>
        <Label for="slug" class="block text-sm font-bold text-gray-700 mb-2">Slug</Label>
        <Input
          type="text"
          id="slug"
          value={slug}
          onkeyup={onSlugChange}
          required
        />
        {#if protectedSlugViolation}
          <Alert color="red" class="mt-2">
              <span class="font-medium">Changing the slug may break links.</span>
          </Alert>
         {/if}
      </div>

      <!-- Content Type -->
      <div>
        <Label for="content-type" class="block text-sm font-bold text-gray-700 mb-2">Content Type</Label>
        <Select id="content-type" bind:value={contentType} items={contentTypeOptions}></Select>
      </div>
  
      <!-- Content -->
      <div>
        <Label for="content" class="block text-sm font-bold text-gray-700 mb-2">Content</Label>
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
              <Toggle id="visible" bind:checked={visible}>
                  <svelte:fragment slot="offLabel">Hidden</svelte:fragment>
                  <span>Visible</span>
              </Toggle>
          </div>
      </div>
      
      <!-- Submit Button -->
      <div class="flex justify-end">
        <SubmitButton savingState={savingState} />
      </div>
    </form>
  </div>