<script lang="ts">
  import { Plus } from "@lucide/svelte";
  import { Button, buttonVariants } from "$lib/components/ui/button/index";
  import * as Dialog from "$lib/components/ui/dialog/index";
  import { Input } from "$lib/components/ui/input/index";
  import Label from "$lib/components/ui/label/label.svelte";
  import { Search } from "@lucide/svelte";

  let profileName = $state<string>("");
  let selectedMods = $state<Set<string>>(new Set());

  let selectedCategory = $state<string>("Cars");
  let searchQuery = $state<string>("");
  let sortField = $state<"name" | "lastModified">("name");
  let sortDirection = $state<"asc" | "desc">("asc");
  let showUnusedMods = $state<boolean>(false);

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
  ]);

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
      if (sortField === "lastModified") {
        comparison =
          new Date(a.lastModified).getTime() -
          new Date(b.lastModified).getTime();
      }
      return sortDirection === "asc" ? comparison : -comparison;
    });

    return result;
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
            </div>
            <!-- Mod List -->
            <div class="max-h-80 min-h-80 overflow-y-auto border mt-2">
              {#each getFilteredMods() as mod}
                <div
                  class="flex items-center px-4 py-2 hover:bg-background-light cursor-pointer"
                  class:selected={selectedMods.has(mod.id)}
                  onclick={() => toggleSelection(mod.id)}
                >
                  <span>{mod.name}</span>
                  <input
                    type="checkbox"
                    class="ml-auto mr-4"
                    checked={selectedMods.has(mod.id)}
                  />
                </div>
              {/each}
            </div>
            <p class="flex justify-end">{modCountLabel(selectedMods.size)}</p>
          </div>
        </div>
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
    </Dialog.Content>
  </form>
</Dialog.Root>
