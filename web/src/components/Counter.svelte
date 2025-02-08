<script lang="ts">
  let count: number = $state(0);
  const increment = () => {
    count += 1;
  };
</script>

<button
  onclick={increment}
  class="rounded bg-blue-500 px-4 py-2 font-bold text-white hover:bg-blue-700"
>
  count is {count}
</button>
