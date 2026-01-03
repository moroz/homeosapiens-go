#!/usr/bin/env -S deno run --allow-write --allow-read

import path from "node:path";
import { fromFileUrl } from "https://deno.land/std@0.208.0/path/mod.ts";

const codes: string[] = [];

const isoCodeFilePath = path.join(
  path.dirname(fromFileUrl(Deno.mainModule)),
  "iso3166.tab",
);

const isoCodeFile = await Deno.readTextFile(isoCodeFilePath);
const lines = isoCodeFile.split("\n");

for (const line of lines) {
  if (line.startsWith("#")) continue;

  const code = line.split(/\s+/)[0];
  code && codes.push(code);
}

function formatCountry(locale: string, isoCode: string): string {
  return new Intl.DisplayNames(locale, { type: "region" }).of(isoCode)!;
}

const options = codes.map((iso) => ({
  value: iso,
  labelPl: formatCountry("pl", iso),
  labelEn: formatCountry("en-GB", iso),
}));

const json = JSON.stringify(options, null, 2);
const resolved = path.resolve(
  path.join(
    path.dirname(fromFileUrl(Deno.mainModule)),
    "../internal/countries/countries.json",
  ),
);

await Deno.writeTextFile(resolved, json);
