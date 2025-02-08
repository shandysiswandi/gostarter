<script lang="ts">
  import { Link } from "svelte-routing";
  import Head from "../components/Head.svelte";
</script>

<Head title="Home" />

<main class="flex h-screen items-center justify-center bg-gray-100">
  <Link
    to="/login"
    class="flex h-40 w-40 items-center justify-center rounded-full bg-violet-600 text-white hover:bg-violet-700 hover:shadow-md"
  >
    <h1 class="text-2xl font-bold">Let's Go</h1>
  </Link>
</main>
