<script lang="ts">
  import AddProfileDialog from "./AddProfileDialog.svelte";

  let {
    profileId = "",
    initialProfileName = "",
    initialCategory = "Cars",
    preselectedMods = new Set<string>(),
    open = $bindable(false),
  } = $props();
</script>

<AddProfileDialog
  editMode={true}
  {profileId}
  {initialProfileName}
  {initialCategory}
  {preselectedMods}
  bind:open
/>
