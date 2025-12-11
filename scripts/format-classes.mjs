#!/usr/bin/env -S deno run --allow-write --allow-sys --allow-env --allow-read
// tools/sort-tailwind.js
// Sorts Tailwind classes inside gomponents Class("...") calls using prettier-plugin-tailwindcss.
import fs from "node:fs";
import path from "node:path";
import prettier from "npm:prettier";
import * as pluginTailwind from "npm:prettier-plugin-tailwindcss";

const classCallRx = /Class\("([^"]+)"\)/g;

async function sortOne(classStr) {
  // Feed into Prettier so the plugin reorders.
  const fake = `<div class="${classStr}"></div>`;
  const formatted = await prettier.format(fake, {
    parser: "html",
    plugins: [pluginTailwind],
  });
  // Extract reordered classes.
  const m = formatted.match(/class="([^"]*)"/);
  return m ? m[1].trim() : classStr;
}

async function processFile(file) {
  const src = fs.readFileSync(file, "utf8");
  let changed = false;
  const out = await replaceAsync(src, classCallRx, async (full, inner) => {
    const sorted = await sortOne(inner);
    if (sorted !== inner) changed = true;
    return `Class("${sorted}")`;
  });
  if (changed) {
    fs.writeFileSync(file, out);
    console.log(`Sorted: ${file}`);
  } else {
    console.log(`No changes: ${file}`);
  }
}

async function replaceAsync(str, regex, asyncFn) {
  const promises = [];
  str.replace(regex, (match, ...groups) => {
    const promise = asyncFn(match, ...groups);
    promises.push(promise);
    return match;
  });
  const results = await Promise.all(promises);
  let i = 0;
  return str.replace(regex, () => results[i++]);
}

async function main() {
  const files = process.argv.slice(2);
  if (!files.length) {
    console.error("usage: node tools/sort-tailwind.js `file.go` [...]");
    process.exit(1);
  }
  for (const f of files) {
    await processFile(path.resolve(f));
  }
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
