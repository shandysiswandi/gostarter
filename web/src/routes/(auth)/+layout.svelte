<script lang="ts">
	let { children } = $props();
</script>

<div class="flex min-h-screen items-center justify-center bg-gray-100 p-4">
	<div class="w-full max-w-md rounded-xl bg-white p-8 shadow-lg">
		{@render children()}
	</div>
</div>
