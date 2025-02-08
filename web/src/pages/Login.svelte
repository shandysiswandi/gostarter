<script lang="ts">
  import { Link, navigate } from "svelte-routing";
  import { onMount } from "svelte";
  import Head from "../components/Head.svelte";

  let email = "";
  let password = "";
  let rememberMe = false;
  let isLoading = false;

  onMount(() => {
    const savedEmail = localStorage.getItem("rememberedEmail");
    if (savedEmail) {
      email = savedEmail;
      rememberMe = true;
    }
  });

  async function handleSubmit(_: SubmitEvent) {
    isLoading = true;

    setTimeout(() => {
      if (rememberMe) localStorage.setItem("rememberedEmail", email);
      else localStorage.removeItem("rememberedEmail");

      isLoading = false;
      email = "";
      password = "";
      rememberMe = false;
      navigate("/dashboard");
    }, 1000);
  }
</script>

<Head title="Sign In" />

<main class="flex min-h-screen items-center justify-center bg-gray-100 p-4">
  <div class="w-full max-w-md rounded-xl bg-white p-8 shadow-lg">
    <h2 class="mb-6 text-center text-2xl font-bold text-gray-800">Sign In</h2>
    <form class="space-y-4" on:submit|preventDefault={handleSubmit}>
      <div>
        <label for="email" class="mb-1 block text-sm font-medium text-gray-700"
          >Email</label
        >
        <input
          type="email"
          id="email"
          bind:value={email}
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
          bind:value={password}
          class="w-full rounded-lg border border-gray-300 px-4 py-2 outline-none transition-all focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500"
          placeholder="••••••••"
        />
      </div>

      <div class="flex items-center justify-between">
        <label class="flex items-center">
          <input
            type="checkbox"
            bind:checked={rememberMe}
            class="rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
          />
          <span class="ml-2 text-sm text-gray-600">Remember me</span>
        </label>
        <Link
          to="/forgot-password"
          class="text-sm text-indigo-600 hover:text-indigo-500"
        >
          Forgot password?
        </Link>
      </div>

      <button
        type="submit"
        disabled={isLoading}
        class="w-full rounded-lg py-2.5 font-medium text-white transition-colors {isLoading
          ? 'cursor-not-allowed bg-indigo-400'
          : 'bg-indigo-600 hover:bg-indigo-700'}"
      >
        {isLoading ? "Loading..." : "Sign In"}
      </button>
    </form>

    <div class="mt-6 text-center text-sm text-gray-600">
      Don't have an account? <Link
        to="/register"
        class="font-medium text-indigo-600 hover:text-indigo-500"
      >
        Sign up
      </Link>
    </div>
  </div>
</main>
