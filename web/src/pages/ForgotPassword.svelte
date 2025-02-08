<script lang="ts">
  import Head from "../components/Head.svelte";

  let email = "";
  let isLoading = false;
  let isSubmitted = false;

  const handleSubmit = (_: SubmitEvent) => {
    isSubmitted = false;
    isLoading = true;
    email = "";

    setTimeout(() => {
      isLoading = false;
      isSubmitted = true;
    }, 1000);
  };
</script>

<Head title="Forgot Password" />

<div class="flex min-h-screen items-center justify-center bg-gray-100 p-4">
  <div class="w-full max-w-md rounded-xl bg-white p-8 shadow-lg">
    <h2 class="mb-2 text-center text-2xl font-bold text-gray-900">
      Forgot Password
    </h2>
    <h4 class="mb-6 text-center text-sm text-gray-600">
      We will send a password reset link to your email.
    </h4>

    {#if isSubmitted}
      <div class="mb-4 rounded-lg bg-green-50 p-4 text-sm text-green-700">
        We have sent you an email with instructions to reset your password.
      </div>
    {/if}

    <form class="space-y-4" on:submit|preventDefault={handleSubmit}>
      <div>
        <label for="email" class="mb-1 block text-sm font-medium text-gray-700">
          Email
        </label>
        <input
          type="email"
          id="email"
          name="email"
          bind:value={email}
          class="w-full rounded-lg border border-gray-300 px-4 py-2 outline-none transition-all focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500"
          placeholder="your@email.com"
          required
        />
      </div>

      <button
        type="submit"
        class="w-full rounded-lg py-2.5 font-medium text-white transition-colors {isLoading
          ? 'cursor-not-allowed bg-indigo-400'
          : 'bg-indigo-600 hover:bg-indigo-700'}"
        disabled={isLoading}
      >
        {isLoading ? "Loading..." : "Send Reset Instructions"}
      </button>
    </form>

    <div class="mt-6 text-center text-sm text-gray-600">
      Don't have an account?
      <a
        href="/register"
        class="font-medium text-indigo-600 hover:text-indigo-500"
      >
        Sign up
      </a>
    </div>
  </div>
</div>
