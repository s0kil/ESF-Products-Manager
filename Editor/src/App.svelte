<script>
  import {onMount} from "svelte";
  import loadAsset from "loadjs";
  import API from '@suddenly/api';

  const api = new API("products", "/");

  // Initialize
  loadAsset(
    [
      "https://cdn.jsdelivr.net/npm/jsuites@2.7.1/dist/jsuites.min.js",
      "https://cdn.jsdelivr.net/npm/jexcel@3.9.1/dist/jexcel.min.js"
    ],
    "assets",
    {numRetries: 7}
  );

  onMount(async () => {
    let products = await api.get("/product").then(async r => await r);
    let productsArray = products.map(item => Object.values(item))

    loadAsset.ready("assets", {
      success: () => Main(productsArray),
      error: () => {
        throw new Error("Failed To Load Required Assets")
      }
    });
  });

  function Main(products) {
    let container = document.getElementById("spreadsheet");
    let containerStyles = getComputedStyle(container)
    let containerSize = {
      // TODO Refactor The Following ???
      width: screen.width,
      height: screen.height - 200
    };

    const spreadSheetColumns = [
      {
        type: "text",
        title: "Title",
      }
    ]

    const spreadSheet = jexcel(container, {
      data: products,
      defaultColWidth: 100,
      tableOverflow: true,
      tableWidth: `${containerSize.width}px`,
      tableHeight: `${containerSize.height}px`,
      columns: spreadSheetColumns,
      locked: true
    });
  }

  $: spreadSheetState = {
    locked: false,
    sheets: [
      "products"
    ]
  }

  const Icon = iconName => `https://cdn.jsdelivr.net/npm/feather-icons@4.26.0/dist/icons/${iconName}.svg`
</script>

<style>
  main {
    width: 100vw;
    height: 100vh;
    overflow: hidden;
  }

  header {
    font-family: "Stalinist One", serif;
  }

  header h1 {
    letter-spacing: 2px;
  }

  #spreadsheet {
  }

  .naive-lock {
    user-select: none;
    pointer-events: none;
  }

  .icon {
    width: 30px;
    height: 30px;
    cursor: pointer;
  }
</style>

<svelte:head>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/modern-normalize@0.6.0/modern-normalize.min.css"/>
  <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Stalinist+One&display=swap"/>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/colors.css@3.0.0/css/colors.min.css"/>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/jsuites@2.7.1/dist/jsuites.min.css"/>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/basscss@8.1.0/css/basscss.min.css"/>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/jexcel@3.9.1/dist/jexcel.min.css"/>
</svelte:head>

<main class="bg-silver">
  <header>
    <div class="flex">
      <div class="mr4">
        <h1 class="p1 m0 h1 black line-height-1">ESF Products Manager</h1>
      </div>
      <div class="flex items-center">
        <div class="lock" on:click="{() => spreadSheetState.locked = !spreadSheetState.locked}">
            {#if spreadSheetState.locked === true}
              <img class="icon block mx-auto" src={Icon("lock")} alt="Lock Spreadsheet">
            {:else}
              <img class="icon block mx-auto" src={Icon("unlock")} alt="Unlock Spreadsheet">
            {/if}
        </div>
      </div>
    </div>
  </header>

  <div id="spreadsheet" class="bg-white" class:naive-lock="{spreadSheetState.locked === true}"></div>

  <nav>
    <div class="sheets">
        {#each spreadSheetState.sheets as sheet}
          <button class="render-sheet btn not-rounded">{sheet}</button>
        {/each}
    </div>
  </nav>
</main>
