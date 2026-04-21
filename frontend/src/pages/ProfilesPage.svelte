<script lang="ts">
  import {
    Search,
    Play,
    Pencil,
    Trash2,
    EllipsisVertical,
    Save,
    Lock,
  } from "@lucide/svelte";
  import { onMount } from "svelte";
  import { EventsOn, EventsOff } from "../../wailsjs/runtime/runtime";
  import { Button } from "$lib/components/ui/button/index";
  import { Input } from "$lib/components/ui/input/index";
  import * as Card from "$lib/components/ui/card/index";
  import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index";
  import SortButton from "../components/SortButton.svelte";
  import AddProfileDialog from "../components/AddProfileDialog.svelte";
  import EditProfileDialog from "../components/EditProfileDialog.svelte";

  import {
    GetData,
    DeleteProfile,
    UpdateProfiles,
  } from "../../wailsjs/go/Main/App";

  interface Profile {
    id: number;
    name: string;
    category: string;
    path: string;
    modCount: number;
    active: boolean;
    auto_created: boolean;
  }

  let data = $state<any>(null);
  let searchQuery = $state("");
  let selectedCategory = $state<string>("All");
  let sortField = $state<"name" | "category" | "modCount" | "active">("name");
  let sortDirection = $state<"asc" | "desc">("asc");

  let profiles = $state<Profile[]>([]);
  let originalProfiles = $state<Profile[]>([]);

  let changedProfiles = $derived(
    profiles.filter((p) => {
      const updated = originalProfiles.find((o) => o.id === p.id);
      return updated && updated.active !== p.active;
    }),
  );

  const categories = ["All", "Cars", "Tracks"];

  const gridTable = "grid grid-cols-[4fr_120px_120px_180px]";

  async function loadProfiles() {
    try {
      const rawData = await GetData();
      data = rawData;
      profiles = data.Profiles.map((p: any) => ({
        id: p.Id,
        name: p.Name,
        path: p.Path,
        category: p.Category,
        modCount: getModCount(p.Id),
        active: p.Active,
        auto_created: p.AutoCreated,
      }));
      originalProfiles = profiles.map((p) => ({ ...p }));
    } catch (error) {
      console.error("Failed to load profiles:", error);
    }
  }

  function getModsInProfile(profileId: number): string[] {
    if (!data || !data.ModProfiles) {
      return [];
    }
    const modsInProfile = data.ModProfiles.filter(
      (mp: any) => mp.ProfileId === profileId,
    );
    return modsInProfile.map((mp: any) => mp.ModDir);
  }

  function getModCount(profileId: number): number {
    return getModsInProfile(profileId).length;
  }

  function getFilteredProfiles() {
    let result = profiles.filter((p) => {
      const matchesSearch = p.name
        .toLowerCase()
        .includes(searchQuery.toLowerCase());
      const matchesCategory =
        selectedCategory === "All" || p.category === selectedCategory;
      return matchesSearch && matchesCategory;
    });

    result.sort((a, b) => {
      let comparison = 0;
      if (sortField === "name") {
        comparison = a.name.localeCompare(b.name);
      } else if (sortField === "category") {
        comparison = a.category.localeCompare(b.category);
      } else if (sortField === "modCount") {
        comparison = a.modCount - b.modCount;
      } else if (sortField === "active") {
        comparison = a.active === b.active ? 0 : a.active ? -1 : 1;
      }
      return sortDirection === "asc" ? comparison : -comparison;
    });

    return result;
  }

  function toggleSort(field: "name" | "category" | "modCount" | "active") {
    if (sortField === field) {
      sortDirection = sortDirection === "asc" ? "desc" : "asc";
    } else {
      sortField = field;
      sortDirection = "asc";
    }
  }

  function toggleActive(profile: Profile) {
    profile.active = !profile.active;
  }

  async function deleteProfile(profile: Profile) {
    try {
      await DeleteProfile(profile.id);
    } catch (error) {
      console.error("Failed to delete profile:", error);
    }
  }

  let editDialogOpen = $state(false);
  let editingProfile = $state<Profile | null>(null);

  function editProfile(profile: Profile) {
    editingProfile = profile;
    editDialogOpen = true;
  }

  async function handleSave() {
    const profileToSave = changedProfiles.map((p) => ({
      Id: p.id,
      Name: p.name,
      Path: p.path,
      Category: p.category,
      Active: p.active,
      AutoCreated: p.auto_created,
    }));

    try {
      await UpdateProfiles(profileToSave);
    } catch (error) {
      console.error("Failed to update profiles:", error);
    }
  }

  onMount(() => {
    loadProfiles();

    EventsOn("data-updated", async () => {
      profiles = [];
      await loadProfiles();
    });

    return () => {
      EventsOff("data-updated");
    };
  });

  $effect(() => {
    if (!editDialogOpen) {
      editingProfile = null;
    }
  });
</script>

