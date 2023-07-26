Promise.allSettled([
    customElements.whenDefined("sl-button"),
]).then(() => {
    document.body.classList.add("ready");
});
