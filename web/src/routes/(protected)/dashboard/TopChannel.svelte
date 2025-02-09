<script lang="ts">
	import { Pencil, Trash } from 'lucide-svelte';

	const tableChannels = [
		{
			id: 1,
			song: 'New York No. 1 Lake Park mNew York No. 1 Lake Park New York No. 1 Lake Park',
			artist: 'John Brown',
			year: '2019'
		},
		{
			id: 2,
			song: 'John Brown',
			artist: 'New York No. 1 Lake Park mNew York No. 1 Lake Park New York No. 1 Lake Park',
			year: '2019'
		},
		{
			id: 1,
			song: 'New York No. 1 Lake Park mNew York No. 1 Lake Park New York No. 1 Lake Park',
			artist: 'John Brown',
			year: '2019'
		},
		{
			id: 2,
			song: 'John Brown',
			artist: 'New York No. 1 Lake Park mNew York No. 1 Lake Park New York No. 1 Lake Park',
			year: '2019'
		}
	];
</script>

<section class="max-w-screen mt-5">
	<div class="rounded-lg border border-gray-100 bg-white p-5 shadow">
		<h4 class="mb-6 text-xl font-bold text-black">Top Channels</h4>

		<table class=" w-full table-fixed divide-y divide-gray-200">
			<thead class="bg-gray-100">
				<tr>
					<th scope="col" class=" px-6 py-3 text-start text-xs font-medium uppercase text-gray-500">
						Song
					</th>
					<th scope="col" class="px-6 py-3 text-start text-xs font-medium uppercase text-gray-500">
						Artist
					</th>
					<th
						scope="col"
						class="w-30 px-6 py-3 text-start text-xs font-medium uppercase text-gray-500"
					>
						Year
					</th>
					<th
						scope="col"
						class="w-30 px-6 py-3 text-center text-xs font-medium uppercase text-gray-500"
					>
						Action
					</th>
				</tr>
			</thead>
			<tbody>
				{#each tableChannels as item}
					<tr class="odd:bg-white even:bg-gray-50">
						<td class="truncate px-6 py-4 text-sm text-gray-800">{item.song}</td>
						<td class="truncate px-6 py-4 text-sm text-gray-800">{item.artist}</td>
						<td class="truncate px-6 py-4 text-sm text-gray-800">{item.year}</td>
						<td class="flex justify-center space-x-4 truncate px-6 py-4 text-sm">
							<Pencil class="h-5 w-5 text-violet-400 hover:cursor-pointer" />
							<Trash class="h-5 w-5 text-red-400 hover:cursor-pointer" />
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</section>
