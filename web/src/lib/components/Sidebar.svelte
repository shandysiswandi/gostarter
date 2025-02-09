<script lang="ts">
	import { fly } from 'svelte/transition';
	import { page } from '$app/state';
	import { menus } from '$lib';
</script>

<aside
	class="fixed left-0 top-0 z-40 mt-16 h-screen w-40 border-r border-gray-200 bg-white sm:w-48 md:w-64"
	in:fly={{ x: -300, duration: 300 }}
	out:fly={{ x: -300, duration: 100 }}
>
	<nav class="mt-1 flex flex-col space-y-1 px-2">
		{#each menus as menu}
			<a
				href={menu.link === page.url.pathname ? 'javascript:' : menu.link}
				class="flex items-center gap-2 rounded-lg p-2 {menu.link === page.url.pathname
					? 'bg-gray-200'
					: 'hover:bg-gray-100'}"
			>
				<menu.icon class="h-5 w-5" />
				<span>{menu.title}</span>
			</a>
		{/each}
	</nav>
</aside>
