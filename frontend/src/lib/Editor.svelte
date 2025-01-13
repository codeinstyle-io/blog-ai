<script lang="ts">
  import { Textarea } from 'flowbite-svelte';
  import EditorToolbar from './EditorToolbar.svelte';

  let { value = $bindable(''), rows = 10, name = '' } = $props();
  let element = $state(null);
</script>

<div>
  <EditorToolbar target={element} />
  <textarea
    bind:this={element}
    bind:value
    {name}
    {rows}
    class="h-48 w-full my-2 p-2 border border-gray-300 focus:border-indigo-500 focus:ring-indigo-500 rounded-md shadow-sm"
  ></textarea>
</div>
