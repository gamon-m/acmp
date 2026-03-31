<script lang="ts">
  import { Search, Plus } from "@lucide/svelte";
  import * as Select from "$lib/components/ui/select/index";
  import { Button } from "$lib/components/ui/button/index";
  import { Input } from "$lib/components/ui/input/index";
  import * as Table from "$lib/components/ui/table/index";
  import SortButton from "../components/SortButton.svelte";

  let searchQuery = $state("");
  let selectedCategory = $state<string>("All");
  let sortField = $state<"name" | "category" | "active">("name");
  let sortDirection = $state<"asc" | "desc">("asc");
  const categories = ["All", "Cars", "Tracks"];

  interface Mod {
    id: string;
    name: string;
    category: string;
    active: boolean;
  }

  const mockMods: Mod[] = $state([
    { id: "1", name: "esda_mk5_public", category: "Cars", active: false },
    { id: "2", name: "alfa_romeo_156", category: "Cars", active: false },
    { id: "3", name: "ariel_atom_500", category: "Cars", active: true },
    { id: "4", name: "acu_okutama-circuit", category: "Tracks", active: true },
    { id: "5", name: "sdv_bmw_e46_330ci_rhd", category: "Cars", active: true },
    {
      id: "6",
      name: "shuto_revival_project_beta",
      category: "Tracks",
      active: false,
    },
    { id: "7", name: "swarm_fluffs_mx5", category: "Cars", active: true },
    { id: "8", name: "The Springs", category: "Tracks", active: false },
    {
      id: "9",
      name: "wdts_toyota_cresta_jzx100",
      category: "Cars",
      active: true,
    },
    {
      id: "10",
      name: "fumi_cp_bmw_e36_330i_sedan",
      category: "Cars",
      active: true,
    },
    { id: "11", name: "union_island", category: "Tracks", active: true },
    {
      id: "12",
      name: "pcp2_fumi_cp24_lex_is200_v8",
      category: "Cars",
      active: false,
    },
  ]);

  function toggleSort(field: "name" | "category" | "active") {
    if (sortField === field) {
      sortDirection = sortDirection === "asc" ? "desc" : "asc";
    } else {
      sortField = field;
      sortDirection = "asc";
    }
  }

  function getFilteredMods() {
    let result = mockMods.filter((m) => {
      const matchesSearch = m.name
        .toLowerCase()
        .includes(searchQuery.toLowerCase());
      const matchesCategory =
        selectedCategory === "All" || m.category === selectedCategory;
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

    <!-- Add mod button -->
    <div>
      <Button class="h-8 min-w-30">
        <Plus class="size-4 mr-2" />
        Add Mod
      </Button>
    </div>
  </div>

  <!-- Table -->
  <div>
    <Table.Root>
      <Table.Header class="bg-background hover:bg-background">
        <Table.Row class="hover:bg-background">
          <Table.Head>
            <SortButton
              displayName="MOD NAME"
              catName="name"
              {toggleSort}
              {sortField}
              {sortDirection}
            ></SortButton>
          </Table.Head>
          <Table.Head class="w-50">
            <div class="flex justify-center">
              <SortButton
                displayName="&emsp;CATEGORY"
                catName="category"
                {toggleSort}
                {sortField}
                {sortDirection}
              ></SortButton>
            </div>
          </Table.Head>
          <Table.Head class="w-50">
            <div class="flex justify-center">
              <SortButton
                displayName="&emsp;ACTIVE"
                catName="active"
                {toggleSort}
                {sortField}
                {sortDirection}
              ></SortButton>
            </div>
          </Table.Head>
        </Table.Row>
      </Table.Header>
      <Table.Body class="bg-input/30">
        {#each getFilteredMods() as mod}
          <Table.Row
            class="border-l-4! {mod.active ? ' border-l-primary' : ''} "
          >
            <Table.Cell class="border-r">{mod.name}</Table.Cell>
            <Table.Cell class="border-r text-center">{mod.category}</Table.Cell>
            <Table.Cell class="text-center"
              >{mod.active ? "Active" : "Inactive"}</Table.Cell
            >
          </Table.Row>
        {/each}
      </Table.Body>
    </Table.Root>
  </div>

  <!-- Empty space -->
  <div class="min-h-20"></div>
</div>
