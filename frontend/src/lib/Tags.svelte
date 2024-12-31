<script lang="ts">
  import { Input, Kbd, CloseButton } from "flowbite-svelte";

  let { tags = $bindable() } = $props();

  function removeItem(item: string) {
    tags = tags.filter((i: string) => i !== item);
  }

  function addItem(event: KeyboardEvent) {
    if (event.key !== "Enter") {
      return;
    }
    const element = event.target as HTMLInputElement;
    if (!element.value) {
      return;
    }

    event.preventDefault();
    event.stopPropagation();

    const item = element.value;
    if (tags.includes(item)) {
      return;
    }

    tags = [...tags, item];
    element.value = "";
  }
</script>

<div>
  <Input type="text" id="tags" name="tags" on:keydown={addItem} class="mb-2" />
  <div class="pt-2" id="selected-tags">
    {#each tags as tag}
      <Kbd class="py-3 px-4 mr-4">
        {tag}
        <CloseButton class="align-middle" on:click={() => removeItem(tag)} />
      </Kbd>
    {/each}
  </div>
</div>
