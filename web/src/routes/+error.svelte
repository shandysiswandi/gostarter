<script lang="ts">
	import { page } from "$app/state";
</script>

<main class="flex h-screen items-center justify-center bg-gray-100">
	<div class="flex items-center bg-gray-100 text-gray-800">
		<div class="container mx-auto flex flex-col items-center justify-center">
			<div class="max-w-md text-center">
				<h2 class="mb-8 text-9xl font-extrabold text-gray-400">
					<span class="sr-only">Error</span>{page?.status}
				</h2>
				<p class="text-2xl font-semibold md:text-3xl">Sorry, we couldn't find this page.</p>
				<p class="mb-8 mt-4 text-gray-500 md:text-lg">
					But don't worry, you can find plenty of other things on our homepage.
				</p>
				<a href="/" class="rounded bg-violet-600 px-8 py-3 font-semibold text-gray-50">
					Back to Home
				</a>
			</div>
		</div>
	</div>
</main>
