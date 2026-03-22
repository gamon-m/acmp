<script lang="ts">
  import { NavigationMenu as NavigationMenuPrimitive } from "bits-ui";
  import { cn } from "$lib/utils.js";

  let {
    ref = $bindable(null),
    class: className,
    ...restProps
  }: NavigationMenuPrimitive.LinkProps = $props();

  const baseStyles =
    "hover:text-muted-foreground px-8 text-2xl font-bold transition-colors";
  const activeStyles =
    "data-[active=true]:text-primary data-[active=true]:underline data-[active=true]:underline-offset-8 data-[active=true]:decoration-3";
</script>

<NavigationMenuPrimitive.Link
  bind:ref
  data-slot="navigation-menu-link"
  class={cn(baseStyles, activeStyles, className)}
  {...restProps}
/>
