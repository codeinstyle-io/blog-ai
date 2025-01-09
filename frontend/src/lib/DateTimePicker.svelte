<script lang="ts">
  import { onMount } from "svelte";
  import { type Action } from "svelte/action";
  import { Datepicker } from "flowbite-svelte";
  import Timepicker from "./Timepicker.svelte";
  import TimezonePicker from "./TimezonePicker.svelte";

  let { value = $bindable(""), timezone = $bindable("UTC") }: { value: string, timezone: string } = $props();
  let internalValue = $state(new Date(value || Date.now()));
  let timeValue = $state("00:00");


  const setTime = (date: Date) => {
    const hours = date.getHours().toString().padStart(2, "0");
    const minutes = date.getMinutes().toString().padStart(2, "0");
    timeValue = `${hours}:${minutes}`;
  };

  const inputObserver: Action = (node) => {
    let mutationObserver = new MutationObserver(function(mutations) {
      mutations.forEach(function(mutation) {
        if(mutation.type === "attributes" && mutation.attributeName === "value") {
          value = node.value;
          internalValue = new Date(value);
          setTime(internalValue);
        }
      });
    });
    mutationObserver.observe(node, { attributes: true });
  };

  const updateDate = (e: CustomEvent) => {
    const selectedDate: Date = new Date(e.detail.valueOf());
    const [hours, minutes] = timeValue.split(":");

    selectedDate.setDate(selectedDate.getDate());
    selectedDate.setMonth(selectedDate.getMonth());
    selectedDate.setFullYear(selectedDate.getFullYear());
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

  onMount(() => {
    setTime(internalValue);
  });
</script>

<style>
  .datepicker-container :global(.inline-block) {
    width: 100%;
  }
</style>

<div class="mt-4">
  <div>
    <h3 class="text-sm font-medium text-gray-900 dark:text-white my-2">Time of publication:</h3>
    <Timepicker onchange={updateTime} value={timeValue} />
    <h3 class="text-sm font-medium text-gray-900 dark:text-white my-2">Timezone:</h3>
    <TimezonePicker bind:timezone={timezone} />
  </div>
  <div class="full text-black dark:text-white">
    <h3 class="text-sm font-medium text-gray-900 dark:text-white my-2">Date of publication:</h3>
    <div class="datepicker-container">
      <Datepicker inline bind:value={internalValue} on:select={updateDate} />
    </div>
  </div>
  <input type="hidden" name="datetime" use:inputObserver />
</div>
