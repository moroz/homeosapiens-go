<script lang="ts">
  import { Temporal } from "@js-temporal/polyfill";

  const pp =
    "https://d3n1g0yg3ja4p3.cloudfront.net/019beef9-ad4c-736f-9bb0-965b59ca21ae.png";
  const locale = "en";

  const host = "Dr Asher Shaikh";
  const title = "What prevents me from moving on?";
  const date = Temporal.PlainDate.from("2026-02-08");
</script>

<div class="w-full h-screen grid place-items-center">
  <div class="aspect-video bg-slate-100 w-160 outline-2 flex p-8 gap-8">
    <div class="h-full aspect-3/4 outline flex overflow-hidden shrink-0">
      <img src={pp} class="w-full h-full object-cover" alt="" />
    </div>
    <div class="flex flex-col gap-4">
      <h1 class="font-bold text-3xl">{host}</h1>
      <h2 class="font-bold text-4xl leading-tight">{title}</h2>
      <h3>{date.toLocaleString(locale, { dateStyle: "long" })}</h3>
    </div>
  </div>
</div>
