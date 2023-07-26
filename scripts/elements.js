Promise.allSettled([
    "sl-button",
    "sl-input",
    "sl-select",
    "sl-option",
    "sl-checkbox",
    "sl-badge"
].map(tag => customElements.whenDefined(tag))).then(() => {
    document.body.classList.add("ready");
});
