<script lang="ts">
  import { fetchLinks, type Link } from "$lib/api";
  import Affiliation from "$lib/partials/Affiliation.svelte";
  import Experience from "$lib/partials/Experience.svelte";
  import Projects from "$lib/partials/Projects.svelte";
  import { onMount } from "svelte";

  let links: Link[] = $state([]);

  onMount(() => {
    fetchLinks().then((linksResp) => {
      links = linksResp.links;
    });
  });
</script>

<div class="w-full min-h-screen bg-background">
  <div class="p-8 max-w-7xl mx-auto flex flex-col md:flex-row justify-between items-center">
    <div class="font-medium text-text">Dane Walker</div>
    <nav class="mt-4 md:mt-0">
      <ul class="flex items-center gap-8">
        <li>
          <a href="#projects" class="text-link hover:underline">Projects</a>
        </li>
        <li>
          <a href="#experience" class="text-link hover:underline">Experience</a>
        </li>
        <li>
          <a href="#affiliations" class="text-link hover:underline"
            >Affiliations</a
          >
        </li>
        <li>
          <a href="#info" class="text-link hover:underline">Info</a>
        </li>
      </ul>
    </nav>
  </div>

  <Projects />

  <Experience />

  <Affiliation />

  <footer id="info">
    <div
      class="w-full border-t border-border mt-16 text-center border-inner-border text-sm text-footnote"
    >
      <div class="w-full max-w-7xl mx-auto flex justify-between items-start md:items-center p-8">
        <div class="w-full flex items-start">
          <span class="text-sm mr-2 text-footnote">[4]</span>
          <h1 class="font-bold text-2xl text-text">Info</h1>
        </div>
        <ul
          class="flex flex-col md:flex-row md:items-center items-end md:gap-8 gap-2 w-full justify-end"
        >
          {#each links as link}
            <li>
              {#if link.type === "file"}
                <a
                  href={link.url}
                  class="hover:text-link hover:underline duration-100"
                  target="_blank"
                  rel="noopener noreferrer"
                  data-umami-event="file-link-click"
                  data-umami-event-file={link.title}
                  >{link.title}
                </a>
              {:else}
                <a
                  href={link.url}
                  class="hover:text-link hover:underline duration-100"
                  data-umami-event="outbound-link-click"
                  data-umami-event-url={link.url}>{link.title}</a
                >
              {/if}
            </li>
          {/each}
        </ul>
      </div>
    </div>
  </footer>
</div>
