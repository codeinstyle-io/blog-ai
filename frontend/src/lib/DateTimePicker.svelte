<script lang="ts">
    import { onMount } from 'svelte';
    import { Datepicker } from 'flowbite-svelte';
    import Timepicker from './Timepicker.svelte';

    let { value = $bindable(new Date()) }: {value: Date} = $props();
    let timeValue = $state('00:00');

    onMount(() => {
        if (!(value instanceof Date)) {
            value = new Date();
        }
        timeValue = `${value.getHours()}:${value.getMinutes()}`;
    })

    const updateDate = (e: CustomEvent) => {
        const selectedDate: Date = e.detail;
        const [hours, minutes] = timeValue.split(':');

        selectedDate.setHours(parseInt(hours));
        selectedDate.setMinutes(parseInt(minutes));

        value = selectedDate;
    }

    const updateTime = (e: Event) => {
        const time: string = (e.target as HTMLInputElement).value;
        const [hours, minutes] = time.split(':');

        value.setHours(parseInt(hours));
        value.setMinutes(parseInt(minutes));

        timeValue = time;
        value = new Date(value);
    }
</script>

<div class="mt-4">
    <div class="flex mb-4">
        <div class="grow text-black dark:text-white">
            <Datepicker inline bind:value={value} on:select={updateDate} />
        </div>
        <div class="mx-2">
            <Timepicker onchange={updateTime} value={timeValue} />
        </div>
    </div>
</div>
