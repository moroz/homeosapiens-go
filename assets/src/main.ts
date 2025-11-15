import "./app.css";

import videojs from "video.js";

document.querySelectorAll(".video-js video").forEach((el) => {
	videojs(el, { controls: true, fill: true });
});
