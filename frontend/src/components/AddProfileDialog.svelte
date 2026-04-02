<script lang="ts">
  import { Button, buttonVariants } from "$lib/components/ui/button/index";
  import Checkbox from "$lib/components/ui/checkbox/checkbox.svelte";
  import * as Dialog from "$lib/components/ui/dialog/index";
  import * as Table from "$lib/components/ui/table/index";
  import { Input } from "$lib/components/ui/input/index";
  import Label from "$lib/components/ui/label/label.svelte";
  import { Plus, Search } from "@lucide/svelte";
  import SortButton from "./SortButton.svelte";

  let profileName = $state<string>("");
  let selectedMods = $state<Set<string>>(new Set());

  let selectedCategory = $state<string>("Cars");
  let searchQuery = $state<string>("");
  let sortField = $state<"name" | "lastModified">("name");
  let sortDirection = $state<"asc" | "desc">("asc");
  let hideUsedMods = $state<boolean>(false);

  interface Mod {
    id: string;
    name: string;
    category: string;
    inProfile: boolean;
    lastModified: string;
  }

  const mockMods: Mod[] = $state([
    {
      id: "1",
      name: "esda_mk5_public",
      category: "Cars",
      inProfile: false,
      lastModified: "2024-05-28",
    },
    {
      id: "2",
      name: "alfa_romeo_156",
      category: "Cars",
      inProfile: false,
      lastModified: "2024-06-03",
    },
    {
      id: "3",
      name: "ariel_atom_500",
      category: "Cars",
      inProfile: true,
      lastModified: "2024-05-15",
    },
    {
      id: "4",
      name: "acu_okutama-circuit",
      category: "Tracks",
      inProfile: true,
      lastModified: "2024-06-10",
    },
    {
      id: "5",
      name: "sdv_bmw_e46_330ci_rhd",
      category: "Cars",
      inProfile: true,
      lastModified: "2024-05-22",
    },
    {
      id: "6",
      name: "shuto_revival_project_beta",
      category: "Tracks",
      inProfile: false,
      lastModified: "2024-06-05",
    },
    {
      id: "7",
      name: "swarm_fluffs_mx5",
      category: "Cars",
      inProfile: true,
      lastModified: "2024-05-19",
    },
    {
      id: "8",
      name: "The Springs",
      category: "Tracks",
      inProfile: false,
      lastModified: "2024-06-12",
    },
    {
      id: "9",
      name: "wdts_toyota_cresta_jzx100",
      category: "Cars",
      inProfile: true,
      lastModified: "2024-05-25",
    },
    {
      id: "10",
      name: "fumi_cp_bmw_e36_330i_sedan",
      category: "Cars",
      inProfile: true,
      lastModified: "2024-06-01",
    },
    {
      id: "11",
      name: "union_island",
      category: "Tracks",
      inProfile: true,
      lastModified: "2024-05-30",
    },
    {
      id: "12",
      name: "pcp2_fumi_cp24_lex_is200_v8",
      category: "Cars",
      inProfile: false,
      lastModified: "2024-06-07",
    },
    {
      id: "13",
      name: "esda_mk5_public",
      category: "Cars",
      inProfile: false,
      lastModified: "2024-05-28",
    },
    {
      id: "14",
      name: "alfa_romeo_156",
      category: "Cars",
      inProfile: false,
      lastModified: "2024-06-03",
    },
    {
      id: "15",
      name: "ariel_atom_500",
      category: "Cars",
      inProfile: true,
      lastModified: "2024-05-15",
    },
    {
      id: "16",
      name: "acu_okutama-circuit",
      category: "Tracks",
      inProfile: true,
      lastModified: "2024-06-10",
    },
    {
      id: "17",
      name: "sdv_bmw_e46_330ci_rhd",
      category: "Cars",
      inProfile: true,
      lastModified: "2024-05-22",
    },
    {
      id: "18",
      name: "shuto_revival_project_beta",
      category: "Tracks",
      inProfile: false,
      lastModified: "2024-06-05",
    },
    {
      id: "19",
      name: "swarm_fluffs_mx5",
      category: "Cars",
      inProfile: true,
      lastModified: "2024-05-19",
    },
    {
      id: "20",
      name: "The Springs",
      category: "Tracks",
      inProfile: false,
      lastModified: "2024-06-12",
    },
    {
      id: "21",
      name: "wdts_toyota_cresta_jzx100",
      category: "Cars",
      inProfile: true,
      lastModified: "2024-05-25",
    },
    {
      id: "22",
      name: "fumi_cp_bmw_e36_330i_sedan",
      category: "Cars",
      inProfile: true,
      lastModified: "2024-06-01",
    },
    {
      id: "23",
      name: "union_island",
      category: "Tracks",
      inProfile: true,
      lastModified: "2024-05-30",
    },
    {
      id: "24",
      name: "pcp2_fumi_cp24_lex_is200_v8",
      category: "Cars",
      inProfile: false,
      lastModified: "2024-06-07",
    },
  ]);

  function getFilteredMods() {
    let result = mockMods.filter((m) => {
      const matchesSearch = m.name
        .toLowerCase()
        .includes(searchQuery.toLowerCase());
      const matchesCategory =
        selectedCategory === "All" || m.category === selectedCategory;
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
  }
</script>

<Dialog.Root
  onOpenChange={(open) => {
    if (!open) resetForm();
  }}
>
  <form>
    <Dialog.Trigger>
      <Button class="h-8 min-w-30">
        <Plus class="size-4 mr-2" />
        Add Profile
      </Button>
    </Dialog.Trigger>
    <Dialog.Content class="sm:max-w-150">
      <Dialog.Header>
        <Dialog.Title>Add Profile</Dialog.Title>
        <Dialog.Description>Create new mod profile.</Dialog.Description>
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
                      <Table.Row onclick={() => toggleSelection(mod.id)}>
                        <Table.Cell class="w-[60%]">{mod.name}</Table.Cell>
                        <Table.Cell class="w-[30%] text-center"
                          >{mod.lastModified}</Table.Cell
                        >
                        <Table.Cell class="w-[10%]">
                          <Checkbox
                            checked={selectedMods.has(mod.id)}
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
            <Button type="submit">Create</Button>
          </Dialog.Footer>
        </div>
      </div>
    </Dialog.Content>
  </form>
</Dialog.Root>
