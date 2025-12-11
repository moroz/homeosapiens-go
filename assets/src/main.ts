import "./app.css";

import videojs from "video.js";

document.querySelectorAll(".video-js").forEach((el) => {
	videojs(el, { controls: true, fluid: true });
});

document.querySelectorAll(".user-dropdown").forEach((el) => {
	const dropdown = el.querySelector(".dropdown");
	if (!dropdown) return;
	let timeout: number;

	el.addEventListener("mouseover", () => {
		clearTimeout(timeout);
		dropdown.classList.remove("hidden");
	});

	el.addEventListener("mouseout", () => {
		timeout = setTimeout(() => {
			dropdown.classList.add("hidden");
		}, 500);
	});
});
