#!/usr/bin/env -S deno run --allow-read --allow-write --allow-sys --allow-env

import { chromium } from "npm:playwright";

chromium.launch();
