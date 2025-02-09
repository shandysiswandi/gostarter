<script lang="ts">
	import Header from '$lib/components/Header.svelte';
	import Sidebar from '$lib/components/Sidebar.svelte';

	let isSidebarOpen: boolean = $state(true);

	let { children } = $props();
</script>

<div class="flex min-h-screen flex-col bg-gray-100">
	<Header handleMenuClick={() => (isSidebarOpen = !isSidebarOpen)} />

	<div class="flex">
		{#if isSidebarOpen}
			<Sidebar />
		{/if}

		<main class="mt-16 w-full p-4" class:ml-64={isSidebarOpen}>
			{@render children()}
		</main>
	</div>
</div>
