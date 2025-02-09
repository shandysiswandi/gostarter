<script lang="ts">
	import { Menu } from 'lucide-svelte';
	import { menus } from '$lib';
	import { page } from '$app/state';

	let title: string = $state('');
	$effect(() => {
		title = menus.find((menu) => menu.link === page.url.pathname)?.title || '';
	});

	let { handleMenuClick }: { handleMenuClick: () => void } = $props();
</script>

<header class="fixed z-50 w-full select-none border-b border-gray-200 bg-white">
	<div class="mx-auto px-4">
		<div class="flex h-16 items-center justify-between">
			<div class="mr-auto flex">
				<Menu class="h-7 w-7 hover:cursor-pointer" onclick={handleMenuClick} />
				<div class="ml-2 flex flex-shrink-0 items-center space-x-2">
					<img class="hidden h-7 w-7 md:block" src="/logo.svg" alt="" />
					<h2 class="text-xl font-bold">{title}</h2>
				</div>
			</div>
			<div class="flex items-center justify-end space-x-6 sm:ml-5">
				<button
					class="flex max-w-xs items-center rounded-full bg-gray-100 hover:cursor-pointer focus:outline-none"
				>
					<img
						class="bg-secondary h-8 w-8 rounded-full object-cover"
						src="https://api.dicebear.com/7.x/croodles/png?seed=Shivam"
						alt=""
					/>
					<span class="ml-2 mr-3 hidden text-sm font-medium md:block"> Username </span>
				</button>
			</div>
		</div>
	</div>
</header>
