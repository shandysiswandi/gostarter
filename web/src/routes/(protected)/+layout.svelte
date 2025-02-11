<script lang="ts">
	import Aside from "$lib/components/layouts/Aside.svelte";
	import Header from "$lib/components/layouts/Header.svelte";
	import Footer from "$lib/components/layouts/Footer.svelte";
	import AsideChild from "$lib/components/layouts/AsideChild.svelte";

	let isSidebarOpen: boolean = $state(false);

	const toggleSidebar = () => {
		isSidebarOpen = !isSidebarOpen;
	};

	let { children } = $props();
</script>

<div class="flex min-h-screen bg-gray-100">
	<Aside />

	{#if isSidebarOpen}
		<button
			aria-label="overlay"
			class="fixed inset-0 z-20 bg-gray-900 opacity-50 transition-opacity lg:hidden"
			onclick={() => (isSidebarOpen = false)}
		></button>
	{/if}

	<div
		class={`ease-linier fixed left-0 top-0 z-30 min-h-screen w-64 select-none border-r border-gray-200 bg-white transition-transform duration-300 lg:hidden ${isSidebarOpen ? "translate-x-0" : "-translate-x-full"}`}
	>
		<AsideChild />
	</div>

	<div class="flex flex-1 flex-col duration-300 ease-linear lg:ml-64">
		<Header toggleSidebar={() => (isSidebarOpen = !isSidebarOpen)} />

		<main class="flex-1 p-4">
			{@render children()}
		</main>

		<Footer />
	</div>
</div>
