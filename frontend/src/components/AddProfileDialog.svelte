<script lang="ts">
  import { Button, buttonVariants } from "$lib/components/ui/button/index";
  import Checkbox from "$lib/components/ui/checkbox/checkbox.svelte";
  import * as Dialog from "$lib/components/ui/dialog/index";
  import * as Table from "$lib/components/ui/table/index";
  import { Input } from "$lib/components/ui/input/index";
  import Label from "$lib/components/ui/label/label.svelte";
  import { Plus, Search } from "@lucide/svelte";
  import SortButton from "./SortButton.svelte";

  import { GetData, SaveProfile } from "../../wailsjs/go/Main/App";

  let {
    editMode = false,
    profileId = "",
    initialProfileName = "",
    initialCategory = "Cars",
    preselectedMods = new Set<string>(),
    open = $bindable(false),
  } = $props();

  let profileName = $state<string>("");
  let selectedMods = $state<Set<string>>(new Set());
  let selectedCategory = $state<string>("Cars");
  let searchQuery = $state<string>("");
  let sortField = $state<"name" | "lastModified">("name");
  let sortDirection = $state<"asc" | "desc">("asc");
  let hideUsedMods = $state<boolean>(false);

  interface Mod {
    dir: string;
    name: string;
    category: string;
    active: boolean;
    inProfile: boolean;
    lastModified: string;
  }

  let mods = $state<Mod[]>([]);

  async function loadMods() {
    try {
      const data = await GetData();
      mods = data.Mods.map((m) => ({
        dir: m.Dir,
        name: m.Name,
        category: m.Category,
        active: m.Active,
        inProfile: m.InProfile,
        lastModified: m.LastModified,
      }));
    } catch (error) {
      console.error("Failed to load mods:", error);
      alert("Failed to load mods. Please try again.");
    }
  }

  function getFormattedDate(dateTimeStr: string) {
    const date = new Date(dateTimeStr);
    return date.toLocaleDateString();
  }

  function getFilteredMods() {
    let result = mods.filter((m) => {
      const matchesSearch = m.name
        .toLowerCase()
        .includes(searchQuery.toLowerCase());
      const matchesCategory =
        selectedCategory === "All" ||
        m.category.toLowerCase() === selectedCategory.toLowerCase();
      const notInProfiles = !hideUsedMods || !m.inProfile;
      return matchesSearch && matchesCategory && notInProfiles;
    });

    result.sort((a, b) => {
      let comparison = 0;
      if (sortField === "name") {
        comparison = a.name.localeCompare(b.name);
      }
      if (sortField === "lastModified") {
        comparison =
          new Date(a.lastModified).getTime() -
          new Date(b.lastModified).getTime();
      }
      return sortDirection === "asc" ? comparison : -comparison;
    });

    return result;
  }

  function toggleSort(field: "name" | "lastModified") {
    if (sortField === field) {
      sortDirection = sortDirection === "asc" ? "desc" : "asc";
    } else {
      sortField = field;
      sortDirection = "asc";
    }
  }

  function toggleSelection(modId: string) {
    if (selectedMods.has(modId)) {
      selectedMods.delete(modId);
    } else {
      selectedMods.add(modId);
    }
    selectedMods = new Set(selectedMods);
  }

  function modCountLabel(count: number) {
    return `${count} mod${count !== 1 ? "s" : ""} selected`;
  }

  function resetForm() {
    profileName = "";
    selectedCategory = "Cars";
    searchQuery = "";
    sortField = "name";
    sortDirection = "asc";
    selectedMods = new Set();
    open = false;
  }

  function initForm() {
    profileName = initialProfileName;
    selectedCategory = initialCategory;
    selectedMods = new Set(preselectedMods);
    searchQuery = "";
    sortField = "name";
    sortDirection = "asc";
    open = true;
  }

  function handleSubmit() {
    if (profileName.trim() === "") {
      alert("Profile name cannot be empty.");
      return;
    }
    const profileData = {
      id: profileId,
      name: profileName,
      category: selectedCategory,
      mods: Array.from(selectedMods),
    };
    SaveProfile(profileData);
    resetForm();
  }

  $effect(() => {
    if (open) {
      loadMods();
      initForm();
    }
  });
</script>

<Dialog.Root
  bind:open
  onOpenChange={(open) => {
    if (!open) resetForm();
  }}
