<script lang="ts">
  import {
    Search,
    Play,
    Pencil,
    Trash2,
    EllipsisVertical,
  } from "@lucide/svelte";
  import { Button } from "$lib/components/ui/button/index";
  import { Input } from "$lib/components/ui/input/index";
  import * as Card from "$lib/components/ui/card/index";
  import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index";
  import * as Select from "$lib/components/ui/select/index";
  import SortButton from "../components/SortButton.svelte";
  import AddProfileDialog from "../components/AddProfileDialog.svelte";

  interface Profile {
    id: string;
    name: string;
    category: string;
    modCount: number;
    active: boolean;
  }

  let searchQuery = $state("");
  let selectedCategory = $state<string>("All");
  let sortField = $state<"name" | "category" | "modCount" | "active">("name");
  let sortDirection = $state<"asc" | "desc">("asc");

  const mockProfiles: Profile[] = $state([
    {
      id: "1",
      name: "Drift",
      category: "Cars",
      modCount: 12,
      active: false,
    },
    {
      id: "2",
      name: "Touge",
      category: "Tracks",
      modCount: 45,
      active: true,
    },
    {
      id: "3",
      name: "Other",
      category: "Cars",
      modCount: 8,
      active: false,
    },
    {
      id: "4",
      name: "Drift",
      category: "Tracks",
      modCount: 167,
      active: false,
    },
    {
      id: "5",
      name: "Drift",
      category: "Cars",
      modCount: 12,
      active: false,
    },
    {
      id: "6",
      name: "Touge",
      category: "Tracks",
      modCount: 45,
      active: true,
    },
    {
      id: "7",
      name: "Other",
      category: "Cars",
      modCount: 8,
      active: false,
    },
    {
      id: "8",
      name: "Drift",
      category: "Tracks",
      modCount: 167,
      active: false,
    },
  ]);

  const categories = ["All", "Cars", "Tracks"];

  const gridTable = "grid grid-cols-[4fr_120px_120px_180px]";

  function getFilteredProfiles() {
    let result = mockProfiles.filter((p) => {
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

  function deleteProfile(profile: Profile) {
    console.log("deleting profile:", profile.name);
  }

  function editProfile(profile: Profile) {
    console.log("deleting profile:", profile.name);
  }
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

      <!-- Category selection -->
      <Select.Root type="single" bind:value={selectedCategory}>
        <Select.Trigger class="w-40 min-h-8"
          >{selectedCategory === "All"
            ? "All Categories"
            : selectedCategory}</Select.Trigger
        >
        <Select.Content>
          {#each categories as category}
            <Select.Item value={category}>
              {category === "All" ? "All Categories" : category}
            </Select.Item>
          {/each}
        </Select.Content>
      </Select.Root>
    </div>

    <!-- Add profile button -->
    <div>
      <AddProfileDialog />
    </div>
  </div>

  <!-- Main table -->
  <div class="space-y-2">
    <!-- Sort table -->
    <div
      class="{gridTable} gap-6 items-center px-5 pb-4 text-xs font-medium text-muted-foreground border-b border-border"
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

    <!-- Profile list -->
    {#each getFilteredProfiles() as profile (profile.id)}
      <Card.Root
        class="bg-input/30 transition-colors border-l-6 hover:bg-accent/50 data-[active=true]:bg-muted data-[active=true]:hover:bg-accent data-[active=true]:border-primary"
        data-active={profile.active ? "true" : "false"}
      >
        <Card.Content class="{gridTable} gap-6 items-center">
          <!-- Profile name -->
          <div class="border-r border-border pr-4">
            <p class="text-sm">{profile.name}</p>
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
            </div>
          </div>
        </Card.Content>
      </Card.Root>
    {/each}
  </div>

  <!-- Empty space -->
  <div class="min-h-20"></div>
</div>
