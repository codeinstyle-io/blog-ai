<script lang="ts">
    import { Datepicker } from 'flowbite-svelte';
    import Timepicker from './Timepicker.svelte';

    let { value = $bindable(new Date()) } = $props();
    let timeValue = $state(`${value.getHours()}:${value.getMinutes()}`);

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

        value.setHours(hours);
        value.setMinutes(minutes);

        timeValue = time;
        value = new Date(value);
    }
</script>

<div class="mt-4">
    <div class="flex mb-4">
        <div class="grow">
            <Datepicker inline bind:value={value} on:select={updateDate} />
        </div>
        <div class="mx-2">
            <Timepicker onchange={updateTime} value={timeValue} />
        </div>
    </div>
</div>