>
  <form>
    {#if !editMode}
      <Dialog.Trigger>
        <Button class="h-8 min-w-30">
          <Plus class="size-4 mr-2" />
          Add Profile
        </Button>
      </Dialog.Trigger>
    {/if}
    <Dialog.Content class="sm:max-w-150">
      <Dialog.Header>
        <Dialog.Title>{editMode ? "Edit Profile" : "Add Profile"}</Dialog.Title>
        <Dialog.Description
          >{editMode
            ? "Save changes to profile."
            : "Create new mod profile."}</Dialog.Description
        >
      </Dialog.Header>
      <div class="flex items-center gap-2">
        <div class="grid flex-1 gap-4">
          <!-- Profile Name Input -->
          <div class="grid gap-2">
            <Label>Profile Name</Label>
            <Input
              type="text"
              placeholder="Enter profile name..."
              required
              class="w-full"
              bind:value={profileName}
            />
          </div>
          <!-- Category Selection -->
          <div class="grid gap-2">
            <Label>Category</Label>
            <div class="flex gap-2">
              <Button
                variant={selectedCategory === "Cars" ? "default" : "outline"}
                onclick={() => (selectedCategory = "Cars")}
              >
                Cars
              </Button>
              <Button
                variant={selectedCategory === "Tracks" ? "default" : "outline"}
                onclick={() => (selectedCategory = "Tracks")}
              >
                Tracks
              </Button>
            </div>
          </div>

          <!-- Mod table -->
          <div class="grid gap-2">
            <Label>Mods</Label>
            <div class="flex items-center gap-4">
              <!-- Search Input -->
              <div class="relative">
                <Search
                  class="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground"
                />
                <Input
                  type="text"
                  placeholder="Search mods..."
                  class="pl-10 min-w-64"
                  bind:value={searchQuery}
                />
              </div>
              <!-- Show unused mods toggle -->
              <div class="flex items-center gap-2">
                <Checkbox
                  id="hide-used-mods"
                  checked={hideUsedMods}
                  onCheckedChange={() => (hideUsedMods = !hideUsedMods)}
                />
                <Label for="hide-used-mods">Only show unused mods</Label>
              </div>
            </div>

            <!-- Mod List -->
            <div class="bg-input/30 border border-input mt-2">
              <!-- Header Table (non-scrollable) -->
              <Table.Root>
                <Table.Header class="border-b border-input">
                  <Table.Row class="hover:bg-transparent">
                    <Table.Head class="w-[60%]">
                      <SortButton
                        displayName="MOD NAME"
                        catName="name"
                        {toggleSort}
                        {sortField}
                        {sortDirection}
                      ></SortButton>
                    </Table.Head>
                    <Table.Head class="w-[30%]">
                      <div class="flex justify-center">
                        <SortButton
                          displayName="&emsp;LAST MODIFIED"
                          catName="lastModified"
                          {toggleSort}
                          {sortField}
                          {sortDirection}
                          center={true}
                        ></SortButton>
                      </div>
                    </Table.Head>
                    <Table.Head class="w-[10%]"></Table.Head>
                  </Table.Row>
                </Table.Header>
              </Table.Root>

              <!-- Scrollable Body Table -->
              <div class="max-h-80 min-h-80 overflow-y-auto">
                <Table.Root>
                  <Table.Body class="border-b border-input">
                    {#each getFilteredMods() as mod}
                      <Table.Row onclick={() => toggleSelection(mod.dir)}>
                        <Table.Cell class="w-[60%]">{mod.name}</Table.Cell>
                        <Table.Cell class="w-[30%] text-center"
                          >{getFormattedDate(mod.lastModified)}</Table.Cell
                        >
                        <Table.Cell class="w-[10%]">
                          <Checkbox
                            checked={selectedMods.has(mod.dir)}
                            readonly
                          />
                        </Table.Cell>
                      </Table.Row>
                    {/each}
                  </Table.Body>
                </Table.Root>
              </div>
            </div>
            <p class="text-right">
              {modCountLabel(selectedMods.size)}
            </p>
          </div>
          <Dialog.Footer>
            <Dialog.Close
              type="button"
              class={buttonVariants({ variant: "outline" })}
            >
              Cancel
            </Dialog.Close>
            <Button onclick={handleSubmit} type="submit"
              >{editMode ? "Save" : "Create"}</Button
            >
          </Dialog.Footer>
        </div>
      </div>
    </Dialog.Content>
  </form>
</Dialog.Root>
