<script lang="ts">
  import { type Affiliation, fetchAffiliations } from "$lib/api";
  import { onMount } from "svelte";

  let hoveredCard: number = $state(0);
  let affiliations: Affiliation[] = $state([]);
  let imagesLoaded = $state(0);
  let imagesReady = $state(false);
  let loadedMap = new Map<string, boolean>();

  const preloadImages = (
    urls: string[],
    timeout: number = 0
  ) => {
    return Promise.all(
      urls.map(
        (url) =>
          new Promise<void>(async (resolve) => {
            const img = new Image();
            img.src = url;
            img.onload = async () => {
              loadedMap.set(url, true);
              imagesLoaded += 1;
              await new Promise((r) => setTimeout(r, timeout)); // fake delay
              resolve();
            };
            img.onerror = async () => {
              loadedMap.set(url, false);
              imagesLoaded += 1;
              await new Promise((r) => setTimeout(r, timeout)); // fake delay
              resolve();
            };
          })
      )
    );
  };
  
  onMount(() => {
    fetchAffiliations().then((affiliationsResp) => {
      affiliations = affiliationsResp.affiliations;
      preloadImages(
        affiliations.map((p) => p.logoURL),
      ).then(() => {
        setTimeout(() => {
          imagesReady = true;
        }, 500);
      });
    });
  })
</script>

<section id="affiliations" class="w-full max-w-7xl mx-auto p-8 text-text">
    <div class="w-full flex items-start">
      <span class="text-sm mr-2 text-footnote">[3]</span>
      <h1 class="font-bold text-2xl">Affiliations</h1>
    </div>

    <div class="flex mt-8">
      {#if !imagesReady && affiliations.length > 0}
        <div class="w-full h-full flex items-center justify-center">
          <div class="grid grid-cols-4 gap-1">
            {#each affiliations as _, i}
              <div
                class="flip-card {i < imagesLoaded
                  ? 'loaded'
                  : ''}"
              >
                <div class="flip-card-inner">
                  <div class="flip-card-front"></div>
                  <div class="flip-card-back"></div>
                </div>
              </div>
            {/each}
          </div>
        </div>
      {:else}
        <div class="grid grid-cols-4 gap-4">
          {#each affiliations as affiliation}
            <div
              class="group w-full aspect-square hover:p-4 p-0 overflow-hidden duration-200 rounded-xl bg-inner"
            >
              <div
                class="w-full aspect-square group-hover:aspect-6/5 duration-200 rounded-none group-hover:rounded-lg flex items-center justify-center text-footnote text-sm overflow-hidden"
              >
                <img
                  src={affiliation.logoURL}
                  alt={affiliation.title}
                  class="w-full h-full object-cover"
                />
              </div>
              <p class="mt-5 font-medium">{affiliation.title}</p>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  </section>