<!-- Main content -->
<div class="p-6 mx-50">
  <!-- Title and menus -->
  <div class="flex items-center justify-between mb-6 gap-4">
    <div class="flex items-center gap-4">
      <!-- Search bar -->
      <div class="relative">
        <Search
          class="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground"
        />
        <Input
          type="text"
          placeholder="Search profiles..."
          class="pl-10 min-w-64 h-8"
          bind:value={searchQuery}
        />
      </div>

      <!-- Category toggle -->
      <div
        class="flex items-center bg-secondary rounded-md p-0.5 border border-input h-8"
      >
        {#each categories as category}
          <Button
            variant={selectedCategory === category ? "default" : "ghost"}
            onclick={() => (selectedCategory = category)}
            class="px-3 text-xs"
          >
            {category === "All" ? "All" : category}
          </Button>
        {/each}
      </div>
    </div>

    <!-- Add profile button -->
    <div>
      <AddProfileDialog />
    </div>
  </div>

  <!-- Edit profile dialog -->
  <EditProfileDialog
    bind:open={editDialogOpen}
    profileId={String(editingProfile?.id) || ""}
    initialProfileName={editingProfile?.name}
    initialCategory={editingProfile?.category}
    preselectedMods={new Set(getModsInProfile(editingProfile?.id || -1))}
  />

  <!-- Main table -->
  <div class="flex flex-col max-h-[calc(100vh-200px)]">
    <!-- Sort table - fixed header -->
    <div
      class="{gridTable} gap-6 items-center px-5 pb-4 text-xs font-medium text-muted-foreground border-b border-border shrink-0"
    >
      <!-- Profile name sort -->
      <div class="flex items-center">
        <SortButton
          displayName="PROFILE NAME"
          catName="name"
          {toggleSort}
          {sortField}
          {sortDirection}
          center={false}
        ></SortButton>
      </div>

      <!-- Category sort -->
      <SortButton
        displayName="&ensp;CATEGORY"
        catName="category"
        {toggleSort}
        {sortField}
        {sortDirection}
      ></SortButton>

      <!-- Mod count sort -->
      <SortButton
        displayName="&ensp;MOD COUNT"
        catName="modCount"
        {toggleSort}
        {sortField}
        {sortDirection}
      ></SortButton>

      <!-- Active sort -->
      <SortButton
        displayName="&ensp;ACTIVE"
        catName="active"
        {toggleSort}
        {sortField}
        {sortDirection}
      ></SortButton>
    </div>

    <!-- Profile list - scrollable content -->
    <div class="overflow-y-auto space-y-2">
      {#each getFilteredProfiles() as profile (profile.id)}
        <Card.Root
          class="bg-input/30 transition-colors border-l-6 hover:bg-accent/50 data-[active=true]:bg-muted data-[active=true]:hover:bg-accent data-[active=true]:border-primary"
          data-active={profile.active ? "true" : "false"}
        >
          <Card.Content class="{gridTable} gap-6 items-center">
            <!-- Profile name -->
            <div class="border-r border-border pr-4 flex items-center gap-2">
              <p class="text-sm">{profile.name}</p>
              {#if profile.auto_created}
                <div
                  class="flex items-center gap-1 text-xs text-muted-foreground bg-secondary px-2 py-0.5 rounded"
                >
                  <Lock class="w-3 h-3" />
                  Auto
                </div>
              {/if}
            </div>

            <!-- Category -->
            <div class="border-r border-border pr-4 text-center">
              <p class="text-sm">{profile.category}</p>
            </div>

            <!-- Mod count -->
            <div class="border-r border-border pr-4 text-center">
              <p class="text-sm">{profile.modCount} mods</p>
            </div>

            <!-- Active button and dropdown menu -->
            <div class="flex items-center">
              <div class="flex flex-1 justify-center px-2">
                <Button
                  variant={profile.active ? "default" : "outline"}
                  size="sm"
                  onclick={() => toggleActive(profile)}
                  class="min-w-30 cursor-pointer h-8"
                >
                  <Play class="w-5 mr-1" />
                  {profile.active ? "Active" : "Activate"}
                </Button>
              </div>

              <!-- Dropdown menu -->
              <div class="flex flex-1 justify-center">
                {#if profile.auto_created}
                  <div class="h-8 w-8"></div>
                {:else}
                  <DropdownMenu.Root>
                    <DropdownMenu.Trigger>
                      <Button variant="ghost" size="sm" class="h-8 w-8">
                        <EllipsisVertical class="size-4" />
                      </Button>
                    </DropdownMenu.Trigger>
                    <DropdownMenu.Content>
                      <DropdownMenu.Item onclick={() => editProfile(profile)}>
                        <Pencil class="size-4 mr-2" />
                        Edit
                      </DropdownMenu.Item>
                      <DropdownMenu.Separator />
                      <DropdownMenu.Item
                        class="text-destructive focus:text-destructive"
                        onclick={() => deleteProfile(profile)}
                      >
                        <Trash2 class="size-4 mr-2" />
                        Delete
                      </DropdownMenu.Item>
                    </DropdownMenu.Content>
                  </DropdownMenu.Root>
                {/if}
              </div>
            </div>
          </Card.Content>
        </Card.Root>
      {/each}
    </div>
  </div>
</div>

{#if changedProfiles.length > 0}
  <div
    class="fixed bottom-6 left-1/2 -translate-x-1/2 z-50 flex items-center gap-3 bg-primary border border-border rounded-lg shadow-lg"
  >
    <Button
      onclick={handleSave}
      variant="default"
      size="sm"
      class="cursor-pointer"
    >
      <Save class="w-4 h-4 mr-2" />
      Save Changes
    </Button>
  </div>
{/if}
