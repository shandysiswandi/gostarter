<script lang="ts">
	import '../app.css';
	import { page } from '$app/state';

	let { children } = $props();
	$effect(() => {
		const title = ['My App', ...page.url.pathname.split('/').slice(1)].filter(Boolean).join(' - ');
		document.title = title;
	});
</script>

{@render children()}
