<script lang="ts">
  import "./layout.css";

  let { children } = $props();
</script>

<svelte:head>
  <title>Dane Walker</title>
  <!-- favicons -->
  <link rel="icon" href="/favicon.ico" />
  <link rel="shortcut icon" href="/favicon.ico" />
  <link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png" />
  <link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png" />
  <link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png" />
  <!-- android -->
  <link
    rel="icon"
    type="image/png"
    sizes="192x192"
    href="/android-chrome-192x192.png"
  />
  <link
    rel="icon"
    type="image/png"
    sizes="512x512"
    href="/android-chrome-512x512.png"
  />
  <!-- manifest -->
  <link rel="manifest" href="/site.webmanifest" />
  <script
    defer
    src="https://danecwalker.com/uevents.js"
    data-website-id="d0a6f2a0-155c-43b2-9498-6696a9346b69"
  ></script>
</svelte:head>

{@render children()}
