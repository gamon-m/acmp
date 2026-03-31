<script lang="ts">
  import Navigation from "./components/Navigation.svelte";
  import ProfilesPage from "./pages/ProfilesPage.svelte";
  import ModsPage from "./pages/ModsPage.svelte";
  import SettingsPage from "./pages/SettingsPage.svelte";

  type Page = "profiles" | "mods";

  let activeTab: Page = "profiles";
</script>

<main class="min-h-screen bg-background">
  <Navigation bind:activeTab />

  {#if activeTab === "profiles"}
    <ProfilesPage />
  {:else if activeTab === "mods"}
    <ModsPage />
  {:else if activeTab === "settings"}
    <SettingsPage />
  {/if}
</main>
