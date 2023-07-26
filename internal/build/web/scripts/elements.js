Promise.allSettled(
    ["sl-button", "sl-input", "sl-select", "sl-option", "sl-checkbox", "sl-badge"]
        .map(e => customElements.whenDefined(e))
).then(() => document.body.classList.add("ready"));
