<script lang="ts">
  import { Link } from "svelte-routing";
  import Head from "../components/Head.svelte";

  let formData = { email: "", password: "" };
  let isLoading = false;

  async function handleSubmit(_: SubmitEvent) {
    isLoading = true;

    setTimeout(() => {
      isLoading = false;
      formData = { email: "", password: "" }; // Reset form
    }, 1000);
  }
</script>

<Head title="Sign Up" />

<main class="flex min-h-screen items-center justify-center bg-gray-100 p-4">
  <div class="w-full max-w-md rounded-xl bg-white p-8 shadow-lg">
    <h2 class="mb-6 text-center text-2xl font-bold text-gray-800">Sign Up</h2>
    <form class="space-y-4" on:submit|preventDefault={handleSubmit}>
      <div>
        <label for="email" class="mb-1 block text-sm font-medium text-gray-700"
          >Email</label
        >
        <input
          type="email"
          id="email"
          bind:value={formData.email}
          class="w-full rounded-lg border border-gray-300 px-4 py-2 outline-none transition-all focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500"
          placeholder="your@email.com"
        />
      </div>

      <div>
        <label
          for="password"
          class="mb-1 block text-sm font-medium text-gray-700">Password</label
        >
        <input
          type="password"
          id="password"
          bind:value={formData.password}
          class="w-full rounded-lg border border-gray-300 px-4 py-2 outline-none transition-all focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500"
          placeholder="••••••••"
        />
      </div>

      <p class="text-xs text-gray-600">
        By clicking Sign Up, you agree to our Terms and Conditions.
      </p>

      <button
        type="submit"
        disabled={isLoading}
        class="w-full rounded-lg py-2.5 font-medium text-white transition-colors {isLoading
          ? 'cursor-not-allowed bg-indigo-400'
          : 'bg-indigo-600 hover:bg-indigo-700'}"
      >
        {isLoading ? "Loading..." : "Sign Up"}
      </button>
    </form>

    <div class="mt-6 text-center text-sm text-gray-600">
      Have an account? <Link
        to="/login"
        class="font-medium text-indigo-600 hover:text-indigo-500"
      >
        Sign in
      </Link>
    </div>
  </div>
</main>
