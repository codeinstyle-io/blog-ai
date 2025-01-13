<script lang="ts">
  import { Button, Spinner } from 'flowbite-svelte';

  let { savingState = 'draft' } = $props();
  let isSaving: boolean = $state(savingState === 'saving');
</script>

<Button
  type="submit"
  disabled={isSaving}
  class="bg-indigo-600 py-2 px-4 text-sm font-bold text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
>
  {#if savingState === 'saving'}
    <Spinner class="me-3" size="4" color="white" />Saving ...
  {:else if savingState === 'saved'}
    Saved !
  {:else}
    Save
  {/if}
</Button>
