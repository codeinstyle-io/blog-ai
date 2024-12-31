<script lang="ts">
  import { onMount } from "svelte";
  import { type Action } from "svelte/action";
  import { Datepicker } from "flowbite-svelte";
  import Timepicker from "./Timepicker.svelte";

  let { value = $bindable("") }: { value: string } = $props();
  let timeValue = $state("00:00");
  let internalValue = $state(new Date(value || Date.now()));

  onMount(() => {
    timeValue = `${internalValue.getHours()}:${internalValue.getMinutes()}`;
  });

  const test: Action = (node) => {
    let mutationObserver = new MutationObserver(function(mutations) {
      mutations.forEach(function(mutation) {
        if(mutation.type === "attributes" && mutation.attributeName === "value") {
          value = node.value;
        }
      });
    });
    mutationObserver.observe(node, { attributes: true });
  };

  const updateDate = (e: CustomEvent) => {
    const selectedDate: Date = e.detail;
    const [hours, minutes] = timeValue.split(":");

    selectedDate.setHours(parseInt(hours));
    selectedDate.setMinutes(parseInt(minutes));

    value = selectedDate.toJSON();
  };

  const updateTime = (e: Event) => {
    const time: string = (e.target as HTMLInputElement).value;
    const [hours, minutes] = time.split(":");

    internalValue.setHours(parseInt(hours));
    internalValue.setMinutes(parseInt(minutes));

    timeValue = time;
    value = new Date(internalValue).toJSON();
  };
</script>

<div class="mt-4">
  <div class="flex mb-4">
    <div class="grow text-black dark:text-white">
      <Datepicker inline bind:value={internalValue} on:select={updateDate} />
    </div>
    <div class="mx-2">
      <Timepicker onchange={updateTime} value={timeValue} />
    </div>
    <input type="hidden" name="datetime" bind:value use:test  />
  </div>
</div>
