<script lang="ts">
	import { page } from "$app/state";
	import { ChevronUp, Home, X, Settings, User, TriangleAlert, Baby } from "lucide-svelte";

	const user = {
		name: "John Doe",
		email: "john.doe@mail.com"
	};

	const menus = [
		{
			link: "/dashboard",
			icon: Home,
			title: "Dashboard",
			hasChild: false,
			children: []
		},
		{
			link: "/setting",
			icon: Settings,
			title: "Setting",
			hasChild: false,
			children: []
		},
		{
			link: "/profile",
			icon: User,
			title: "Profile",
			hasChild: false,
			children: []
		},
		{
			link: "#",
			icon: null,
			title: "Others",
			hasChild: true,
			children: [
				{
					link: "/child",
					icon: Baby,
					title: "Child"
				},
				{
					link: "/not-found",
					icon: TriangleAlert,
					title: "Not Found"
				}
			]
		}
	];
</script>

<div class="flex h-screen flex-col justify-between">
	<div class="p-4">
		<div class="flex justify-between space-x-2">
			<span
				class="flex h-10 w-full items-center justify-center rounded-lg bg-gray-100 text-xs text-gray-600"
			>
				Logo
			</span>
			<button
				class="flex h-10 w-12 cursor-pointer items-center justify-center rounded-full text-xs text-gray-600 hover:bg-gray-200 lg:hidden"
			>
				<X class="h-5 w-5" />
			</button>
		</div>

		<!-- add scrollable -->
		<ul
			class="mt-6 max-h-[calc(100vh-150px)] space-y-1 overflow-y-auto scroll-smooth"
			style="scrollbar-width: none;"
		>
			{#each menus as menu}
				{#if menu.hasChild}
					<li>
						<details class="group">
							<summary
								class="flex cursor-pointer items-center justify-between rounded-lg px-4 py-2 text-gray-500 hover:bg-gray-100 hover:text-gray-700"
							>
								<span class="text-sm font-medium"> {menu.title} </span>
								<span class="shrink-0 transition duration-300 group-open:-rotate-180">
									<ChevronUp class="size-5" />
								</span>
							</summary>

							<ul class="mt-2 space-y-1 px-4">
								{#each menu.children as child}
									<li>
										<a
											href={child.link}
											class:bg-gray-100={child.link === page.url.pathname}
											class:text-gray-700={child.link === page.url.pathname}
											class="flex items-center rounded-lg px-4 py-2 text-sm font-medium text-gray-500 hover:bg-gray-100 hover:text-gray-700"
										>
											<child.icon class="size-5" />
											<span class="ml-2"> {child.title} </span>
										</a>
									</li>
								{/each}
							</ul>
						</details>
					</li>
				{:else}
					<li>
						<a
							href={menu.link}
							class="flex items-center rounded-lg px-4 py-2 text-sm font-medium text-gray-500 hover:bg-gray-100 hover:text-gray-700"
							class:bg-gray-100={menu.link === page.url.pathname}
							class:text-gray-700={menu.link === page.url.pathname}
						>
							<menu.icon class="size-5" />
							<span class="ml-2"> {menu.title} </span>
						</a>
					</li>
				{/if}
			{/each}
		</ul>
	</div>

	<div class="sticky inset-x-0 bottom-0">
		<a href="/profile" class="flex items-center gap-2 bg-gray-50 p-3">
			<img alt="" src="/profile.avif" class="size-10 rounded-full object-cover" />

			<div>
				<p class="text-xs">
					<strong class="block truncate font-medium text-gray-600">{user.name}</strong>

					<span class="truncate text-gray-500"> {user.email} </span>
				</p>
			</div>
		</a>
	</div>
</div>
