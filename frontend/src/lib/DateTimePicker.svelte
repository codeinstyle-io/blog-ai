<script>
    import { Datepicker, Label } from 'flowbite-svelte';
    import Timepicker from './Timepicker.svelte';

    let { value = $bindable() } = $props();
</script>

<div class="mt-4">
    <Label for="date" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">
        Select date and time of publication:
    </Label>

    <div class="flex items-center mb-4">
        <Datepicker bind:value={value} />
        <Timepicker bind:value={value} />
    </div>
    
</div>
