#!/usr/bin/env -S deno run --allow-read --allow-write --allow-sys --allow-env --allow-run

import { chromium } from "npm:playwright";

const browser = await chromium.launch();
browser.newPage()
