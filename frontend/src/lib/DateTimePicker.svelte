<script lang="ts">
  import { onMount } from "svelte";
  import { type Action } from "svelte/action";
  import { Datepicker } from "flowbite-svelte";
  import Timepicker from "./Timepicker.svelte";

  let { value = $bindable("") }: { value: string } = $props();
  let internalValue = $state(new Date(value || Date.now()));
  let timeValue = $state("00:00");


  const setTime = (date: Date) => {
    const hours = date.getUTCHours().toString().padStart(2, "0");
    const minutes = date.getUTCMinutes().toString().padStart(2, "0");
    timeValue = `${hours}:${minutes}`;
  };

  const inputObserver: Action = (node) => {
    let mutationObserver = new MutationObserver(function(mutations) {
      mutations.forEach(function(mutation) {
        console.log(mutation);
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

    selectedDate.setUTCDate(selectedDate.getUTCDate());
    selectedDate.setUTCMonth(selectedDate.getUTCMonth());
    selectedDate.setUTCFullYear(selectedDate.getUTCFullYear());
    selectedDate.setUTCHours(parseInt(hours));
    selectedDate.setUTCMinutes(parseInt(minutes));

    value = selectedDate.toJSON();
  };

  const updateTime = (e: Event) => {
    const time: string = (e.target as HTMLInputElement).value;
    const [hours, minutes] = time.split(":");

    internalValue.setUTCHours(parseInt(hours));
    internalValue.setUTCMinutes(parseInt(minutes));

    timeValue = time;
    value = new Date(internalValue).toJSON();
  };

  onMount(() => {
    setTime(internalValue);
  });
</script>

<div class="mt-4">
  <div class="flex mb-4">
    <div class="grow text-black dark:text-white">
      <Datepicker inline bind:value={internalValue} on:select={updateDate} />
    </div>
    <div class="mx-2">
      <Timepicker onchange={updateTime} value={timeValue} />
    </div>
    <input type="hidden" name="datetime" use:inputObserver />
  </div>
</div>
