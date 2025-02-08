<script lang="ts">
  import Head from "../components/Head.svelte";

  let formData = { newPassword: "", confirmPassword: "" };
  let isLoading = false;

  const handleSubmit = (_: SubmitEvent) => {
    isLoading = true;

    setTimeout(() => {
      isLoading = false;
      formData = { newPassword: "", confirmPassword: "" }; // Reset form
    }, 1000);
  };
</script>

<Head title=" Reset your password" />

<div class="flex min-h-screen items-center justify-center bg-gray-100 p-4">
  <div class="w-full max-w-md rounded-xl bg-white p-8 shadow-lg">
    <h2 class="mb-2 text-center text-2xl font-bold text-gray-900">
      Reset your password
    </h2>
    <p class="mb-6 text-center text-sm text-gray-600">
      Enter a new password below to change your password
    </p>

    <form class="space-y-4" on:submit|preventDefault={handleSubmit}>
      <div>
        <label
          for="newPassword"
          class="mb-1 block text-sm font-medium text-gray-700"
        >
          New Password
        </label>
        <input
          type="password"
          id="newPassword"
          name="newPassword"
          bind:value={formData.newPassword}
          class="w-full rounded-lg border border-gray-300 px-4 py-2 outline-none transition-all focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500"
          placeholder="••••••••"
          required
        />
      </div>

      <div>
        <label
          for="confirmPassword"
          class="mb-1 block text-sm font-medium text-gray-700"
        >
          Confirm Password
        </label>
        <input
          type="password"
          id="confirmPassword"
          name="confirmPassword"
          bind:value={formData.confirmPassword}
          class="w-full rounded-lg border border-gray-300 px-4 py-2 outline-none transition-all focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500"
          placeholder="••••••••"
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
        {isLoading ? "Loading..." : "Change Password"}
      </button>
    </form>
  </div>
</div>
