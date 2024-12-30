<script lang="ts">
  import { onMount } from "svelte";
  import {
    Alert,
    Input,
    Label,
    Textarea,
    Toggle,
    Select,
  } from "flowbite-svelte";

  import Tags from "../lib/Tags.svelte";
  import DateTimePicker from "../lib/DateTimePicker.svelte";
  import SubmitButton from "../lib/SubmitButton.svelte";

  import { type Posts } from "../utils/types/posts";
  import { slugify } from "../utils/text";
  import type { SavingStates } from "../utils/types/common";

  let originalSlug = $state("");
  let protectedSlug = $state(false);
  let protectedSlugViolation = $state(false);
  let publish = $state("immediately");

  let {
    title = "",
    tags = [],
    excerpt = "",
    content = "",
    visible = false,
    publishedAt = "",
    slug = "",
    savingState = "draft",
    onSubmit = (data: any, done: (savingState: SavingStates) => void) => {
      done("saved");
    },
  }: Posts = $props();

  const publishOptions = [
    { value: "immediately", name: "Immediately" },
    { value: "scheduled", name: "Scheduled" },
  ];

  onMount(() => {
    if (slug !== "") {
      protectedSlug = true;
      originalSlug = slug;
    }

    if (publishedAt) {
      publish = "scheduled";
    }
  });

  function updateSlug() {
    if (protectedSlug) {
      return;
    }
    slug = slugify(title);
  }

  function onSlugChange(e: Event) {
    const target = e.target as HTMLInputElement;
    protectedSlugViolation = target.value !== originalSlug && protectedSlug;
  }

  function handleSubmit(e: Event) {
    e.preventDefault();

    onSubmit(
      {
        title,
        tags,
        excerpt,
        content,
        visible,
        publishedAt,
        slug,
      },
      (newSavingState: SavingStates) => {
        savingState = newSavingState;
      },
    );
  }
</script>

<div>
  <form onsubmit={handleSubmit} class="space-y-6">
    <!-- Title -->
    <div>
      <Label for="title" class="block text-sm font-bold text-gray-700 mb-2"
        >Title</Label
      >
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
      <Label for="slug" class="block text-sm font-bold text-gray-700 mb-2"
        >Slug</Label
      >
      <Input
        type="text"
        id="slug"
        bind:value={slug}
        onkeyup={onSlugChange}
        required
      />
      {#if protectedSlugViolation}
        <Alert color="red" class="mt-2">
          <span class="font-medium">Changing the slug may break links.</span>
        </Alert>
      {/if}
    </div>

    <!-- Tags -->
    <div>
      <Label for="tags" class="block text-sm font-bold text-gray-700 mb-2"
        >Tags</Label
      >
      <Tags bind:tags />
    </div>

    <!-- Excerpt -->
    <div>
      <Label for="excerpt" class="block text-sm font-bold text-gray-700 mb-2"
        >Excerpt</Label
      >
      <Textarea id="excerpt" bind:value={excerpt} rows={4}></Textarea>
    </div>

    <!-- Content -->
    <div>
      <Label for="content" class="block text-sm font-bold text-gray-700 mb-2"
        >Content</Label
      >
      <Textarea id="content" bind:value={content} rows={10}></Textarea>
    </div>

    <div class="grid gap-4 sm:grid-cols-2 sm:gap-6">
      <!-- Visibility Toggle -->
      <div>
        <Label for="visible" class="block text-sm font-bold text-gray-700 mb-4"
          >Visibility</Label
        >
        <Toggle id="visible" bind:checked={visible}>
          <svelte:fragment slot="offLabel">Hidden</svelte:fragment>
          <span>Visible</span>
        </Toggle>
      </div>
      <!-- Publish Status -->
      <div>
        <Label for="publish" class="block text-sm font-bold text-gray-700 mb-2"
          >Publish</Label
        >
        <Select id="publish" bind:value={publish} items={publishOptions}
        ></Select>

        {#if publish === "scheduled"}
          <div class="mt-4">
            <span
              class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
            >
              Select date and time of publication:
            </span>
            <DateTimePicker bind:value={publishedAt} />
          </div>
        {/if}
      </div>
    </div>

    <!-- Submit Button -->
    <div class="flex justify-end">
      <SubmitButton {savingState} />
    </div>
  </form>
</div>
