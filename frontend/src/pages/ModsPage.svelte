<script lang="ts">
  import { Search, Plus } from "@lucide/svelte";
  import { Button } from "$lib/components/ui/button/index";
  import { Input } from "$lib/components/ui/input/index";
  import * as Table from "$lib/components/ui/table/index";
  import SortButton from "../components/SortButton.svelte";
  import { GetData } from "../../wailsjs/go/Main/App";

  let searchQuery = $state("");
  let selectedCategory = $state<string>("All");
  let sortField = $state<"name" | "category" | "active">("name");
  let sortDirection = $state<"asc" | "desc">("asc");
  const categories = ["All", "Cars", "Tracks"];

  interface Mod {
    dir: string;
    name: string;
    category: string;
    active: boolean;
    inProfile: boolean;
    lastModified: string;
  }

  let mods = $state<Mod[]>([]);

  function toggleSort(field: "name" | "category" | "active") {
    if (sortField === field) {
      sortDirection = sortDirection === "asc" ? "desc" : "asc";
    } else {
      sortField = field;
      sortDirection = "asc";
    }
  }

  function getFilteredMods() {
    let result = mods.filter((m) => {
      const matchesSearch = m.name
        .toLowerCase()
        .includes(searchQuery.toLowerCase());
      const matchesCategory =
        selectedCategory === "All" ||
        m.category.toLowerCase() === selectedCategory.toLowerCase();
      return matchesSearch && matchesCategory;
    });

    result.sort((a, b) => {
      let comparison = 0;
      if (sortField === "name") {
        comparison = a.name.localeCompare(b.name);
      }
      if (sortField === "category") {
        comparison = a.category.localeCompare(b.category);
      }
      if (sortField === "active") {
        comparison = a.active === b.active ? 0 : a.active ? -1 : 1;
      }
      return sortDirection === "asc" ? comparison : -comparison;
    });

    return result;
  }

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

  $effect(() => {
    loadMods();
  });
</script>

<div class="p-6 mx-50">
  <div class="flex items-center justify-between mb-6 gap-4">
    <div class="flex items-center gap-4">
      <!-- Search bar -->
      <div class="relative">
        <Search
          class="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground"
        />
        <Input
          type="text"
          placeholder="Search mods..."
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

    <!-- Add mod button -->
    <div>
      <Button class="h-8 min-w-30">
        <Plus class="size-4 mr-2" />
        Add Mod
      </Button>
    </div>
  </div>

  <!-- Table -->
  <div class="bg-input/30">
    <!-- Header Table (non-scrollable) -->
    <div class="bg-background">
      <Table.Root>
        <Table.Header>
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
            <Table.Head class="w-[20%]">
              <div class="flex justify-center">
                <SortButton
                  displayName="CATEGORY"
                  catName="category"
                  {toggleSort}
                  {sortField}
                  {sortDirection}
                ></SortButton>
              </div>
            </Table.Head>
            <Table.Head class="w-[20%]">
              <div class="flex justify-center">
                <SortButton
                  displayName="ACTIVE"
                  catName="active"
                  {toggleSort}
                  {sortField}
                  {sortDirection}
                ></SortButton>
              </div>
            </Table.Head>
          </Table.Row>
        </Table.Header>
      </Table.Root>
    </div>

    <!-- Scrollable Body Table -->
    <div
      class="max-h-[calc(100vh-280px)] overflow-y-auto border border-input border-t-0"
    >
      <Table.Root>
        <Table.Body>
          {#each getFilteredMods() as mod}
            <Table.Row
              class="border-l-4! {mod.active ? ' border-l-primary' : ''} "
            >
              <Table.Cell class="w-[60%] border-r">{mod.name}</Table.Cell>
              <Table.Cell class="w-[20%] border-r text-center"
                >{mod.category}</Table.Cell
              >
              <Table.Cell class="w-[20%] text-center"
                >{mod.active ? "Active" : "Inactive"}</Table.Cell
              >
            </Table.Row>
          {/each}
        </Table.Body>
      </Table.Root>
    </div>
  </div>
</div>
