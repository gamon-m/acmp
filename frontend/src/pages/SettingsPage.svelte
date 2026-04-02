<script lang="ts">
  import { FolderOpen, Save } from "@lucide/svelte";
  import { Button } from "$lib/components/ui/button/index";
  import { Input } from "$lib/components/ui/input/index";
  import { Checkbox } from "$lib/components/ui/checkbox/index";
  import * as Card from "$lib/components/ui/card/index";
  import * as Dialog from "$lib/components/ui/dialog/index";
  import { Label } from "$lib/components/ui/label/index";
  import { Separator } from "$lib/components/ui/separator/index";

  let assettoCorsaModFolder = $state("");
  let assettoCorsaRootFolder = $state("");
  let automaticProfiles = $state(false);
  let resetDialogOpen = $state(false);

  async function selectModFolder() {}

  async function selectRootFolder() {}

  async function saveSettings() {}

  async function resetApp() {
    assettoCorsaModFolder = "";
    assettoCorsaRootFolder = "";
    automaticProfiles = false;
    resetDialogOpen = false;
  }

  async function loadSettings() {}

  // Load settings on mount
  $effect(() => {
    loadSettings();
  });
</script>

<div class="p-6 mx-auto max-w-2xl pb-20">
  <h2 class="text-2xl font-semibold mb-4">Settings</h2>

  <div class="space-y-6">
    <!-- Assetto Corsa Folders -->
    <Card.Root>
      <Card.Header>
        <Card.Title>Assetto Corsa Folders</Card.Title>
        <Card.Description>
          Configure the paths to your Assetto Corsa installation
        </Card.Description>
      </Card.Header>
      <Card.Content class="space-y-4">
        <!-- Mod Folder -->
        <div class="space-y-2">
          <Label for="mod-folder">Assetto Corsa Mod Folder</Label>
          <div class="flex gap-2">
            <Input
              id="mod-folder"
              type="text"
              placeholder="Select mod folder..."
              value={assettoCorsaModFolder}
              class="flex-1"
              readonly
            />
            <Button variant="outline" onclick={selectModFolder}>
              <FolderOpen class="w-4 h-4 mr-2" />
              Browse
            </Button>
          </div>
        </div>

        <!-- Root Folder -->
        <div class="space-y-2">
          <Label for="root-folder">Assetto Corsa Root Folder</Label>
          <div class="flex gap-2">
            <Input
              id="root-folder"
              type="text"
              placeholder="Select root folder..."
              value={assettoCorsaRootFolder}
              class="flex-1"
              readonly
            />
            <Button variant="outline" onclick={selectRootFolder}>
              <FolderOpen class="w-4 h-4 mr-2" />
              Browse
            </Button>
          </div>
        </div>
      </Card.Content>
    </Card.Root>

    <!-- Profile Settings -->
    <Card.Root>
      <Card.Header>
        <Card.Title>Profile Settings</Card.Title>
        <Card.Description>
          Configure automatic profile management
        </Card.Description>
      </Card.Header>
      <Card.Content>
        <div class="flex items-center space-x-3">
          <Checkbox
            id="automatic-profiles"
            checked={automaticProfiles}
            onCheckedChange={(checked) =>
              (automaticProfiles = checked as boolean)}
          />
          <Label
            for="automatic-profiles"
            class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
          >
            Enable automatic profiles
          </Label>
        </div>
      </Card.Content>
    </Card.Root>

    <!-- Actions -->
    <div class="flex gap-4">
      <Button onclick={saveSettings}>
        <Save class="w-4 h-4 mr-2" />
        Save Settings
      </Button>
    </div>

    <Separator />

    <!-- Danger Zone -->
    <Card.Root class="border-destructive">
      <Card.Header>
        <Card.Title class="text-destructive">Danger Zone</Card.Title>
        <Card.Description>
          Irreversible actions that can affect your data
        </Card.Description>
      </Card.Header>
      <Card.Content>
        <Dialog.Root bind:open={resetDialogOpen}>
          <Dialog.Trigger asChild>
            <Button variant="destructive">Reset App</Button>
          </Dialog.Trigger>
          <Dialog.Content>
            <Dialog.Header>
              <Dialog.Title>Reset Application</Dialog.Title>
              <Dialog.Description>
                This will delete all data from the database and reset all
                settings to default. This action cannot be undone.
              </Dialog.Description>
            </Dialog.Header>
            <Dialog.Footer>
              <Button
                variant="outline"
                onclick={() => (resetDialogOpen = false)}
              >
                Cancel
              </Button>
              <Button variant="destructive" onclick={resetApp}>Reset</Button>
            </Dialog.Footer>
          </Dialog.Content>
        </Dialog.Root>
      </Card.Content>
    </Card.Root>
  </div>
</div>
