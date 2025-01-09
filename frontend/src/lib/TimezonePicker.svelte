<script>
  import { onMount } from "svelte";
  import { Select } from "flowbite-svelte";
  import timezones from "../utils/timezones";

  let { timezone = $bindable("Africa/Casablanca") } = $props();


  onMount(() => {
    const tz = new Intl.DateTimeFormat().resolvedOptions().timeZone;
    timezone = tz;
  });

</script>

<Select
  bind:value={timezone}
  items={timezones}
  class="w-full my-2"
/>
