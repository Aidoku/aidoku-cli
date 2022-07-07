/** @type {import('tailwindcss').Config} */
module.exports = {
	content: ["./src/**/*.{html,js,svelte,ts}", "./index.html"],
	theme: {
		extend: {
			colors: {
				"aidoku-color": "#ff375f",
				"light-tint": "rgba(255,55,95,.06)",
			},
			backgroundImage: {
				"button-gradient": "linear-gradient(rgba(255, 55, 95, 0.02) 0 0)",
			},
			keyframes: {
				fade: {
					from: { opacity: 1 },
					to: { opacity: 0.25 },
				},
			},
			animation: {
				fade: "fade 1s linear infinite",
			},
		},
	},
	plugins: [],
};
