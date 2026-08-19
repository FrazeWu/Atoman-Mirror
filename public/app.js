const servicePages = [...document.querySelectorAll("[data-page]")];
const serviceLinks = [...document.querySelectorAll("[data-service]")];
const serviceIds = new Set(servicePages.map((page) => page.dataset.page));
const defaultService = "docker-hub";

function currentService() {
	const service = window.location.hash.slice(1);
	return serviceIds.has(service) ? service : defaultService;
}

function renderService() {
	const service = currentService();

	servicePages.forEach((page) => {
		page.hidden = page.dataset.page !== service;
	});

	serviceLinks.forEach((link) => {
		const active = link.dataset.service === service;
		link.toggleAttribute("aria-current", active);
	});

	document.title = `${serviceLinks.find((link) => link.dataset.service === service).textContent} | Atoman Mirror`;
}

function copyCommand(button) {
	const code = button.parentElement.querySelector("code").textContent.trim();
	return navigator.clipboard.writeText(code);
}

serviceLinks.forEach((link) => {
	link.addEventListener("click", () => {
		window.setTimeout(renderService, 0);
	});
});

document.querySelectorAll(".copy").forEach((button) => {
	button.addEventListener("click", async () => {
		try {
			await copyCommand(button);
			button.textContent = "已复制";
			button.dataset.copied = "true";
			window.setTimeout(() => {
				button.textContent = "复制";
				delete button.dataset.copied;
			}, 1600);
		} catch {
			button.textContent = "复制失败";
			window.setTimeout(() => {
				button.textContent = "复制";
			}, 1600);
		}
	});
});

window.addEventListener("hashchange", renderService);

if (!window.location.hash || !serviceIds.has(window.location.hash.slice(1))) {
	window.history.replaceState(null, "", `#${defaultService}`);
}

renderService();
